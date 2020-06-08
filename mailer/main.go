package main

import (
	"bytes"
	"context"
	pb "github.com/NikolayOskin/go-trello-clone/mailer/mailerpkg"
	"github.com/mailgun/mailgun-go"
	"google.golang.org/grpc"
	"html/template"
	"log"
	"net"
	"os"
	"time"
)

type GRPCServer struct{}

var mg *mailgun.MailgunImpl
var tpl *template.Template

type EmailConfirmMessage struct {
	VerificationCode string
	tplname          string
}

func (s *GRPCServer) SendEmail(ctx context.Context, r *pb.EmailRequest) (*pb.EmailResponse, error) {
	m := &EmailConfirmMessage{r.Code, "signup-confirm.html"}

	var buf bytes.Buffer
	err := tpl.ExecuteTemplate(&buf, m.tplname, m)
	if err != nil {
		log.Fatalf("failed to execute template: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()

	if err = sendToMailgun("Thanks for registration!", buf.String(), r.Email); err != nil {
		return &pb.EmailResponse{Sent: false}, nil
	}

	return &pb.EmailResponse{Sent: true}, nil
}

func sendToMailgun(s string, body string, toEmail string) error {
	message := mg.NewMessage(os.Getenv("MAILER_SENDER"), s, "", toEmail)
	message.SetHtml(body)
	_, _, err := mg.Send(message)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	mg = mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_API_KEY"))
	tpl = template.Must(template.New("").ParseGlob("./templates/*.html"))
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
