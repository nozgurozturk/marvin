package app

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nozgurozturk/marvin/pkg/errors"
	"github.com/nozgurozturk/marvin/server/entity"
	"github.com/nozgurozturk/marvin/server/internal/config"
	"strings"
	"time"
)

func verifyToken(t string, tokenType string) (*jwt.Token, *errors.AppError) {
	cnf := config.Get().HTTP
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if tokenType == "Access" {
			return []byte(cnf.AccessSecret), nil
		}

		if tokenType == "Refresh" {
			return []byte(cnf.RefreshSecret), nil
		}

		if tokenType == "Sub" {
			return []byte(cnf.SubSecret), nil
		}
		return nil, fmt.Errorf("anexpected token type: %s", tokenType)
	})
	if err != nil {
		return nil, errors.Unauthorized("Unauthorized user")
	}
	return token, nil
}

func ValidateToken(t string, tType string) (*jwt.Token, *errors.AppError) {
	token, err := verifyToken(t, tType)
	if err != nil {
		if token != nil {
			return token, err
		}
		return token, errors.Unauthorized("Unauthorized user")
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, errors.Unauthorized("Unauthorized user")
	}
	return token, nil
}

func CreateSubToken(s *entity.SubscriberDTO) (*entity.SubToken, *errors.AppError){
	var err error
	cnf := config.Get().HTTP
	sub := &entity.SubToken{}

	sub.Expires = time.Now().Add(time.Hour * time.Duration(cnf.SubExpire)).Unix()
	sub.Uuid = uuid.New().String()
	sub.RepoID = s.RepoID
	sub.Email = s.Email
	sub.SubID = *s.ID

	subClaims := jwt.MapClaims{
		"exp": sub.Expires,
		"uuid": sub.Uuid,
		"subID": sub.SubID,
		"repoID": sub.RepoID,
		"email": sub.Email,
	}

	subToken := jwt.NewWithClaims(jwt.SigningMethodHS256, subClaims)

	sub.Token, err = subToken.SignedString([]byte(cnf.SubSecret))
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}

	return sub, nil
}

func ExtractSubTokenMetaData(t string) (*entity.SubToken, *errors.AppError) {
	token, err := verifyToken(t, "Sub")
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, err
	}

	return &entity.SubToken{
		Uuid:   claims["uuid"].(string),
		RepoID: claims["repoID"].(string),
		SubID: claims["subID"].(string),
		Email: claims["email"].(string),
	}, nil
}

func CreateToken(u *entity.UserDTO) (*entity.TokenDetails, *errors.AppError) {
	var err error

	uniqueID := uuid.New().String()

	cnf := config.Get().HTTP
	rt := &entity.Token{}
	rt.Expires = time.Now().Add(time.Hour * time.Duration(cnf.RefreshExpire)).Unix()
	rt.Uuid = uniqueID
	rt.UserID = u.ID
	rt.Email = u.Email
	rt.Authorized = u.IsConfirmed

	rtClaims := jwt.MapClaims{
		"exp":    rt.Expires,
		"uuid":   rt.Uuid,
		"userID": rt.UserID,
		"email":  rt.Email,
		"authorized": rt.Authorized,
	}

	rtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

	rt.Token, err = rtToken.SignedString([]byte(cnf.RefreshSecret))
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}

	at := &entity.Token{}
	at.Expires = time.Now().Add(time.Minute * time.Duration(cnf.AccessExpire)).Unix()
	at.Uuid = uniqueID
	at.UserID = u.ID
	at.Email = u.Email
	at.Authorized = u.IsConfirmed

	atClaims := jwt.MapClaims{
		"exp":    at.Expires,
		"uuid":   at.Uuid,
		"userID": at.UserID,
		"email":  at.Email,
		"authorized": at.Authorized,
	}

	atToken := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	at.Token, err = atToken.SignedString([]byte(cnf.AccessSecret))
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}

	return &entity.TokenDetails{
		AccessToken:  at,
		RefreshToken: rt,
	}, nil
}

func ExtractToken(c *fiber.Ctx) (string, *errors.AppError) {
	var token string

	authHeader := c.Get("Authorization")
	authHeaderArr := strings.Split(authHeader, " ")

	if len(authHeaderArr) == 2 {
		// [0] bearer, [1] token
		token = authHeaderArr[1]
	}

	if token == "" {
		return "", errors.Unauthorized("Unauthorized user")
	}

	return token, nil
}

func ExtractTokenMetaData(t string) (*entity.AccessDetail, *errors.AppError) {
	token, err := verifyToken(t, "Access")
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, err
	}

	return &entity.AccessDetail{
		Uuid:   claims["uuid"].(string),
		UserID: claims["userID"].(string),
	}, nil
}
