package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"

	"github.com/sample/sample-server/dto"
	"github.com/sample/sample-server/model"
)

func (api *API) InitUser() {
	jwtAuth := jwtauth.New(api.Config.TOKEN_SIGN_METHOD, []byte(api.Config.TOKEN_SIGN_KEY), []byte(api.Config.TOKEN_VERIFY_KEY))

	api.Router.Route("/users", func(route chi.Router) {
		route.Method("POST", "/", api.APIHandler(api.CreateUser))
		route.Method("POST", "/login", api.APIHandler(api.Login))

		route.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(jwtAuth))
			r.Use(Authenticator)

			route.Method("GET", "/", api.APIHandler(api.GetAllUsers))
			route.Method("GET", "/{userId}", api.APIHandler(api.GetUser))
			route.Method("PUT", "/", api.APIHandler(api.UpdateUser))
			route.Method("DELETE", "/{userId}", api.APIHandler(api.DeleteUser))
		})
	})
}

func (api *API) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	user.Prepare()
	if err = user.Validate(); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	createdUser, err := api.App.User().CreateUser(user)
	if err != nil {
		api.Logger.WithError(err).Fatal("Error in creating user")
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.JSON(w, r, createdUser)
}

func (api *API) Login(w http.ResponseWriter, r *http.Request) {
	var authdetails dto.Authentication
	err := json.NewDecoder(r.Body).Decode(&authdetails)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	token, err := api.App.User().Login(authdetails)
	if err != nil {
		api.Logger.WithError(err)
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.JSON(w, r, token)
}

func (api *API) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := api.App.User().GetAllUsers()

	if err != nil {
		api.Logger.WithError(err)
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.JSON(w, r, users)
}

func (api *API) GetUser(w http.ResponseWriter, r *http.Request) {
	var userId = chi.URLParam(r, "userId")
	if userId == "" {
		render.Render(w, r, ErrInvalidRequest(fmt.Errorf("Error missing UserId")))
		return
	}

	user, err := api.App.User().GetUser(userId)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.JSON(w, r, user)
}

func (api *API) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err = user.Validate(); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	updatedUser, err := api.App.User().UpdateUser(user)
	if err != nil {
		api.Logger.WithError(err).Fatal("Error in updating user")
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.JSON(w, r, updatedUser)
}

func (api *API) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var userId = chi.URLParam(r, "userId")

	err := api.App.User().DeleteUser(userId)
	if err != nil {
		api.Logger.WithError(err).Fatal("Error in deleting user")
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	ReturnStatusOK(w)
}
