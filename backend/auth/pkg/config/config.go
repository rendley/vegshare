package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

// Загрузчик конфигов

// Config - структура для хранение всех конфигов.
type Config struct {
	HTTP     HTTPConfig     `yaml:"http"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Redis    RedisConfig    `yaml:"redis"`
}

type HTTPConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	DBName   string `yaml:"dbname"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"sslmode"`
}

type JWTConfig struct {
	Secret    string `yaml:"secret"`
	ExpiresIn string `yaml:"expires_in"`
}

type RedisConfig struct {
	URL string `yaml:"url"`
	TTL int    `yaml:"ttl"`
}

// Load загружает конфиг из YAML-файла.
func Load() *Config {
	// Путь к файлу конфига (относительно корня проекта).
	configPath := filepath.Join("configs", "config.yaml")

	// Читаем файл.
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	// Декодируем YAML в структуру.
	// Принимаеm YAML-данные (в виде []byte).
	// Заполняет структуру Config значениями из YAML.
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Failed ro parse YAML config: %v", err)
	}

	// Валидация (пример).
	if cfg.JWT.Secret == "" {
		log.Fatal("JWT secret is required")
	}

	return &cfg
}
