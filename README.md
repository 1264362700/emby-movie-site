# Emby Cinema

基于 **Wails + Go + React + TypeScript** 的跨平台 Emby 桌面客户端。目标是将私人 Emby 媒体库做成沉浸式桌面影院，视觉方向参考 LibreTV 的深色、玻璃质感和海报式布局。

## 当前进度
- [x] Wails 桌面应用入口与 React/Vite 前端骨架
- [x] 深色玻璃风格连接页
- [x] Go 端调用 Emby 用户名密码认证接口
- [ ] 安全保存登录凭据与服务器配置
- [ ] 首页媒体库、最近添加、继续观看
- [ ] 分类、搜索、详情和剧集列表
- [ ] 播放器集成、字幕/音轨选择、播放进度同步
- [ ] macOS / Windows 安装包

## 架构
- Wails v2：桌面窗口与 Go/前端绑定
- Go：Emby API、配置管理和系统能力
- React + Vite + TypeScript：桌面 UI
- Emby Server：媒体文件、元数据、用户权限和播放状态

客户端不提供影视片源，也不存储媒体文件；需要连接你自己的 Emby Server。

## 开发环境
安装 Go、Node.js 和 Wails v2 CLI。在项目根目录执行：
```bash
wails dev
```
构建前端：
```bash
cd frontend
npm install
npm run build
cd ..
wails build
```

## 安全说明
当前连接页将密码通过本地 Wails 绑定交给 Go，再由 Go 请求 Emby。开发阶段不记录密码或 AccessToken 日志，也暂不持久化登录状态。正式版应使用操作系统凭据存储保存 Token，并支持退出和清除凭据。建议优先使用 HTTPS 或可信局域网连接。

## 项目状态
当前是桌面端基础工程阶段。媒体库完整浏览和播放器尚未实现，不能视为完整 Emby 客户端。