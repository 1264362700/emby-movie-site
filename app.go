package main

import (
 "context"
 "errors"
 "net/url"
 "strings"
 "github.com/1264362700/emby-movie-site/internal/emby"
)

type App struct {
 ctx context.Context
 client *emby.Client
 token string
 serverURL string
 userID string
 serverName string
 userName string
}
func NewApp() *App { return &App{client:emby.NewClient()} }
func (a *App) startup(ctx context.Context) { a.ctx=ctx }

type ConnectResult struct {
 ServerName string `json:"serverName"`
 UserName string `json:"userName"`
 UserID string `json:"userId"`
}
type MediaItem struct {
 ID string `json:"id"`
 Name string `json:"name"`
 Type string `json:"type"`
 Overview string `json:"overview"`
 Year int `json:"year"`
 Rating float64 `json:"rating"`
 ImageURL string `json:"imageUrl"`
 SeriesName string `json:"seriesName"`
}
type HomeData struct {
 ServerName string `json:"serverName"`
 UserName string `json:"userName"`
 Views []emby.View `json:"views"`
 Latest []MediaItem `json:"latest"`
 Resume []MediaItem `json:"resume"`
}
func (a *App) Connect(serverURL, username, password string) (ConnectResult,error) {
 serverURL=strings.TrimRight(strings.TrimSpace(serverURL),"/")
 if serverURL==""||username==""||password=="" { return ConnectResult{},errors.New("请填写 Emby 地址、用户名和密码") }
 result,err:=a.client.Authenticate(serverURL,username,password);if err!=nil{return ConnectResult{},err}
 a.token,a.serverURL,a.userID=result.Token,serverURL,result.UserID
 a.serverName,a.userName=result.ServerName,result.UserName
 return ConnectResult{ServerName:result.ServerName,UserName:result.UserName,UserID:result.UserID},nil
}
func imageURL(serverURL,token,itemID,tag string) string {
 if tag=="" { return "" }
 q:=url.Values{"maxWidth":{"480"},"quality":{"85"},"tag":{tag},"api_key":{token}}
 return strings.TrimRight(serverURL,"/")+"/Items/"+url.PathEscape(itemID)+"/Images/Primary?"+q.Encode()
}
func mapItems(items []emby.Item,serverURL,token string) []MediaItem {
 out:=make([]MediaItem,0,len(items))
 for _,item:=range items {
  tag:="";if item.ImageTags!=nil { tag=item.ImageTags["Primary"] }
  out=append(out,MediaItem{ID:item.ID,Name:item.Name,Type:item.Type,Overview:item.Overview,Year:item.ProductionYear,Rating:item.CommunityRating,ImageURL:imageURL(serverURL,token,item.ID,tag),SeriesName:item.SeriesName})
 }
 return out
}
func (a *App) LoadHome() (HomeData,error) {
 if a.token==""||a.userID=="" { return HomeData{},errors.New("请先连接 Emby 服务器") }
 views,err:=a.client.Views(a.serverURL,a.token,a.userID);if err!=nil{return HomeData{},err}
 latestQuery:=url.Values{"Recursive":{"true"},"IncludeItemTypes":{"Movie,Series"},"SortBy":{"DateCreated"},"SortOrder":{"Descending"},"Limit":{"18"},"Fields":{"Overview,ProductionYear,CommunityRating,DateCreated"}}
 latest,err:=a.client.Items(a.serverURL,a.token,a.userID,latestQuery);if err!=nil{return HomeData{},err}
 resumeQuery:=url.Values{"Recursive":{"true"},"Filters":{"IsResumable"},"IncludeItemTypes":{"Movie,Episode"},"SortBy":{"DatePlayed"},"SortOrder":{"Descending"},"Limit":{"12"},"Fields":{"Overview,ProductionYear,CommunityRating,PlaybackPositionTicks,RunTimeTicks"}}
 resume,resumeErr:=a.client.Items(a.serverURL,a.token,a.userID,resumeQuery);if resumeErr!=nil { resume=[]emby.Item{} }
 return HomeData{ServerName:a.serverName,UserName:a.userName,Views:views,Latest:mapItems(latest,a.serverURL,a.token),Resume:mapItems(resume,a.serverURL,a.token)},nil
}
