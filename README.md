# Emby Cinema

基于 **Wails v2 + Go + React + TypeScript** 的跨平台 Emby 桌面客户端。目标是将私人 Emby 媒体库做成沉浸式桌面影院，视觉方向参考 LibreTV 的深色、玻璃质感和海报式布局。

## 当前进度

- [x] Wails 桌面应用入口与 React/Vite 前端骨架
- [x] 深色玻璃风格桌面连接页
- [x] Go 端 Emby 用户名密码认证
- [x] 读取 Emby 媒体库分类
- [x] 首页最近添加电影/电视剧
- [x] 首页继续观看（可恢复播放的电影/单集）
- [x] 首页海报、标题、年份和评分展示
- [x] GitHub Actions 前端构建与 Go 测试流程
- [ ] Go 图片代理（当前图片 URL 暂时包含 Emby api_key）
- [ ] 电影/电视剧独立浏览、筛选、分页和搜索
- [ ] 详情页、季/集导航
- [ ] 播放器集成、字幕/音轨、播放进度同步
- [ ] 操作系统凭据安全存储
- [ ] macOS / Windows 安装包与自动更新

## 架构

- Wails v2：桌面窗口与 Go/前端绑定
- Go：Emby API、会话状态和系统能力
- React + Vite + TypeScript：桌面 UI
- Emby Server：媒体文件、元数据、用户权限和播放状态

客户端不提供影视片源，也不存储媒体文件；需要连接你自己的 Emby Server。

## 开发环境

需要 Go、Node.js、Wails v2 CLI，以及 Wails 对应平台的系统依赖。项目根目录执行：

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

Emby 用户名密码通过本地 Wails 绑定交给 Go，Token 暂存在 Go 进程内存中，不写入日志或本地配置。当前海报 URL 使用 Emby 的 `api_key` 参数供 WebView 加载图片；后续应改为 Go 图片代理，避免 Token 出现在前端 URL。建议优先使用 HTTPS 或可信局域网连接。

## 项目状态

当前已具备桌面工程骨架、Emby 登录和首页数据读取能力，但尚不是完整媒体播放器。电影/电视剧浏览、搜索、详情、播放控制和字幕等功能仍在开发计划中。
