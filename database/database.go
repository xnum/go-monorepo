package database

import (
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Predefined database type.
const (
	Default string = "default"
)

// database defines database instance.
type database struct {
	Opt ConnectOption
	db  *gorm.DB
}

type databaseStore struct {
	sync.Mutex
	dbs map[string]*database
}

var store databaseStore

// Initialize inits singleton.
func Initialize(typ string, opt ConnectOption) {
	store.Lock()
	defer store.Unlock()

	if _, ok := store.dbs[typ]; ok {
		log.Printf("warning: database type:%v was initialized", typ)
		return
	}

	db := &database{}
	db.Opt = opt
	db.Open()

	if store.dbs == nil {
		store.dbs = map[string]*database{}
	}

	store.dbs[typ] = db
}

// Finalize closes singleton.
func Finalize() {
	store.Lock()
	defer store.Unlock()

	for key, db := range store.dbs {
		db.Close()
		delete(store.dbs, key)
	}
}

// GetDB gets db from singleton.
func GetDB(typ string) *gorm.DB {
	store.Lock()
	defer store.Unlock()
	if _, ok := store.dbs[typ]; !ok {
		log.Panicf("db uninited bad type:%v", typ)
	}
	return store.dbs[typ].GetDB()
}

// AutoMigrate migrates table.
func AutoMigrate(typ string, models []any) {
	store.Lock()
	defer store.Unlock()
	if !store.dbs[typ].Opt.Testing {
		panic("migrate to code base just for testing")
	}

	db := store.dbs[typ].db
	for _, m := range models {
		err := db.AutoMigrate(m)
		if err != nil {
			log.Panicf("AutoMigrate(%T): %v", m, err)
		}
	}
}

var (
	defaultLogger = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: 180 * time.Millisecond,
			LogLevel:      logger.Warn,
			Colorful:      true,
		},
	)
)

// Open opens database connection.
func (db *database) Open() {
	if db.db != nil {
		return
	}

	dialector := db.Opt.Dialector()
	if dialector == nil {
		log.Panicf(
			"gorm driver open dialector fail, connect str: (%v)",
			db.Opt.ConnStr(),
		)
	}

	if db.Opt.Silence {
		db.Opt.Config.Logger = logger.Discard
	} else {
		db.Opt.Config.Logger = defaultLogger
	}
	conn, err := gorm.Open(dialector, &db.Opt.Config)
	if err != nil {
		log.Panicf("sql.Open(%v): %v", db.Opt.ConnStr(), err)
	}

	if db.Opt.Dialect == "postgres" {
		err := conn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
		if err != nil {
			log.Panicf("install UUID %v", err)
		}
	}
	realConn, err := conn.DB()
	if err != nil {
		log.Panicf("can't not get real connection %v", err)
	}

	realConn.SetConnMaxLifetime(10 * time.Minute)
	realConn.SetMaxIdleConns(20)
	realConn.SetMaxOpenConns(20)

	db.db = conn
}

// GetDB get gorm db instance.
func (db *database) GetDB() *gorm.DB {
	if db.db == nil {
		panic("database is not initialized.")
	}

	return db.db.Session(&gorm.Session{})
}

// Close closes db connection.
func (db *database) Close() {
	realConn, err := db.db.DB()
	if err == nil {
		realConn.Close()
	}
	db.db = nil
}
