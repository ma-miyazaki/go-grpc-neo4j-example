package main

import (
	"log"
	"net"

	"github.com/ma-miyazaki/go-grpc-neo4j-example/pb/employee"
	"github.com/ma-miyazaki/go-grpc-neo4j-example/server/infrastracture/persistence"
	"github.com/ma-miyazaki/go-grpc-neo4j-example/server/interface/handler"
	"github.com/ma-miyazaki/go-grpc-neo4j-example/server/usecase"

	"google.golang.org/grpc"
)

const port = ":50051"

// type ServerUnary struct {
// 	pb.UnimplementedCalcServer
// }

// func (s *ServerUnary) Sum(ctx context.Context, in *pb.SumRequest) (*pb.SumReply, error) {
// 	a := in.GetA()
// 	b := in.GetB()
// 	fmt.Println(a, b)
// 	reply := fmt.Sprintf("%d + %d = %d", a, b, a+b)
// 	return &pb.SumReply{
// 		Message: reply,
// 	}, nil
// }

func createEmployeeServer() employee.EmployeeServiceServer {
	employeeRepository := persistence.NewEmployeeRepository()
	employeeUseCase := usecase.NewEmployeeUseCase(employeeRepository)
	return handler.NewEmployeeHandler(employeeUseCase)
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()
	// pb.RegisterCalcServer(s, &ServerUnary{})
	employee.RegisterEmployeeServiceServer(s, createEmployeeServer())

	defer persistence.CloseNeo4jDriver()
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return
	}
}
