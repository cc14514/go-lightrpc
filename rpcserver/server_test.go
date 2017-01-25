package rpcserver

import (
	"testing"

	"github.com/alecthomas/log4go"
	"github.com/tidwall/gjson"
)

func TestLog4go(t *testing.T) {
	logServiceMap(nil)
	//log4go.AddFilter("stdout", log4go.DEBUG, log4go.NewConsoleLogWriter())                  //输出到控制台,级别为DEBUG
	//log4go.AddFilter("file", log4go.DEBUG, log4go.NewFileLogWriter("/tmp/test.log", false)) //输出到文件,级别为DEBUG,文件名为test.log,每次追加该原文件
	//log4go.LoadConfiguration("log.xml")//使用加载配置文件,类似与java的log4j.propertites
	log4go.Debug(">>>>>>>> %s -- %s", "213", "sad")
	dict := map[string]struct{}{"foo": struct{}{}, "bar": struct{}{}, "hello": struct{}{}}
	for itm := range dict {
		log4go.Debug(itm)
	}
	defer log4go.Close() //注:如果不是一直运行的程序,请加上这句话,否则主线程结束后,也不会输出和log到日志文件
}

func TestPaserMethodName(t *testing.T) {
	s := paserMethodName("getUserinfo")
	t.Log(s)
}

func TestGjson(t *testing.T) {
	j := []byte(`{"token":"123456","service":"user","method":"login","params":{"username":"foo","password":"123123"}}`)
	r := gjson.GetBytes(j, "params")
	t.Log("r.Value()=", r.Value())
	t.Log("r.Map()=", r.Map())
	t.Log("r.String()=", r.String())
	if gjson.GetBytes(j, "hello").String() == "null" {
		t.Log("ok : empty value (Result.String()) is \"null\" str")
	}
	r = gjson.Get("hello world", "foo")
	t.Log("r =", r)
}

type UserService struct{}
type UserVo struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (self *UserService) LoginInterface(vo interface{}) Success {
	return Success{
		Success: true,
		Entity:  vo,
	}
}
func (self *UserService) LoginMap(vo map[string]interface{}) Success {
	return Success{
		Success: true,
		Entity:  vo,
	}
}
func (self *UserService) LoginString(vo string) Success {
	return Success{
		Success: true,
		Entity:  vo,
	}
}
func (self *UserService) LoginStruct(vo UserVo) Success {
	return Success{
		Success: true,
		Entity:  vo,
	}
}

var (
	serviceReg ServiceReg = ServiceReg{
		Namespace: "user",
		Version:   "0.0.1",
		Service:   &UserService{},
	}
)

func Test_InputMap(t *testing.T) {
	body := `{"token":"123456","service":"user","method":"LoginMap","params":{"username":"foo","password":"123123"}}`
	success := &Success{}
	executeMethod(serviceReg, body, success)
	t.Log(success)
}
func Test_InputInterface(t *testing.T) {
	body := `{"token":"123456","service":"user","method":"LoginInterface","params":{"username":"foo","password":"123123"}}`
	success := &Success{}
	executeMethod(serviceReg, body, success)
	t.Log(success)
}
func Test_InputString(t *testing.T) {
	body := `{"token":"123456","service":"user","method":"LoginString","params":{"username":"foo","password":"123123"}}`
	success := &Success{}
	executeMethod(serviceReg, body, success)
	t.Log(success)
}
func Test_InputStruct(t *testing.T) {
	body := `{"token":"123456","service":"user","method":"LoginStruct","params":{"username":"foo","password":"123123"}}`
	success := &Success{}
	executeMethod(serviceReg, body, success)
	t.Log(body)
	t.Log(success)

}
