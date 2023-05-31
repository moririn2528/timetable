package usecase

import (
	"crypto/sha256"
	"os"
	"time"

	"timetable/errors"

	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

func hashPassword(pass string) string {
	pass += os.Getenv("PASSWORD_SALT")
	sh := sha256.Sum256([]byte(pass))
	return string(sh[:])
}

func createToken(user User) (string, error) {
	token := jwt.New(jwt.GetSigningMethod("HS256"))

	token.Claims = jwt.MapClaims{
		"iss":  "moririn2528/timetable",
		"sub":  "access",
		"aud":  "user",
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
		"nbf":  time.Now().Unix(),
		"lat":  time.Now().Unix(),
		"jti":  user.Id,
		"name": user.Name,
	}

	var secretKey = os.Getenv("JWT_SECRET_KEY")
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func Login(user User, password string) (string, error) {
	password = hashPassword(password)
	err := Db_any.FindUser(user, password)
	if err != nil {
		return "", errors.ErrorWrap(err)
	}
	logger.Infof("login success: id = %v, name = %v", user.Id, user.Name)
	token, err := createToken(user)
	if err != nil {
		return "", errors.ErrorWrap(err)
	}
	return token, nil
}

func VerifyToken(token_string string) (*jwt.Token, error) {
	token, err := jwt.Parse(token_string, func(*jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	logger.Debug(token.Claims)
	return token, nil
}

func Register(user User, password string) error {
	password = hashPassword(password)
	err := Db_any.InsertUser(user, password)
	if err != nil {
		return errors.ErrorWrap(err)
	}
	logger.Infof("register success: id = %v, name = %v", user.Id, user.Name)
	return nil
}
