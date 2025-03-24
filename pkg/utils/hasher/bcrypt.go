package hasher

import (
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
	cost int
}

func NewBcryptHasher(cfg *viper.Viper) *BcryptHasher {
	return &BcryptHasher{cost: cfg.GetInt("auth.bcrypt_cost")}
}

func (h *BcryptHasher) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (h *BcryptHasher) ComparePasswords(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
