package tokens

import (
	"encoding/base64"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"

	//"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

var SECRET_KEY string = "shv2w7b"

type PayloadData struct {
	User_id string
	Uniq    string

	jwt.StandardClaims
}

type RefreshPayloadData struct {
	Uniq string
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

	uniqSessionId := uuid.New().String()

	refreshClaims := &RefreshPayloadData{

		Uniq: uniqSessionId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS512, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}
	refreshTokenBase64 := base64.StdEncoding.EncodeToString([]byte(refreshToken))

	claims := &PayloadData{

		User_id: uid,
		Uniq:    uniqSessionId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(1)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshTokenBase64
}
