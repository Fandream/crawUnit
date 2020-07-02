package services

import (
	"crawunit/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)



func MakeToken(user *model.User) (tokenStr string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.Id,
		"username": user.Username,
		"nbf":      time.Now().Unix(),
	})

	return token.SignedString([]byte(viper.GetString("server.secret")))

}






func GetAdminUser() (user *model.User, err error) {
	u, err := model.GetUserByUsername("admin")
	if err != nil {
		return user, err
	}
	return &u, nil
}
