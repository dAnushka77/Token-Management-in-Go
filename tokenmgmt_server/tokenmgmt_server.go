package main

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	pb "example.com/go-tokenmgmt-grpc/tokenmgmt"
	"google.golang.org/grpc"
)

var (
	port       = flag.Int("port", 50051, "The server port")
	tokenDB    = make(map[string]tokenData)
	map_access sync.Mutex
)

type domainStruct struct {
	low  uint64
	mid  uint64
	high uint64
}

type stateStruct struct {
	partialValue uint64
	finalValue   uint64
}

type tokenData struct {
	id     string
	name   string
	domain []domainStruct
	state  []stateStruct
}

type TokenManagementServer struct {
	pb.UnimplementedTokenManagementServer
}

func Hash(name string, nonce uint64) uint64 {
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s %d", name, nonce)))
	return binary.BigEndian.Uint64(hasher.Sum(nil))
}

func getPartialValue(ipArray []uint64, low uint64, mid uint64) uint64 {
	min := ipArray[0]
	parValue := low
	for i := 0; i < len(ipArray); i++ {
		if ipArray[i] < min {
			min = ipArray[i]
			parValue = low
		}
		low = low + 1
	}
	return parValue
}

func (s *TokenManagementServer) Create(ctx context.Context, in *pb.Token_Create_I) (*pb.Token_Create_O, error) {
	log.Printf("Create Module Received: %v", in.GetId())
	map_access.Lock()
	tokenDB[in.GetId()] = tokenData{id: in.GetId(), state: []stateStruct{{partialValue: 0, finalValue: 0}}}
	status := "Sucess"
	fmt.Printf(`Current token information
	ID: %v
	Name: %v
	Domain[low, mid, high]: %d
	State[partial value, full value: %d`, tokenDB[in.GetId()].id, tokenDB[in.GetId()].name, tokenDB[in.GetId()].domain, tokenDB[in.GetId()].state)
	fmt.Println("List of ID's all the tokens:")
	idList := make([]string, 0, len(tokenDB))
	for i := range tokenDB {
		idList = append(idList, i)
	}
	fmt.Println(idList)
	map_access.Unlock()
	return &pb.Token_Create_O{Status: status}, nil
}

func (s *TokenManagementServer) Drop(ctx context.Context, in *pb.Token_Drop_I) (*pb.Token_Drop_O, error) {
	log.Printf("Drop Module Received: %v", in.GetId())
	//fmt.Println(tokenDB)
	map_access.Lock()
	fmt.Printf(`Current token information
	ID: %v
	Name: %v
	Domain[low, mid, high]: %d
	State[partial value, full value: %d`, tokenDB[in.GetId()].id, tokenDB[in.GetId()].name, tokenDB[in.GetId()].domain, tokenDB[in.GetId()].state)
	delete(tokenDB, in.GetId())
	fmt.Println("List of ID's all the tokens:")
	idList := make([]string, 0, len(tokenDB))
	for i := range tokenDB {
		idList = append(idList, i)
	}
	fmt.Println(idList)
	map_access.Unlock()
	return &pb.Token_Drop_O{}, nil
}

func (s *TokenManagementServer) Write(ctx context.Context, in *pb.Token_Write_I) (*pb.Token_Write_O, error) {
	log.Printf("Write Module Received: %v", in.GetId())
	map_access.Lock()
	tokenDB[in.GetId()] = tokenData{id: in.GetId(), name: in.GetName(), domain: []domainStruct{{in.GetLow(), in.GetMid(), in.GetHigh()}}}
	status := "Write Module Sucess"
	var x_array []uint64
	for i := in.GetLow(); i < in.GetMid(); i++ {
		h := Hash(in.GetName(), i)
		x_array = append(x_array, h)
	}
	partial_value := getPartialValue(x_array, in.GetLow(), in.GetMid())
	tokenDB[in.GetId()] = tokenData{id: in.GetId(), name: in.GetName(), domain: []domainStruct{{in.GetLow(), in.GetMid(), in.GetHigh()}}, state: []stateStruct{{partialValue: partial_value, finalValue: 0}}}
	//fmt.Println(tokenDB)
	fmt.Printf(`Current token information
	ID: %v
	Name: %v
	Domain[low, mid, high]: %d
	State[partial value, full value: %d`, tokenDB[in.GetId()].id, tokenDB[in.GetId()].name, tokenDB[in.GetId()].domain, tokenDB[in.GetId()].state)
	fmt.Println("List of ID's all the tokens:")
	idList := make([]string, 0, len(tokenDB))
	for i := range tokenDB {
		idList = append(idList, i)
	}
	fmt.Println(idList)
	map_access.Unlock()
	return &pb.Token_Write_O{Status: status, Partial_Value: partial_value}, nil
}

func (s *TokenManagementServer) Read(ctx context.Context, in *pb.Token_Read_I) (*pb.Token_Read_O, error) {
	log.Printf("Read Module Received: %v", in.GetId())
	map_access.Lock()
	rName := tokenDB[in.GetId()].name
	rLow := tokenDB[in.GetId()].domain[0].low
	rMid := tokenDB[in.GetId()].domain[0].mid
	rHigh := tokenDB[in.GetId()].domain[0].high
	var x_array []uint64
	for i := rLow; i < rMid; i++ {
		h := Hash(rName, i)
		x_array = append(x_array, h)
	}
	partial_value1 := getPartialValue(x_array, rLow, rMid)

	var y_array []uint64
	for i := rMid; i < rHigh; i++ {
		h := Hash(rName, i)
		y_array = append(y_array, h)
	}
	partial_value2 := getPartialValue(y_array, rMid, rHigh)
	final_value := partial_value1
	if partial_value2 < partial_value1 {
		final_value = partial_value2
	}
	//fmt.Println(partial_value1, partial_value2)
	tokenDB[in.GetId()] = tokenData{id: in.GetId(), name: rName, domain: []domainStruct{{rLow, rMid, rHigh}}, state: []stateStruct{{partialValue: partial_value1, finalValue: final_value}}}
	fmt.Printf(`Current token information
	ID: %v
	Name: %v
	Domain[low, mid, high]: %d
	State[partial value, full value: %d`, tokenDB[in.GetId()].id, tokenDB[in.GetId()].name, tokenDB[in.GetId()].domain, tokenDB[in.GetId()].state)
	fmt.Println("List of ID's all the tokens:")
	idList := make([]string, 0, len(tokenDB))
	for i := range tokenDB {
		idList = append(idList, i)
	}
	fmt.Println(idList)
	map_access.Unlock()
	return &pb.Token_Read_O{Final_Value: final_value}, nil
}

func main() {

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen : %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTokenManagementServer(s, &TokenManagementServer{})
	log.Printf("server listening at: %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
