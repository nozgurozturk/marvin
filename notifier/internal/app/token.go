package app

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/nozgurozturk/marvin/notifier/entity"
	"github.com/nozgurozturk/marvin/notifier/internal/config"
	"time"
)


func CreateToken(s *entity.SubscriberDTO) (*entity.SubToken, error) {
	var err error
	cnf := config.Get().HTTP
	sub := &entity.SubToken{}
	sub.Expires = time.Now().Add(time.Hour * time.Duration(cnf.SubExpire)).Unix()
	sub.Uuid = uuid.New().String()
	sub.RepoID = s.RepoID
	sub.Email = s.Email
	sub.SubID = *s.ID

	subClaims := jwt.MapClaims{
		"exp":    sub.Expires,
		"uuid":   sub.Uuid,
		"subID":  sub.SubID,
		"repoID": sub.RepoID,
		"email":  sub.Email,
	}

	subToken := jwt.NewWithClaims(jwt.SigningMethodHS256, subClaims)

	sub.Token, err = subToken.SignedString([]byte(cnf.SubSecret))
	if err != nil {
		return nil, err
	}

	return sub, nil
}
