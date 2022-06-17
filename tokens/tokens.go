package tokens

import (
	"github.com/golang-jwt/jwt"
	"log"
	"time"
)

var SECRET_KEY string = "shv2w7b"

type PayloadData struct {
	User_id string
	jwt.StandardClaims
}

// GenerateAllTokens generates both teh detailed token and refresh token
func GenerateAllTokens(uid string) (signedToken string, signedRefreshToken string) {
	claims := &PayloadData{

		User_id: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &PayloadData{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS512, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken
}
