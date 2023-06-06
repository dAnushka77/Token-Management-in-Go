# Token-Management-in-Go

CMSC 621 – Advanced Operating Systems
Project 2 – Client-Server Token Manager
Done by – Anushka Dhekne (VD19739)

Steps to run the project –

1. Open the terminal and run the command given below to get start the server.

go run tokenmgmt_server/tokenmgmt_server.go -port 50051
If it displays listening at: [::]: 50051 it means that the 50051 port has been successfully passed through the terminal.

2) Open another terminal and run the command given below to start the client.

go run tokenmgmt_client/tokenmgmt_client.go -create -id 1234 -host localhost -port 50051 

3) Run the below given command at the client side.

go run tokenmgmt_client/tokenmgmt_client.go -write -id 1234 -name abc -low 0 -mid 10 -high 100 -host localhost -port 50051

4) Run on client side the below given command.

go run tokenmgmt_client/tokenmgmt_client.go -read -id 1234 -host localhost -port 50051

5) Run the given command at the client side.

go run tokenmgmt_client/tokenmgmt_client.go -drop 1234 -host localhost -port 50051    

References used –
1)	Overview | Protocol Buffers Documentation (protobuf.dev)
2)	Quick start | Go | gRPC
3)	Protocol Buffer Compiler Installation | gRPC
4)	Protocol Buffer Basics: Go | Protocol Buffers Documentation (protobuf.dev)
5)	Basics tutorial | Go | gRPC
6)	https://grpc.io/docs/languages/go/quickstart/#prerequisites
7)	https://grpc.io/docs/languages/go/basics/
