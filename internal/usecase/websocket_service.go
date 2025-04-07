package usecase

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/personal/task-management/internal/domain"
	"github.com/personal/task-management/internal/repositories"
)

type WebSocketService interface {
	// Connection management
	HandleConnection(conn *websocket.Conn, userID string)

	// Room operations
	CreateDirectRoom(userID1, userID2 string) (*domain.Room, error)
	CreateGroupRoom(name string, userIDs []string) (*domain.Room, error)
	JoinRoom(roomID, userID string) error
	LeaveRoom(roomID, userID string) error

	// Message operations
	SendDirectMessage(senderID, receiverID, content string) error
	SendGroupMessage(roomID, userID, content string) error
	SendFileMessage(roomID, userID, fileURL, fileName string, fileSize int64, fileType string) error
	SendImageMessage(roomID, userID, imageURL, thumbnailURL string) error
	SendVideoMessage(roomID, userID, videoURL, thumbnailURL string, duration int) error
	SendAudioMessage(roomID, userID, audioURL string, duration int) error
	SendTypingIndicator(roomID, userID string) error
	MarkMessageAsRead(roomID, userID, messageID string) error
	PinMessage(roomID, messageID string) error
	UnpinMessage(roomID, messageID string) error

	// Room management
	ListRooms(userID string) ([]*domain.Room, error)
	ArchiveRoom(roomID, userID string) error
	UnarchiveRoom(roomID, userID string) error
	MuteRoom(roomID, userID string) error
	UnmuteRoom(roomID, userID string) error
	UpdateRoomInfo(roomID, name, description, avatarURL string) error

	// History and status
	GetRoomHistory(roomID string, limit, offset int) ([]domain.WebSocketMessage, error)
	GetUnreadCount(roomID, userID string) (int, error)

	// Notification operations
	SendTaskUpdateNotification(userID, taskID, taskTitle, taskStatus string) error
	SendMentionNotification(userID, senderID, content string) error
	SendSystemNotification(userID, title, content string) error
	MarkNotificationAsRead(notificationID string) error
	GetUnreadNotificationCount(userID string) (int, error)
}

type websocketService struct {
	hub      *domain.Hub
	roomRepo repositories.ChatRepository
	mu       sync.RWMutex
}

func NewWebSocketService(roomRepo repositories.ChatRepository) WebSocketService {
	hub := &domain.Hub{
		Rooms:         make(map[string]*domain.Room),
		Connections:   make(map[string]*domain.Connection),
		Register:      make(chan *domain.Connection),
		Unregister:    make(chan *domain.Connection),
		Broadcast:     make(chan domain.WebSocketMessage),
		DirectMessage: make(chan domain.WebSocketMessage),
	}

	service := &websocketService{
		hub:      hub,
		roomRepo: roomRepo,
	}

	go service.runHub()
	return service
}

func (s *websocketService) runHub() {
	for {
		select {
		case conn := <-s.hub.Register:
			s.mu.Lock()
			s.hub.Connections[conn.UserID] = conn
			s.mu.Unlock()

		case conn := <-s.hub.Unregister:
			s.mu.Lock()
			delete(s.hub.Connections, conn.UserID)
			if conn.RoomID != "" {
				room, exists := s.hub.Rooms[conn.RoomID]
				if exists {
					for i, userID := range room.Users {
						if userID == conn.UserID {
							room.Users = append(room.Users[:i], room.Users[i+1:]...)
							break
						}
					}
				}
			}
			s.mu.Unlock()

		case message := <-s.hub.DirectMessage:
			s.mu.RLock()
			if targetConn, exists := s.hub.Connections[message.TargetID]; exists {
				targetConn.Send <- message
			}
			s.mu.RUnlock()

		case message := <-s.hub.Broadcast:
			s.mu.RLock()
			if message.RoomID != "" {
				// Group message
				room, exists := s.hub.Rooms[message.RoomID]
				if exists {
					for _, userID := range room.Users {
						if conn, exists := s.hub.Connections[userID]; exists {
							conn.Send <- message
						}
					}
					room.LastMessage = &domain.Message{
						ID:        message.ID,
						RoomID:    message.RoomID,
						UserID:    message.UserID,
						Content:   message.Content,
						Type:      message.Type,
						CreatedAt: message.Timestamp,
						UpdatedAt: message.Timestamp,
					}
				}
			} else if message.Type == domain.MessageTypeTaskUpdate {
				for _, conn := range s.hub.Connections {
					conn.Send <- message
				}
			}
			s.mu.RUnlock()
		}
	}
}

