package main

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var database *gorm.DB

/*
* models is an interface array for database models.
* In order for the model to be migrated, it must be added to this array.
 */
var models = []interface{}{
	&UserModel{},
	&PostModel{},
	&CategoryModel{},
	&PostCommentModel{},
	&ActionModel{},
}

func init() {
	var dialector gorm.Dialector

	if loadedConfig.Database.UseSqlite {
		dialector = sqlite.Open(loadedConfig.Database.SqlitePath)
	} else {
		dialector = mysql.Open(loadedConfig.Database.DSN)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(models...)
	if err != nil {
		log.Fatal(err)
	}

	database = db
}
