package store

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"github.com/sample/sample-server/model"
)

type UserStore struct {
	Db     *gorm.DB
	Logger *logrus.Logger
}

func (s *Store) User() *UserStore {
	if s.user == nil {
		s.user = &UserStore{
			Db:     s.Db,
			Logger: s.Logger,
		}
	}
	return s.user
}

func (us *UserStore) GetAll() ([]model.User, error) {
	var users []model.User
	result := us.Db.Table("Users").Find(&users)

	if result.Error != nil {
		return users, result.Error
	}

	return users, nil
}

func (us *UserStore) Get(id string) (model.User, error) {
	var user model.User
	us.Db.Table("Users").Where("id = ?", id).First(&user)

	if user.ID == "" {
		return user, fmt.Errorf("no User found")
	}

	return user, nil
}

func (us *UserStore) GetByEmail(email string) (model.User, error) {
	var user model.User
	us.Db.Table("Users").Where("email = ?", email).First(&user)

	if user.ID == "" {
		return user, fmt.Errorf("no User found")
	}

	return user, nil
}

func (us *UserStore) Create(user model.User) (model.User, error) {
	var db = us.Db.Table("Users").Create(&user)

	if db.Error != nil {
		us.Logger.WithError(db.Error).Fatal("Error creating user", user.Username)
		return model.User{}, db.Error
	}

	return user, nil
}

func (us *UserStore) Update(user model.User) (model.User, error) {
	var db = us.Db.Table("Users").Where("id = ?", user.ID).Updates(model.User{
		Name:      user.Name,
		Email:     user.Email,
		Username:  user.Username,
		Password:  user.Password,
		CompanyId: user.CompanyId,
		Deleted:   user.Deleted,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	})

	if db.Error != nil {
		us.Logger.WithError(db.Error).Fatal("Error updating user", user.Name)
		return model.User{}, db.Error
	}

	return user, nil
}

func (us *UserStore) Delete(userId string) error {
	var db = us.Db.Table("Users").Where("id = ?", userId).Delete(&model.User{})

	if db.Error != nil {
		us.Logger.WithError(db.Error).Fatal("Error deleting user", userId)
		return db.Error
	}

	return nil
}
