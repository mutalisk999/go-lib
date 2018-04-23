package jsonrpc_cli

import (
	"fmt"
	"testing"
)

func TestAll(t *testing.T) {
	result, err := RpcRequestSimple("127.0.0.1", 8888, 3, "info", nil)
	if err == nil {
		fmt.Println(result)
	} else {
		fmt.Println(err)
	}

	rpcClient := new(RpcClient)
	err = rpcClient.RpcConnect("127.0.0.1", 8888, 3)
	if err == nil {
		result, err = rpcClient.RpcRequest("info", nil)
		if err == nil {
			fmt.Println(result)
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
	rpcClient.RpcDisConnect()
}
