package domain

import (
	"errors"
	"time"
)

// Room represents a chat room
type Room struct {
	ID             string         `json:"id" gorm:"primaryKey"`
	Name           string         `json:"name"`
	Type           string         `json:"type"` // "direct" or "group"
	Description    string         `json:"description,omitempty"`
	AvatarURL      string         `json:"avatar_url,omitempty"`
	Users          []string       `json:"users" gorm:"-"`
	LastMessage    *Message       `json:"last_message,omitempty" gorm:"-"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	IsArchived     bool           `json:"is_archived"`
	IsMuted        bool           `json:"is_muted"`
	UnreadCount    map[string]int `json:"unread_count" gorm:"type:jsonb"`
	PinnedMessages []string       `json:"pinned_messages" gorm:"type:text[]"`
}

// Message represents a chat message
type Message struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	RoomID       string    `json:"room_id"`
	UserID       string    `json:"user_id"`
	Content      string    `json:"content"`
	Type         string    `json:"type"`
	FileURL      string    `json:"file_url,omitempty"`
	FileName     string    `json:"file_name,omitempty"`
	FileSize     int64     `json:"file_size,omitempty"`
	FileType     string    `json:"file_type,omitempty"`
	ThumbnailURL string    `json:"thumbnail_url,omitempty"`
	Duration     int       `json:"duration,omitempty"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// RoomUser represents the relationship between rooms and users
type RoomUser struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	RoomID    string    `json:"room_id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MessageStatus represents the status of a message for a specific user
type MessageStatus struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	MessageID string    `json:"message_id"`
	UserID    string    `json:"user_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Notification represents a system notification
type Notification struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Data      string    `json:"data,omitempty"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// WebSocketMessage represents a message sent over WebSocket
type WebSocketMessage struct {
	Type         string    `json:"type"`
	ID           string    `json:"id,omitempty"`
	RoomID       string    `json:"room_id,omitempty"`
	UserID       string    `json:"user_id,omitempty"`
	TargetID     string    `json:"target_id,omitempty"`
	Content      string    `json:"content,omitempty"`
	FileURL      string    `json:"file_url,omitempty"`
	FileName     string    `json:"file_name,omitempty"`
	FileSize     int64     `json:"file_size,omitempty"`
	FileType     string    `json:"file_type,omitempty"`
	ThumbnailURL string    `json:"thumbnail_url,omitempty"`
	Duration     int       `json:"duration,omitempty"`
	MessageID    string    `json:"message_id,omitempty"`
	Status       string    `json:"status,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
}

// Hub maintains active connections and broadcasts messages
type Hub struct {
	Rooms         map[string]*Room
	Connections   map[string]*Connection
	Register      chan *Connection
	Unregister    chan *Connection
	Broadcast     chan WebSocketMessage
	DirectMessage chan WebSocketMessage
}

// Connection represents a WebSocket connection
type Connection struct {
	ID     string
	UserID string
	RoomID string
	Send   chan WebSocketMessage
	Hub    *Hub
}

// Message types
const (
	MessageTypeText       = "text"
	MessageTypeFile       = "file"
	MessageTypeImage      = "image"
	MessageTypeVideo      = "video"
	MessageTypeAudio      = "audio"
	MessageTypeTyping     = "typing"
	MessageTypeRead       = "read"
	MessageTypeTaskUpdate = "task_update"
	MessageTypeMention    = "mention"
	MessageTypeSystem     = "system"
)

// Message statuses
const (
	MessageStatusSent      = "sent"
	MessageStatusDelivered = "delivered"
	MessageStatusRead      = "read"
)

// Room types
const (
	RoomTypeDirect = "direct"
	RoomTypeGroup  = "group"
)

// Notification types
const (
	NotificationTypeTaskUpdate = "task_update"
	NotificationTypeMention    = "mention"
	NotificationTypeSystem     = "system"
)

// Error constants
var (
	ErrRoomNotFound    = errors.New("room not found")
	ErrUserNotInRoom   = errors.New("user not in room")
	ErrInvalidMessage  = errors.New("invalid message")
	ErrInvalidRoomType = errors.New("invalid room type")
)
