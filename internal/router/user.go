package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/1gkx/salary/internal/session"
	"github.com/1gkx/salary/internal/store"
	templates "github.com/1gkx/salary/internal/template"
	"github.com/gorilla/mux"
)

func setUserRouters() {
	// Users
	r.Handle("/admin/users", authRequireHandlerWrap(userList)).Methods("GET")
	r.Handle("/admin/users/new", authRequireHandlerWrap(userNew)).Methods("GET")
	r.Handle("/admin/users/{id:[0-9]+}", authRequireHandlerWrap(userprofile)).Methods("GET")
	r.Handle("/admin/users", authRequireHandlerWrap(userAdd)).Methods("POST")
	r.Handle("/admin/users", authRequireHandlerWrap(userUpdate)).Methods("PUT")
	r.Handle("/admin/users", authRequireHandlerWrap(userRemove)).Methods("DELETE")
}

func userprofile(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	u := store.FindByID(vars["id"])
	j := map[string]interface{}{
		"user": session.GetUser(r),
		"data": u,
	}
	templates.Templates.ExecuteTemplate(w, "view", j)
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

	// fmt.Printf("Request: %+v\n", r)
	// u := new(store.User)
	// _ = json.NewDecoder(r.Body).Decode(&u)
	// fmt.Printf("User: %+v\n", u)

	// fmt.Printf("User update: %+v\n", u.user)

	// if err := store.AddUser(u.user); err != nil {
	// 	w.WriteHeader(501)
	// 	json.NewEncoder(w).Encode(err.Error())
	// 	return
	// }
	// w.WriteHeader(201)
	return
}

func userUpdate(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	u := new(store.User)
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		w.WriteHeader(501)
		fmt.Printf("{\"status\": \"%s\"}", err.Error())
		json.NewEncoder(w).Encode(
			fmt.Sprintf("{\"status\": \"%s\"}", err.Error()),
		)
		return
	}

	fmt.Printf("{\"status\": \"%+v\"}", u)

	// if u.user.ComparePass(u.newpass) {
	// 	u.user.Password = u.newpass
	// }

	// if err := store.UpdateUser(u); err != nil {
	// 	w.WriteHeader(501)
	// 	fmt.Printf("{\"status\": \"%s\"}", err.Error())
	// 	json.NewEncoder(w).Encode(
	// 		fmt.Sprintf("{\"status\": \"%s\"}", err.Error()),
	// 	)
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		fmt.Sprintf("{\"status\": \"OK\"}"),
	)
	return
}

func userRemove(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	u := new(store.User)
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		w.WriteHeader(501)
		json.NewEncoder(w).Encode(err)
		return
	}

	// if err := store.DeleteUser(u); err != nil {
	// 	w.WriteHeader(501)
	// 	json.NewEncoder(w).Encode("{ status: Fail }")
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("{\"status\": \"OK\"}")
	return
}
