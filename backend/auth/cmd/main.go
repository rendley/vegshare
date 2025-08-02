package main

import (
	"fmt"
	"github.com/rendley/auth/pkg/config"
)

func main() {
	// Загружаем конфиги (порт, секреты)
	cfg := config.Load()
	fmt.Printf("Config: %+v\n", cfg)

}
