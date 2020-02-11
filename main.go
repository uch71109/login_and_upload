package main

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"login_and_upload/file"
	"login_and_upload/model"
	rrest "login_and_upload/rest"
)

const (
	concurrentJobs = 100
)

func main() {
	fileOp := file.NewFile(concurrentJobs)

	db, err := gorm.Open("sqlite3", "file::memory:")
	if err != nil {
		log.Fatal(err)
	}
	dataOp := model.New(db)
	if err := dataOp.Init(); err != nil {
		log.Fatal(err)
	}

	// admin user
	dataOp.User.Create("admin@gmail.com", "admin@gmail.com", model.RoleAdmin)
	// normal user
	dataOp.User.Create("user1@gmail.com", "user1@gmail.com", model.RoleNormal)

	rest := rrest.RegisterRest(dataOp, fileOp)
	if err := rest.Serve(); err != nil {
		log.Fatal(err)
	}
}
