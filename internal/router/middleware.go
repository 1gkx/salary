package router

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"

	"github.com/1gkx/salary/internal/session"
)

func requireCookieAuth(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	if c, _ := session.Get(r); len(c.ID) > 0 {

		fmt.Printf("Middleware session: %+v\n", c)

		isAuth, _ := c.Values["isAuth"].(bool)
		isVerify, _ := c.Values["isVeryfy"].(bool)

		if isAuth && isVerify {
			next.ServeHTTP(w, r)
			return
		}
	}
	session.Delete(r, w)
	w.Header().Set("Cache-Control", "No-Cache")
	http.Redirect(w, r, "/login", 301)
}

func authRequireHandlerWrap(fn func(http.ResponseWriter, *http.Request, http.HandlerFunc)) *negroni.Negroni {
	return negroni.New(
		negroni.HandlerFunc(requireCookieAuth),
		negroni.HandlerFunc(fn),
	)
}
