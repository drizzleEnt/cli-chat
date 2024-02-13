package root

import (
	"context"
	"log"
	"os"

	"github.com/drizzleent/cli-chat/cmd/root/initcmd"
	"github.com/drizzleent/cli-chat/cmd/root/token"
	chat "github.com/drizzleent/cli-chat/pkg/chat_v1"
	login "github.com/drizzleent/cli-chat/pkg/login_v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
)

var refreshToken string
var accessToken string

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
	Short: "connecting",
}

var connectChatCmd = &cobra.Command{
	Use:   "chat",
	Short: "connecting to chat with id",
	Run: func(cmd *cobra.Command, args []string) {
		chatId, err := cmd.Flags().GetString("chatId")
		if err != nil {
			log.Fatalf("failed to get chat id: %v", err)
		}

		token, err := token.ReadRefresh()
		if err != nil {
			log.Fatalf("failed to get token")
		}
		conn := ConnectChatServer()
		defer conn.Close()
		client := chat.NewChatV1Client(conn)
		ctx := context.Background()
		md := metadata.New(map[string]string{"Authrization": "Bearer " + token})
		ctx = metadata.NewOutgoingContext(ctx, md)
		err = ConnectChat(ctx, client, chatId)
		if err != nil {
			log.Fatalf("failed to connect chat: %v", err)
		}
		log.Printf("user connected to chat %s", chatId)
	},
}

var connectExistUserCmd = &cobra.Command{
	Use:   "user",
	Short: "getting refresh token with username",
	Run: func(cmd *cobra.Command, args []string) {
		logStr, err := cmd.Flags().GetString("username")
		if err != nil {
			log.Fatalf("failed to get login: %v", err)
		}

		passwrdStr, err := cmd.Flags().GetString("passwrd")
		if err != nil {
			log.Fatalf("failed to get password: %v", err)
		}

		conn := ConnectLoginServer()
		defer conn.Close()

		ctx := context.Background()
		client := login.NewLoginV1Client(conn)
		resp, err := client.Login(ctx, &login.LoginRequest{
			Info: &login.Login{
				Username: logStr,
				Password: passwrdStr,
			},
		})
		if err != nil {
			log.Fatalf("failed to login: %v", err)
		}
		refreshToken = resp.GetRefreshToken()
		file, err := os.Create("bin/token.txt")
		if err != nil {
			log.Fatalf("failed to create or open file: %v", err)
		}
		defer file.Close()
		_, err = file.Write([]byte(refreshToken))
		if err != nil {
			log.Fatalf("failed to write in file: %v", err)
		}
		log.Printf("your refresh token is %s\n", refreshToken)
	},
}

var createUserCmd = &cobra.Command{
	Use:   "user",
	Short: "creating new user",
	Run: func(cmd *cobra.Command, args []string) {
		conn := ConnectChatServer()
		defer conn.Close()
		client := chat.NewChatV1Client(conn)
		ctx := context.Background()

		usernameStr, err := cmd.Flags().GetString("username")
		if err != nil {
			log.Fatalf("failed to get username %v", err.Error())
		}
		passwordStr, err := cmd.Flags().GetString("passwrd")
		if err != nil {
			log.Fatalf("failed to get password %v", err.Error())
		}

		resp, err := client.Create(ctx, &chat.CreateRequest{
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
		conn := ConnectChatServer()
		defer conn.Close()

		client := chat.NewChatV1Client(conn)
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

	connectCmd.AddCommand(connectChatCmd)
	connectCmd.AddCommand(connectExistUserCmd)

	createCmd.AddCommand(createUserCmd)
	createCmd.AddCommand(createChatCmd)

	deleteCmd.AddCommand(deleteUserCmd)

	connectChatCmd.Flags().StringP("chatId", "i", "", "Chat id")
	err := connectChatCmd.MarkFlagRequired("chatId")
	if err != nil {
		log.Fatalf("failed to mark chatId flag as required %v", err.Error())
	}

	initcmd.ConnectUserExistFlags(connectExistUserCmd)

	initcmd.CreateUserFlags(createUserCmd)

	deleteUserCmd.Flags().StringP("username", "u", "", "User name")
	err = deleteUserCmd.MarkFlagRequired("username")
	if err != nil {
		log.Fatalf("failed to mark username flag as required %v", err.Error())
	}
}
