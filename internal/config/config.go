package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Cache    CacheConfig
	Storage  StorageConfig
	LLM      LLMConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string
	Port int
	Mode      string // debug, release, test
	RateLimit float64
	RateBurst int
	JWTSecret string
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int // seconds
	ConnMaxIdleTime int // seconds
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Enabled     bool
	MaxCost     int64 // 最大内存占用（字节）
	NumCounters int64
	BufferItems int64
	TTL         int // 缓存过期时间（秒）
}

// StorageConfig 存储配置
type StorageConfig struct {
	UploadPath string
	MaxSize    int64 // 最大上传大小（字节）
}

// LLMConfig 大模型配置
type LLMConfig struct {
	BaseURL string
	APIKey  string
	Model   string
}

// Load 从环境变量加载配置
func Load() *Config {
	// 加载 .env 文件（如果存在）
	if err := godotenv.Load(); err != nil {
		log.Println("未找到 .env 文件，使用系统环境变量")
	}
	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnvInt("SERVER_PORT", 8080),
			Mode: getEnv("SERVER_MODE", "release"),
			RateLimit: float64(getEnvInt("SERVER_RATE_LIMIT", 100)),
			RateBurst: getEnvInt("SERVER_RATE_BURST", 200),
			JWTSecret: getEnv("JWT_SECRET", ""),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnvInt("DB_PORT", 5432),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", ""),
			Database:        getEnv("DB_NAME", "litecms"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 10),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvInt("DB_CONN_MAX_LIFETIME", 1800),
			ConnMaxIdleTime: getEnvInt("DB_CONN_MAX_IDLE_TIME", 300),
		},
		Cache: CacheConfig{
			Enabled:     getEnvBool("CACHE_ENABLED", true),
			MaxCost:     getEnvInt64("CACHE_MAX_COST", 50*1024*1024), // 50MB
			NumCounters: getEnvInt64("CACHE_NUM_COUNTERS", 100000),
			BufferItems: getEnvInt64("CACHE_BUFFER_ITEMS", 64),
			TTL:         getEnvInt("CACHE_TTL", 300), // 5分钟
		},
		Storage: StorageConfig{
			UploadPath: getEnv("STORAGE_UPLOAD_PATH", "./static/uploads"),
			MaxSize:    getEnvInt64("STORAGE_MAX_SIZE", 5*1024*1024), // 5MB
		},
		LLM: LLMConfig{
			BaseURL: getEnv("LLM_BASE_URL", ""),
			APIKey:  getEnv("LLM_API_KEY", ""),
			Model:   getEnv("LLM_MODEL", "gpt-3.5-turbo"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.ParseInt(value, 10, 64); err == nil {
			return i
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return defaultValue
}
