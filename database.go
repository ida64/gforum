package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var database *gorm.DB

/*
* Models is an interface array for database models.
* In order for the model to be migrated, it must be added to this array.
 */
var Models = []interface{}{
	&UserModel{},
	&PostModel{},
	&CategoryModel{},
	&PostCommentModel{},
	&ActionModel{},
}

func init() {
	db, err := gorm.Open(sqlite.Open(loadedConfig.Database.Path), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(Models...)
	if err != nil {
		log.Fatal(err)
	}

	database = db
}
