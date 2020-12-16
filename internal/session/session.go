package session

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/wader/gormstore"

	"github.com/1gkx/salary/internal/store"
)

var S *gormstore.Store

const (
	cookieName   = "__ssid"
	cookieAgeMax = 60 * 60 * 1 // 1 hours
)

type sms struct {
	Code    string
	Expired string
}

type SessionStruct struct {
	*sessions.Session
	// User     store.User
	// Sms      sms
	// IsAuth   bool
	// IsVerify bool
}

// Start ...
func Start() {
	S = gormstore.NewOptions(
		store.GetEnginie(),
		gormstore.Options{
			TableName:       "sessions",
			SkipCreateTable: false,
		},
		[]byte("secret-hash-key"),
	)

	S.SessionOpts.Secure = true
	S.SessionOpts.HttpOnly = true
	S.SessionOpts.MaxAge = cookieAgeMax
}

func Get(r *http.Request) (*sessions.Session, error) {
	return S.Get(r, cookieName)
}

// GetUser ...
func GetUser(r *http.Request) *store.User {
	c, _ := S.Get(r, cookieName)
	email, _ := c.Values["user"].(string)

	if u, err := store.FindByEmail(email); err == nil {
		return u
	}
	return nil
}

func IsAdmin(r *http.Request) bool {
	c, _ := S.Get(r, cookieName)
	isAdmin, _ := c.Values["isAdmin"].(bool)
	return isAdmin
}

// func (c *sessions.Session) IsAuth() bool {
// 	isAdmin, _ := c.Values["isAuth"].(bool)
// 	return isAdmin
// }

// func Reset(r *http.Request, w http.ResponseWriter) {
// 	vs := map[string]interface{}{
// 		"sms_code":   nil,
// 		"user":       nil,
// 		"expired_at": nil,
// 		"isAuth":     false,
// 		"isVeryfy":   false,
// 	}
// 	Add(r, w, vs)
// }

func Delete(r *http.Request, w http.ResponseWriter) {

	S.MaxAge(-1)
	st, err := S.New(r, cookieName)
	if err != nil {
		fmt.Println(err.Error())
	}
	if err = S.Save(r, w, st); err != nil {
		fmt.Println(err.Error())
	}
	S.Cleanup()
	S.MaxAge(cookieAgeMax)
}

// Add ...
func Add(r *http.Request, w http.ResponseWriter, val map[string]interface{}) error {

	st, err := S.Get(r, cookieName)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if len(st.ID) == 0 {
		st, err = S.New(r, cookieName)
		if err != nil {
			fmt.Println(err.Error())
		}
		// fmt.Printf("Session nil: %+v\n", st)
	}

	for key, value := range val {
		st.Values[key] = value
	}

	if err = S.Save(r, w, st); err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil

}

// CheckAuth Проверка логина и пароля
func CheckAuth(r *http.Request) bool {
	var (
		u   *store.User
		err error
	)
	if u, err = store.FindByEmail(r.FormValue("email")); err != nil {
		fmt.Printf("auth failed: %s\n", err.Error())
		return false
	}

	return u.ComparePass(r.FormValue("password"))
}

func CheckSms(r *http.Request) bool {

	c, _ := S.Get(r, cookieName)

	// Проверяем, что пользователь авторизован
	if isAuth, _ := c.Values["isAuth"].(bool); !isAuth {
		return false
	}

	// Проверяем, что смс совпадает
	if sms_code, _ := c.Values["sms_code"].(string); sms_code != r.FormValue("sms") {
		return false
	}

	t, _ := c.Values["expiried_at"].(string)
	expired_at, _ := time.Parse(time.RFC3339, t)

	if dur := expired_at.Sub(time.Now()).Minutes(); dur >= 2 {
		return false
	}

	return true
}
