package store

func migrate() {

	// Миграция схем
	x.Debug().AutoMigrate(&User{}, &Client{})

}
