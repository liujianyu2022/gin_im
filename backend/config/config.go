package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App   AppConfig   `yaml:"app"`
	MySQL MySQLConfig `yaml:"mysql"`
	Redis RedisConfig `yaml:"redis"`
	JWT   JWTConfig   `yaml:"jwt"` // 确保有这个标签
}

type AppConfig struct {
	Name string
	Port string
	Mode string
}

type MySQLConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	PoolSize int
}

// 新增JWT配置结构体
type JWTConfig struct {
	Secret     string `yaml:"secret"`
	ExpireTime int    `yaml:"expire_time"` // 过期时间(小时)
}

var (
	once     sync.Once
	instance *Config
)

func LoadConfig(path string) *Config {
	once.Do(func() {

		// 读取文件内容
		data, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("Error reading config file: %v", err)
		}

		// 解析YAML
		instance = &Config{}
		if err := yaml.Unmarshal(data, instance); err != nil {
			log.Fatalf("Error parsing YAML config: %v", err)
		}
	})

	return instance
}

func (c *Config) GetMySQLDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.MySQL.User,
		c.MySQL.Password,
		c.MySQL.Host,
		c.MySQL.Port,
		c.MySQL.DBName,
	)
}

func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}
