package main

import (
	"logconnection/conf"
	"logconnection/consul"
	"logconnection/proto/server"
)

func main() {
	conf.InitConfig()
	go consul.RegisterServer()
	proto.GrpcRegisterLogServer()

}
