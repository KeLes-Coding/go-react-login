package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// 配置数据库连接参数
const (
	host     = "localhost"
	port     = 5433
	user     = "keles"
	password = "c05022007"
	dbname   = "login_demo"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Env struct {
	db *sql.DB
}

func CreateUser(db *sql.DB, user User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		"INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", user.Username, user.Email, string(hashedPassword),
	)
	if err != nil {
		return err
	}
	return nil
}

func (env *Env) signupHandler(w http.ResponseWriter, r *http.Request) {
	var u User

	// 从请求体中读取 JSON 数据
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	err = CreateUser(env.db, u)
	if err != nil {
		http.Error(w, "创建用户失败", http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("用户创建成功"))
	}
}

func (env *Env) loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds LoginRequest
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	var storeHash string
	var storeUsername string

	err = env.db.QueryRow("SELECT username, password FROM users WHERE email = $1", creds.Email).Scan(&storeUsername, &storeHash)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "用户不存在", http.StatusUnauthorized)
			return
		}
		fmt.Println("数据库查询错误:", err)
		http.Error(w, "服务器错误", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storeHash), []byte(creds.Password))
	if err != nil {
		http.Error(w, "密码错误", http.StatusUnauthorized)
		return
	}

	// 生成 JWT 令牌
	// 1. 设置过期时间
	expirationTime := time.Now().Add(5 * time.Minute)

	// 2. 创建 Claims
	claims := &Claims{
		Username: storeUsername,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// 3. 使用 HS256 算法和我们的密钥进行签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Println("❌ Token 生成错误详情:", err)
		http.Error(w, "生成 Token 失败", http.StatusInternalServerError)
		return
	}

	// 4. 将 Token 发送回客户端
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

func (env *Env) welcomeHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "未提供 Token", http.StatusUnauthorized)
		return
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrTokenExpired {
			http.Error(w, "Token 已过期", http.StatusUnauthorized)
			return
		}
		http.Error(w, "无效的 Token", http.StatusUnauthorized)
		return
	}

	if !token.Valid {
		http.Error(w, "Token 无效", http.StatusUnauthorized)
		return
	}
	w.Write([]byte(fmt.Sprintf("Welcome, %s!", claims.Username)))
}

// enableCORS 是一个中间件，用来给响应头添加跨域允许的标记
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 允许任何来源访问 (在生产环境中，这里应该换成具体的 React 域名)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// 允许的请求方法
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		// 允许的请求头 (特别要注意 Authorization，因为我们要传 Token)
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// 浏览器在发送 POST 之前，会先发一个 OPTIONS 请求来“探路”
		// 如果是 OPTIONS 请求，直接返回 OK，不用进入业务逻辑
		if r.Method == "OPTIONS" {
			return
		}

		// 如果不是 OPTIONS，就继续执行原本的处理函数
		next(w, r)
	}
}

func main() {
	// 1. 构建连接字符库
	// sslmode=disable 是因为本地开发通常没有配置 SSL 证书
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// --- 调试代码 START ---
	// 打印出实际生成的连接字符串，检查 dbname 是否正确
	fmt.Println("------------ 调试信息 ------------")
	fmt.Println("正在尝试连接，配置信息如下：")
	fmt.Println(psqlInfo)
	fmt.Println("----------------------------------")
	// --- 调试代码 END ---

	// 2. 打开数据库连接
	db, err := sql.Open("postgres", psqlInfo)
	env := &Env{db: db}
	if err != nil {
		log.Fatal(err)
	}
	// 函数结束时关闭数据库连接
	defer db.Close()

	// 3. Ping
	err = db.Ping()
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}
	fmt.Println("成功连接到数据库")

	// 4. 定义建表 SQL
	// 注意 password 我们给 255 长度，因为稍后我们会用 bcrypt 加密，生成的哈希值较长
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL, 
			password TEXT NOT NULL
		);`

	// 5. 执行 SQL 创建表
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("创建表失败：%v", err)
	}

	// 6. 插入测试用户
	// u := User{
	// 	Username: "admin",
	// 	Password: "123456",
	// }
	// err = CreateUser(db, u)
	// if err != nil {
	// 	log.Fatalf("创建用户失败：%v", err)
	// } else {
	// 	fmt.Println("用户 admin 已成功创建")
	// }
	// fmt.Println("成功: users 表已就绪 (如果表不存在则已创建)")

	// 7. 注册路由
	http.HandleFunc("/signup", enableCORS(env.signupHandler))
	http.HandleFunc("/login", enableCORS(env.loginHandler))
	http.HandleFunc("/welcome", enableCORS(env.welcomeHandler))

	// 8. 启动监听
	fmt.Println("服务器正在监听 :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
