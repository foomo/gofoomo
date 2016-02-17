package handler

import (
	"log"
	"testing"

	"github.com/foomo/gofoomo/rpc"
)

type TestService struct {
}

func NewTestService() *TestService {
	t := new(TestService)
	return t
}

func getTestRPC() *RPC {
	return NewRPC(NewTestService(), "/services/test.php")
}

func (t *TestService) Test(arguments []string, argumentMap map[string]interface{}, reply *rpc.MethodReply) {
	reply.Value = true
}

func TestHandlesMethod(t *testing.T) {
	r := getTestRPC()
	if r.handlesMethod("Test") == false {
		t.Fail()
	}
	if r.handlesMethod("testi") == true {
		t.Fail()
	}
}

func TestGetApplicationPath(t *testing.T) {
	p := getTestRPC().getApplicationPath("/services/test.php/Foomo.Services.RPC/serve/test")
	if p != "test" {
		t.Fatal("i do not like this path", p)
	}
}

func TestHandlesPath(t *testing.T) {
	r := getTestRPC()
	if r.handlesPath("/services/test.php/Foomo.Services.RPC/serve/test") == false {
		t.Fatal("/services/test.php/Foomo.Services.RPC/serve/test")
	}
	if r.handlesPath("/services/test.php/Foomo.Services.RPC/serve/test/foo") == false {
		t.Fatal("/services/test.php/Foomo.Services.RPC/serve/test/foo")
	}
	if r.handlesPath("/services/test.php/Foomo.Services.RPC/serve/testi/foo") == true {
		t.Fatal("/services/test.php/Foomo.Services.RPC/serve/testi/foo")
	}
}

func TestExtractArguments(t *testing.T) {
	r := getTestRPC()
	args := r.extractArguments("/services/test.php/Foomo.Services.RPC/serve/test/%C3%BCb%C3%A4l/B%C3%A4r")
	if len(args) != 2 {
		t.Fatal("wrong args length", args)
	}
	if args[0] != "übäl" {
		t.Fatal("no übäl")
	}
	if args[1] != "Bär" {
		t.Fatal("where is the bear")
	}
}

func TestCallServiceObject(t *testing.T) {
	r := getTestRPC()
	var argumentMap map[string]interface{}
	var arguments []string
	reply := &rpc.MethodReply{}
	r.callServiceObject("Test", arguments, argumentMap, reply)
	if reply.Value != true {
		log.Println(reply)
		t.Fail()
	}
}
