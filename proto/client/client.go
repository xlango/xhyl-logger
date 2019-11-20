package logc

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
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
	go sendLog(pb.Level_INFO, content, NodeName)
}
func Debug(content string) {
	go sendLog(pb.Level_DEBUG, content, NodeName)
}
func Error(content string) {
	go sendLog(pb.Level_ERROR, content, NodeName)
}
func Warn(content string) {
	go sendLog(pb.Level_WARN, content, NodeName)
}
func Fatal(content string) {
	go sendLog(pb.Level_FATAL, content, NodeName)
}
func Off(content string) {
	go sendLog(pb.Level_OFF, content, NodeName)
}
func Trace(content string) {
	go sendLog(pb.Level_TRACE, content, NodeName)
}
func All(content string) {
	go sendLog(pb.Level_ALL, content, NodeName)
}

func sendLog(level pb.Level, content, nodeName string) {
	conn, err := grpc.Dial(Address, grpc.WithInsecure())

	if err != nil {
		fmt.Println("log grpc did not connect: ", err)
	}

	defer conn.Close()

	c := pb.NewLoggerClient(conn)

	_, err = c.Info(context.Background(), &pb.RequestInfo{Level: level, Content: content, NodeName: nodeName})

	if err != nil {
		fmt.Println("log grpc could not greet: ", err)
	}

}
