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

		token, err := token.ReadAccess()
		if err != nil {
			log.Fatalf("failed to get token")
		}
		conn := ConnectChatServer()
		defer conn.Close()
		client := chat.NewChatV1Client(conn)
		ctx := context.Background()
		md := metadata.New(map[string]string{"authorization": "Bearer " + token})
		ctx = metadata.NewOutgoingContext(ctx, md)
		res, err := client.GetName(ctx, &emptypb.Empty{})
		if err != nil {
			log.Fatalf("failed to get username: %v", err)
		}
		err = ConnectChat(ctx, client, chatId, res.GetName())
		if err != nil {
			log.Fatalf("failed to connect chat: %v", err)
		}
		log.Printf("user connected to chat %s", chatId)
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

		conn := ConnectLoginServer()
		defer conn.Close()

		ctx := context.Background()
		client := login.NewLoginV1Client(conn)
		loginResp, err := client.Login(ctx, &login.LoginRequest{
			Info: &login.Login{
				Username: logStr,
				Password: passwrdStr,
			},
		})
		if err != nil {
			log.Fatalf("failed to login: %v", err)
		}
		refreshToken := loginResp.GetRefreshToken()
		err = token.CreateRefresh(refreshToken)
		if err != nil {
			log.Fatalf("failed to create refreshtoken file: %v", err)
		}
		accsesResp, err := client.GetAccesToken(ctx, &login.GetAccessTokenRequest{
			RefreshToken: refreshToken,
		})
		if err != nil {
			log.Fatalf("failed to get access token: %v", err)
		}

		accessToken := accsesResp.GetAccessToken()
		err = token.CreateAccess(accessToken)
		if err != nil {
			log.Fatalf("failed to create acceasstoken file: %v", err)
		}
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
		accessToken, err := token.ReadAccess()
		if err != nil {
			log.Fatalf("failed to read token file: %v", err)
		}

		conn := ConnectChatServer()
		defer conn.Close()

		client := chat.NewChatV1Client(conn)
		ctx := context.Background()
		md := metadata.New(map[string]string{"authorization": "Bearer " + accessToken})
		ctx = metadata.NewOutgoingContext(ctx, md)
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
