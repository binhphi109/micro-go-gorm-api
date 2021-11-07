package model

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Company struct {
	ID        string     `gorm:"primary_key" json:"id"`
	Name      string     `json:"name"`
	Deleted   bool       `json:"deleted"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

func (c *Company) Prepare() {
	c.ID = uuid.NewString()
	c.Name = html.EscapeString(strings.TrimSpace(c.Name))
	c.Deleted = false
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}

func (c *Company) Validate() error {
	if c.Name == "" {
		return errors.New("required Name")
	}

	return nil
}
