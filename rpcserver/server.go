package rpcserver

import (
	"encoding/json"
	"fmt"
	"github.com/alecthomas/log4go"
	"github.com/rs/cors"
	"net/http"
	"reflect"
)

var (
	//仅在判断参数类型时使用
	_token TOKEN
	this   *Rpcserver
)

type Rpcserver struct {
	Port           int
	ServiceMap     map[string]ServiceReg
	CheckToken     func(token TOKEN) bool
	AllowedMethods []string
}

func (self *Rpcserver) makeCors() *cors.Cors {
	log4go.Debug("StartServer port->%s ; allow_method->%s", self.Port, self.AllowedMethods)
	return cors.New(cors.Options{
		AllowedMethods: self.AllowedMethods,
	})
}

func (self *Rpcserver) StartServer() (e error) {
	this = self
	logServiceMap(self.ServiceMap)
	if len(self.AllowedMethods) == 0 {
		self.AllowedMethods = []string{"POST", "GET"}
	}
	c := self.makeCors()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/", handlerFunc)
	h := c.Handler(mux)
	host := fmt.Sprintf(":%d", self.Port)
	log4go.Debug("host = %s", host)
	http.ListenAndServe(host, h)
	return e
}

func logServiceMap(m map[string]ServiceReg) {
	for k, _ := range m {
		log4go.Debug("<<<< service_reg_map >>>> : %s", k)
	}
	if m == nil {
		log4go.Debug("<<<< service_reg_map_empty >>>>")
	}
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	var (
		success *Success = &Success{Success: true}
		ibody   interface{}
	)
	r.ParseForm()
	body := r.FormValue("body")
	if len(body) == 0 {
		success.Success = false
	} else if err := json.Unmarshal([]byte(body), &ibody); err == nil {
		mbody := ibody.(map[string]interface{})
		log4go.Debug("mbody = %s", mbody)
		var token TOKEN = ""
		if tokenParam, ok := mbody["token"]; ok {
			token = TOKEN(tokenParam.(string))
		}
		if service, ok := mbody["service"]; !ok {
			success.Error("1002", "service error or notfound")
		} else if method, ok := mbody["method"]; !ok {
			success.Error("1002", "method error or notfound")
		} else {
			s := service.(string)
			m := method.(string)
			//TODO check service version
			serviceReg, ok := this.ServiceMap[s]
			if ok {
				serviceObj := serviceReg.Service
				refService := reflect.ValueOf(serviceObj)
				refMethod := refService.MethodByName(m)
				log4go.Debug("refService = %s", refService)
				log4go.Debug("refMethod = %s", refMethod)
				auth := false
				if refMethod.IsValid() {
					//TODO check method input/output of define
					rmt := refMethod.Type()
					inArr := make([]reflect.Value, rmt.NumIn())
					for i := 0; i < rmt.NumIn(); i++ {
						in := rmt.In(i)
						var _token TOKEN
						log4go.Debug("in = %s", in)
						if in == reflect.TypeOf(_token) {
							log4go.Debug("TODO: AuthFilter ========>")
							inArr[i] = reflect.ValueOf(token)
							auth = true
						} else {
							if params, ok := mbody["params"]; ok {
								inArr[i] = reflect.ValueOf(params)
							}
						}
					}
					runservice := func() {
						rtn := refMethod.Call(inArr)[0].Interface().(Success)
						log4go.Debug("rtn =", rtn)
						success = &rtn
					}
					if auth {
						if this.CheckToken(token) {
							runservice()
						} else {
							success.Error("1003", fmt.Sprintf("error token"))
						}
					} else {
						runservice()
					}
				} else {
					success.Error("1002", fmt.Sprintf("method notfond ;;; %s", m))
				}
			}
		}
	} else {
		log4go.Debug("mbody_err =", err)
		success.Error("1001", fmt.Sprintf("Params format error ;;; %s", err))
	}
	success.ResponseAsJson(w)
}
