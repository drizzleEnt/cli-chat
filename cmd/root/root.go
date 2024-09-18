package root

import (
	"context"
	"log"
	"os"

	"github.com/drizzleent/cli-chat/cmd/root/initcmd"
	"github.com/drizzleent/cli-chat/cmd/root/md"
	"github.com/drizzleent/cli-chat/cmd/root/server"
	"github.com/drizzleent/cli-chat/cmd/root/token"
	chat "github.com/drizzleent/cli-chat/pkg/chat_v1"
	auth "github.com/drizzleent/cli-chat/pkg/user_v2"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"
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
	Use:   "conn",
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

		accessToken, err := token.ReadAccess()
		if err != nil {
			log.Fatalf("failed to get access token")
		}

		conn := server.ConnectChat()
		defer conn.Close()

		client := chat.NewChatV1Client(conn)
		ctx := md.CreateNewAuthMd(accessToken)
		res, err := client.GetName(ctx, &emptypb.Empty{})
		if err != nil {
			log.Fatalf("failed to get username: %v \n", err)
		}
		err = ConnectChat(ctx, client, chatId, res.GetName())
		if err != nil {
			log.Fatalf("failed to connect chat: %v \n", err)
		}
	},
}

var loginUserCmd = &cobra.Command{
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
		err = token.ReciveTokens(logStr, passwrdStr)
		if err != nil {
			log.Fatalf("failed to get tokens user %v", err)
		}
	},
}

var createUserCmd = &cobra.Command{
	Use:   "user",
	Short: "creating new user",
	Run: func(cmd *cobra.Command, args []string) {
		conn := server.ConnectLogin()
		defer conn.Close()
		client := auth.NewUserV1Client(conn)
		ctx := context.Background()

		usernameStr, err := cmd.Flags().GetString("username")
		if err != nil {
			log.Fatalf("failed to get username %v", err.Error())
		}
		passwordStr, err := cmd.Flags().GetString("passwrd")
		if err != nil {
			log.Fatalf("failed to get password %v", err.Error())
		}

		resp, err := client.Create(ctx, &auth.CreateRequest{
			Info: &auth.UserCreate{
				UserUpdate: &auth.UserUpdate{
					Name:  usernameStr,
					Email: "",
					Role:  0,
				},
				Password: passwordStr,
			},
		})
		if err != nil {
			log.Fatalf("failed to create user %v", err.Error())
		}
		err = token.ReciveTokens(usernameStr, passwordStr)
		if err != nil {
			log.Fatalf("failed to get tokens user %v", err)
		}

		log.Printf("user %s created with id:%v", usernameStr, resp.GetId())
	},
}

var createChatCmd = &cobra.Command{
	Use:   "chat",
	Short: "creating new chat",
	Run: func(cmd *cobra.Command, args []string) {
		accessToken, err := token.ReadAccess()
		if err != nil {
			log.Fatalf("failed to read token file: %v", err)
		}

		conn := server.ConnectChat()
		defer conn.Close()

		client := chat.NewChatV1Client(conn)
		ctx := md.CreateNewAuthMd(accessToken)
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
	connectCmd.AddCommand(loginUserCmd)

	createCmd.AddCommand(createUserCmd)
	createCmd.AddCommand(createChatCmd)

	deleteCmd.AddCommand(deleteUserCmd)

	connectChatCmd.Flags().StringP("chatId", "i", "", "Chat id")
	err := connectChatCmd.MarkFlagRequired("chatId")
	if err != nil {
		log.Fatalf("failed to mark chatId flag as required %v", err.Error())
	}

	initcmd.ConnectUserExistFlags(loginUserCmd)

	initcmd.CreateUserFlags(createUserCmd)

	deleteUserCmd.Flags().StringP("username", "u", "", "User name")
	err = deleteUserCmd.MarkFlagRequired("username")
	if err != nil {
		log.Fatalf("failed to mark username flag as required %v", err.Error())
	}
}
