package db

import (
	"log"
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/db/migrator"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	ErrDatabaseConnection = errors.New("failed to establish connection to the database")
	ErrDatabaseMigration  = errors.New("failed to perform database migration")
)

// GetInstance retrieves instance of the database, performing initial connection if needed.
var GetInstance = sync.OnceValue[*gorm.DB](func() *gorm.DB {
	db, err := connect()
	if err != nil {
		log.Fatalln(err)
	}

	return db
})

// Init initializes database instance, performing migrations first of all.
func Init() {
	instance := GetInstance()

	if err := migrateDatabase(instance); err != nil {
		log.Fatalln(errors.Wrap(err, ErrDatabaseMigration.Error()))
	}
}

// connect establishes database connection.
func connect() (*gorm.DB, error) {
	retryTicker := time.NewTimer(config.GetDatabaseConnectionRetryDelay())

	var (
		connection *gorm.DB
		err        error
	)

	for range retryTicker.C {
		retryTicker.Stop()

		connection, err = gorm.Open(sqlite.Open(config.GetDatabaseName()), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			log.Println(errors.Wrap(err, ErrDatabaseConnection.Error()).Error())

			retryTicker.Reset(config.GetDatabaseConnectionRetryDelay())
			continue
		}

		break
	}

	return connection, nil
}

// migrateDatabase performs database migration.
func migrateDatabase(src *gorm.DB) error {
	goose.SetBaseFS(migrator.Migrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	db, err := src.DB()
	if err != nil {
		return err
	}

	err = goose.Up(db, "migration")
	if err != nil {
		return err
	}

	return src.AutoMigrate()
}
