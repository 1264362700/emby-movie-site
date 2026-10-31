import { useState } from "react";
import { Clapperboard, Search, House, Film, Tv, Settings, Play, Server, LoaderCircle } from "lucide-react";

declare global {
  interface Window {
    go?: { main?: { App?: { Connect: (url: string, username: string, password: string) => Promise<{serverName:string;userName:string;userId:string}> } } };
  }
}

export default function App() {
  const [server, setServer] = useState("");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [busy, setBusy] = useState(false);
  const [connected, setConnected] = useState("");
  const [error, setError] = useState("");

  async function connect(e: React.FormEvent) {
    e.preventDefault(); setBusy(true); setError(""); setConnected("");
    try {
      const api = window.go?.main?.App;
      if (!api) throw new Error("当前不是 Wails 桌面运行环境，请先启动桌面开发模式");
      const result = await api.Connect(server, username, password);
      setConnected(`已连接到 ${result.serverName}，用户：${result.userName}`);
    } catch (err) { setError(err instanceof Error ? err.message : "连接失败"); }
    finally { setBusy(false); }
  }

  return <div className="shell">
    <aside className="sidebar">
      <div className="brand"><span className="brand-icon"><Clapperboard size={21}/></span><span>EMBER<span className="brand-light"> CINEMA</span></span></div>
      <div className="nav-label">媒体空间</div>
      <a className="nav active"><House size={18}/> 首页</a>
      <a className="nav"><Film size={18}/> 电影</a>
      <a className="nav"><Tv size={18}/> 电视剧</a>
      <div className="sidebar-bottom"><a className="nav"><Settings size={18}/> 设置</a><div className="server-state"><span className="dot"/> 尚未连接服务器</div></div>
    </aside>
    <main className="main">
      <header className="topbar"><div className="crumb">你的私人影院 <span>/</span> 首页</div><button className="icon-button"><Search size={19}/></button></header>
      <section className="hero">
        <div className="hero-copy"><div className="eyebrow">YOUR PERSONAL CINEMA</div><h1>把热爱，<br/><em>留在私人影院。</em></h1><p>连接你的 Emby 媒体服务器，继续探索收藏的电影与剧集。</p><a className="hero-link" href="#connect"><Play size={15} fill="currentColor"/> 开始连接 <span>↗</span></a></div>
        <div className="hero-art"><div className="orb orb-one"/><div className="orb orb-two"/><div className="art-card"><Clapperboard size={58} strokeWidth={1}/><span>YOUR LIBRARY<br/>YOUR CINEMA</span></div></div>
        <div className="hero-index">01 <span>/ 03</span></div>
      </section>
      <section className="connect-section" id="connect">
        <div className="section-heading"><div><div className="eyebrow">GET STARTED</div><h2>连接 Emby 服务器</h2><p>使用你的 Emby 服务器地址和账户登录。</p></div><div className="secure"><Server size={15}/> 本地直连 · 不经过第三方</div></div>
        <form className="connect-card" onSubmit={connect}>
          <label>服务器地址<input value={server} onChange={e=>setServer(e.target.value)} placeholder="例如：http://192.168.1.100:8096" required /></label>
          <div className="form-row"><label>用户名<input value={username} onChange={e=>setUsername(e.target.value)} placeholder="Emby 用户名" required /></label><label>密码<input type="password" value={password} onChange={e=>setPassword(e.target.value)} placeholder="Emby 密码" required /></label></div>
          <button className="connect-button" disabled={busy}>{busy?<LoaderCircle className="spin" size={17}/>:<Server size={17}/>} {busy?"正在连接…":"连接服务器"}</button>
          {connected && <p className="success">{connected}</p>}{error && <p className="error">{error}</p>}
          <p className="privacy">密码仅用于本次登录请求；当前版本尚未将登录状态持久化保存。</p>
        </form>
      </section>
      <footer>EMBER CINEMA <span>·</span> A PERSONAL MEDIA EXPERIENCE</footer>
    </main>
  </div>;
}
