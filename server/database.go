package main

import (
	"log"
	"reflect"

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

func transferDatabase(source *gorm.DB, destination *gorm.DB) error {
	var transferableModels = []interface{}{
		&[]UserModel{},
		&[]PostModel{},
		&[]CategoryModel{},
		&[]PostCommentModel{},
		&[]ActionModel{},
	}

	for _, model := range transferableModels {
		v := reflect.ValueOf(model).Elem()
		source.Find(model)

		for i := 0; i < v.Len(); i++ {
			item := v.Index(i).Addr().Interface()

			err := destination.Create(item).Error
			if err != nil {
				log.Println(err)
			}
		}
	}

	return nil
}

func init() {
	var dialector gorm.Dialector

	if loadedConfig.Database.UseSqlite {
		dialector = sqlite.Open(loadedConfig.Database.SqlitePath)
	} else {
		dialector = mysql.Open(loadedConfig.Database.DSN)
	}

	var config = &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	db, err := gorm.Open(dialector, config)
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(models...)
	if err != nil {
		log.Fatal(err)
	}

	if loadedConfig.Database.UseSqlite {
		if loadedConfig.Database.MigrateToSqlite {
			mysqlDialector := mysql.Open(loadedConfig.Database.DSN)

			mysqlDatabase, err := gorm.Open(mysqlDialector, config)
			if err != nil {
				log.Fatal(err)
			}

			err = transferDatabase(mysqlDatabase, db)
			if err != nil {
				log.Fatal(err)
			}

			loadedConfig.Database.MigrateToSqlite = false

			err = saveConfig("config.yaml", loadedConfig)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	database = db
}
