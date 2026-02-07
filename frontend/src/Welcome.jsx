import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';

export default function Welcome() {
    const [message, setMessage] = useState("正在验证身份...");
    const navigate = useNavigate();

    useEffect(() => {
        // 1. 从本地存储获取 Token
        const token = localStorage.getItem("token");

        // 2. 如果没有 Token，直接跳转回登录页
        if (!token) {
            alert("🔒 请先登录！");
            navigate("/login");
            return;
        }

        // 3. 发送请求，带上 Token
        const fetchWelcome = async () => {
            try {
                const response = await fetch("http://localhost:8080/welcome", {
                    method: "GET",
                    headers: {
                        // 👇 关键点：把 Token 放在这里传给后端
                        "Authorization": `Bearer ${token}`,
                    },
                });

                if (response.ok) {
                    const text = await response.text();
                    setMessage(text); // 显示 "Welcome, <用户名>!"
                } else {
                    // 如果 Token 过期或无效
                    alert("🚫 会话已过期，请重新登录");
                    localStorage.removeItem("token"); // 清理掉无效的 token
                    navigate("/login");
                }
            } catch (error) {
                console.error("请求错误:", error);
                setMessage("🔌 无法连接服务器");
            }
        };

        fetchWelcome();
    }, [navigate]); // 空依赖数组表示只在组件加载时执行一次

    return (
        <div style={{ textAlign: 'center' }}>
            <h1>{message}</h1>
            {/* 加个退出按钮方便测试 */}
            <button 
                className="btn-primary" 
                style={{ maxWidth: '200px', marginTop: '2rem' }}
                onClick={() => {
                    localStorage.removeItem("token");
                    navigate("/login");
                }}
            >
                退出登录
            </button>
        </div>
    );
}