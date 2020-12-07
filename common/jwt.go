package common

import (
	"com/mittacy/gomeet/config"
	"com/mittacy/gomeet/model"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// GenerateToken 生成token函数
func GenerateToken(session *model.Session) (string, error) {
	// 获取参数
	now := time.Now()
	num, err := config.Cfg.Section("jwt").Key("expire").Int()
	if err != nil {
		return "", err
	}
	expireTime := now.Add(time.Duration(num) * time.Hour)

	name := config.Cfg.Section("jwt").Key("tokenName").String()

	// 创建Token
	claims := model.Session {
		Sno: session.Sno,
		Phone: session.Phone,
		IsAdmin: session.IsAdmin,
		IsRoot: session.IsRoot,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    name,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := []byte(config.Cfg.Section("jwt").Key("secret").String())

	return tokenClaims.SignedString(jwtSecret)
}

// ParseToken 解析token返回保存信息
func ParseToken(token string) (*model.Session, error) {
	jwtSecret := []byte(config.Cfg.Section("jwt").Key("secret").String())
	tokenClaims, err := jwt.ParseWithClaims(token, &model.Session{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if session, ok := tokenClaims.Claims.(*model.Session); ok && tokenClaims.Valid {
			return session, nil
		}
	}
	return nil, err
}
