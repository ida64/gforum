package database

import (
	"log"
	"reflect"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	config "gforum/config"
)

var Database *gorm.DB

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

	if config.LoadedConfig.Database.UseSqlite {
		dialector = sqlite.Open(config.LoadedConfig.Database.SqlitePath)
	} else {
		dialector = mysql.Open(config.LoadedConfig.Database.DSN)
	}

	var gormConfig = &gorm.Config{}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(models...)
	if err != nil {
		log.Fatal(err)
	}

	if config.LoadedConfig.Database.UseSqlite {
		if config.LoadedConfig.Database.MigrateToSqlite {
			mysqlDialector := mysql.Open(config.LoadedConfig.Database.DSN)

			mysqlDatabase, err := gorm.Open(mysqlDialector, gormConfig)
			if err != nil {
				log.Fatal(err)
			}

			err = transferDatabase(mysqlDatabase, db)
			if err != nil {
				log.Fatal(err)
			}

			config.LoadedConfig.Database.MigrateToSqlite = false

			err = config.SaveConfig("config.yaml", config.LoadedConfig)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	Database = db
}
