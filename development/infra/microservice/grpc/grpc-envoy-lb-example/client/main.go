package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"

	pb "grpc-envoy-lb/pb"
)

var (
	ipport     string
	numRequest int
)

func init() {
	flag.StringVar(&ipport, "addr", "127.0.0.1:8080", "server ip port")
	flag.IntVar(&numRequest, "n", 100, "total request")
}

func main() {
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, err := getConn(ctx, ipport)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)

	msgCh := make(chan string, numRequest)
	done := make(chan string)
	go func() {
		statistic := map[string]int{}
		for {
			select {
			case name, ok := <-msgCh:
				if !ok {
					done <- fmt.Sprint(statistic)
					return
				}
				statistic[name] += 1
			}
		}
	}()

	for i := 0; i < numRequest; i++ {
		resp, err := client.Greet(ctx, &pb.GreetRequest{})
		if err != nil {
			panic(err.Error())
		}
		msgCh <- strings.TrimPrefix(resp.Message, "Message from ")
	}
	close(msgCh)

	result := <-done
	fmt.Println(result)
}

func getConn(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	callOpts := []grpc.CallOption{}
	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(callOpts...),
		grpc.WithInsecure(),
	}

	return grpc.DialContext(ctx, addr, dialOpts...)
}
