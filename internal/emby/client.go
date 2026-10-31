package emby

import (
	"bytes"
	"encoding/json"
	"io"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct { http *http.Client }
func NewClient() *Client { return &Client{http: &http.Client{Timeout: 20 * time.Second}} }

type authRequest struct { Username string; Pw string }
type authResponse struct {
	AccessToken string
	User struct { ID string; Name string; ServerName string }
}
type AuthResult struct { Token, UserID, UserName, ServerName string }

type Item struct {
	ID string
	Name string
	Type string
	Overview string
	ProductionYear int
	CommunityRating float64
	ImageTags map[string]string
	SeriesName string
	ParentIndexNumber int
	IndexNumber int
}
type ItemResult struct { Items []Item; TotalRecordCount int }
type View struct { ID string; Name string; CollectionType string; ImageTags map[string]string }
type ViewResult struct { Items []View }

func (c *Client) Authenticate(serverURL, username, password string) (AuthResult, error) {
	serverURL = strings.TrimRight(strings.TrimSpace(serverURL), "/")
	body, _ := json.Marshal(authRequest{Username: username, Pw: password})
	req, err := http.NewRequest(http.MethodPost, serverURL+"/Users/AuthenticateByName", bytes.NewReader(body))
	if err != nil { return AuthResult{}, fmt.Errorf("无效的服务器地址: %w", err) }
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Emby-Authorization", "MediaBrowser Client=\"Emby Cinema\", Device=\"Desktop\", DeviceId=\"emby-cinema-desktop\", Version=\"0.1.0\"")
	resp, err := c.http.Do(req)
	if err != nil { return AuthResult{}, fmt.Errorf("无法连接 Emby 服务器: %w", err) }
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 { return AuthResult{}, fmt.Errorf("Emby 登录失败，HTTP %d", resp.StatusCode) }
	var data authResponse
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil { return AuthResult{}, fmt.Errorf("解析登录响应失败: %w", err) }
	if data.AccessToken == "" || data.User.ID == "" { return AuthResult{}, fmt.Errorf("Emby 未返回有效登录令牌或用户信息") }
	name := data.User.ServerName
	if name == "" { name = "Emby" }
	return AuthResult{Token:data.AccessToken, UserID:data.User.ID, UserName:data.User.Name, ServerName:name}, nil
}

func (c *Client) get(serverURL, token, path string, query url.Values, target interface{}) error {
	if query == nil { query = url.Values{} }
	query.Set("api_key", token)
	u := strings.TrimRight(serverURL, "/") + path + "?" + query.Encode()
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil { return err }
	req.Header.Set("X-Emby-Token", token)
	req.Header.Set("X-Emby-Authorization", "MediaBrowser Client=\"Emby Cinema\", Device=\"Desktop\", DeviceId=\"emby-cinema-desktop\", Version=\"0.1.0\"")
	resp, err := c.http.Do(req)
	if err != nil { return fmt.Errorf("请求 Emby 数据失败: %w", err) }
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 { return fmt.Errorf("Emby API 返回 HTTP %d", resp.StatusCode) }
	if err = json.NewDecoder(resp.Body).Decode(target); err != nil { return fmt.Errorf("解析 Emby 数据失败: %w", err) }
	return nil
}

func (c *Client) Views(serverURL, token, userID string) ([]View, error) {
	var result ViewResult
	err := c.get(serverURL, token, "/Users/"+url.PathEscape(userID)+"/Views", nil, &result)
	return result.Items, err
}
func (c *Client) Items(serverURL, token, userID string, query url.Values) ([]Item, error) {
	var result ItemResult
	err := c.get(serverURL, token, "/Users/"+url.PathEscape(userID)+"/Items", query, &result)
	return result.Items, err
}
func (c *Client) ImageURL(serverURL, token, itemID, tag string) string {
	if tag == "" { return "" }
	q := url.Values{"maxWidth":{"480"}, "quality":{"85"}, "tag":{tag}, "api_key":{token}}
	return strings.TrimRight(serverURL, "/")+"/Items/"+url.PathEscape(itemID)+"/Images/Primary?"+q.Encode()
}

func (c *Client) Item(serverURL, token, userID, itemID string) (Item, error) {
	var item Item
	err := c.get(serverURL, token, "/Users/"+url.PathEscape(userID)+"/Items/"+url.PathEscape(itemID), nil, &item)
	return item, err
}

func (c *Client) ImageBytes(serverURL, token, itemID, tag string) ([]byte, string, error) {
	if tag == "" { return nil, "", fmt.Errorf("该项目没有主海报") }
	q := url.Values{"maxWidth":{"720"}, "quality":{"88"}, "tag":{tag}}
	u := strings.TrimRight(serverURL, "/")+"/Items/"+url.PathEscape(itemID)+"/Images/Primary?"+q.Encode()
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil { return nil, "", err }
	req.Header.Set("X-Emby-Token", token)
	resp, err := c.http.Do(req)
	if err != nil { return nil, "", err }
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 { return nil, "", fmt.Errorf("获取海报失败，HTTP %d", resp.StatusCode) }
	data, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil { return nil, "", err }
	return data, resp.Header.Get("Content-Type"), nil
}
