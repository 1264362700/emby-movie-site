package emby

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client struct { http *http.Client }

func NewClient() *Client {
	return &Client{http: &http.Client{Timeout: 15 * time.Second}}
}

type authRequest struct {
	Username string `json:"Username"`
	Pw string `json:"Pw"`
}
type authResponse struct {
	AccessToken string `json:"AccessToken"`
	User struct {
		ID string `json:"Id"`
		Name string `json:"Name"`
		ServerName string `json:"ServerName"`
	} `json:"User"`
}
type AuthResult struct {
	Token string
	UserID string
	UserName string
	ServerName string
}

func (c *Client) Authenticate(serverURL, username, password string) (AuthResult, error) {
	serverURL = strings.TrimRight(strings.TrimSpace(serverURL), "/")
	body, _ := json.Marshal(authRequest{Username: username, Pw: password})
	req, err := http.NewRequest(http.MethodPost, serverURL+"/Users/AuthenticateByName", bytes.NewReader(body))
	if err != nil { return AuthResult{}, fmt.Errorf("无效的服务器地址: %w", err) }
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Emby-Authorization", `MediaBrowser Client="Emby Cinema", Device="Desktop", DeviceId="emby-cinema-desktop", Version="0.1.0"`)
	resp, err := c.http.Do(req)
	if err != nil { return AuthResult{}, fmt.Errorf("无法连接 Emby 服务器: %w", err) }
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return AuthResult{}, fmt.Errorf("Emby 登录失败，服务器返回 HTTP %d", resp.StatusCode)
	}
	var data authResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil { return AuthResult{}, fmt.Errorf("解析 Emby 登录响应失败: %w", err) }
	if data.AccessToken == "" || data.User.ID == "" { return AuthResult{}, fmt.Errorf("Emby 未返回有效的登录令牌或用户信息") }
	name := data.User.ServerName
	if name == "" { name = "Emby" }
	return AuthResult{Token:data.AccessToken, UserID:data.User.ID, UserName:data.User.Name, ServerName:name}, nil
}
