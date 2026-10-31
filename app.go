package main

import (
 "context"
 "encoding/base64"
 "errors"
 "fmt"
 "net/url"
 "strings"
 "github.com/1264362700/emby-movie-site/internal/emby"
)

type App struct {
 ctx context.Context
 client *emby.Client
 token, serverURL, userID, serverName, userName string
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
 ImageTag string `json:"imageTag"`
 SeriesName string `json:"seriesName"`
}
type MediaPage struct {
 Items []MediaItem `json:"items"`
 Total int `json:"total"`
 Start int `json:"start"`
 Limit int `json:"limit"`
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
func (a *App) requireSession() error {
 if a.token==""||a.userID=="" { return errors.New("请先连接 Emby 服务器") }
 return nil
}
func mapItems(items []emby.Item) []MediaItem {
 out:=make([]MediaItem,0,len(items))
 for _,item:=range items {
  tag:="";if item.ImageTags!=nil { tag=item.ImageTags["Primary"] }
  out=append(out,MediaItem{ID:item.ID,Name:item.Name,Type:item.Type,Overview:item.Overview,Year:item.ProductionYear,Rating:item.CommunityRating,ImageTag:tag,SeriesName:item.SeriesName})
 }
 return out
}
func (a *App) LoadHome() (HomeData,error) {
 if err:=a.requireSession();err!=nil{return HomeData{},err}
 views,err:=a.client.Views(a.serverURL,a.token,a.userID);if err!=nil{return HomeData{},err}
 latestQuery:=url.Values{"Recursive":{"true"},"IncludeItemTypes":{"Movie,Series"},"SortBy":{"DateCreated"},"SortOrder":{"Descending"},"Limit":{"18"},"Fields":{"Overview,ProductionYear,CommunityRating,DateCreated"}}
 latest,err:=a.client.Items(a.serverURL,a.token,a.userID,latestQuery);if err!=nil{return HomeData{},err}
 resumeQuery:=url.Values{"Recursive":{"true"},"Filters":{"IsResumable"},"IncludeItemTypes":{"Movie,Episode"},"SortBy":{"DatePlayed"},"SortOrder":{"Descending"},"Limit":{"12"},"Fields":{"Overview,ProductionYear,CommunityRating,PlaybackPositionTicks,RunTimeTicks"}}
 resume,resumeErr:=a.client.Items(a.serverURL,a.token,a.userID,resumeQuery);if resumeErr!=nil{resume=[]emby.Item{}}
 return HomeData{ServerName:a.serverName,UserName:a.userName,Views:views,Latest:mapItems(latest),Resume:mapItems(resume)},nil
}
func (a *App) Browse(kind, parentID string, start, limit int, sortBy, sortOrder string) (MediaPage,error) {
 if err:=a.requireSession();err!=nil{return MediaPage{},err}
 if limit<1||limit>100{limit=30};if start<0{start=0}
 q:=url.Values{"Recursive":{"true"},"StartIndex":{fmt.Sprint(start)},"Limit":{fmt.Sprint(limit)},"Fields":{"Overview,ProductionYear,CommunityRating,DateCreated"}}
 if parentID!=""{q.Set("ParentId",parentID)}
 if kind=="Movie"{q.Set("IncludeItemTypes","Movie")}else if kind=="Series"{q.Set("IncludeItemTypes","Series")}else{q.Set("IncludeItemTypes","Movie,Series")}
 if sortBy==""{sortBy="SortName"};if sortOrder==""{sortOrder="Ascending"}
 q.Set("SortBy",sortBy);q.Set("SortOrder",sortOrder)
 items,err:=a.client.Items(a.serverURL,a.token,a.userID,q);if err!=nil{return MediaPage{},err}
 return MediaPage{Items:mapItems(items),Total:len(items),Start:start,Limit:limit},nil
}
func (a *App) Search(term string, start, limit int) (MediaPage,error) {
 if err:=a.requireSession();err!=nil{return MediaPage{},err}
 term=strings.TrimSpace(term);if term==""{return MediaPage{Items:[]MediaItem{},Total:0,Start:0,Limit:limit},nil}
 if limit<1||limit>100{limit=30};if start<0{start=0}
 q:=url.Values{"Recursive":{"true"},"SearchTerm":{term},"IncludeItemTypes":{"Movie,Series,Episode"},"StartIndex":{fmt.Sprint(start)},"Limit":{fmt.Sprint(limit)},"SortBy":{"SortName"},"SortOrder":{"Ascending"},"Fields":{"Overview,ProductionYear,CommunityRating"}}
 items,err:=a.client.Items(a.serverURL,a.token,a.userID,q);if err!=nil{return MediaPage{},err}
 return MediaPage{Items:mapItems(items),Total:len(items),Start:start,Limit:limit},nil
}
func (a *App) GetImageData(itemID, tag string) (string,error) {
 if err:=a.requireSession();err!=nil{return "",err}
 data,mime,err:=a.client.ImageBytes(a.serverURL,a.token,itemID,tag);if err!=nil{return "",err}
 if mime==""{mime="image/jpeg"}
 return "data:"+mime+";base64,"+base64.StdEncoding.EncodeToString(data),nil
}
func (a *App) GetItem(itemID string) (MediaItem,error) {
 if err:=a.requireSession();err!=nil{return MediaItem{},err}
 item,err:=a.client.Item(a.serverURL,a.token,a.userID,itemID);if err!=nil{return MediaItem{},err}
 return mapItems([]emby.Item{item})[0],nil
}
