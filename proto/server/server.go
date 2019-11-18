package proto

import (
	context2 "context"
	"fmt"
	"log"
	"logconnection/conf"
	"logconnection/conf/es"
	"logconnection/proto/model"
	"net"
	"strconv"

	"google.golang.org/grpc"
	pb "logconnection/proto/model"
)

type LogServer struct{}

func (s *LogServer) Info(c context2.Context, req *model.RequestInfo) (*model.ResposeResult, error) {
	err := es.InsertLog(req.NodeName, "node", req.Content)
	if err != nil {
		return nil, err
	}
	return &model.ResposeResult{
		Status: 0,
	}, nil
}

func GrpcRegisterLogServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", strconv.Itoa(conf.GlobalConfig.Port)))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterLoggerServer(s, &LogServer{})
	log.Println("logger rpc服务已经开启")
	s.Serve(lis)
}
