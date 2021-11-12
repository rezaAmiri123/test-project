package article

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// Article model
type Article struct {
	UUID        string `json:"uuid"`
	UserUUID    string `json:"user_uuid"`
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Body        string `json:"body"`
}

func (a Article) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Title, validation.Required),
		validation.Field(&a.Body, validation.Required),
	)
}
