import { BrowserRouter, Routes, Route, Link } from 'react-router-dom';
import Signup from './Signup';
import Login from './Login';
import Welcome from './Welcome';

function App() {
  return (
    <BrowserRouter>
      {/* 整个应用包裹在一个卡片容器里 */}
      <div className="app-container">
        
        {/* 顶部导航 */}
        <nav className="nav-header">
          <Link to="/signup" className="nav-link">注册</Link>
          <Link to="/login" className="nav-link">登录</Link>
          <Link to="/welcome" className="nav-link">欢迎页</Link>
        </nav>

        {/* 页面内容区域 */}
        <div className="content">
          <Routes>
            <Route path="/" element={<Login />} />
            <Route path="/signup" element={<Signup />} />
            <Route path="/login" element={<Login />} />
            <Route path="/welcome" element={<Welcome />} />
          </Routes>
        </div>

      </div>
    </BrowserRouter>
  );
}

export default App;