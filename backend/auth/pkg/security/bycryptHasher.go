package security

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
	cost int // Сложность хеширования (от 4 до 31) Рекомендуемое значение: 10-12
}

func NewBcryptHasher(cost int) *BcryptHasher {
	return &BcryptHasher{cost: cost}
}

func (h *BcryptHasher) Hash(password string) (string, error) {
	// 1. Генерируем "соль" и хушируем пароль
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("hashing failed: %w", err)
	}
	return string(bytes), nil
}

func (h *BcryptHasher) Check(hashedPassword, inputPassword string) bool {
	// Сравниваем хеш с введенным паролем
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(inputPassword),
	)
	return err == nil
}