func (s *websocketService) HandleConnection(conn *websocket.Conn, userID string) {
	connection := &domain.Connection{
		ID:     userID,
		UserID: userID,
		Send:   make(chan domain.WebSocketMessage),
		Hub:    s.hub,
	}

	s.hub.Register <- connection

	go s.writePump(conn, connection)
	go s.readPump(conn, connection)
}

func (s *websocketService) CreateDirectRoom(userID1, userID2 string) (*domain.Room, error) {
	room := &domain.Room{
		ID:        generateRoomID(),
		Type:      domain.RoomTypeDirect,
		Users:     []string{userID1, userID2},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.roomRepo.CreateRoom(room); err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.hub.Rooms[room.ID] = room
	s.mu.Unlock()

	return room, nil
}

func (s *websocketService) CreateGroupRoom(name string, userIDs []string) (*domain.Room, error) {
	room := &domain.Room{
		ID:        generateRoomID(),
		Name:      name,
		Type:      domain.RoomTypeGroup,
		Users:     userIDs,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.roomRepo.CreateRoom(room); err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.hub.Rooms[room.ID] = room
	s.mu.Unlock()

	return room, nil
}

func (s *websocketService) JoinRoom(roomID, userID string) error {
	room, err := s.roomRepo.GetRoom(roomID)
	if err != nil {
		return err
	}

	if room == nil {
		return domain.ErrRoomNotFound
	}

	room.Users = append(room.Users, userID)
	if err := s.roomRepo.UpdateRoom(room); err != nil {
		return err
	}

	s.mu.Lock()
	s.hub.Rooms[roomID] = room
	s.mu.Unlock()

	return nil
}

func (s *websocketService) LeaveRoom(roomID, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, exists := s.hub.Rooms[roomID]
	if !exists {
		return errors.New("room not found")
	}

	for i, id := range room.Users {
		if id == userID {
			room.Users = append(room.Users[:i], room.Users[i+1:]...)
			if err := s.roomRepo.RemoveUserFromRoom(roomID, userID); err != nil {
				return err
			}
			return nil
		}
	}

	return domain.ErrUserNotInRoom
}

func (s *websocketService) SendDirectMessage(senderID, receiverID, content string) error {
	// Create or get direct room
	room, err := s.roomRepo.GetRoom(generateDirectRoomID(senderID, receiverID))
	if err != nil {
		return err
	}

	if room == nil {
		room = &domain.Room{
			ID:        generateDirectRoomID(senderID, receiverID),
			Type:      domain.RoomTypeDirect,
			Users:     []string{senderID, receiverID},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := s.roomRepo.CreateRoom(room); err != nil {
			return err
		}
	}

	// Create message
	message := &domain.Message{
		ID:        generateMessageID(),
		RoomID:    room.ID,
		UserID:    senderID,
		Content:   content,
		Type:      domain.MessageTypeText,
		Status:    domain.MessageStatusSent,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.roomRepo.CreateMessage(message); err != nil {
		return err
	}

	// Update room's last message
	room.LastMessage = message
	if err := s.roomRepo.UpdateRoom(room); err != nil {
		return err
	}

	// Send message to receiver
	wsMessage := domain.WebSocketMessage{
		Type:      domain.MessageTypeText,
		ID:        message.ID,
		RoomID:    room.ID,
		UserID:    senderID,
		TargetID:  receiverID,
		Content:   content,
		Timestamp: time.Now(),
	}

	s.hub.DirectMessage <- wsMessage
	return nil
}

func (s *websocketService) SendGroupMessage(roomID, userID, content string) error {
	room, err := s.roomRepo.GetRoom(roomID)
	if err != nil {
		return err
	}

	if room == nil {
		return domain.ErrRoomNotFound
	}

	// Create message
	message := &domain.Message{
		ID:        generateMessageID(),
		RoomID:    roomID,
		UserID:    userID,
		Content:   content,
		Type:      domain.MessageTypeText,
		Status:    domain.MessageStatusSent,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.roomRepo.CreateMessage(message); err != nil {
		return err
	}

	// Update room's last message
	room.LastMessage = message
	if err := s.roomRepo.UpdateRoom(room); err != nil {
		return err
	}

	// Send message to all room users
	wsMessage := domain.WebSocketMessage{
		Type:      domain.MessageTypeText,
		ID:        message.ID,
		RoomID:    roomID,
		UserID:    userID,
		Content:   content,
		Timestamp: time.Now(),
	}

	s.hub.Broadcast <- wsMessage
	return nil
}

func (s *websocketService) SendFileMessage(roomID, userID, fileURL, fileName string, fileSize int64, fileType string) error {
	message := &domain.Message{
		ID:        generateMessageID(),
		RoomID:    roomID,
		UserID:    userID,
		Type:      domain.MessageTypeFile,
		FileURL:   fileURL,
		FileName:  fileName,
		FileSize:  fileSize,
		FileType:  fileType,
		Status:    domain.MessageStatusSent,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.roomRepo.CreateMessage(message); err != nil {
		return err
	}

	wsMessage := domain.WebSocketMessage{
		Type:      domain.MessageTypeFile,
		ID:        message.ID,
		RoomID:    roomID,
		UserID:    userID,
		FileURL:   fileURL,
		FileName:  fileName,
		FileSize:  fileSize,
		FileType:  fileType,
		Timestamp: time.Now(),
	}

	s.hub.Broadcast <- wsMessage
	return nil
}

func (s *websocketService) SendImageMessage(roomID, userID, imageURL, thumbnailURL string) error {
	message := &domain.Message{
		ID:           generateMessageID(),
		RoomID:       roomID,
		UserID:       userID,
		Type:         domain.MessageTypeImage,
		FileURL:      imageURL,
		ThumbnailURL: thumbnailURL,
		Status:       domain.MessageStatusSent,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.roomRepo.CreateMessage(message); err != nil {
		return err
	}

	wsMessage := domain.WebSocketMessage{
		Type:         domain.MessageTypeImage,
		ID:           message.ID,
		RoomID:       roomID,
		UserID:       userID,
		FileURL:      imageURL,
		ThumbnailURL: thumbnailURL,
		Timestamp:    time.Now(),
	}

	s.hub.Broadcast <- wsMessage
	return nil
}

func (s *websocketService) SendVideoMessage(roomID, userID, videoURL, thumbnailURL string, duration int) error {
	message := &domain.Message{
		ID:           generateMessageID(),
		RoomID:       roomID,
		UserID:       userID,
		Type:         domain.MessageTypeVideo,
		FileURL:      videoURL,
		ThumbnailURL: thumbnailURL,
		Duration:     duration,
		Status:       domain.MessageStatusSent,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.roomRepo.CreateMessage(message); err != nil {
		return err
	}

	wsMessage := domain.WebSocketMessage{
		Type:         domain.MessageTypeVideo,
		ID:           message.ID,
		RoomID:       roomID,
		UserID:       userID,
		FileURL:      videoURL,
		ThumbnailURL: thumbnailURL,
		Duration:     duration,
		Timestamp:    time.Now(),
	}

	s.hub.Broadcast <- wsMessage
	return nil
}

func (s *websocketService) SendAudioMessage(roomID, userID, audioURL string, duration int) error {
	message := &domain.Message{
		ID:        generateMessageID(),
		RoomID:    roomID,
		UserID:    userID,
		Type:      domain.MessageTypeAudio,
		FileURL:   audioURL,
		Duration:  duration,
		Status:    domain.MessageStatusSent,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.roomRepo.CreateMessage(message); err != nil {
		return err
	}

	wsMessage := domain.WebSocketMessage{
		Type:      domain.MessageTypeAudio,
		ID:        message.ID,
		RoomID:    roomID,
		UserID:    userID,
		FileURL:   audioURL,
		Duration:  duration,
		Timestamp: time.Now(),
	}

	s.hub.Broadcast <- wsMessage
	return nil
}

func (s *websocketService) SendTypingIndicator(roomID, userID string) error {
	message := domain.WebSocketMessage{
		Type:      domain.MessageTypeTyping,
		RoomID:    roomID,
		UserID:    userID,
		Timestamp: time.Now(),
	}

	s.hub.Broadcast <- message
	return nil
}

func (s *websocketService) MarkMessageAsRead(roomID, userID, messageID string) error {
	// Update message status in database
	status := &domain.MessageStatus{
		ID:        generateMessageStatusID(),
		MessageID: messageID,
		UserID:    userID,
		Status:    domain.MessageStatusRead,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.roomRepo.UpdateMessageStatus(status); err != nil {
		return err
	}

	// Update unread count for the room
	room, err := s.roomRepo.GetRoom(roomID)
	if err != nil {
		return err
	}

	if room == nil {
		return domain.ErrRoomNotFound
	}

	if room.UnreadCount == nil {
		room.UnreadCount = make(map[string]int)
	}
	room.UnreadCount[userID] = 0

	if err := s.roomRepo.UpdateRoom(room); err != nil {
		return err
	}

	// Send read receipt
	message := domain.WebSocketMessage{
		Type:      domain.MessageTypeRead,
		RoomID:    roomID,
		UserID:    userID,
		MessageID: messageID,
		Status:    domain.MessageStatusRead,
		Timestamp: time.Now(),
	}

	s.hub.Broadcast <- message
	return nil
}

func (s *websocketService) PinMessage(roomID, messageID string) error {
	room, err := s.roomRepo.GetRoom(roomID)
	if err != nil {
		return err
	}

	if room == nil {
		return domain.ErrRoomNotFound
	}

	// Check if message is already pinned
	for _, pinnedID := range room.PinnedMessages {
		if pinnedID == messageID {
			return nil // Message is already pinned
		}
	}

	room.PinnedMessages = append(room.PinnedMessages, messageID)
	return s.roomRepo.UpdateRoom(room)
}

func (s *websocketService) UnpinMessage(roomID, messageID string) error {
	room, err := s.roomRepo.GetRoom(roomID)
	if err != nil {
		return err
	}

	if room == nil {
		return domain.ErrRoomNotFound
	}

	// Remove message from pinned messages
	for i, pinnedID := range room.PinnedMessages {
		if pinnedID == messageID {
			room.PinnedMessages = append(room.PinnedMessages[:i], room.PinnedMessages[i+1:]...)
			return s.roomRepo.UpdateRoom(room)
		}
	}

	return nil // Message was not pinned
}

func (s *websocketService) ArchiveRoom(roomID, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, exists := s.hub.Rooms[roomID]
	if !exists {
		return domain.ErrRoomNotFound
	}

	room.IsArchived = true
	if err := s.roomRepo.UpdateRoom(room); err != nil {
		return err
	}

	return nil
}

func (s *websocketService) UnarchiveRoom(roomID, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, exists := s.hub.Rooms[roomID]
	if !exists {
		return domain.ErrRoomNotFound
	}

	room.IsArchived = false
	if err := s.roomRepo.UpdateRoom(room); err != nil {
		return err
	}

	return nil
}

func (s *websocketService) MuteRoom(roomID, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, exists := s.hub.Rooms[roomID]
	if !exists {
		return domain.ErrRoomNotFound
	}

	room.IsMuted = true
	if err := s.roomRepo.UpdateRoom(room); err != nil {
		return err
	}

	return nil
}

func (s *websocketService) UnmuteRoom(roomID, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, exists := s.hub.Rooms[roomID]
	if !exists {
		return domain.ErrRoomNotFound
	}

	room.IsMuted = false
	if err := s.roomRepo.UpdateRoom(room); err != nil {
		return err
	}

	return nil
}

func (s *websocketService) GetUnreadCount(roomID, userID string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	room, exists := s.hub.Rooms[roomID]
	if !exists {
		return 0, domain.ErrRoomNotFound
	}

	if room.UnreadCount == nil {
		return 0, nil
	}

	return room.UnreadCount[userID], nil
}

func (s *websocketService) UpdateRoomInfo(roomID, name, description, avatarURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, exists := s.hub.Rooms[roomID]
	if !exists {
		return domain.ErrRoomNotFound
	}

	if name != "" {
		room.Name = name
	}
	if description != "" {
		room.Description = description
	}
	if avatarURL != "" {
		room.AvatarURL = avatarURL
	}

	if err := s.roomRepo.UpdateRoom(room); err != nil {
		return err
	}

	return nil
}

func (s *websocketService) ListRooms(userID string) ([]*domain.Room, error) {
	rooms, err := s.roomRepo.ListUserRooms(userID)
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (s *websocketService) GetRoomHistory(roomID string, limit, offset int) ([]domain.WebSocketMessage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.hub.Rooms[roomID]
	if !exists {
		return nil, domain.ErrRoomNotFound
	}

	messages, err := s.roomRepo.GetRoomMessages(roomID, limit, offset)
	if err != nil {
		return nil, err
	}

	wsMessages := make([]domain.WebSocketMessage, len(messages))
	for i, msg := range messages {
		wsMessages[i] = domain.WebSocketMessage{
			Type:         msg.Type,
			ID:           msg.ID,
			RoomID:       msg.RoomID,
			UserID:       msg.UserID,
			Content:      msg.Content,
			FileURL:      msg.FileURL,
			FileName:     msg.FileName,
			FileSize:     msg.FileSize,
			FileType:     msg.FileType,
			ThumbnailURL: msg.ThumbnailURL,
			Duration:     msg.Duration,
			Status:       msg.Status,
			Timestamp:    msg.CreatedAt,
		}
	}

	return wsMessages, nil
}

func (s *websocketService) writePump(conn *websocket.Conn, c *domain.Connection) {
	defer func() {
		conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			json.NewEncoder(w).Encode(message)
		}
	}
}

func (s *websocketService) readPump(conn *websocket.Conn, c *domain.Connection) {
	defer func() {
		s.hub.Unregister <- c
		conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var wsMessage domain.WebSocketMessage
		if err := json.Unmarshal(message, &wsMessage); err != nil {
			log.Printf("error unmarshaling message: %v", err)
			continue
		}

		switch wsMessage.Type {
		case domain.RoomTypeDirect:
			s.hub.DirectMessage <- wsMessage
		case domain.RoomTypeGroup:
			s.hub.Broadcast <- wsMessage
		default:
			s.hub.Broadcast <- wsMessage
		}
	}
}

func generateRoomID() string {
	return time.Now().Format("20060102150405") + "_" + time.Now().Format("000000000")
}

func generateMessageID() string {
	return time.Now().Format("20060102150405") + "_" + time.Now().Format("000000000")
}

func generateMessageStatusID() string {
	return time.Now().Format("20060102150405") + "_" + time.Now().Format("000000000")
}

func generateDirectRoomID(userID1, userID2 string) string {
	if userID1 < userID2 {
		return userID1 + "_" + userID2
	}
	return userID2 + "_" + userID1
}

// Notification methods
func (s *websocketService) SendTaskUpdateNotification(userID, taskID, taskTitle, taskStatus string) error {
	notification := &domain.Notification{
		ID:        generateNotificationID(),
		UserID:    userID,
		Type:      domain.NotificationTypeTaskUpdate,
		Title:     "Task Update",
		Content:   taskTitle + " status changed to " + taskStatus,
		Data:      `{"task_id": "` + taskID + `", "task_title": "` + taskTitle + `", "task_status": "` + taskStatus + `"}`,
		IsRead:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.roomRepo.CreateNotification(notification); err != nil {
		return err
	}

	message := domain.WebSocketMessage{
		Type:      domain.MessageTypeTaskUpdate,
		ID:        notification.ID,
		UserID:    userID,
		Content:   notification.Content,
		Timestamp: time.Now(),
	}

	s.hub.DirectMessage <- message
	return nil
}

func (s *websocketService) SendMentionNotification(userID, senderID, content string) error {
	notification := &domain.Notification{
		ID:        generateNotificationID(),
		UserID:    userID,
		Type:      domain.NotificationTypeMention,
		Title:     "You were mentioned",
		Content:   content,
		Data:      `{"sender_id": "` + senderID + `"}`,
		IsRead:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.roomRepo.CreateNotification(notification); err != nil {
		return err
	}

	message := domain.WebSocketMessage{
		Type:      domain.MessageTypeMention,
		ID:        notification.ID,
		UserID:    userID,
		Content:   notification.Content,
		Timestamp: time.Now(),
	}

	s.hub.DirectMessage <- message
	return nil
}

func (s *websocketService) SendSystemNotification(userID, title, content string) error {
	notification := &domain.Notification{
		ID:        generateNotificationID(),
		UserID:    userID,
		Type:      domain.NotificationTypeSystem,
		Title:     title,
		Content:   content,
		IsRead:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.roomRepo.CreateNotification(notification); err != nil {
		return err
	}

	message := domain.WebSocketMessage{
		Type:      domain.MessageTypeSystem,
		ID:        notification.ID,
		UserID:    userID,
		Content:   notification.Content,
		Timestamp: time.Now(),
	}

	s.hub.DirectMessage <- message
	return nil
}

func (s *websocketService) MarkNotificationAsRead(notificationID string) error {
	return s.roomRepo.MarkNotificationAsRead(notificationID)
}

func (s *websocketService) GetUnreadNotificationCount(userID string) (int, error) {
	return s.roomRepo.GetUnreadNotificationCount(userID)
}

func generateNotificationID() string {
	return time.Now().Format("20060102150405") + "_" + time.Now().Format("000000000")
}
