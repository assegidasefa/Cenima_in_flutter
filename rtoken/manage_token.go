package rtoken

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/joocosta/bloctrial/model"
	"log"
)

type Service struct{
	privateKey []byte
}

func NewToken(privateKey []byte) Service{
	return Service{
		privateKey: privateKey,
	}
}

// JwtClaim adds email as a claim to the token
type CustomJwtClaim struct {
	User model.User
	jwt.StandardClaims
}

func (t *Service) GenerateToken(claims jwt.Claims) (string, error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(t.privateKey)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return tokenString, nil
}
func (t *Service) ValidateToken(signedToken string) (bool, error){
	token, err := jwt.ParseWithClaims(signedToken, &CustomJwtClaim{}, func(token *jwt.Token) (i interface{}, e error) {
		return t.privateKey, nil
	})
	if err != nil{
		return false, err
	}

	if _, ok := token.Claims.(*CustomJwtClaim); !ok || token.Valid{
		return false, err
	}

	return true, nil
}

func (t *Service) GetClaims(signedToken string) (*CustomJwtClaim, error) {
	token, err := jwt.ParseWithClaims(signedToken, &CustomJwtClaim{}, func(token *jwt.Token) (i interface{}, e error) {
		return t.privateKey, nil
	})

	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*CustomJwtClaim)
	if !ok || !token.Valid {
		return claims, err
	}
	return claims, err
}
