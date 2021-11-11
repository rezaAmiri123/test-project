package user

import (
	"regexp"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmptyPassword = errors.New("empty password")
)

// User is user model
type User struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
}

func NewUser(username, password, email, bio, image string) *User {
	return &User{
		UUID:     uuid.New().String(),
		Username: username,
		Password: password,
		Email:    email,
		Bio:      bio,
		Image:    image,
	}
}

// Validate validates fields of user model
func (u User) Validate() error {
	err := validation.ValidateStruct(&u,
		validation.Field(
			&u.Username,
			validation.Required,
			validation.Match(regexp.MustCompile("[a-zA-Z0-9]+")),
		),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required),
		validation.Field(&u.UUID, validation.Required, is.UUID),
	)
	if err != nil {
		return errors.Wrap(err, "cannot validate user")
	}
	return nil
}

// HashPassword makes password field crypted
func (u *User) SetUUID() error {
	if u.UUID ==""{
		u.UUID = uuid.New().String()
	}
	return nil
}

// HashPassword makes password field crypted
func (u *User) HashPassword() error {
	if len(u.Password) == 0 {
		return ErrEmptyPassword
	}
	h, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "can't generate password")
	}
	u.Password = string(h)
	return nil
}

// CheckPassword checks user password correct
func (u *User) CheckPassword(plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain))
	return err == nil
}

// HidePassword hide user password
func (u *User) HidePassword() {
	u.Password = ""
}
