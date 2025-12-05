package atylabjwt

import (
	"context"
	"fmt"
	"time"

	firebase "firebase.google.com/go"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/option"
)

type JwtSvcInterface interface {
	CreateJwt(config *JwtConfig) (string, error)
	Validate(jwtSecret string, jwtToken string) error
	GetUUID() string
	GetEmail() string
}

type JwtSvcStruct struct {
	uuid  string
	email string
}

func NewJwtSvc() *JwtSvcStruct {
	return &JwtSvcStruct{}
}

type JwtConfig struct {
	Key   []byte
	Uuid  string
	Email string
	Exp   time.Time
}

func (s *JwtSvcStruct) CreateJwt(config *JwtConfig) (string, error) {
	claims := jwt.MapClaims{
		"sub":   fmt.Sprintf("user%s", config.Uuid),
		"email": config.Email,
		"exp":   config.Exp.Add(time.Hour * 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.Key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *JwtSvcStruct) Validate(
	jwtSecret string,
	jwtToken string,
) error {
	if err := s.originalValidateToken(jwtSecret, jwtToken); err == nil {
		fmt.Print("jwt use OriginalValidateToken\n")
		return nil
	}
	if err := s.firebaseValidateToken(jwtToken); err == nil {
		fmt.Print("jwt use FirebaseValidateToken\n")
		return nil
	}
	return fmt.Errorf("invalid jwt token")
}

func (j *JwtSvcStruct) originalValidateToken(
	jwtSecret string,
	jwtToken string,
) error {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		fmt.Println("jwtSecret:", string([]byte(jwtSecret)))
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return fmt.Errorf("invalid jwt token: %v", err)
	}

	j.uuid = claims["sub"].(string)
	j.email = claims["email"].(string)

	fmt.Print("jwt use OriginalValidateToken\n")

	return nil
}

func (j *JwtSvcStruct) firebaseValidateToken(
	jwtToken string,
) error {
	opt := option.WithCredentialsFile("secret/firebase_credentials.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return err
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		return err
	}

	token, err := client.VerifyIDToken(context.Background(), jwtToken)
	if err != nil {
		return err
	}
	j.uuid = token.UID
	j.email = token.Claims["email"].(string)

	fmt.Print("jwt use FirebaseValidateToken\n")

	return nil
}

func (j *JwtSvcStruct) GetUUID() string {
	return j.uuid
}

func (j *JwtSvcStruct) GetEmail() string {
	return j.email
}
