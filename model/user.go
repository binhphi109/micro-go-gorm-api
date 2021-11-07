package model

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        string     `gorm:"primary_key" json:"id"`
	Name      string     `json:"name"`
	Email     string     `gorm:"unique" json:"email"`
	Username  string     `json:"username"`
	Password  string     `json:"password"`
	CompanyId *string    `json:"companyId"`
	Deleted   bool       `json:"deleted"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

func (u *User) Prepare() {
	u.ID = uuid.NewString()
	u.Name = html.EscapeString(strings.TrimSpace(u.Name))
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.Password = html.EscapeString(strings.TrimSpace(u.Password))
	u.Deleted = false
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("required Name")
	}

	if u.Email == "" {
		return errors.New("required Email")
	}

	if u.Username == "" {
		return errors.New("required Username")
	}

	if u.Password == "" {
		return errors.New("required Password")
	}

	return nil
}
