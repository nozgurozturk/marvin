package entity

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUp struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func AuthToUserDTO(user *SignUp) *UserDTO {
	return &UserDTO{
		Name:        user.Name,
		Email:       user.Email,
		Password:    user.Password,
	}
}