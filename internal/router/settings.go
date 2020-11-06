package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/1gkx/salary/internal/conf"
)

func setSettingRouters() {
	r.Handle("/admin/settings", authRequireHandlerWrap(getSettings)).Methods("GET")
	r.Handle("/admin/settings", authRequireHandlerWrap(setSettings)).Methods("POST")
}

func getSettings(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	responce(w, r, "settings", conf.Cfg)
}

func setSettings(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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
	// templates.Templates.ExecuteTemplate(w, "settings", conf.Cfg)
}
