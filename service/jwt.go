package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// 自定义Claims
type CustomClaims struct {
	Userid   int32  `json:"id"`       // 用户id
	Username string `json:"username"` // 用户昵称
	jwt.RegisteredClaims
}

// 令牌有效期，7天有效
const TokenExpireDuration = time.Hour * 24

// 密钥
var CustomSecret = []byte("89b7203b81a0273d")

// 生成jwt
func GenToken(userid int32, username string) (token, rToken string, err error) {

	// 创建声明
	claims := CustomClaims{
		userid,
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			Issuer:    "Sytalk Monitor",
		},
	}
	// 用指定的签名方法来创建token
	token, err = jwt.NewWithClaims(jwt.SigningMethodES256, claims).SignedString(CustomSecret)

	if err != nil {
		return
	}

	// 生成refresh token
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodES256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(TokenExpireDuration * 7).Unix(), // 过期时间，7天有效,
		Issuer:    "Sytalk Monitor",
	}).SignedString(CustomSecret)

	// 返回token和refresh token
	return
}

// 解析jwt
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 因为自定义了签名方法，所以用 parseWithClaims来解析
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return CustomSecret, nil
	})
	// 原生方法如下
	// token, err := jwt.Parse(tokenString, func (token *jwt.Token) (i interface{}, err error)) {
	//   return CustomSecret, nil
	// }
	if err != nil {
		return nil, err
	}
	// 对token对象中的claim进行类型断言
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("令牌解析错误")
}

// 刷新jwt令牌
func RefreshToken(token, rToken string) (newToken, newRToken string, err error) {
	// 解析当前的rtoken
	_, err = jwt.Parse(rToken, func(token *jwt.Token) (i interface{}, err error) {
		return CustomSecret, nil
	})

	if err != nil {
		return
	}

	var claims CustomClaims

	_, err = jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (i interface{}, err error) {
		return CustomSecret, nil
	})

	// 错误类型断言
	v, _ := err.(*jwt.ValidationError)

	// 如果错误类型是已过期，那么直接重新验发
	if v.Errors == jwt.ValidationErrorExpired {
		return GenToken(claims.Userid, claims.Username)
	}

	return
}
