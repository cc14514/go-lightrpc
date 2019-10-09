package rpcserver

import (
	"encoding/json"
	"fmt"
	"github.com/rs/cors"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
)

var (
	//仅在判断参数类型时使用
	_token TOKEN
	this   *Rpcserver
	server *http.Server
)

type Rpcserver struct {
	// url , 默认 /api/
	Pattern        string
	Port           int
	ServiceMap     map[string]ServiceReg
	CheckToken     func(token TOKEN) bool
	AllowedMethods []string
	ServeMux       *http.ServeMux
}

func (self *Rpcserver) makeCors() *cors.Cors {
	log.Println(fmt.Sprintf("StartDefaultServer port->%d ; allow_method->%s", self.Port, self.AllowedMethods))
	return cors.Default()
	/*return cors.New(cors.Options{
		AllowedMethods: self.AllowedMethods,
	})*/
}

func (self *Rpcserver) StopServer() (e error) {
	if server != nil {
		e = server.Shutdown(nil)
		log.Println(fmt.Sprintf("<%v> lightrpc-server-shutdown , see you ...", e))
	}
	return e
}

func (self *Rpcserver) StartServer() (e error) {
	this = self
	logServiceMap(self.ServiceMap)
	if len(self.AllowedMethods) == 0 {
		self.AllowedMethods = []string{"POST", "GET"}
	}
	c := self.makeCors()
	if self.ServeMux == nil {
		self.ServeMux = http.NewServeMux()
	}
	if self.Pattern != "" {
		self.ServeMux.HandleFunc(self.Pattern, handlerFunc)
	} else {
		self.ServeMux.HandleFunc("/api/", handlerFunc)
	}
	self.ServeMux.HandleFunc("/", handlerFunc)
	h := c.Handler(self.ServeMux)
	host := fmt.Sprintf(":%d", self.Port)
	log.Println("host =", host)
	//http.ListenAndServe(host, h)
	//server := &http.Server{Addr: host, Handler: h}
	//server.ListenAndServe()
	http.ListenAndServe(host, h)

	return e
}

func logServiceMap(m map[string]ServiceReg) {
	for k, _ := range m {
		log.Println("<<<< service_reg_map >>>> :", k)
	}
	if m == nil {
		log.Println("<<<< service_reg_map_empty >>>>")
	}
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	success := &Success{Success: true}
	var body = ""
	switch r.Method {
	case http.MethodGet:
		r.ParseForm()
		body = r.FormValue("body")
	case http.MethodPost:
		d, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("getbody_err", err)
		} else {
			body = string(d)
		}
	}

	if body == "" {
		success.Success = false
		success.Error("1000", "params of body notfound")
	} else {
		if serviceRes := gjson.Get(body, "service"); serviceRes.String() == "null" {
			success.Error("1002", "service error or notfound")
		} else if methodRes := gjson.Get(body, "method"); methodRes.String() == "null" {
			success.Error("1002", "method error or notfound")
		} else {
			s := serviceRes.String()
			//TODO check service version
			serviceReg, ok := this.ServiceMap[s]
			if ok {
				executeMethod(serviceReg, body, success)
			} else {
				success.Success = false
				success.Error("1003", "service not reg")
			}
		}
	}
	success.ResponseAsJson(w)
}

func executeMethod(serviceReg ServiceReg, body string, success *Success) {
	token := getToken(body)
	sn := gjson.Get(body, "sn").String()
	methodName := gjson.Get(body, "method").String()
	methodName = paserMethodName(methodName)
	serviceObj := serviceReg.Service
	refService := reflect.ValueOf(serviceObj)
	refMethod := refService.MethodByName(methodName)
	auth := false
	if refMethod.IsValid() {
		rmt := refMethod.Type()
		inArr := make([]reflect.Value, rmt.NumIn())
		for i := 0; i < rmt.NumIn(); i++ {
			in := rmt.In(i)
			var _token TOKEN
			if in == reflect.TypeOf(_token) {
				log.Println("TODO: AuthFilter ========>")
				inArr[i] = reflect.ValueOf(token)
				auth = true
			} else if kind := in.Kind().String(); kind == "interface" || kind == "map" {
				//interface{} 和 map[string]interface{} 就是序列化成 map 传递
				target := map[string]interface{}{}
				if paramsRes := gjson.Get(body, "params"); paramsRes.String() != "null" {
					json.Unmarshal([]byte(paramsRes.String()), &target)
				}
				inArr[i] = reflect.ValueOf(target)
			} else if kind == "string" {
				// 字符串则直接传递 json 字符串
				if paramsRes := gjson.Get(body, "params"); paramsRes.String() != "null" {
					inArr[i] = reflect.ValueOf(paramsRes.String())
				} else {
					inArr[i] = reflect.ValueOf("")
				}
			} else {
				//TODO 2016-12-06 : 非常遗憾，当前版本还不能支持此功能
				// 否则反射成 in 指定的 struct 类型
				//if paramsRes := gjson.Get(body, "params"); paramsRes.String() != "null" {
				//	inVal := reflect.New(in).Interface()
				//	json.Unmarshal([]byte(paramsRes.String()), inVal)
				//	inArr[i] = reflect.ValueOf(&inVal)
				//}
				success.Success = false
				success.Error("TODO", "not support struct yet")
				return
			}
		}
		runservice := func() {
			rtn := refMethod.Call(inArr)[0].Interface().(Success)
			success.Sn = sn
			success.Success = rtn.Success
			success.Entity = rtn.Entity
		}
		if auth {
			if this.CheckToken == nil || this.CheckToken(token) {
				runservice()
			} else {
				success.Error("1003", fmt.Sprintf("error token"))
			}
		} else {
			runservice()
		}
	} else {
		success.Error("1002", fmt.Sprintf("method notfond ;;; %s", methodName))
	}
}

func paserMethodName(s string) string {
	b0 := s[0]
	s0 := string(b0)
	s0 = strings.ToUpper(s0)
	sn := s[1:len(s)]
	s = s0 + sn
	return s
}

func getToken(body string) TOKEN {
	var token TOKEN
	if tokenRes := gjson.Get(body, "token"); tokenRes.String() != "null" {
		token = TOKEN(tokenRes.String())
	}
	return token
}
