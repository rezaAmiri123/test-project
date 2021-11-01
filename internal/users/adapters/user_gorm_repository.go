package adapters

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/rezaAmiri123/test-project/internal/users/domain/user"
)

// User is user model
type GORMUserModel struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
}

func (m *GORMUserModel) protoDomainUser() *user.User {
	return &user.User{
		Username: m.Username,
		Password: m.Password,
		Email:    m.Email,
		Bio:      m.Bio,
		Image:    m.Image,
	}
}

func (m *GORMUserModel) protoGORMUser(user *user.User) {
	m.Username = user.Username
	m.Password = user.Password
	m.Email = user.Email
	m.Bio = user.Bio
	m.Image = user.Image
}

type GORMUserRepository struct {
	db *gorm.DB
}

func NewGORMUserRepository(db *gorm.DB) *GORMUserRepository {
	if db == nil {
		panic("missing gorm db")
	}
	return &GORMUserRepository{db: db}
}

func (r *GORMUserRepository) Create(ctx context.Context, user *user.User) error {
	var gormUser GORMUserModel
	gormUser.protoGORMUser(user)
	err := r.db.Create(gormUser).Error
	if err != nil {
		return errors.Wrap(err, "cannot create user")
	}
	return nil
}

func (r *GORMUserRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	var gormUser GORMUserModel
	if err := r.db.Where(GORMUserModel{Username: username}).First(&gormUser).Error; err != nil {
		return nil, errors.Wrap(err, "cannot find user")
	}
	return gormUser.protoDomainUser(), nil
}
