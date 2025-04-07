package migrations

import (
	"github.com/personal/task-management/internal/domain"
	"gorm.io/gorm"
)

func MigrateChatTables(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&domain.Room{},
		&domain.Message{},
		&domain.RoomUser{},
		&domain.MessageStatus{},
	); err != nil {
		return err
	}

	return nil
}
