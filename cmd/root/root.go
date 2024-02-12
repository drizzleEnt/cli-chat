package root

import (
	"context"
	"fmt"
	"log"
	"os"

	desc "github.com/drizzleent/cli-chat/pkg/chat_v1"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chat",
	Short: "cli-chat",
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create ",
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "deleting",
}

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "connecting to chat with id",
	Run: func(cmd *cobra.Command, args []string) {
		chatId, err := cmd.Flags().GetString("chatId")
		if err != nil {
			log.Fatalf("failed to get chat id: %v", err)
		}

		client := ConnectChatServer()
		ctx := context.Background()
		var username string
		log.Print("Enter login: ")
		_, err = fmt.Fscan(os.Stdin, &username)
		if err != nil {
			log.Fatalf("failed to get username: %v", err)
		}
		err = ConnectChat(ctx, client, chatId, username)
		if err != nil {
			log.Fatalf("failed to connect chat: %v", err)
		}
		log.Printf("user %s connected to chat %s", username, chatId)
	},
}

// var connectExistUser = &cobra.Command{
// 	Use:   "user",
// 	Short: "connecting to chat with username",
// }

var createUserCmd = &cobra.Command{
	Use:   "user",
	Short: "creating new user",
	Run: func(cmd *cobra.Command, args []string) {
		client := ConnectChatServer()
		ctx := context.Background()

		usernameStr, err := cmd.Flags().GetString("username")
		if err != nil {
			log.Fatalf("failed to get username %v", err.Error())
		}
		passwordStr, err := cmd.Flags().GetString("passwrd")
		if err != nil {
			log.Fatalf("failed to get password %v", err.Error())
		}

		resp, err := client.Create(ctx, &desc.CreateRequest{
			Username: usernameStr,
			Password: passwordStr,
		})
		if err != nil {
			log.Fatalf("failed to create user %v", err.Error())
		}

		log.Printf("user %s created with id:%v", usernameStr, resp.GetId())
	},
}

var createChatCmd = &cobra.Command{
	Use:   "chat",
	Short: "creating new chat",
	Run: func(cmd *cobra.Command, args []string) {
		client := ConnectChatServer()

		ctx := context.Background()
		chatId, err := createChat(ctx, client)

		if err != nil {
			log.Fatalf("failed to create chat: %v", err)
		}

		log.Printf("Chat created %s\n", chatId)
	},
}

var deleteUserCmd = &cobra.Command{
	Use:   "user",
	Short: "delete user",
	Run: func(cmd *cobra.Command, args []string) {
		usernameStr, err := cmd.Flags().GetString("username")
		if err != nil {
			log.Fatalf("failed to get username %v", err.Error())
		}
		log.Printf("user %s deleted", usernameStr)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(connectCmd)

	createCmd.AddCommand(createUserCmd)
	createCmd.AddCommand(createChatCmd)

	deleteCmd.AddCommand(deleteUserCmd)

	connectCmd.Flags().StringP("chatId", "i", "", "Chat id")
	err := connectCmd.MarkFlagRequired("chatId")
	if err != nil {
		log.Fatalf("failed to mark chatId flag as required %v", err.Error())
	}

	createUserCmd.Flags().StringP("username", "u", "", "User name")
	err = createUserCmd.MarkFlagRequired("username")
	if err != nil {
		log.Fatalf("failed to mark username flag as required %v", err.Error())
	}
	createUserCmd.Flags().StringP("passwrd", "p", "", "Password")
	err = createUserCmd.MarkFlagRequired("passwrd")
	if err != nil {
		log.Fatalf("failed to mark password flag as required %v", err.Error())
	}

	deleteUserCmd.Flags().StringP("username", "u", "", "User name")
	err = deleteUserCmd.MarkFlagRequired("username")
	if err != nil {
		log.Fatalf("failed to mark username flag as required %v", err.Error())
	}
}
