package app

import (
	"fmt"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/sirupsen/logrus"

	"github.com/sample/sample-server/core"
	"github.com/sample/sample-server/dto"
	"github.com/sample/sample-server/model"
	"github.com/sample/sample-server/store"
	"github.com/sample/sample-server/utils"
)

type UserApp struct {
	Store  *store.Store
	Logger *logrus.Logger
	Config *core.Config
}

func (a *App) User() *UserApp {
	if a.user == nil {
		a.user = &UserApp{
			Store:  a.Store,
			Logger: a.Logger,
			Config: a.Config,
		}
	}
	return a.user
}

func (ua *UserApp) Login(authdetails dto.Authentication) (*dto.Token, error) {
	authUser, err := ua.Store.User().GetByEmail(authdetails.Email)
	if err != nil {
		return &dto.Token{}, err
	}

	if authUser.Email == "" {
		return &dto.Token{}, fmt.Errorf("error Username or Password is incorrect")
	}

	check := utils.CheckPasswordHash(authdetails.Password, authUser.Password)

	if !check {
		return &dto.Token{}, fmt.Errorf("error Username or Password is incorrect")
	}

	jwtAuth := jwtauth.New(ua.Config.TOKEN_SIGN_METHOD, []byte(ua.Config.TOKEN_SIGN_KEY), []byte(ua.Config.TOKEN_VERIFY_KEY))

	claims := make(map[string]interface{})
	claims["authorized"] = true
	claims["email"] = authUser.Email
	// claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	_, tokenString, err := jwtAuth.Encode(claims)
	if err != nil {
		return &dto.Token{}, fmt.Errorf("error Failed to generate token")
	}

	token := &dto.Token{
		Email:       authdetails.Email,
		TokenString: tokenString,
	}
	return token, nil
}

func (ua *UserApp) GetAllUsers() ([]model.User, error) {
	users, err := ua.Store.User().GetAll()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		user.Password = ""
	}

	return users, nil
}

func (ua *UserApp) GetUser(userId string) (*model.User, error) {
	user, err := ua.Store.User().Get(userId)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return &user, nil
}

func (ua *UserApp) CreateUser(user model.User) (*model.User, error) {
	oldUser, err := ua.Store.User().GetByEmail(user.Email)
	if err != nil {
		return nil, err
	}

	//checks if email is already register or not
	if oldUser.Email != "" {
		return nil, fmt.Errorf("email already in use")
	}

	hash, err := utils.GenerateHashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("error in password hash")
	}

	user.Password = string(hash)

	//insert user details in database
	newUser, err := ua.Store.User().Create(user)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (ua *UserApp) UpdateUser(user model.User) (model.User, error) {
	oldUser, err := ua.Store.User().Get(user.ID)
	if err != nil {
		return user, err
	}

	oldUser.Name = user.Name
	oldUser.Username = user.Username
	oldUser.Email = user.Email
	if oldUser.Password != user.Password {
		hash, err := utils.GenerateHashPassword(user.Password)
		if err != nil {
			return user, fmt.Errorf("error in password hash")
		}

		oldUser.Password = hash
	}
	oldUser.CompanyId = user.CompanyId
	oldUser.UpdatedAt = time.Now()

	//update user details in database
	oldUser, err = ua.Store.User().Update(oldUser)
	if err != nil {
		return user, err
	}

	return oldUser, nil
}

func (ua *UserApp) DeleteUser(userId string) error {
	user, err := ua.Store.User().Get(userId)
	if err != nil {
		return err
	}

	timeNow := time.Now()
	user.Deleted = true
	user.DeletedAt = &timeNow

	//soft delete user in database
	_, err = ua.Store.User().Update(user)
	if err != nil {
		return err
	}

	return nil
}
