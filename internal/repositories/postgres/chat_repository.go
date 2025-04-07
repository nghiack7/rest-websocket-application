package postgres

import (
	"time"

	"github.com/personal/task-management/internal/domain"
	"github.com/personal/task-management/internal/repositories"
	"gorm.io/gorm"
)

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) repositories.ChatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) CreateRoom(room *domain.Room) error {
	return r.db.Create(room).Error
}

func (r *chatRepository) GetRoom(roomID string) (*domain.Room, error) {
	var room domain.Room
	err := r.db.First(&room, "id = ?", roomID).Error
	if err != nil {
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
	err := r.db.Joins("JOIN room_users ON room_users.room_id = rooms.id").
		Where("room_users.user_id = ?", userID).
		Order("rooms.updated_at DESC").
		Find(&rooms).Error
	return rooms, err
}

func (r *chatRepository) CreateMessage(message *domain.Message) error {
	return r.db.Create(message).Error
}

func (r *chatRepository) GetMessage(messageID string) (*domain.Message, error) {
	var message domain.Message
	err := r.db.First(&message, "id = ?", messageID).Error
	if err != nil {
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
	err := r.db.Where("room_id = ?", roomID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, err
}

func (r *chatRepository) AddUserToRoom(roomID, userID string) error {
	roomUser := &domain.RoomUser{
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
	err := r.db.Model(&domain.RoomUser{}).
		Where("room_id = ?", roomID).
		Pluck("user_id", &userIDs).Error
	return userIDs, err
}

func (r *chatRepository) UpdateMessageStatus(status *domain.MessageStatus) error {
	return r.db.Save(status).Error
}

func (r *chatRepository) GetMessageStatus(messageID, userID string) (*domain.MessageStatus, error) {
	var status domain.MessageStatus
	err := r.db.First(&status, "message_id = ? AND user_id = ?", messageID, userID).Error
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (r *chatRepository) CreateNotification(notification *domain.Notification) error {
	return r.db.Create(notification).Error
}

func (r *chatRepository) GetNotification(notificationID string) (*domain.Notification, error) {
	var notification domain.Notification
	err := r.db.First(&notification, "id = ?", notificationID).Error
	if err != nil {
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
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error
	return notifications, err
}

func (r *chatRepository) MarkNotificationAsRead(notificationID string) error {
	return r.db.Model(&domain.Notification{}).
		Where("id = ?", notificationID).
		Update("is_read", true).Error
}

func (r *chatRepository) GetUnreadNotificationCount(userID string) (int, error) {
	var count int64
	err := r.db.Model(&domain.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return int(count), err
}
