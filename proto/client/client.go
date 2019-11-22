package logc

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "logconnection/proto/model"
)

var (
	address  string
	nodeName string
)

func init() {
	address = "127.0.0.1:5021"
	nodeName = "unknown"
}

func Info(content string) {
	go sendLog(pb.Level_INFO, content, nodeName)
}
func Debug(content string) {
	go sendLog(pb.Level_DEBUG, content, nodeName)
}
func Error(content string) {
	go sendLog(pb.Level_ERROR, content, nodeName)
}
func Warn(content string) {
	go sendLog(pb.Level_WARN, content, nodeName)
}
func Fatal(content string) {
	go sendLog(pb.Level_FATAL, content, nodeName)
}
func Off(content string) {
	go sendLog(pb.Level_OFF, content, nodeName)
}
func Trace(content string) {
	go sendLog(pb.Level_TRACE, content, nodeName)
}
func All(content string) {
	go sendLog(pb.Level_ALL, content, nodeName)
}

func sendLog(level pb.Level, content, nodeName string) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())

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

func SetLogcAddress(addr string) {
	address = addr
}

func SetLogcNodeName(node string) {
	nodeName = node
}
