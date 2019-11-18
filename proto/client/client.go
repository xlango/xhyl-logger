package main

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	pb "logconnection/proto/model"
	"net/http"
)

const (
	address = "192.168.10.190:5021"
)

func main() {
	http.HandleFunc("/logger/client/test", testClient)
	err := http.ListenAndServe(fmt.Sprintf(":%d", 5020), nil)
	fmt.Println(err)
}

func testClient(writer http.ResponseWriter, request *http.Request) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	c := pb.NewLoggerClient(conn)

	r, err := c.Info(context.Background(), &pb.RequestInfo{Content: "12345677", NodeName: "logtest1"})

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Println(r.Status)
	writer.Write([]byte(fmt.Sprintf("%d", r.Status)))
}
