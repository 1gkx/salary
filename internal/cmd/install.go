package cmd

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/1gkx/salary/internal/conf"
	"github.com/gorilla/mux"
	"github.com/urfave/cli"
	// "github.com/1gkx/salary/internal/template"
)

var Install = cli.Command{
	Name:        "install",
	Usage:       "Install and configuration web server",
	Description: `Description`,
	Action:      install,
}

var T *template.Template

func install(c *cli.Context) {

	T = template.Must(template.New("").ParseFiles(
		"templates/install.html", "templates/status/400.html",
	))

	router := mux.NewRouter().StrictSlash(true)
	router.NotFoundHandler = http.HandlerFunc(errorhendler)
	router.HandleFunc("/", getInstallForm).Methods("GET")
	router.HandleFunc("/", setSettings2).Methods("POST")
	publicFolder := http.FileServer(http.Dir("./public"))
	router.PathPrefix("/js/").Handler(publicFolder)
	router.PathPrefix("/css/").Handler(publicFolder)
	router.PathPrefix("/img/").Handler(publicFolder)

	if err := http.ListenAndServeTLS(":443", "conf/cert/cert.pem", "conf/cert/key.pem", router); err != nil {
		log.Fatal(err)
	}

}
func errorhendler(w http.ResponseWriter, r *http.Request) {
	responce(w, r, "status/40x", nil)
}

func getInstallForm(w http.ResponseWriter, r *http.Request) {
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

func responce(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	T.ExecuteTemplate(w, tmpl,
		map[string]interface{}{
			// "user": session.GetUser(r),
			"data": data,
		},
	)
}
