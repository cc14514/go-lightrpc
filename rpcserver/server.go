package rpcserver

import (
	"encoding/json"
	"fmt"
	"github.com/rs/cors"
	"log"
	"net/http"
	"reflect"
)

var (
	//仅在判断参数类型时使用
	_token TOKEN
	logger *log.Logger
)

func init() {
	logger, _ = NewLogger("/tmp", "server.log", "eng")
}

type Rpcserver struct {
	Port            int
	ServiceMap      map[string]ServiceReg
	CheckToken      func(token TOKEN) bool
	AllowedMethods  []string
	LogDir, LogFile string
}

func (self *Rpcserver) StartServer() (e error) {
	if len(self.LogDir) > 0 && len(self.LogFile) > 0 {
		logger, _ = NewLogger(self.LogDir, self.LogFile, "eng")
	}
	rpcport := self.Port
	if len(self.AllowedMethods) == 0 {
		self.AllowedMethods = []string{"POST", "GET"}
	}
	logger.Print("StartServer port -->", rpcport, " ; allow_method -->", self.AllowedMethods)
	for k, _ := range self.ServiceMap {
		logger.Print("service_reg_map : ", k)
	}
	c := cors.New(cors.Options{
		AllowedMethods: self.AllowedMethods,
	})
	mux := http.NewServeMux()
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		var (
			success *Success = &Success{Success: true}
			ibody   interface{}
		)
		r.ParseForm()
		logger.Print("url =", r.URL, "method =", r.Method)
		logger.Print("r.Form", r.Form)
		logger.Print("r.PostForm", r.PostForm)
		logger.Print("r.URL.Query", r.URL.Query())
		body := r.FormValue("body")
		logger.Print("body =", body, []byte(body))
		if len(body) == 0 {
			success.Success = false
		} else if err := json.Unmarshal([]byte(body), &ibody); err == nil {
			mbody := ibody.(map[string]interface{})
			logger.Print("mbody", mbody)
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
				serviceReg, ok := self.ServiceMap[s]
				if ok {
					serviceObj := serviceReg.Service
					refService := reflect.ValueOf(serviceObj)
					refMethod := refService.MethodByName(m)
					logger.Print("refService=", refService)
					logger.Print("refMethod=", refMethod)
					auth := false
					if refMethod.IsValid() {
						//TODO check method input/output of define
						rmt := refMethod.Type()
						inArr := make([]reflect.Value, rmt.NumIn())
						for i := 0; i < rmt.NumIn(); i++ {
							in := rmt.In(i)
							var _token TOKEN
							logger.Print("in:", in)
							if in == reflect.TypeOf(_token) {
								logger.Print("TODO: AuthFilter")
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
							logger.Print("rtn =", rtn)
							success = &rtn
						}
						if auth {
							if self.CheckToken(token) {
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
			logger.Print("mbody_err =", err)
			success.Error("1001", fmt.Sprintf("Params format error ;;; %s", err))
		}
		success.ResponseAsJson(w)
	})
	h := c.Handler(mux)
	host := fmt.Sprintf(":%d", rpcport)
	logger.Print("host = ", host)
	err := http.ListenAndServe(host, h)
	logger.Print(err)
	return e
}
