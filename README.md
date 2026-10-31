# Emby Cinema

基于 **Wails v2 + Go + React + TypeScript** 的跨平台 Emby 桌面客户端。界面采用深色影院风格和玻璃质感，目标是把私人 Emby 媒体库做成独立桌面应用。

## 已实现

- Wails 桌面工程骨架与 React/Vite 前端
- Emby 用户名/密码登录
- 首页：媒体库分类、最近添加、继续观看
- 电影与电视剧独立浏览
- 列表分页（每页 30 项，读取 Emby 总记录数）
- Emby 全库搜索（电影、电视剧、单集）
- 点击媒体卡片查看详情
- 海报由 Go 后端携带 Token 请求，再以 Base64 图片数据交给前端；前端不再持有含 Token 的图片 URL
- GitHub Actions：前端构建与 Go 测试

## 尚未实现

- 播放器集成（MPV）、播放控制和播放进度回写
- 字幕与音轨选择
- 电视剧季/集导航及完整详情信息
- 操作系统凭据存储、登录状态持久化
- macOS / Windows 安装包和自动更新
- 更完善的分类筛选、排序选项

## 架构

- Wails v2：桌面窗口与 Go/前端绑定
- Go：Emby API、会话状态、海报代理
- React + Vite + TypeScript：桌面 UI
- Emby Server：媒体文件、元数据、用户权限和播放状态

客户端本身不提供影视片源，也不存储媒体文件；需要连接你自己的 Emby Server。

## 开发

需要 Go、Node.js、Wails v2 CLI，以及 Wails 对应平台的系统依赖。

```bash
wails dev
```

前端单独构建：

```bash
cd frontend
npm install
npm run build
```

桌面应用构建：

```bash
cd ..
wails build
```

## 安全说明

登录 Token 暂存在 Go 进程内存中，不写入前端状态或本地配置。海报由 Go 端携带 `X-Emby-Token` 请求，并以 data URL 交给前端，避免在图片请求 URL 中暴露 Token。建议优先使用 HTTPS 或可信局域网连接。当前版本退出应用后需要重新登录。

## 项目状态

目前已经具备 Emby 连接、首页、电影/电视剧列表、搜索、分页和基础详情弹层；尚未完成媒体播放，因此还不是完整可用的播放器客户端。
