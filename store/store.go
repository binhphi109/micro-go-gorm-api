package store

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type IStore interface {
	User() UserStore
	Company() CompanyStore
}

type Store struct {
	company *CompanyStore
	user    *UserStore
	Db      *gorm.DB
	Logger  *logrus.Logger
}

func New(db *gorm.DB, logger logrus.Logger) *Store {
	return &Store{
		Db:     db,
		Logger: &logger,
	}
}
