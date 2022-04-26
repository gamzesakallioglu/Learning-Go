package db

import (
	"time"

	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/pkg/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPsqlDB(cfg *config.Config) *gorm.DB {

	connectionDsn := cfg.DBConfig.DataSourceName
	db, err := gorm.Open(postgres.Open(connectionDsn), &gorm.Config{})
	if err != nil {
		zap.L().Fatal("cannot connect to db", zap.Error(err)) // Exit the program if db cannot be connected
	}

	sqlDB, err := db.DB()
	if err != nil {
		zap.L().Fatal("cannot get db from postgresql", zap.Error(err)) // Exit the program if db cannot be open
	}

	if err := sqlDB.Ping(); err != nil {
		zap.L().Fatal("cannot ping db", zap.Error(err)) // Exit the program if db cannot be pinged
	}

	// set const variables
	sqlDB.SetMaxOpenConns(cfg.DBConfig.MaxOpen)                                     // max number of open conns to the db
	sqlDB.SetMaxIdleConns(cfg.DBConfig.MaxIdle)                                     // max number of conns in the idle conn pool
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.DBConfig.MaxLifetime) * time.Second) // max time a conn can be reused
	//

	return db
}
