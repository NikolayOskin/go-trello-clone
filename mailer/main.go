package main

import (
	"context"
	pb "github.com/NikolayOskin/go-trello-clone/mailer/mailerpkg"
	"github.com/mailgun/mailgun-go"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"
)

type GRPCServer struct{}

func (s *GRPCServer) SendEmail(ctx context.Context, in *pb.EmailRequest) (*pb.EmailResponse, error) {
	mg := mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_API_KEY"))

	sender := "sender@example.com"
	subject := "Welcome!"
	body := in.Code
	recipient := in.Email

	message := mg.NewMessage(sender, subject, body, recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := mg.Send(message)
	if err != nil {
		return nil, err
	}

	return &pb.EmailResponse{Sent: true}, nil
}

func main() {
	lis, err := net.Listen("tcp", os.Getenv("MAILER_SERVICE_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMailerServer(grpcServer, &GRPCServer{})
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
}
