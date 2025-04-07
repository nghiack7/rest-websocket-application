package repositories

import (
	"time"

	"github.com/personal/task-management/internal/domain"
	"gorm.io/gorm"
)

type ChatRepository interface {
	// Room operations
	CreateRoom(room *domain.Room) error
	GetRoom(roomID string) (*domain.Room, error)
	UpdateRoom(room *domain.Room) error
	DeleteRoom(roomID string) error
	ListUserRooms(userID string) ([]*domain.Room, error)

	// Message operations
	CreateMessage(message *domain.Message) error
	GetMessage(messageID string) (*domain.Message, error)
	UpdateMessage(message *domain.Message) error
	DeleteMessage(messageID string) error
	GetRoomMessages(roomID string, limit, offset int) ([]*domain.Message, error)

	// Room user operations
	AddUserToRoom(roomID, userID string) error
	RemoveUserFromRoom(roomID, userID string) error
	GetRoomUsers(roomID string) ([]string, error)

	// Message status operations
	UpdateMessageStatus(status *domain.MessageStatus) error
	GetMessageStatus(messageID, userID string) (*domain.MessageStatus, error)

	// Notification operations
	CreateNotification(notification *domain.Notification) error
	GetNotification(notificationID string) (*domain.Notification, error)
	UpdateNotification(notification *domain.Notification) error
	DeleteNotification(notificationID string) error
	GetUserNotifications(userID string, limit, offset int) ([]*domain.Notification, error)
	MarkNotificationAsRead(notificationID string) error
	GetUnreadNotificationCount(userID string) (int, error)
}

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) CreateRoom(room *domain.Room) error {
	return r.db.Create(room).Error
}

func (r *chatRepository) GetRoom(roomID string) (*domain.Room, error) {
	var room domain.Room
	if err := r.db.First(&room, "id = ?", roomID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &room, nil
}

func (r *chatRepository) UpdateRoom(room *domain.Room) error {
	return r.db.Save(room).Error
}

func (r *chatRepository) DeleteRoom(roomID string) error {
	return r.db.Delete(&domain.Room{}, "id = ?", roomID).Error
}

func (r *chatRepository) ListUserRooms(userID string) ([]*domain.Room, error) {
	var rooms []*domain.Room
	if err := r.db.Where("id IN (SELECT room_id FROM room_users WHERE user_id = ?)", userID).Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *chatRepository) CreateMessage(message *domain.Message) error {
	return r.db.Create(message).Error
}

func (r *chatRepository) GetMessage(messageID string) (*domain.Message, error) {
	var message domain.Message
	if err := r.db.First(&message, "id = ?", messageID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &message, nil
}

func (r *chatRepository) UpdateMessage(message *domain.Message) error {
	return r.db.Save(message).Error
}

func (r *chatRepository) DeleteMessage(messageID string) error {
	return r.db.Delete(&domain.Message{}, "id = ?", messageID).Error
}

func (r *chatRepository) GetRoomMessages(roomID string, limit, offset int) ([]*domain.Message, error) {
	var messages []*domain.Message
	if err := r.db.Where("room_id = ?", roomID).Order("created_at DESC").Limit(limit).Offset(offset).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *chatRepository) AddUserToRoom(roomID, userID string) error {
	roomUser := &domain.RoomUser{
		ID:        time.Now().Format("20060102150405") + "_" + time.Now().Format("000000000"),
		RoomID:    roomID,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return r.db.Create(roomUser).Error
}

func (r *chatRepository) RemoveUserFromRoom(roomID, userID string) error {
	return r.db.Delete(&domain.RoomUser{}, "room_id = ? AND user_id = ?", roomID, userID).Error
}

func (r *chatRepository) GetRoomUsers(roomID string) ([]string, error) {
	var userIDs []string
	if err := r.db.Model(&domain.RoomUser{}).Where("room_id = ?", roomID).Pluck("user_id", &userIDs).Error; err != nil {
		return nil, err
	}
	return userIDs, nil
}

func (r *chatRepository) UpdateMessageStatus(status *domain.MessageStatus) error {
	return r.db.Save(status).Error
}

func (r *chatRepository) GetMessageStatus(messageID, userID string) (*domain.MessageStatus, error) {
	var status domain.MessageStatus
	if err := r.db.First(&status, "message_id = ? AND user_id = ?", messageID, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &status, nil
}

func (r *chatRepository) CreateNotification(notification *domain.Notification) error {
	return r.db.Create(notification).Error
}

func (r *chatRepository) GetNotification(notificationID string) (*domain.Notification, error) {
	var notification domain.Notification
	if err := r.db.First(&notification, "id = ?", notificationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &notification, nil
}

func (r *chatRepository) UpdateNotification(notification *domain.Notification) error {
	return r.db.Save(notification).Error
}

func (r *chatRepository) DeleteNotification(notificationID string) error {
	return r.db.Delete(&domain.Notification{}, "id = ?", notificationID).Error
}

func (r *chatRepository) GetUserNotifications(userID string, limit, offset int) ([]*domain.Notification, error) {
	var notifications []*domain.Notification
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Offset(offset).Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *chatRepository) MarkNotificationAsRead(notificationID string) error {
	return r.db.Model(&domain.Notification{}).Where("id = ?", notificationID).Update("is_read", true).Error
}

func (r *chatRepository) GetUnreadNotificationCount(userID string) (int, error) {
	var count int64
	if err := r.db.Model(&domain.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
