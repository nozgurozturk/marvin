package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
	"time"
)

type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Email       string             `json:"email" bson:"email"`
	IsConfirmed bool               `json:"isConfirmed" bson:"isConfirmed"`
	Password    string             `json:"password" bson:"password"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
}

type UserDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	IsConfirmed bool   `json:"isConfirmed"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type UserResponse struct {
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
}

func ToUser(userDTO *UserDTO) *User {

	id, _ := primitive.ObjectIDFromHex(userDTO.ID)

	return &User{
		ID:          id,
		Name:        userDTO.Name,
		Email:       userDTO.Email,
		IsConfirmed: userDTO.IsConfirmed,
		Password:    userDTO.Password,
	}
}

func ToUserDTO(user *User) *UserDTO {
	return &UserDTO{
		ID:          user.ID.Hex(),
		Name:        user.Name,
		Email:       user.Email,
		IsConfirmed: user.IsConfirmed,
		Password:    user.Password,
	}
}

func ToUserResponse(user *UserDTO) *UserResponse {
	return &UserResponse{
		Name:  user.Name,
		Email: user.Email,
	}
}

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword string, userPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userPassword))
}

// Validates user's email
func ValidateEmail(email string) string {
	email = strings.TrimSpace(email)

	if email == "" {
		return "Email is required"
	}
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegex.MatchString(email) {
		return "Email is not valid"
	}
	return ""
}

// Validates user's  [name, password, email]
func ValidateUser(user *UserDTO) string {

	user.Name = strings.TrimSpace(user.Name)
	if user.Name == "" {
		return "Name is required"
	}

	if len(user.Name) < 3 {
		return "Name must be minimum 3 character long"
	}

	user.Password = strings.TrimSpace(user.Password)
	if user.Password == "" {
		return "Password is required"
	}
	if len(user.Password) < 8 {
		return "Password must be minimum 8 character long"
	}

	user.Email = strings.TrimSpace(user.Email)
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if user.Email == "" {
		return "Email is required"
	}
	if !emailRegex.MatchString(user.Email) {
		return "Email is not valid"
	}

	return ""
}
