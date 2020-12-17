package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/1gkx/salary/internal/session"
	"github.com/1gkx/salary/internal/store"
	templates "github.com/1gkx/salary/internal/template"
	"github.com/gorilla/mux"
)

func setUserRouters() {
	// Users
	r.Handle("/admin/users", authRequireHandlerWrap(userList)).Methods("GET")
	r.Handle("/admin/users/{id:[0-9]+}", authRequireHandlerWrap(userprofile)).Methods("GET")
	r.Handle("/admin/users", authRequireHandlerWrap(userAdd)).Methods("POST")
	r.Handle("/admin/users", authRequireHandlerWrap(userUpdate)).Methods("PUT")
	r.Handle("/admin/users", authRequireHandlerWrap(userRemove)).Methods("DELETE")
}

func userprofile(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	u := store.FindByID(uint(id))

	isNew := false
	if uint(id) == 0 {
		isNew = true
	}

	w.Header().Set("Cache-Control", "No-Cache")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.Templates.ExecuteTemplate(w, "user",
		map[string]interface{}{
			"isNew":     isNew,
			"isProfile": false,
			"user":      session.GetUser(r),
			"data":      u,
			"redirect":  "/admin/users",
		},
	)
	return

}

func userList(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var page = ""
	params := r.URL.Query()
	if len(params) > 0 {
		page = params["page"][0]
	} else {
		page = "1"
	}
	users := store.FindUserLimit(page)

	j := map[string]interface{}{
		"user": session.GetUser(r),
		"data": users,
	}
	templates.Templates.ExecuteTemplate(w, "list", j)
}

func userNew(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	j := map[string]interface{}{
		"user": session.GetUser(r),
		"data": store.FindUser(),
	}
	templates.Templates.ExecuteTemplate(w, "new", j)
}

func userAdd(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	u := new(store.User)
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		RespAPI(http.StatusInternalServerError, w, err)
		return
	}

	fmt.Printf("User: %+v\n", u)

	if err := store.AddUser(u); err != nil {
		RespAPI(http.StatusInternalServerError, w, err)
		return
	}
	w.WriteHeader(201)
	return
}

func userUpdate(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	tmpUser := new(store.User)
	if err := json.NewDecoder(r.Body).Decode(&tmpUser); err != nil {
		RespAPI(http.StatusInternalServerError, w, err)
		return
	}

	u := store.FindByID(tmpUser.ID)
	fmt.Println(u)
	if u == nil {
		RespAPI(http.StatusInternalServerError, w, "user not found")
		return
	}

	if err := store.UpdateUser(tmpUser); err != nil {
		RespAPI(http.StatusInternalServerError, w, err)
		return
	}

	RespAPI(http.StatusOK, w, "OK")
	return
}

func userRemove(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	u := new(store.User)
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		RespAPI(http.StatusInternalServerError, w, err)
		return
	}

	if err := store.DeleteUserByID(u.ID); err != nil {
		RespAPI(http.StatusInternalServerError, w, "Fail")
		return
	}

	RespAPI(http.StatusOK, w, "OK")
	return
}

func RespAPI(code int, w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(
		fmt.Sprintf("{\"status\": \"%s\"}", data),
	)
}
