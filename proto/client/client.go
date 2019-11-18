package logc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	pb "logconnection/proto/model"
)

var (
	Address string
)

func init() {
	Address = "127.0.0.1:5021"
}

func Info(nodeName, content string) {
	sendLog(pb.Level_INFO, content, nodeName)
}
func Debug(nodeName, content string) {
	sendLog(pb.Level_DEBUG, content, nodeName)
}
func Error(nodeName, content string) {
	sendLog(pb.Level_ERROR, content, nodeName)
}
func Warn(nodeName, content string) {
	sendLog(pb.Level_WARN, content, nodeName)
}
func Fatal(nodeName, content string) {
	sendLog(pb.Level_FATAL, content, nodeName)
}
func Off(nodeName, content string) {
	sendLog(pb.Level_OFF, content, nodeName)
}
func Trace(nodeName, content string) {
	sendLog(pb.Level_TRACE, content, nodeName)
}
func All(nodeName, content string) {
	sendLog(pb.Level_ALL, content, nodeName)
}

func sendLog(level pb.Level, content, nodeName string) {
	conn, err := grpc.Dial(Address, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("log grpc did not connect: %v", err)
	}

	defer conn.Close()

	c := pb.NewLoggerClient(conn)

	_, err = c.Info(context.Background(), &pb.RequestInfo{Level: level, Content: content, NodeName: nodeName})

	if err != nil {
		log.Fatalf("log grpc could not greet: %v", err)
	}

}
