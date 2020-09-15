package jsonrpc_cli

import (
	"errors"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
	"time"
)

type RpcClient struct {
	rpcConn *net.Conn
	rpcCli  *rpc.Client
}

func (g *RpcClient) RpcConnect(serverAddr string, serverPort uint16, timeOut float64) error {
	client, err := net.DialTimeout("tcp", serverAddr+":"+strconv.Itoa(int(serverPort)),
		time.Duration(timeOut*1000*1000*1000))
	if err != nil {
		return err
	}

	rpcClient := jsonrpc.NewClient(client)
	if rpcClient == nil {
		return errors.New("invalid rpc client")
	}

	g.rpcConn = &client
	g.rpcCli = rpcClient
	return nil
}

func (g *RpcClient) RpcDisConnect() {
	if g.rpcCli != nil {
		_ = g.rpcCli.Close()
		g.rpcCli = nil
	}

	if g.rpcConn != nil {
		_ = (*g.rpcConn).Close()
		g.rpcConn = nil
	}
}

func (g RpcClient) RpcRequest(rpcFuncName string, rpcArg interface{}) (interface{}, error) {
	if g.rpcConn == nil {
		return nil, errors.New("invalid rpc connection")
	}

	if g.rpcCli == nil {
		return nil, errors.New("invalid rpc client")
	}

	var replyObj interface{}
	err := g.rpcCli.Call(rpcFuncName, rpcArg, &replyObj)

	if err != nil {
		return nil, err
	} else {
		return replyObj, nil
	}
}

func RpcRequestSimple(serverAddr string, serverPort uint16, timeOut float64,
	rpcFuncName string, rpcArg interface{}) (interface{}, error) {
	rpcClient := new(RpcClient)
	err := rpcClient.RpcConnect(serverAddr, serverPort, timeOut)
	if err != nil {
		return nil, err
	}

	defer rpcClient.RpcDisConnect()
	return rpcClient.RpcRequest(rpcFuncName, rpcArg)
}
