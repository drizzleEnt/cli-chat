package token

import (
	"context"
	"fmt"

	"github.com/drizzleent/cli-chat/cmd/root/server"
	login "github.com/drizzleent/cli-chat/pkg/login_v1"
)

func ReciveTokens(logStr, passwrdStr string) error {
	conn := server.ConnectLogin()
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
		return fmt.Errorf("failed to login: %v", err)
	}

	refreshToken := loginResp.GetRefreshToken()
	err = CreateRefresh(refreshToken)
	if err != nil {
		return fmt.Errorf("failed to create refreshtoken file: %v", err)
	}
	accsesResp, err := client.GetAccesToken(ctx, &login.GetAccessTokenRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return fmt.Errorf("failed to get access token: %v", err)
	}

	accessToken := accsesResp.GetAccessToken()
	err = CreateAccess(accessToken)
	if err != nil {
		return fmt.Errorf("failed to create acceasstoken file: %v", err)
	}
	return nil

}
