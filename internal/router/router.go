package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/1gkx/salary/internal/session"
	"github.com/1gkx/salary/internal/store"
	templates "github.com/1gkx/salary/internal/template"
)

var r *mux.Router

func NewRouter() *mux.Router {

	r = mux.NewRouter().StrictSlash(true)

	setPublicFolder()
	setUserRouters()
	setSignInRouters()
	setSettingRouters()
	SetClientHendler()

	r.Handle("/", authRequireHandlerWrap(index)).Methods("GET")
	r.Handle("/user", authRequireHandlerWrap(profile)).Methods("GET")

	r.NotFoundHandler = http.HandlerFunc(errorhendler)

	return r
}

func setPublicFolder() {
	publicFolder := http.FileServer(http.Dir("./public"))
	r.PathPrefix("/js/").Handler(publicFolder)
	r.PathPrefix("/css/").Handler(publicFolder)
	r.PathPrefix("/img/").Handler(publicFolder)
}

func index(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var page = ""
	params := r.URL.Query()
	if len(params) > 0 {
		page = params["page"][0]
	} else {
		page = "1"
	}
	cls := store.GetClients(page)
	responce(w, r, "home", cls)
}

func errorhendler(w http.ResponseWriter, r *http.Request) {
	responce(w, r, "status/40x", nil)
}

func profile(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Set("Cache-Control", "No-Cache")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.Templates.ExecuteTemplate(w, "user",
		map[string]interface{}{
			"isNew":     false,
			"isProfile": true,
			"user":      session.GetUser(r),
			"data":      session.GetUser(r),
			"redirect":  "/user",
		},
	)
	return
}

func responceJson(code int, w http.ResponseWriter, data map[string]interface{}) {
	w.WriteHeader(code)
	w.Header().Set("Cache-Control", "No-Cache")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": code,
		"data": data,
	})
}

func responce(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	w.Header().Set("Cache-Control", "No-Cache")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.Templates.ExecuteTemplate(w, tmpl,
		map[string]interface{}{
			"user": session.GetUser(r),
			"data": data,
		},
	)
}
