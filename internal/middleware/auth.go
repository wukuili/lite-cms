package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

const maxVisitors = 10000 // 限流器map最大容量，防止分布式攻击匸尽内存

var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex
)

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		mode := os.Getenv("SERVER_MODE")
		if mode == "release" {
			log.Fatal("生产环境必须设置 JWT_SECRET 环境变量")
		}
		log.Println("警告: 未设置 JWT_SECRET，使用默认密钥（仅限开发环境）")
		secret = "lite-cms-secret-key-change-in-production"
	}
	jwtSecret = []byte(secret)

	go cleanupVisitors()
}

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	
	v, exists := visitors[ip]
	if !exists {
		// 容量检查：超过上限时先紧急清理
		if len(visitors) >= maxVisitors {
			now := time.Now()
			for ip, v := range visitors {
				if now.Sub(v.lastSeen) > time.Minute {
					delete(visitors, ip)
				}
			}
			// 如果清理后仍超过限制，返回一个已耗尽的limiter来拒绝新IP
			if len(visitors) >= maxVisitors {
				return rate.NewLimiter(0, 0) // 拒绝
			}
		}
		// 允许每秒10个请求，突发容量为20
		v = &visitor{limiter: rate.NewLimiter(10, 20)}
		visitors[ip] = v
	}
	v.lastSeen = time.Now()
	return v.limiter
}

func cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

// JWT密钥（生产环境应从配置读取）
var jwtSecret []byte

// Claims JWT声明
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// AuthRequired 认证中间件
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			tokenString, _ = c.Cookie("token")
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			c.Abort()
			return
		}

		// 移除Bearer前缀
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*Claims); ok {
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("role", claims.Role)
		} else {
			// 如果 *Claims 转型失败，记录原因
			log.Printf("Token claims 转型失败: type=%T, claims=%+v", token.Claims, token.Claims)
		}

		c.Next()
	}
}

// AdminRequired 管理员权限中间件
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "权限验证失败"})
			c.Abort()
			return
		}
		
		roleStr := strings.ToLower(role.(string))
		if roleStr != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// GenerateToken 生成JWT Token
func GenerateToken(userID uint, username, role string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return gin.Recovery()
}

// RateLimit 简单限流中间件
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := getLimiter(c.ClientIP())
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "请求太频繁，请稍后再试"})
			c.Abort()
			return
		}
		c.Next()
	}
}
