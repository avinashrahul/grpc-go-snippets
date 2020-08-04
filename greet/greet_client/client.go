package main

import (
	"context"
	"github.com/personal/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("unable to connect to server %v", err)
	}

	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	//fmt.Printf("create client %f", c)
	//doUnary(c)
	//doSum(c)
	//doServerStream(c)
	//primeNumberStream(c)
	//doClientStreaming(c)
	//doComputeAverageClientStreaming(c)
	doBiDiStreaming(c)
	doSquareRoot

}

func doUnary(conn greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Rahul1",
			LastName: "Gudimetla",
		},
	}

	res, err := conn.Greet(context.Background(), req)

	if err != nil {
		log.Fatalf("error %v", req)
	}

	log.Printf("response %v", res.Result)
}

func doSum(conn greetpb.GreetServiceClient) {
	req := &greetpb.SumRequest{
		VarA: 12,
		VarB: 15,
	}

	res, err := conn.Sum(context.Background(), req)

	if err != nil {
		log.Fatalf("error %v", req)
	}

	log.Printf("response %v", res.Result)

}

func doServerStream(conn greetpb.GreetServiceClient) {

	req := greetpb.GreetManyTimesRequest{
		FirstName: "Rahul23",
		LastName: "Gudimetla",
	}
	res, err := conn.GreetManyTimes(context.Background(), &req)

	if err != nil {
		log.Fatalf("Error in doServerStream request %v", err)
	}

	for {
		resp, err := res.Recv()

		if err == io.EOF {
			break;
		}

		if err != nil {
			log.Fatalf("Error in doServerStream response %v", err)
		}

		log.Printf("Response is %v", resp.Result)
	}


}

func primeNumberStream(conn greetpb.GreetServiceClient) {

	req := greetpb.PrimeNumberDecompositionRequest {
		Number: 210,
	}
	res, err := conn.PrimeNumberDecomposition(context.Background(), &req)

	if err != nil {

	}

	for {
		resp, err := res.Recv()

		if err == io.EOF {
			break;
		}

		if err != nil {

		}

	log.Printf("Response is %v", resp.Result)
	}
}

func doClientStreaming(conn greetpb.GreetServiceClient) {

	req := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			FirstName: "Rahul",
		},
		&greetpb.LongGreetRequest{
			FirstName: "Rahul1",
		},
		&greetpb.LongGreetRequest{
			FirstName: "Rahul2",
		},
		&greetpb.LongGreetRequest{
			FirstName: "Rahul3",
		},
	}
	stream, err := conn.LongGreet(context.Background())

	if err != nil {
		log.Fatalf("Error is %v", err)
	}

	for _, req := range req {
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, _ := stream.CloseAndRecv()

	log.Println("outout is %v", res)
	//fmt.Println("output is %v", res)

}

func doComputeAverageClientStreaming(conn greetpb.GreetServiceClient) {

	requests := []*greetpb.ComputeAverageRequest{
		&greetpb.ComputeAverageRequest{
			Number: 1,
		},
		&greetpb.ComputeAverageRequest{
			Number: 2,
		},
		&greetpb.ComputeAverageRequest{
			Number: 4,
		},
	}

	stream, _ := conn.ComputeAverage(context.Background())

	for _, req := range requests {
		stream.Send(req)
	}

	res, _ := stream.CloseAndRecv()

	log.Println("Response is %v", res)

}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	log.Println("started")

	stream, err := c.GreetEveryone(context.Background())

	if err != nil {
		log.Fatalf("Error in client %v", err)
	}

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			FirstName: "Rahul",
		},
		&greetpb.GreetEveryoneRequest{
			FirstName: "Rahul1",
		},
		&greetpb.GreetEveryoneRequest{
			FirstName: "Rahul2",
		},
		&greetpb.GreetEveryoneRequest{
			FirstName: "Rahul3",
		},
	}

	waitc := make(chan struct{})

	// func to send
	go func() {

		for _, req := range requests {
			sendErr := stream.Send(req)
			if sendErr != nil {
				return
			}
		}

		stream.CloseSend()
	}()



	//func to receive

	go func() {
		for {
			resp, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("error in receive client %v", err)
				break
			}


			log.Println("Response received from Server - %v", resp.Result)
		}
		close(waitc)
	}()

	<-waitc
}
