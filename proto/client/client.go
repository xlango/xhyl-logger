package logc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	pb "logconnection/proto/model"
)

var (
	Address  string
	NodeName string
)

func init() {
	Address = "127.0.0.1:5021"
	NodeName = "unknown"
}

func Info(content string) {
	sendLog(pb.Level_INFO, content, NodeName)
}
func Debug(content string) {
	sendLog(pb.Level_DEBUG, content, NodeName)
}
func Error(content string) {
	sendLog(pb.Level_ERROR, content, NodeName)
}
func Warn(content string) {
	sendLog(pb.Level_WARN, content, NodeName)
}
func Fatal(content string) {
	sendLog(pb.Level_FATAL, content, NodeName)
}
func Off(content string) {
	sendLog(pb.Level_OFF, content, NodeName)
}
func Trace(content string) {
	sendLog(pb.Level_TRACE, content, NodeName)
}
func All(content string) {
	sendLog(pb.Level_ALL, content, NodeName)
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
