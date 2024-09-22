package usersService

import (
	"fmt"
	"time"

	"github.com/alperklc/the-zula/service/infrastructure/auth"
	"github.com/alperklc/the-zula/service/infrastructure/cache"
	mqpublisher "github.com/alperklc/the-zula/service/infrastructure/messageQueue/publisher"
)

type UsersService interface {
	RefreshUserInCache(id string) error
	GetUser(id string) (User, error)
	UpdateUser(id, clientId, email, firstName, lastname, displayName string, language, theme *string) error
}

type datasources struct {
	Auth        auth.AuthClient
	Cache       cache.Cache[User]
	mqpublisher mqpublisher.MessagePublisher
}

func NewService(a *auth.AuthClient, c *cache.Cache[User], mqp mqpublisher.MessagePublisher) UsersService {
	return &datasources{
		Auth:        *a,
		Cache:       *c,
		mqpublisher: mqp,
	}
}

func (d *datasources) RefreshUserInCache(id string) error {
	d.Cache.Reset(id)
	a, err := d.GetUser(id)
	fmt.Println(a)
	return err
}

func (d *datasources) GetUser(id string) (User, error) {
	obj := d.Cache.Read(id)
	if obj != nil {
		return *obj, nil
	}
	user, errGetUser := d.Auth.GetUser(id)
	if errGetUser != nil {
		return User{}, errGetUser
	}
	theme, errGetTheme := d.Auth.GetUserMetadata(id, "theme")
	if errGetTheme != nil {
		return User{}, errGetTheme
	}
	creationDate, creationDateParseErr := time.Parse(time.RFC3339Nano, user.Details.CreationDate)
	if creationDateParseErr != nil {
		return User{}, creationDateParseErr
	}
	changeDate, changeDateParseErr := time.Parse(time.RFC3339Nano, user.Details.ChangeDate)
	if changeDateParseErr != nil {
		return User{}, changeDateParseErr
	}

	response := User{
		ID:          user.ID,
		CreatedAt:   creationDate,
		UpdatedAt:   changeDate,
		FirstName:   user.Human.Profile.FirstName,
		LastName:    user.Human.Profile.LastName,
		DisplayName: user.Human.Profile.DisplayName,
		Email:       user.Human.Email.Email,
		Language:    user.Human.Profile.PreferredLanguage,
		Theme:       theme,
	}

	d.Cache.Write(id, response)

	return response, nil
}

func (d *datasources) UpdateUser(id, clientId, email, firstName, lastname, displayName string, language, theme *string) error {
	if theme != nil {
		errGetTheme := d.Auth.SetUserMetadata(id, "theme", *theme)
		if errGetTheme != nil {
			return errGetTheme
		}
	}
	errSetUser := d.Auth.SetUserProfile(id, auth.UserInput{FirstName: firstName, LastName: lastname, DisplayName: displayName, PreferredLanguage: *language})
	if errSetUser != nil {
		return errSetUser
	}

	go d.mqpublisher.Publish(mqpublisher.UserUpdated(id, clientId, nil))

	return nil
}
