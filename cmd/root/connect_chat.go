package root

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"time"

	desc "github.com/drizzleent/cli-chat/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConnectChat(ctx context.Context, client desc.ChatV1Client, chatId string) error {
	stream, err := client.ConnectChat(ctx, &desc.ConnectChatRequest{
		ChatId:   chatId,
		Username: "username",
	})

	if err != nil {
		return err
	}

	go func() {
		for {
			message, errRecv := stream.Recv()
			if errRecv == io.EOF {
				return
			}

			if errRecv != nil {
				log.Printf("failed to receive message from stream: %v", errRecv)
				return
			}
			// if message.GetFrom() == username {
			// 	continue
			// }

			log.Printf("[%v]-[from %s]: %s", message.GetCreatedAt().AsTime().Format(time.RFC3339), message.GetFrom(), message.GetText())
		}
	}()

	for {
		reader := bufio.NewReader(os.Stdin)
		msg, err := reader.ReadString('\n')

		if err != nil {
			log.Println("failed to read message: ", err)
			return err
		}
		_, err = client.SendMessage(ctx, &desc.SendMessageRequest{
			ChatId: chatId,
			Message: &desc.Message{
				From:      "username",
				Text:      msg,
				CreatedAt: timestamppb.Now(),
			},
		})

		if err != nil {
			log.Println("failed to send message: ", err)
			return err
		}
	}

}
