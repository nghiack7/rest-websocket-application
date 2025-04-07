package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/personal/task-management/internal/delivery/rest/dtos"
	"github.com/personal/task-management/internal/usecase"
	"github.com/personal/task-management/pkg/utils/jwt"
)

// ChatHandler handles chat-related HTTP requests
type ChatHandler struct {
	wsService usecase.WebSocketService

	jwtService jwt.JWTTokenServicer
}

// NewChatHandler creates a new ChatHandler instance
func NewChatHandler(wsService usecase.WebSocketService, jwtService jwt.JWTTokenServicer) *ChatHandler {
	return &ChatHandler{
		wsService:  wsService,
		jwtService: jwtService,
	}
}

// CreateDirectRoom godoc
// @Summary Create a direct chat room between two users
// @Description Creates a new direct chat room between the authenticated user and another user
// @Tags chat
// @Accept json
// @Produce json
// @Param request body dtos.CreateDirectRoomRequest true "Create Direct Room Request"
// @Success 200 {object} interface{} "Room created successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /chat/direct [post]
func (h *ChatHandler) CreateDirectRoom(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req dtos.CreateDirectRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	room, err := h.wsService.CreateDirectRoom(userID, req.UserID2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(room)
}

// CreateGroupRoom godoc
// @Summary Create a group chat room
// @Description Creates a new group chat room with multiple users
// @Tags chat
// @Accept json
// @Produce json
// @Param request body dtos.CreateGroupRoomRequest true "Create Group Room Request"
// @Success 200 {object} interface{} "Room created successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /chat/group [post]
func (h *ChatHandler) CreateGroupRoom(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateGroupRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	room, err := h.wsService.CreateGroupRoom(req.Name, req.UserIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(room)
}

// ListRooms godoc
// @Summary List all chat rooms for the authenticated user
// @Description Returns a list of all chat rooms the authenticated user is a member of
// @Tags chat
// @Produce json
// @Success 200 {array} interface{} "List of chat rooms"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /chat/rooms [get]
func (h *ChatHandler) ListRooms(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	rooms, err := h.wsService.ListRooms(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(rooms)
}

// GetRoomHistory godoc
// @Summary Get chat room history
// @Description Retrieves the message history for a specific chat room
// @Tags chat
// @Produce json
// @Param roomId path string true "Room ID"
// @Param limit query integer false "Number of messages to return" default(50)
// @Param offset query integer false "Number of messages to skip" default(0)
// @Success 200 {object} interface{} "Room history"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /chat/rooms/{roomId}/history [get]
func (h *ChatHandler) GetRoomHistory(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	room, err := h.wsService.GetRoomHistory(roomID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(room)
}

// JoinRoom godoc
// @Summary Join a chat room
// @Description Adds the authenticated user to a chat room
// @Tags chat
// @Param roomId path string true "Room ID"
// @Success 200 "Successfully joined room"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /chat/rooms/{roomId}/join [post]
func (h *ChatHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	roomID := chi.URLParam(r, "roomId")

	if err := h.wsService.JoinRoom(roomID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// LeaveRoom godoc
// @Summary Leave a chat room
// @Description Removes the authenticated user from a chat room
// @Tags chat
// @Param roomId path string true "Room ID"
// @Success 200 "Successfully left room"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /chat/rooms/{roomId}/leave [post]
func (h *ChatHandler) LeaveRoom(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	roomID := chi.URLParam(r, "roomId")

	if err := h.wsService.LeaveRoom(roomID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateRoom godoc
// @Summary Update chat room information
// @Description Updates the name, description, or avatar of a chat room
// @Tags chat
// @Accept json
// @Param roomId path string true "Room ID"
// @Param request body dtos.UpdateRoomRequest true "Update Room Request"
// @Success 200 "Room updated successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /chat/rooms/{roomId} [put]
func (h *ChatHandler) UpdateRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")

	var req dtos.UpdateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.wsService.UpdateRoomInfo(roomID, req.Name, req.Description, req.AvatarURL); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetMessages godoc
// @Summary Get messages from a chat room
// @Description Retrieves messages from a specific chat room with pagination
// @Tags chat
// @Produce json
// @Param roomId path string true "Room ID"
// @Param limit query integer false "Number of messages to return" default(50)
// @Param offset query integer false "Number of messages to skip" default(0)
// @Success 200 {array} interface{} "List of messages"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /chat/rooms/{roomId}/messages [get]
func (h *ChatHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	messages, err := h.wsService.GetRoomHistory(roomID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(messages)
}

// SendMessage godoc
// @Summary Send a message to a chat room
// @Description Sends a message to a specific chat room
// @Tags chat
// @Accept json
// @Param roomId path string true "Room ID"
// @Param request body dtos.SendMessageRequest true "Send Message Request"
// @Success 200 "Message sent successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /chat/rooms/{roomId}/messages [post]
func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	roomID := chi.URLParam(r, "roomId")

	var req dtos.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var err error
	switch req.Type {
	case "text":
		err = h.wsService.SendGroupMessage(roomID, userID, req.Content)
	case "file":
		err = h.wsService.SendFileMessage(roomID, userID, req.FileURL, "", 0, "")
	case "image":
		err = h.wsService.SendImageMessage(roomID, userID, req.FileURL, "")
	case "video":
		err = h.wsService.SendVideoMessage(roomID, userID, req.FileURL, "", 0)
	case "audio":
		err = h.wsService.SendAudioMessage(roomID, userID, req.FileURL, 0)
	default:
		err = h.wsService.SendGroupMessage(roomID, userID, req.Content)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// MarkMessageAsRead godoc
// @Summary Mark a message as read
// @Description Marks a specific message as read by the authenticated user
// @Tags chat
// @Param roomId path string true "Room ID"
// @Param messageId path string true "Message ID"
// @Success 200 "Message marked as read"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /chat/rooms/{roomId}/messages/{messageId}/read [post]
func (h *ChatHandler) MarkMessageAsRead(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	roomID := chi.URLParam(r, "roomId")
	messageID := chi.URLParam(r, "messageId")

	if err := h.wsService.MarkMessageAsRead(roomID, userID, messageID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// PinMessage godoc
// @Summary Pin a message in a chat room
// @Description Pins a specific message in a chat room
// @Tags chat
// @Param roomId path string true "Room ID"
// @Param messageId path string true "Message ID"
// @Success 200 "Message pinned successfully"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /chat/rooms/{roomId}/messages/{messageId}/pin [post]
func (h *ChatHandler) PinMessage(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")
	messageID := chi.URLParam(r, "messageId")

	if err := h.wsService.PinMessage(roomID, messageID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UnpinMessage godoc
// @Summary Unpin a message in a chat room
// @Description Unpins a specific message in a chat room
// @Tags chat
// @Param roomId path string true "Room ID"
// @Param messageId path string true "Message ID"
// @Success 200 "Message unpinned successfully"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /chat/rooms/{roomId}/messages/{messageId}/unpin [post]
func (h *ChatHandler) UnpinMessage(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")
	messageID := chi.URLParam(r, "messageId")

	if err := h.wsService.UnpinMessage(roomID, messageID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ArchiveRoom godoc
// @Summary Archive a chat room
// @Description Archives a specific chat room for the authenticated user
// @Tags chat
// @Param roomId path string true "Room ID"
// @Success 200 "Room archived successfully"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /chat/rooms/{roomId}/archive [post]
func (h *ChatHandler) ArchiveRoom(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	roomID := chi.URLParam(r, "roomId")

	if err := h.wsService.ArchiveRoom(roomID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UnarchiveRoom godoc
// @Summary Unarchive a chat room
// @Description Unarchives a specific chat room for the authenticated user
// @Tags chat
// @Param roomId path string true "Room ID"
// @Success 200 "Room unarchived successfully"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /chat/rooms/{roomId}/unarchive [post]
func (h *ChatHandler) UnarchiveRoom(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	roomID := chi.URLParam(r, "roomId")

	if err := h.wsService.UnarchiveRoom(roomID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// MuteRoom godoc
// @Summary Mute a chat room
// @Description Mutes notifications for a specific chat room
// @Tags chat
// @Param roomId path string true "Room ID"
// @Success 200 "Room muted successfully"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /chat/rooms/{roomId}/mute [post]
func (h *ChatHandler) MuteRoom(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	roomID := chi.URLParam(r, "roomId")

	if err := h.wsService.MuteRoom(roomID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UnmuteRoom godoc
// @Summary Unmute a chat room
// @Description Unmutes notifications for a specific chat room
// @Tags chat
// @Param roomId path string true "Room ID"
// @Success 200 "Room unmuted successfully"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /chat/rooms/{roomId}/unmute [post]
func (h *ChatHandler) UnmuteRoom(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	roomID := chi.URLParam(r, "roomId")

	if err := h.wsService.UnmuteRoom(roomID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
