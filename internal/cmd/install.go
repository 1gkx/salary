package cmd

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/1gkx/salary/internal/conf"
	"github.com/1gkx/salary/internal/store"
	"github.com/gorilla/mux"
	"github.com/urfave/cli"
	// "github.com/1gkx/salary/internal/template"
)

var (
	MsgSuccess         = map[string]interface{}{"status": "запрос успешно выполнен"}
	MsgDBNameFail      = map[string]interface{}{"status": "некоректное имя базы данных"}
	MsgUserAlredyExist = map[string]interface{}{"status": "пользователь уже существует"}
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

	if err := conf.Read(); err != nil {
		panic(err)
	}

	if err := json.NewDecoder(r.Body).Decode(&conf.Cfg); err != nil {
		fmt.Printf("Error decode conf: %s\n", err.Error())
		responceAPI(http.StatusInternalServerError, w, err.Error())
		// w.WriteHeader(501)
		// json.NewEncoder(w).Encode(err.Error())
		return
	}

	if conf.Cfg.Database.Driver == "sqlite3" {
		db := strings.Split(conf.Cfg.Database.Path, ".")
		conf.Cfg.Database.Path = fmt.Sprintf("data/%s.db", db[0])
		if len(db) > 2 || len(db) == 0 {
			fmt.Printf("database name incorrect")
			responceAPI(http.StatusInternalServerError, w, MsgDBNameFail)
			// w.WriteHeader(501)
			// json.NewEncoder(w).Encode("database name incorrect")
			return
		}
		if len(db) == 1 {
			conf.Cfg.Database.Path = fmt.Sprintf("data/%s.db", db[0])
		}
	}
	fmt.Printf("Config: %+v\n", conf.Cfg)
	conf.Save()

	if err := store.Initialize(); err != nil {
		panic(err)
	}
	defer store.GetEnginie().Close()

	admin := &store.User{
		FirstName: "Администратор",
		LastName:  "Системы",
		Email:     conf.Cfg.Admin.Email,
		Phone:     "+7 (999) 618-51-15",
		Password:  conf.Cfg.Admin.Password,
		Admin:     "true",
	}
	if err := store.AddUser(admin); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		responceAPI(http.StatusOK, w, MsgSuccess)
		// responceAPI(http.StatusInternalServerError, w, MsgUserAlredyExist)
		// json.NewEncoder(w).Encode(err.Error())
		return
	}

	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(MsgSuccess)
	responceAPI(http.StatusOK, w, MsgSuccess)
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

func responceAPI(code int, w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
