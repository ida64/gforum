package main

import (
	"log"

	"gorm.io/driver/mysql"
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
	db, err := gorm.Open(mysql.Open(loadedConfig.Database.DSN), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(Models...)
	if err != nil {
		log.Fatal(err)
	}

	database = db
}
