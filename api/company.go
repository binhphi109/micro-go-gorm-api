package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"

	"github.com/sample/sample-server/model"
)

func (api *API) InitCompany() {
	jwtAuth := jwtauth.New(api.Config.TOKEN_SIGN_METHOD, []byte(api.Config.TOKEN_SIGN_KEY), []byte(api.Config.TOKEN_VERIFY_KEY))

	api.Router.Route("/companies", func(route chi.Router) {
		route.Use(jwtauth.Verifier(jwtAuth))
		route.Use(Authenticator)

		route.Method("GET", "/", api.APIHandler(api.GetAllCompanies))
		route.Method("GET", "/{companyId}", api.APIHandler(api.GetCompany))
		route.Method("POST", "/", api.APIHandler(api.CreateCompany))
		route.Method("PUT", "/", api.APIHandler(api.UpdateCompany))
		route.Method("DELETE", "/{companyId}", api.APIHandler(api.DeleteCompany))
	})
}

func (api *API) GetAllCompanies(w http.ResponseWriter, r *http.Request) {
	companies, err := api.App.Company().GetAllCompanies()

	if err != nil {
		api.Logger.WithError(err)
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.JSON(w, r, companies)
}

func (api *API) GetCompany(w http.ResponseWriter, r *http.Request) {
	var companyId = chi.URLParam(r, "companyId")
	if companyId == "" {
		render.Render(w, r, ErrInvalidRequest(fmt.Errorf("Error missing CompanyId")))
		return
	}

	company, err := api.App.Company().GetCompany(companyId)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.JSON(w, r, company)
}

func (api *API) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var company model.Company
	err := json.NewDecoder(r.Body).Decode(&company)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	company.Prepare()
	if err = company.Validate(); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	createdCompany, err := api.App.Company().CreateCompany(company)
	if err != nil {
		api.Logger.WithError(err).Fatal("Error in creating company")
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.JSON(w, r, createdCompany)
}

func (api *API) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	var company model.Company
	err := json.NewDecoder(r.Body).Decode(&company)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err = company.Validate(); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	updatedCompany, err := api.App.Company().UpdateCompany(company)
	if err != nil {
		api.Logger.WithError(err).Fatal("Error in updating company")
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.JSON(w, r, updatedCompany)
}

func (api *API) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	var companyId = chi.URLParam(r, "companyId")

	err := api.App.Company().DeleteCompany(companyId)
	if err != nil {
		api.Logger.WithError(err).Fatal("Error in deleting company")
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	ReturnStatusOK(w)
}
