package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	MongoDB  MongoDBConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Storage  StorageConfig
	Video    VideoConfig
	RateLimit RateLimitConfig
	CORS     CORSConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type MongoDBConfig struct {
	URI      string
	Database string
}

type RedisConfig struct {
	URL      string
	Password string
}

type JWTConfig struct {
	Secret             string
	Expiry             time.Duration
	RefreshTokenExpiry time.Duration
}

type StorageConfig struct {
	Provider        string // "s3" or "minio"
	AWSRegion       string
	AWSAccessKey    string
	AWSSecretKey    string
	S3Bucket        string
	S3Endpoint      string
	MinioEndpoint   string
	MinioAccessKey  string
	MinioSecretKey  string
	MinioUseSSL     bool
}

type VideoConfig struct {
	MaxSizeMB          int
	MaxDurationSeconds int
	AllowedFormats     []string
}

type RateLimitConfig struct {
	Requests int
	Window   time.Duration
}

type CORSConfig struct {
	AllowedOrigins []string
}

var AppConfig *Config

func Load() *Config {
	// Load .env file if it exists
	_ = godotenv.Load()

	jwtExpiry, _ := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))
	refreshExpiry, _ := time.ParseDuration(getEnv("REFRESH_TOKEN_EXPIRY", "168h"))
	rateLimitWindow, _ := time.ParseDuration(getEnv("RATE_LIMIT_WINDOW", "60s"))

	maxSizeMB, _ := strconv.Atoi(getEnv("MAX_VIDEO_SIZE_MB", "100"))
	maxDuration, _ := strconv.Atoi(getEnv("MAX_VIDEO_DURATION_SECONDS", "180"))
	rateLimitReqs, _ := strconv.Atoi(getEnv("RATE_LIMIT_REQUESTS", "100"))

	AppConfig = &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGODB_DATABASE", "magicchat"),
		},
		Redis: RedisConfig{
			URL:      getEnv("REDIS_URL", "redis://localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		JWT: JWTConfig{
			Secret:             getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production"),
			Expiry:             jwtExpiry,
			RefreshTokenExpiry: refreshExpiry,
		},
		Storage: StorageConfig{
			Provider:       getEnv("STORAGE_PROVIDER", "minio"),
			AWSRegion:      getEnv("AWS_REGION", "us-east-1"),
			AWSAccessKey:   getEnv("AWS_ACCESS_KEY_ID", ""),
			AWSSecretKey:   getEnv("AWS_SECRET_ACCESS_KEY", ""),
			S3Bucket:       getEnv("S3_BUCKET", "magicchat-videos"),
			S3Endpoint:     getEnv("S3_ENDPOINT", ""),
			MinioEndpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
			MinioAccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			MinioSecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
			MinioUseSSL:    getEnv("MINIO_USE_SSL", "false") == "true",
		},
		Video: VideoConfig{
			MaxSizeMB:          maxSizeMB,
			MaxDurationSeconds: maxDuration,
			AllowedFormats:     strings.Split(getEnv("ALLOWED_VIDEO_FORMATS", "mp4,mov,avi,webm"), ","),
		},
		RateLimit: RateLimitConfig{
			Requests: rateLimitReqs,
			Window:   rateLimitWindow,
		},
		CORS: CORSConfig{
			AllowedOrigins: strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"), ","),
		},
	}

	log.Println("Configuration loaded successfully")
	return AppConfig
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
