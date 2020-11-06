package store

import (
	// "fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/1gkx/salary/internal/conf"
)

var x *gorm.DB

func Initialize() error {
	var err error

	if conf.Cfg.Database.Driver == "sqlite3" {
		x, err = gorm.Open(conf.Cfg.Database.Driver, conf.Cfg.Database.Path)
		if err != nil {
			return err
		}
	}

	if conf.Cfg.Database.Driver == "postgres" {
		x, err = gorm.Open(conf.Cfg.Database.Driver, conf.Cfg.Database.Path)
		if err != nil {
			return err
		}
	}

	if conf.Cfg.Database.Driver == "mysql" {
		x, err = gorm.Open(conf.Cfg.Database.Driver, conf.Cfg.Database.Path)
		if err != nil {
			return err
		}
	}

	migrate()

	return nil
}

func GetEnginie() *gorm.DB {
	return x
}

// func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
// 	return func(db *gorm.DB) *gorm.DB {
// 		page, _ := strconv.Atoi(r.URL.Query()["page"])
// 		if page == 0 {
// 			page = 1
// 		}

// 		pageSize, _ := strconv.Atoi(r.URL.Query()["page_size"])
// 		switch {
// 		case pageSize > 100:
// 			pageSize = 100
// 		case pageSize <= 0:
// 			pageSize = 10
// 		}

// 		offset := (page - 1) * pageSize
// 		return db.Offset(offset).Limit(pageSize)
// 	}
// }
