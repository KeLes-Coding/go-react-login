package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-react-login/backend/internal/models"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// 定义一个密钥 (真实项目中应从环境变量获取)
var JwtKey = []byte("my_secret_key")

// Handler 结构体用于持有数据库连接
type Handler struct {
	DB *sql.DB
}

// 辅助函数：创建用户
func (h *Handler) createUser(user models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = h.DB.Exec(
		"INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", user.Username, user.Email, string(hashedPassword),
	)
	return err
}

// SignupHandler 处理注册
func (h *Handler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	err = h.createUser(u)
	if err != nil {
		http.Error(w, "创建用户失败", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("用户创建成功"))
}

// LoginHandler 处理登录
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	var storeHash string
	var storeUsername string

	err = h.DB.QueryRow("SELECT username, password FROM users WHERE email = $1", creds.Email).Scan(&storeUsername, &storeHash)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "用户不存在", http.StatusUnauthorized)
			return
		}
		http.Error(w, "服务器错误", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storeHash), []byte(creds.Password))
	if err != nil {
		http.Error(w, "密码错误", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &models.Claims{
		Username: storeUsername,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		http.Error(w, "生成 Token 失败", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

// WelcomeHandler 受保护的路由
func (h *Handler) WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 Authorization 头部信息
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "未提供 Token", http.StatusUnauthorized)
		return
	}

	// 2. 移除 "Bearer " 前缀，提取纯 Token
	// 前端通常发送格式为: "Bearer eyJhbGci..."
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "Token 格式错误 (需要 Bearer 前缀)", http.StatusUnauthorized)
		return
	}
	tokenString := parts[1]

	// 3. 解析 Token
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil || !token.Valid {
		// 调试建议：如果还报错，可以 fmt.Println(err) 查看具体原因
		fmt.Println("❌ Token 验证失败:", err)
		http.Error(w, "Token 无效或已过期", http.StatusUnauthorized)
		return
	}

	w.Write([]byte(fmt.Sprintf("Welcome, %s!", claims.Username)))
}
