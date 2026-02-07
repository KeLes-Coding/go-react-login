import { useState } from 'react'; // 1. 引入 useState
import { Link } from 'react-router-dom';

export default function Signup() {
    
    // 2. 定义状态： [变量名, 修改变量的方法] = useState(初始值)
    const [username, setUsername] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");

    const handleSubmit = async (e) => {
        e.preventDefault();
        console.log("准备提交注册:", { username, email, password });

        try {
            // 发送请求给 Go 后端
            const response = await fetch("http://localhost:8080/signup", {
                method: "POST", // 告诉后端：我要创建新数据
                headers: {
                    "Content-Type": "application/json", // 告诉后端：我发的是 JSON 格式
                },
                // 把我们的 state 数据转换成 JSON 字符串发过去
                body: JSON.stringify({ username, email, password }),
            });

            // 检查后端返回的状态码
            if (response.ok) {
                alert("🎉 注册成功！");
                // 这里以后可以加个自动跳转到登录页
            } else {
                // 如果后端返回 400 或 500
                const errorText = await response.text();
                alert("❌ 注册失败: " + errorText);
            }

        } catch (error) {
            // 如果连网都连不上（比如后端没开）
            console.error("请求错误:", error);
            alert("🔌 无法连接到服务器，请检查后端是否启动");
        }
    }

    return (
        <div>
            <h2>注册新账户</h2>
            <form onSubmit={handleSubmit}>
                <div className="form-group">
                    <label>用户名</label>
                    <input 
                        type="text" 
                        placeholder="比如: keles" 
                        value={username} // 3. 绑定值
                        onChange={(e) => setUsername(e.target.value)} // 4. 监听输入变化
                    />
                </div>

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

                <button type="submit" className="btn-primary">立即注册</button>
            </form>

            <p>
                已经有账号了？ <Link to="/login" style={{color: 'var(--primary-color)'}}>直接登录</Link>
            </p>
        </div>
    );
}