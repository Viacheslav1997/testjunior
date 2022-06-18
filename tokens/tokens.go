package tokens

import (
	"encoding/base64"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

var SECRET_KEY string = "shv2w7b"

type PayloadData struct {
	User_id string
	jwt.StandardClaims
}

func Bcrypt(pass string) (Encrypted string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), 12)
	if err != nil {
		log.Panic(err)
		return
	}
	return string(hash)
}

func GenerateAllTokens(uid string) (signedToken string, signedRefreshToken string) {
	claims := &PayloadData{

		User_id: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	//refreshClaims := &PayloadData{
	//	StandardClaims: jwt.StandardClaims{
	//		ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
	//	},
	//}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString([]byte(SECRET_KEY))
	refreshToken := uuid.New()
	refreshTokenBase64 := base64.StdEncoding.EncodeToString([]byte(refreshToken.String()))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshTokenBase64
}
