package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Password string
	Role     string
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(hashedPass, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password))
	return err == nil
}
