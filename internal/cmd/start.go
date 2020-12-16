package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/urfave/cli"

	"github.com/1gkx/salary/internal/conf"
	"github.com/1gkx/salary/internal/router"
	"github.com/1gkx/salary/internal/session"
	"github.com/1gkx/salary/internal/store"
	"github.com/1gkx/salary/internal/template"
)

var Start = cli.Command{
	Name:        "start",
	Usage:       "Start web server",
	Description: `Description`,
	Action:      runWeb,
}

func runWeb(c *cli.Context) {

	if err := conf.Read(); err != nil {
		panic(err)
	}
	if err := store.Initialize(); err != nil {
		panic(err)
	}

	session.Start()
	template.InitTemplate()

	router := router.NewRouter()

	fmt.Printf("Server start at %s port\n", ":8000")
	if err := http.ListenAndServeTLS(":443", "conf/cert/cert.pem", "conf/cert/key.pem", router); err != nil {
		log.Fatal(err)
	}
	// if conf.Prod() {
	// 	// Для production версии
	// 	if err := http.ListenAndServe(":8000", router); err != nil {
	// 		log.Fatal(err)
	// 	}
	// } else {
	// 	// Для local версии. Т.к. куки шифрованные, необходимо https соединение. Добавить в host: 127.0.0.1 salary.pskb.ad
	// 	// Сертификаты для localhost в папке conf/cert
	// 	if err := http.ListenAndServeTLS(":443", "conf/cert/cert.pem", "conf/cert/key.pem", router); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

}
