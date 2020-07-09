package mailer

import (
	"context"
	pb "github.com/NikolayOskin/go-trello-clone/mailer/src"
	"google.golang.org/grpc"
	"log"
	"os"
)

var client pb.MailerClient

func Start() {
	conn, err := grpc.DialContext(
		context.Background(),
		os.Getenv("MAILER_SERVICE"),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}
	client = pb.NewMailerClient(conn)
}

func Client() pb.MailerClient {
	return client
}
