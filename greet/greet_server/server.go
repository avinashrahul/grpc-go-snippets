package main

import (
	"context"
	"fmt"
	"github.com/personal/grpc-go/greet/greetpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"

	//"github.com/prometheus/common/log"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"time"
	"log"
)

type server struct {

}


func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	firstName := req.GetGreeting().GetFirstName()

	result := "Hello" + firstName

	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func (*server) Sum(ctx context.Context, req *greetpb.SumRequest) (*greetpb.SumResponse, error) {
	a := req.VarA
	b := req.VarB

	result := a + b

	res := &greetpb.SumResponse{
		Result: result,
	}

	return res, nil;
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	firstName := req.FirstName;

	for i := 0; i < 10; i++ {
		result := "Hello" + firstName + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}

		stream.Send(res);

		time.Sleep(1000 * time.Millisecond)
	}
	return nil;
}

func (*server) PrimeNumberDecomposition(req *greetpb.PrimeNumberDecompositionRequest, stream greetpb.GreetService_PrimeNumberDecompositionServer) error {
	number := req.Number
	divisor := int32(2)

	for number > 1 {
		newNum := number % divisor
		if number == 1 {
			break;
		} else if newNum == 0 {
			stream.Send(&greetpb.PrimeNumberDecompositionResponse{
				Result: divisor,
			})
			log.Printf("divisor is %v", divisor)

			number = number / divisor
			} else {
			divisor = divisor + 1

		}
	}
	return nil;
}

func (s *server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	result := ""
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}

		if err != nil {

		}

		name := req.FirstName

		result += "Hello!" + name

	}

}

func (s * server) ComputeAverage(stream greetpb.GreetService_ComputeAverageServer) error {

	result := int32(0)
	count := 0

	for {

		req, err := stream.Recv()

		if err == io.EOF {
			finalResult := float32(result)/float32(count)
			return stream.SendAndClose(&greetpb.ComputeAverageResponse{
				Result: finalResult,
			})
		}

		if err != nil {

		}

		count += 1
		result += req.Number
		log.Println("Result is %v", result)
	}
}

func (s * server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("received request")

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error in server %v", err)
			return err
		}

		name := req.FirstName

		sendErr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: name,
		})

		if sendErr != nil {
			log.Fatalf("Error in sending response to client %v", err)
			return sendErr
		}

	}
}

func(s *server) DoSquareRoot(ctx context.Context, req *greetpb.SquareRootRequest) (*greetpb.SquareRootResponse, error) {
	number := req.Number

	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %v", number))
	}


	return &greetpb.SquareRootResponse{
		Response: float32(number),

	}, nil
}


func main() {
	fmt.Println("Hello Go!")

	lis, err := net.Listen("tcp", "0.0.0.0:50051" )

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s:= grpc.NewServer()

	greetpb.RegisterGreetServiceServer(s, &server{})

	if err:= s.Serve(lis); err != nil {
		log.Fatalf("failed to serve %v", err)
	}


}
