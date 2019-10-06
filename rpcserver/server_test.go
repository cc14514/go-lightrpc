package rpcserver

import (
	"testing"

	"github.com/tidwall/gjson"
)

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
