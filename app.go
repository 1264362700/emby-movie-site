package main

import (
	"context"
	"errors"
	"strings"

	"github.com/1264362700/emby-movie-site/internal/emby"
)

type App struct {
	ctx context.Context
	client *emby.Client
	token string
}

func NewApp() *App { return &App{client: emby.NewClient()} }
func (a *App) startup(ctx context.Context) { a.ctx = ctx }

type ConnectResult struct {
	ServerName string `json:"serverName"`
	UserName string `json:"userName"`
	UserID string `json:"userId"`
}

func (a *App) Connect(serverURL, username, password string) (ConnectResult, error) {
	serverURL = strings.TrimRight(strings.TrimSpace(serverURL), "/")
	if serverURL == "" || username == "" || password == "" {
		return ConnectResult{}, errors.New("请填写 Emby 地址、用户名和密码")
	}
	result, err := a.client.Authenticate(serverURL, username, password)
	if err != nil { return ConnectResult{}, err }
	a.token = result.Token
	return ConnectResult{ServerName: result.ServerName, UserName: result.UserName, UserID: result.UserID}, nil
}
