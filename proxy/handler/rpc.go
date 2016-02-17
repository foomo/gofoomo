package handler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/foomo/gofoomo/proxy/utils"
	"github.com/foomo/gofoomo/rpc"
)

// RPC helps you to hijack foomo rpc services. Actually it is even
// better, you can hijack them method by method.
//
// 	f := gofoomo.NewFoomo("/var/www/myApp", "test")
//	p := proxy.NewProxy(f, "http://test.myapp")
//	service := NewFooService()
//	rpcHandler := handler.NewRPC(service, "/foomo/modules/MyModule/services/foo.php")
//	f.AddHandler(rpcHandler)
//
// Happy service hijacking!
type RPC struct {
	path          string
	serviceObject interface{}
}

// NewRPC rpc constructor, path is the path in the url, that you intend to hijack
func NewRPC(serviceObject interface{}, path string) *RPC {
	rpc := new(RPC)
	rpc.path = path
	rpc.serviceObject = serviceObject
	return rpc
}

func (r *RPC) getApplicationPath(path string) string {
	return path[len(r.path+"/Foomo.Services.RPC/serve")+1:]
}

func (r *RPC) getMethodFromPath(path string) string {
	parts := strings.Split(r.getApplicationPath(path), "/")
	if len(parts) > 0 {
		return strings.ToUpper(parts[0][0:1]) + parts[0][1:]
	}
	return ""
}

func (r *RPC) handlesMethod(methodName string) bool {
	return reflect.ValueOf(r.serviceObject).MethodByName(methodName).IsValid()
}

func (r *RPC) handlesPath(path string) bool {
	return strings.HasPrefix(path, r.path) && r.handlesMethod(r.getMethodFromPath(path))
}

// HandlesRequest implementation of request handler interface
func (r *RPC) HandlesRequest(incomingRequest *http.Request) bool {
	return incomingRequest.Method == "POST" && r.handlesPath(incomingRequest.URL.Path)
}

func (r *RPC) callServiceObjectWithHTTPRequest(incomingRequest *http.Request) (reply *rpc.MethodReply) {
	reply = &rpc.MethodReply{}
	path := incomingRequest.RequestURI
	argumentMap := extractPostData(incomingRequest)
	methodName := r.getMethodFromPath(path)
	arguments := r.extractArguments(path)
	r.callServiceObject(methodName, arguments, argumentMap, reply)
	return reply
}

func (r *RPC) extractArguments(path string) (args []string) {
	for _, value := range strings.Split(r.getApplicationPath(path), "/")[1:] {
		unescapedArg, err := url.QueryUnescape(value)
		if err != nil {
			panic(err)
		}
		args = append(args, unescapedArg)
	}
	return args
}

func (r *RPC) callServiceObject(methodName string, arguments []string, argumentMap map[string]interface{}, reply *rpc.MethodReply) {
	reflectionArgs := []reflect.Value{}
	reflectionArgs = append(reflectionArgs, reflect.ValueOf(arguments), reflect.ValueOf(argumentMap), reflect.ValueOf(reply))
	reflect.ValueOf(r.serviceObject).MethodByName(methodName).Call(reflectionArgs)
}

func extractPostData(incomingRequest *http.Request) map[string]interface{} {
	body, err := ioutil.ReadAll(incomingRequest.Body)
	if err != nil {
		panic(err)
	}
	if len(body) > 0 {
		return jsonDecode(body).(map[string]interface{})
	}
	return make(map[string]interface{})
}

func jsonDecode(jsonData []byte) (data interface{}) {
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		panic(err)
	} else {
		return data
	}
}

func jsonEncode(data interface{}) []byte {
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	} else {
		return b
	}
}

func (r *RPC) ServeHTTP(w http.ResponseWriter, incomingRequest *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := r.callServiceObjectWithHTTPRequest(incomingRequest)

	err := utils.ServeCompressed(w, incomingRequest, func(writer io.Writer) error {
		return json.NewEncoder(writer).Encode(response)
	})
	if err != nil {
		panic(err)
	}
}
