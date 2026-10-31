import { useState, useEffect, type FormEvent } from "react";
import { Clapperboard, Search, House, Film, Tv, Settings, Play, Server, LoaderCircle, Wifi, Star, RotateCcw } from "lucide-react";

type MediaItem = { id:string; name:string; type:string; overview:string; year:number; rating:number; imageUrl:string; seriesName:string };
type HomeData = { serverName:string; userName:string; views:Array<{ID:string;Name:string;CollectionType:string}>; latest:MediaItem[]; resume:MediaItem[] };
type MediaPage = {items:MediaItem[];total:number;start:number;limit:number};
declare global {
 interface Window {
  go?: { main?: { App?: {
   Connect:(url:string,username:string,password:string)=>Promise<{serverName:string;userName:string;userId:string}>;
   LoadHome:()=>Promise<HomeData>;
   Browse:(kind:string,parentID:string,start:number,limit:number,sortBy:string,sortOrder:string)=>Promise<MediaPage>;
   Search:(term:string,start:number,limit:number)=>Promise<MediaPage>;
   GetImageData:(id:string,tag:string)=>Promise<string>;
   GetItem:(id:string)=>Promise<MediaItem>;
  } } };
 }
}

function MediaCard({item,onOpen}: {item:MediaItem;onOpen:(item:MediaItem)=>void}) {
 const [poster,setPoster]=useState("");
 useEffect(()=>{let active=true;if(item.imageTag){window.go?.main?.App?.GetImageData(item.id,item.imageTag).then(src=>{if(active)setPoster(src)}).catch(()=>{});}return ()=>{active=false}},[item.id,item.imageTag]);
 return <article className="media-card" onClick={()=>onOpen(item)}>
  <div className="poster">{poster ? <img src={poster} loading="lazy" alt={item.name}/> : <div className="poster-fallback"><Clapperboard size={28}/></div>}
   {item.rating>0 && <span className="rating"><Star size={11} fill="currentColor"/> {item.rating.toFixed(1)}</span>}
  </div>
  <div className="media-title">{item.name}</div><div className="media-meta">{item.year || "年份未知"}{item.type==="Series"?" · 剧集":item.type==="Episode"?" · 单集":" · 电影"}</div>
 </article>;
}
function MediaRow({title,subtitle,items,onOpen}: {title:string;subtitle:string;items:MediaItem[];onOpen:(item:MediaItem)=>void}) {
 return <section className="media-section"><div className="media-heading"><div><h2>{title}</h2><p>{subtitle}</p></div><span>{items.length} 部</span></div>
  {items.length ? <div className="media-grid">{items.map(item=><MediaCard key={item.id} item={item} onOpen={onOpen}/>)}</div> : <div className="empty-state">暂时没有内容</div>}
 </section>;
}

