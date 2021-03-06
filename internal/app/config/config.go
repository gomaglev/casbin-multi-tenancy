package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/koding/multiconfig"
)

var (
	C    = new(Config)
	once sync.Once
)

// MustLoad load configs
func MustLoad(fpaths ...string) {
	once.Do(func() {
		loaders := []multiconfig.Loader{
			&multiconfig.TagLoader{},
			&multiconfig.EnvironmentLoader{},
		}

		for _, fpath := range fpaths {
			if strings.HasSuffix(fpath, "toml") {
				loaders = append(loaders, &multiconfig.TOMLLoader{Path: fpath})
			}
			if strings.HasSuffix(fpath, "json") {
				loaders = append(loaders, &multiconfig.JSONLoader{Path: fpath})
			}
			if strings.HasSuffix(fpath, "yaml") {
				loaders = append(loaders, &multiconfig.YAMLLoader{Path: fpath})
			}
		}

		m := multiconfig.DefaultLoader{
			Loader:    multiconfig.MultiLoader(loaders...),
			Validator: multiconfig.MultiValidator(&multiconfig.RequiredValidator{}),
		}
		m.MustLoad(C)
	})
}

// PrintWithJSON prints configs in JSON format
func PrintWithJSON() {
	if C.PrintConfig {
		b, err := json.MarshalIndent(C, "", " ")
		if err != nil {
			os.Stdout.WriteString("[CONFIG] JSON marshal error: " + err.Error())
			return
		}
		os.Stdout.WriteString(string(b) + "\n")
	}
}

// LogHook 日志钩子
type LogHook string

// IsGorm is gorm log hook
func (h LogHook) IsGorm() bool {
	return h == "gorm"
}

// IsMongo is mongo log hook
func (h LogHook) IsMongo() bool {
	return h == "mongo"
}

// Log log params
type Log struct {
	Level         int
	Format        string
	Output        string
	OutputFile    string
	EnableHook    bool
	HookLevels    []string
	Hook          LogHook
	HookMaxThread int
	HookMaxBuffer int
}

// LogGormHook
type LogGormHook struct {
	DBType       string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
	Table        string
}

// LogMongoHook
type LogMongoHook struct {
	Collection string
}

// Gateway
type Gateway struct {
	Host            string
	Port            int
	CertFile        string
	KeyFile         string
	ShutdownTimeout int
	PathPrefix      string
	Enable          bool
}

// Monitor
type Monitor struct {
	Enable    bool
	Addr      string
	ConfigDir string
}

// Captcha
type Captcha struct {
	Store       string
	Length      int
	Width       int
	Height      int
	RedisDB     int
	RedisPrefix string
}

// RateLimiter
type RateLimiter struct {
	Enable  bool
	Count   int64
	RedisDB int
}

// CORS
type CORS struct {
	Enable           bool
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	MaxAge           int
}

// GZIP
type GZIP struct {
	Enable             bool
	ExcludedExtentions []string
	ExcludedPaths      []string
}

// Gorm
type Gorm struct {
	Debug             bool
	DBType            string
	MaxLifetime       int
	MaxOpenConns      int
	MaxIdleConns      int
	TablePrefix       string
	EnableAutoMigrate bool
	Timeout           time.Duration
}

// MySQL
type MySQL struct {
	Host       string
	Port       int
	User       string
	Password   string
	DBName     string
	Parameters string
}

// DSN
func (a MySQL) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		a.User, a.Password, a.Host, a.Port, a.DBName, a.Parameters)
}

// Postgres
type Postgres struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// DSN
func (a Postgres) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		a.Host, a.Port, a.User, a.DBName, a.Password, a.SSLMode)
}

// Sqlite3
type Sqlite3 struct {
	Path string
}

// DSN
func (a Sqlite3) DSN() string {
	return a.Path
}

// Mongo
type Mongo struct {
	URI              string
	Database         string
	Timeout          time.Duration
	CollectionPrefix string
}

type GRPC struct {
	Host            string
	Port            int
	CertFile        string
	KeyFile         string
	RateLimitCount  int
	ShutdownTimeout time.Duration
}

// Config
type Config struct {
	RunMode     string
	WWW         string
	Swagger     bool
	PrintConfig bool
	GRPC        GRPC
	Gateway     Gateway
	Interceptor Interceptor
	Monitor     Monitor
	BasicAuth   BasicAuth
	Authorizer  Authorizer

	Log          Log
	LogGormHook  LogGormHook
	LogMongoHook LogMongoHook
	RateLimiter  RateLimiter
	CORS         CORS
	Redis        Redis
	Gorm         Gorm
	MySQL        MySQL
	Postgres     Postgres
	Sqlite3      Sqlite3
	Mongo        Mongo
	UniqueID     struct {
		Type      string
		Snowflake struct {
			Node  int64
			Epoch int64
		}
	}
	DefaultLang string
}

type BasicAuth struct {
	User      string
	Password  string
	AuthToken string
}

type Redis struct {
	Host        string
	Port        int
	Auth        string
	Key         string
	Expire      string
	ExpireT     time.Duration
	DialTimeout time.Duration
}

type ES struct {
	Hosts         string
	DeviceIndex   string
	SensorIndex   string
	ReadingsIndex string
}

type Interceptor struct {
	EnableRateLimit      bool
	EnableLogrus         bool
	EnableRecovery       bool
	EnableAuthentication bool
	EnableAuthorization  bool
}

type Authorizer struct{}

// Enforce resources
func (a *Authorizer) Enforce() func(ctx context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		// Not used in current version
		return ctx, nil
	}
}
