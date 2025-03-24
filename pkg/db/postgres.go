package db

import (
	"fmt"

	"github.com/personal/task-management/internal/domain/task"
	"github.com/personal/task-management/internal/domain/user"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
	db *gorm.DB
}

func ConnectDB(config *viper.Viper) *PostgresDB {
	dns := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.GetString("database.host"),
		config.GetInt("database.port"),
		config.GetString("database.user"),
		config.GetString("database.password"),
		config.GetString("database.name"),
		config.GetString("database.ssl_mode"))
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return &PostgresDB{db: db}
}

func (db *PostgresDB) Close() {
	sqlDB, err := db.db.DB()
	if err != nil {
		panic("failed to close database")
	}

	sqlDB.Close()
}

func (db *PostgresDB) GetDB() *gorm.DB {
	return db.db
}

func (db *PostgresDB) MigrateDB() {
	db.db.AutoMigrate(&user.User{}, &task.Task{}) // basic migration
}