export default function App() {
 const [server,setServer]=useState(""); const [username,setUsername]=useState(""); const [password,setPassword]=useState("");
 const [busy,setBusy]=useState(false); const [connected,setConnected]=useState(""); const [error,setError]=useState("");
 const [home,setHome]=useState<HomeData|null>(null); const [loadingHome,setLoadingHome]=useState(false);
 async function loadHome() {
  const api=window.go?.main?.App;if(!api)throw new Error("当前不是 Wails 桌面运行环境");
  setLoadingHome(true);
  try { const data=await api.LoadHome();setHome(data); }
  finally { setLoadingHome(false); }
 }
 async function connect(e:FormEvent) {
  e.preventDefault();setBusy(true);setError("");setConnected("");setHome(null);
  try {
   const api=window.go?.main?.App;if(!api)throw new Error("请通过 Wails 启动桌面客户端");
   const result=await api.Connect(server,username,password);
   setConnected("已连接："+result.serverName+" · "+result.userName);
   await loadHome();
  } catch(err) { setError(err instanceof Error?err.message:"连接失败"); }
  finally { setBusy(false); }
 }
 return <div className="shell">
  <aside className="sidebar">
   <div className="brand"><span className="brand-icon"><Clapperboard size={21}/></span><span>EMBER<span className="brand-light"> CINEMA</span></span></div>
   <div className="nav-label">媒体空间</div><a className="nav active"><House size={18}/> 首页</a><a className="nav"><Film size={18}/> 电影</a><a className="nav"><Tv size={18}/> 电视剧</a>
   <div className="sidebar-bottom"><a className="nav"><Settings size={18}/> 设置</a><div className="server-state"><span className={home?"dot online":"dot"}/>{home?"已连接 Emby":"尚未连接服务器"}</div></div>
  </aside>
  <main className="main">
   <header className="topbar"><div className="crumb">你的私人影院 <span>/</span> 首页</div><div className="top-actions">{home&&<span className="connected-pill"><Wifi size={13}/> {home.serverName}</span>}<button className="icon-button" title="搜索"><Search size={19}/></button></div></header>
   {!home ? <>
    <section className="hero">
     <div className="hero-copy"><div className="eyebrow">YOUR PERSONAL CINEMA</div><h1>把热爱，<br/><em>留在私人影院。</em></h1><p>连接你的 Emby 媒体服务器，开始探索私人收藏。</p><a className="hero-link" href="#connect"><Play size={15} fill="currentColor"/> 开始连接 <span>↗</span></a></div>
     <div className="hero-art"><div className="orb orb-one"/><div className="orb orb-two"/><div className="art-card"><Clapperboard size={58} strokeWidth={1}/><span>YOUR LIBRARY<br/>YOUR CINEMA</span></div></div>
     <div className="hero-index">01 <span>/ 03</span></div>
    </section>
    <section className="connect-section" id="connect">
     <div className="section-heading"><div><div className="eyebrow">GET STARTED</div><h2>连接 Emby 服务器</h2><p>使用你的 Emby 服务器地址和账户登录。</p></div><div className="secure"><Server size={15}/> 本地直连 · 不经过第三方</div></div>
     <form className="connect-card" onSubmit={connect}>
      <label>服务器地址<input value={server} onChange={e=>setServer(e.target.value)} placeholder="例如：http://192.168.1.100:8096" required/></label>
      <div className="form-row"><label>用户名<input value={username} onChange={e=>setUsername(e.target.value)} placeholder="Emby 用户名" required/></label><label>密码<input type="password" value={password} onChange={e=>setPassword(e.target.value)} placeholder="Emby 密码" required/></label></div>
      <button className="connect-button" disabled={busy}>{busy?<LoaderCircle className="spin" size={17}/>:<Server size={17}/>} {busy?"正在连接并读取媒体库…":"连接服务器"}</button>
      {connected&&<p className="success">{connected}</p>}{error&&<p className="error">{error}</p>}
      <p className="privacy">密码仅用于登录请求；当前版本暂不持久化保存登录状态。</p>
     </form>
    </section>
   </> : <>
    <section className="welcome-banner"><div><div className="eyebrow">WELCOME BACK</div><h1>今晚，想看点什么？</h1><p>已连接到 {home.serverName}，你的私人片库已就绪。</p></div><button className="refresh-button" onClick={()=>loadHome()} disabled={loadingHome}>{loadingHome?<LoaderCircle className="spin" size={16}/>:<RotateCcw size={16}/>} 刷新媒体库</button></section>
    <section className="library-stats">{home.views.map(v=><div className="library-chip" key={v.ID}><Film size={15}/><span>{v.Name}</span></div>)}{home.views.length===0&&<span className="muted">服务器尚未返回媒体库分类</span>}</section>
    {home.resume.length>0&&<MediaRow title="继续观看" subtitle="接着上次的进度，继续你的故事" items={home.resume}/>}
    <MediaRow title="最近添加" subtitle="你的媒体库刚刚上新" items={home.latest}/>
    {error&&<p className="error">{error}</p>}
   </>}
   <footer>EMBER CINEMA <span>·</span> A PERSONAL MEDIA EXPERIENCE</footer>
  </main>
 </div>;
}
