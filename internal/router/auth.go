package router

import (
	"net/http"

	"github.com/1gkx/salary/internal/conf"
	"github.com/1gkx/salary/internal/session"
	"github.com/1gkx/salary/internal/store"
	"github.com/1gkx/salary/internal/utils"
)

var (
	ErrorAuthFail     = map[string]interface{}{"message": "Ошибка авторизации"}
	ErrorUserNotFound = map[string]interface{}{"message": "Пользователь не найден"}
	ErrorServer       = map[string]interface{}{"message": "Ошибка на сервере банка. Повторите попытку позже"}
)

func setSignInRouters() {
	r.HandleFunc("/login", login).Methods("GET")
	r.HandleFunc("/login", signin).Methods("POST")
	r.HandleFunc("/verify", verify).Methods("POST")
	r.HandleFunc("/logout", logout).Methods("GET")
}

func login(w http.ResponseWriter, r *http.Request) {
	responce(w, r, "login", nil)
}

func signin(w http.ResponseWriter, r *http.Request) {

	if !session.CheckAuth(r) {
		responceJson(401, w, ErrorAuthFail)
		return
	}

	u, err := store.FindByEmail(r.FormValue("email"))
	if err != nil {
		responceJson(401, w, ErrorUserNotFound)
		return
	}
	res, err := utils.Post(u.GetPhoneNumber(), conf.Cfg.Methods["SMS"])
	if err != nil {
		responceJson(501, w, ErrorServer)
		return
	}

	vs := map[string]interface{}{
		"sms_code":   res.GetSmsCode(),
		"user":       u.Email,
		"expired_at": res.GetExpiredSmsCode(),
		"isAuth":     true,
		"isVeryfy":   false,
	}
	session.Add(r, w, vs)

	responceJson(http.StatusOK, w, map[string]interface{}{
		"auth": true,
	})
	return
}

func verify(w http.ResponseWriter, r *http.Request) {

	if !session.CheckSms(r) {
		responceJson(401, w, map[string]interface{}{
			"message": "Смс код не подтвержден",
		})
		return
	}

	vs := map[string]interface{}{
		"isVeryfy": true,
	}
	session.Add(r, w, vs)

	responceJson(200, w, map[string]interface{}{
		"verify": true,
	})
	return
}

func logout(w http.ResponseWriter, r *http.Request) {

	session.Reset(r, w)

	http.Redirect(w, r, "/login", 301)
	return
}
