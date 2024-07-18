package auth

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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

func NewAuthClient(authURL, clientId, clientSecret string) *AuthClient {
	return &AuthClient{
		authURL:      authURL,
		clientId:     clientId,
		clientSecret: clientSecret,
		client:       &http.Client{},
	}
}

func (c *AuthClient) Authenticate() error {
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("scope", "openid profile email urn:zitadel:iam:org:project:id:zitadel:aud")
	form.Add("client_id", c.clientId)
	form.Add("client_secret", c.clientSecret)

	req, newRequestErr := http.NewRequest("POST", fmt.Sprintf("%soauth/v2/token", c.authURL), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if newRequestErr != nil {
		return newRequestErr
	}

	resp, httpPostErr := c.client.Do(req)
	if httpPostErr != nil {
		return httpPostErr
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("auth_failed")
	}

	var response AuthResponse
	errDecode := json.NewDecoder(resp.Body).Decode(&response)
	if errDecode != nil {
		return errDecode
	}

	defer resp.Body.Close()

	c.token = response.AccessToken
	c.expiresAt = time.Now().Add((time.Second - 1) * time.Duration(response.ExpiresIn))

	return nil
}

func (c *AuthClient) TokenValid() bool {
	if c.token == "" {
		return false
	}

	return c.expiresAt.After(time.Now())
}

func (c *AuthClient) GetUser(userId string) (*User, error) {
	if !c.TokenValid() {
		err := c.Authenticate()
		if err != nil {
			return nil, err
		}
	}

	req, newRequestErr := http.NewRequest("GET", fmt.Sprintf("%s/management/v1/users/%s", c.authURL, userId), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if newRequestErr != nil {
		return nil, newRequestErr
	}

	resp, httpGetErr := c.client.Do(req)
	if httpGetErr != nil || resp.StatusCode != http.StatusOK {
		return nil, httpGetErr
	}

	var response UserResponse
	errDecode := json.NewDecoder(resp.Body).Decode(&response)
	if errDecode != nil {
		return nil, errDecode
	}

	defer resp.Body.Close()

	return &response.User, nil
}

func (c *AuthClient) GetUserMetadata(userId, key string) (string, error) {
	if !c.TokenValid() {
		err := c.Authenticate()
		if err != nil {
			return "", err
		}
	}

	req, newRequestErr := http.NewRequest("GET", fmt.Sprintf("%s/management/v1/users/%s/metadata/%s", c.authURL, userId, key), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if newRequestErr != nil {
		return "", newRequestErr
	}

	resp, httpGetErr := c.client.Do(req)
	if httpGetErr != nil || resp.StatusCode != http.StatusOK {
		return "", httpGetErr
	}

	var response MetadataResponse
	errDecode := json.NewDecoder(resp.Body).Decode(&response)
	if newRequestErr != nil {
		return "", errDecode
	}

	defer resp.Body.Close()

	if response.Metadata.Value == "" {
		return "", nil
	}

	decodedMetadataValue, errDecode := base64.StdEncoding.DecodeString(response.Metadata.Value)
	if errDecode != nil {
		return "", errDecode
	}

	return string(decodedMetadataValue), nil
}

func (c *AuthClient) SetUserProfile(userId string, userinput UserInput) error {
	if !c.TokenValid() {
		err := c.Authenticate()
		if err != nil {
			return err
		}
	}

	requestBody, marshalErr := json.Marshal(userinput)
	if marshalErr != nil {
		return marshalErr
	}

	req, newRequestErr := http.NewRequest("PUT", fmt.Sprintf("%s/management/v1/users/%s/profile", c.authURL, userId), bytes.NewBuffer(requestBody))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if newRequestErr != nil {
		return newRequestErr
	}

	resp, httpGetErr := c.client.Do(req)
	if httpGetErr != nil || resp.StatusCode != http.StatusOK {
		return httpGetErr
	}
	defer resp.Body.Close()

	return nil
}

func (c *AuthClient) SetUserMetadata(userId, key, value string) error {
	if !c.TokenValid() {
		err := c.Authenticate()
		if err != nil {
			return err
		}
	}

	decodedMetadataValue := base64.StdEncoding.EncodeToString([]byte(value))

	requestBody, marshalErr := json.Marshal(MetadataInput{Value: decodedMetadataValue})
	if marshalErr != nil {
		return marshalErr
	}

	req, newRequestErr := http.NewRequest("POST", fmt.Sprintf("%s/management/v1/users/%s/metadata/%s", c.authURL, userId, key), bytes.NewBuffer(requestBody))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if newRequestErr != nil {
		return newRequestErr
	}

	resp, httpGetErr := c.client.Do(req)
	if httpGetErr != nil || resp.StatusCode != http.StatusOK {
		return httpGetErr
	}

	defer resp.Body.Close()

	return nil
}
