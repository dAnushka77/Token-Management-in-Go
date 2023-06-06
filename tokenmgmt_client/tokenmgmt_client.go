package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "example.com/go-tokenmgmt-grpc/tokenmgmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	exeCreate = flag.Bool("create", false, "Create a token?")
	exeWrite  = flag.Bool("write", false, "Write a token?")
	exeRead   = flag.Bool("read", false, "Read a token?")
	exeDrop   = flag.Bool("drop", false, "Drop a token?")
	cmdId     = flag.String("id", "0000", "Token ID")
	cmdName   = flag.String("name", "no name", "Token Name")
	cmdLow    = flag.Uint64("low", 0, "Domain low")
	cmdMid    = flag.Uint64("mid", 0, "Domain mid")
	cmdHigh   = flag.Uint64("high", 0, "Domain high")
	cmdHost   = flag.String("host", "localhost", "host for client")
	cmdPort   = flag.String("port", "50051", "port for client")
)

func main() {
	flag.Parse()
	address := *cmdHost + ":" + *cmdPort
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not onnect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTokenManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if *exeCreate {
		op1, err := c.Create(ctx, &pb.Token_Create_I{Id: *cmdId})
		if err != nil {
			log.Fatalf("could not create token: %v", err)
		}
		log.Printf("Create Module Response: %v", op1.Status)
	}

	if *exeDrop {
		op2, err := c.Drop(ctx, &pb.Token_Drop_I{Id: *cmdId})
		if err != nil {
			log.Fatalf("could not drop token: %v", err)
		}
		i := 1
		if i == 0 {
			op2.Reset()
		}
		log.Printf("Token dropped")
	}

	if *exeWrite {
		op3, err := c.Write(ctx, &pb.Token_Write_I{Id: *cmdId, Name: *cmdName, Low: *cmdLow, Mid: *cmdMid, High: *cmdHigh})
		if err != nil {
			log.Fatalf("could not write user: %v", err)
		}
		log.Printf(`Write Module Response: %v
		Partial Value: %d`, op3.Status, op3.Partial_Value)
	}

	if *exeRead {
		op4, err := c.Read(ctx, &pb.Token_Read_I{Id: *cmdId})
		if err != nil {
			log.Fatalf("could not read token: %v", err)
		}
		log.Printf("Read Module Response: %v", op4.Final_Value)
	}

}
