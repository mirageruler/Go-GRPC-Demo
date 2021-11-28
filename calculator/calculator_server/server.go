package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/mirageruler/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum func was invoked with %v\n", req)

	firstNum := req.GetFirstNumber()
	secondNum := req.GetSecondNumber()
	result := firstNum + secondNum
	res := calculatorpb.SumResponse{
		SumResult: result,
	}

	return &res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("PrimeNumberDecomposition func was invoked with %v\n", req)

	number := req.GetNumber()
	divisor := int64(2)

	for number > 1 {
		if number%divisor == 0 {
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			}
			stream.Send(res)
			time.Sleep(time.Second)
			number = number / divisor
		} else {
			divisor++
			fmt.Printf("Divisor has increased to %v\n", divisor)
		}
	}

	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("ComputeAverage func was invoked with a streaming request\n")

	sum := 0
	count := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we have finished reading the client stream
			average := float64(sum / count)
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: average,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		sum += int(req.GetNumber())
		count++
	}

}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("FindMaximum func was invoked with a streaming request\n")

	var maxNumber int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error while reading client stream: %v", err)
			return err
		}
		number := req.GetNumber()
		if number > maxNumber {
			maxNumber = number
		}

		err = stream.Send(&calculatorpb.FindMaximumResponse{
			MaxNumber: maxNumber,
		})
		if err != nil {
			log.Fatalf("error while sending data to client: %v", err)
		}
	}
}

func main() {
	fmt.Println("Calculator Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	addr := lis.Addr()
	bt1, _ := json.Marshal(addr)
	fmt.Println("ADDR: ", string(bt1))

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
