package entity

type Token struct {
	Token      string
	Uuid       string
	Expires    int64
	UserID     string
	Email      string
	Authorized bool
}

// Public Token that is consumed by confirmation mail and update notify preferences
type SubToken struct {
	Token   string
	Uuid    string
	Expires int64
	SubID   string
	RepoID  string
	Email   string
}

// User by session db that is authorize user's session
type AccessDetail struct {
	Uuid   string
	UserID string
}

type TokenDetails struct {
	AccessToken  *Token
	RefreshToken *Token
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func TokenDetailsToResponse(details *TokenDetails) *TokenResponse {
	return &TokenResponse{
		AccessToken:  details.AccessToken.Token,
		RefreshToken: details.RefreshToken.Token,
	}
}
