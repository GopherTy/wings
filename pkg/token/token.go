package token

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gopherty/wings/common/conf"
)

// Token jwt token
type Token struct {
	Access    string // access token
	Refresh   string // refresh token
	AtExpires int64  // access token expire time
	RtExpires int64  // refresh token expire time
}

// New generate a token
func New(userID uint64) (t *Token, err error) {
	t = new(Token)
	t.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	t.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()

	atClaims := &jwt.MapClaims{
		"authorized": true,
		"user_id":    userID,
		"exp":        t.AtExpires,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	t.Access, err = accessToken.SignedString([]byte(conf.Instance().Secret.Access))
	if err != nil {
		return
	}

	rtClaims := &jwt.MapClaims{
		"authorized": true,
		"user_id":    userID,
		"exp":        t.RtExpires,
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	t.Refresh, err = refreshToken.SignedString([]byte(conf.Instance().Secret.Refresh))
	if err != nil {
		return
	}
	return
}
