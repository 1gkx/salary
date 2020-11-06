package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/1gkx/salary/internal/conf"
	"github.com/gorilla/mux"
)

func RouterInstall() *mux.Router {

	r = mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", install).Methods("GET")
	r.HandleFunc("/", setSettings2).Methods("POST")

	r.NotFoundHandler = http.HandlerFunc(errorhendler)

	return r
}

func install(w http.ResponseWriter, r *http.Request) {
	responce(w, r, "install", nil)
}

func setSettings2(w http.ResponseWriter, r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&conf.Cfg); err != nil {
		w.WriteHeader(501)
		fmt.Printf("Error: %s\n", err.Error())
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	conf.Save()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"status":  "success",
		"message": "Данные успешно сохранены",
	})
	return
}
