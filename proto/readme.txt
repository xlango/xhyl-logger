docker build -t grpc/server:v1 server/*
docker build -t grpc/client:v1 client/*
docker run -d --restart=unless-stopped --name grpcserver -p 50001:50001 grpc/server:v1
docker run -d --restart=unless-stopped --name grpcclient -p 7001:7001 grpc/client:v1


curl http://localhost:7001/grpc/client/test