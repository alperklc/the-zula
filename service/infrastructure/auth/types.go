package auth

import (
	"net/http"
	"time"
)

type AuthClient struct {
	authURL      string
	clientId     string
	clientSecret string
	token        string
	expiresAt    time.Time
	client       *http.Client
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

type UserResponse struct {
	User User `json:"user"`
}

type User struct {
	ID                 string   `json:"id"`
	Details            Details  `json:"details"`
	State              string   `json:"state"`
	UserName           string   `json:"userName"`
	LoginNames         []string `json:"loginNames"`
	PreferredLoginName string   `json:"preferredLoginName"`
	Human              Human    `json:"human"`
}

type Details struct {
	Sequence      string `json:"sequence"`
	CreationDate  string `json:"creationDate"`
	ChangeDate    string `json:"changeDate"`
	ResourceOwner string `json:"resourceOwner"`
}

type Human struct {
	Profile Profile `json:"profile"`
	Email   Email   `json:"email"`
}

type Profile struct {
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	DisplayName       string `json:"displayName"`
	PreferredLanguage string `json:"preferredLanguage"`
}

type Email struct {
	Email           string `json:"email"`
	IsEmailVerified bool   `json:"isEmailVerified"`
}

type MetadataResponse struct {
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MetadataInput struct {
	Value string `json:"value"`
}

type UserInput struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	DisplayName string `json:"displayName"`
}
