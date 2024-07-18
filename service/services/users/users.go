package users

import (
	"time"

	"github.com/alperklc/the-zula/service/infrastructure/auth"
)

type UsersService interface {
	GetUser(id string) (User, error)
	UpdateUser(id, firstName, lastname, displayName string, language, theme *string) error
}

type datasources struct {
	Auth auth.AuthClient
}

func NewService(a *auth.AuthClient) UsersService {
	return &datasources{
		Auth: *a,
	}
}

func (d *datasources) GetUser(id string) (User, error) {
	user, err := d.Auth.GetUser(id)

	theme, errGetTheme := d.Auth.GetUserMetadata(id, "theme")
	if errGetTheme != nil {
		return User{}, errGetTheme
	}
	language, errGetLanguage := d.Auth.GetUserMetadata(id, "language")
	if errGetLanguage != nil {
		return User{}, errGetLanguage
	}

	creationDate, creationDateParseErr := time.Parse(time.RFC3339Nano, user.Details.CreationDate)
	if creationDateParseErr != nil {
		return User{}, creationDateParseErr
	}
	changeDate, changeDateParseErr := time.Parse(time.RFC3339Nano, user.Details.ChangeDate)
	if changeDateParseErr != nil {
		return User{}, changeDateParseErr
	}

	return User{
		ID:          user.ID,
		CreatedAt:   creationDate,
		UpdatedAt:   changeDate,
		FirstName:   user.Human.Profile.FirstName,
		LastName:    user.Human.Profile.LastName,
		DisplayName: user.Human.Profile.DisplayName,
		Email:       user.Human.Email.Email,
		Language:    language,
		Theme:       theme,
	}, err
}

func (d *datasources) UpdateUser(id, firstName, lastname, displayName string, language, theme *string) error {
	if theme != nil {
		errGetTheme := d.Auth.SetUserMetadata(id, "theme", *theme)
		if errGetTheme != nil {
			return errGetTheme
		}
	}
	if language != nil {
		errGetLanguage := d.Auth.SetUserMetadata(id, "language", *language)
		if errGetLanguage != nil {
			return errGetLanguage
		}
	}
	errSetUser := d.Auth.SetUserProfile(id, auth.UserInput{FirstName: firstName, LastName: lastname, DisplayName: displayName})
	if errSetUser != nil {
		return errSetUser
	}

	return nil
}
