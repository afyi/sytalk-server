package model

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name" json:"name" validate:"required"`
	Password  string             `bson:"pass" json:"pass" validate:"required"`
	Nickname  string             `bson:"nick" json:"nick"`
	Avatar    string             `bson:"avatar" json:"avatar"`
	Salt      string             `bson:"salt" json:"salt"`
	CreatedAt time.Time          `bson:"CreatedAt" json:"CreatedAt"`
	UpdatedAt time.Time          `bson:"UpdatedAt" json:"UpdatedAt"`
}

// 生成盐
func (User *User) GetSalt() string {
	b := make([]byte, 2)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		panic("生成盐失败")
	}
	return fmt.Sprintf("%x", b)
}

// 加密
func (User *User) GetSysMd5(password, salt string) string {

	m5 := md5.New()

	pass := []byte(password + salt)

	m5.Write(pass)

	return hex.EncodeToString(m5.Sum(nil))
}
