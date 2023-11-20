package database

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"gorm.io/gorm"
)

var (
	// PostgresOpt is default connection option for postgres.
	PostgresOpt = ConnectOption{
		Dialect: "postgres",
		Host:    "localhost",
		DBName:  "testing",
		Port:    5432,
		User:    "tester",
		Pass:    "aaaa1234",
	}
)

func randomDBName() string {
	return fmt.Sprintf("testing_%v", time.Now().UnixNano())
}

// TestingInitialize creates new db for testing.
func TestingInitialize(typ string, opt ConnectOption) {
	opt.Config.DisableForeignKeyConstraintWhenMigrating = true
	opt.Testing = true

	Initialize(typ, opt)
	dbName := randomDBName()

	db := GetDB(typ)
	err := db.Exec("CREATE DATABASE " + dbName).Error
	if err != nil {
		log.Panicln(err)
	}

	opt.DBName = dbName
	log.Println("use db name:", dbName)

	store.Lock()
	for key, db := range store.dbs {
		if key != typ {
			continue
		}

		db.Close()
		delete(store.dbs, key)
	}
	store.Unlock()

	Initialize(typ, opt)
}

// TestingFinalize cleanups testing data.
func TestingFinalize() {
	// store.Lock()
	// for _, db := range store.dbs {
	// XXX pq: cannot drop the currently open database if you aren't superusers
	// or owner.
	// err := db.GetDB().Exec("DROP DATABASE " + db.Opt.DBName).Error
	// }
	// store.Unlock()

	Finalize()
}

// DeleteCreatedEntities drop all created data
func DeleteCreatedEntities(db *gorm.DB) func() {
	hookName := "cleanupHook"

	models := make([]any, 0)
	// Setup the onCreate Hook
	db.Callback().
		Create().
		After("gorm:create").
		Register(hookName, func(db *gorm.DB) {
			switch db.Statement.ReflectValue.Kind() {
			case reflect.Slice, reflect.Array:
				db.Statement.CurDestIndex = 0
				for index := 0; index < db.Statement.ReflectValue.Len(); index++ {
					elem := reflect.Indirect(
						db.Statement.ReflectValue.Index(index),
					)
					models = append(models, elem.Addr().Interface())
				}
			case reflect.Struct:
				models = append(
					models,
					db.Statement.ReflectValue.Addr().Interface(),
				)
			}
		})

	return func() {
		for _, model := range models {
			err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).
				Debug().Unscoped().Delete(model).Error
			if err != nil {
				panic(err)
			}
		}
	}
}
