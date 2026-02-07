package main

import (
	"fmt"
	"go-react-login/backend/internal/database"
	"go-react-login/backend/internal/handlers"
	"go-react-login/backend/internal/middleware"
	"log"
	"net/http"
	// 引入我们拆分出去的包
	// 请将 login-demo 替换为你 go.mod 中的 module 名
)

func main() {
	// 1. 初始化数据库
	db := database.InitDB()
	defer db.Close()

	// 2. 初始化 Handler 实例
	h := &handlers.Handler{DB: db}

	// 3. 注册路由 (使用中间件包裹)
	http.HandleFunc("/signup", middleware.EnableCORS(h.SignupHandler))
	http.HandleFunc("/login", middleware.EnableCORS(h.LoginHandler))
	http.HandleFunc("/welcome", middleware.EnableCORS(h.WelcomeHandler))

	// 4. 启动服务
	fmt.Println("服务器正在监听 :8080 (重构版)")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
