package root

import (
	"log"

	desc "github.com/drizzleent/cli-chat/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectChatServer() desc.ChatV1Client {
	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("failed to connect server: %v", err)
	}
	//defer conn.Close()
	client := desc.NewChatV1Client(conn)

	return client
}
