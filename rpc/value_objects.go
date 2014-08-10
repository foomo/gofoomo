package rpc

// from php class Foomo\Services\RPC\Protocol\Call\MethodCall
// serializing a method call
type MethodCall struct {
	// id of the method call
	Id string `json:"id"`
	// name of the method to be called
	Method string `json:"method"`
	// the method call arguments
	Arguments []struct {
		Name  string      `json:"name"`
		Value interface{} `json:"value"`
	} `json:"arguments"`
}

// from php class Foomo\Services\RPC\Protocol\Reply\MethodReply
// reply to a method call
type MethodReply struct {
	// id of the method call
	Id string `json:"id"`
	// return value
	Value interface{} `json:"value"`
	// server side exception
	Exception interface{} `json:"exception"`
	// messages from the server
	// possibly many of them
	// possibly many types
	Messages interface{} `json:"messages"`
}
