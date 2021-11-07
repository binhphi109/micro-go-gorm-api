package api

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwt"

	"github.com/sample/sample-server/core"
)

type Handler struct {
	HandleFunc  func(http.ResponseWriter, *http.Request)
	HandlerName string
	Config      *core.Config
}

func GetHandlerName(h func(http.ResponseWriter, *http.Request)) string {
	handlerName := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	pos := strings.LastIndex(handlerName, ".")
	if pos != -1 && len(handlerName) > pos {
		handlerName = handlerName[pos+1:]
	}
	return handlerName
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func (api *API) APIHandler(h func(http.ResponseWriter, *http.Request)) http.Handler {
	handler := &Handler{
		HandleFunc:  h,
		HandlerName: GetHandlerName(h),
		Config:      api.Config,
	}

	return handler
}

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

// AdminOnly middleware restricts access to just administrators.
func (api *API) AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin, ok := r.Context().Value("acl.admin").(bool)
		if !ok || !isAdmin {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
