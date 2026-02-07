import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';

export default function Login() {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    
    // 使用 hook 来进行页面跳转
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        console.log("准备登录:", { email, password });

        try {
            const response = await fetch("http://localhost:8080/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ email, password }),
            });

            if (response.ok) {
                // 1. 解析数据
                const data = await response.json();
                
                // 👇 新增：看看后端到底给了啥？
                console.log("🔍 调试 - 后端返回的数据:", data);

                // 2. 检查是否有 token 字段
                if (data.token) {
                    localStorage.setItem("token", data.token);
                    console.log("✅ Token 已保存到本地:", data.token); // 确认保存动作
                    
                    alert("🎉 登录成功！");
                    navigate("/welcome");
                } else {
                    console.error("❌ 严重错误: 后端返回了 200 OK，但数据里没有 token 字段！");
                    alert("登录异常：未收到令牌");
                }
            } else {
                const errorText = await response.text();
                alert("❌ 登录失败: " + errorText);
            }

        } catch (error) {
            console.error("请求错误:", error);
            alert("🔌 无法连接到服务器");
        }
    }

    return (
        <div>
            <h2>登录你的账户</h2>
            <form onSubmit={handleSubmit}>
                <div className="form-group">
                    <label>邮箱</label>
                    <input 
                        type="email" 
                        placeholder="example@mail.com" 
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                    />
                </div>

                <div className="form-group">
                    <label>密码</label>
                    <input 
                        type="password" 
                        placeholder="**" 
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                    />
                </div>

                <button type="submit" className="btn-primary">登录</button>
            </form>

            <p>
                还没有账号？ <Link to="/signup" style={{color: 'var(--primary-color)'}}>去注册</Link>
            </p>
        </div>
    );
}