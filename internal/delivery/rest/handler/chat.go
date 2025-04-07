package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/personal/task-management/internal/usecase"
	"github.com/personal/task-management/pkg/utils/jwt"
)

type ChatHandler struct {
	wsService usecase.WebSocketService

	jwtService jwt.JWTTokenServicer
}

func NewChatHandler(wsService usecase.WebSocketService, jwtService jwt.JWTTokenServicer) *ChatHandler {
	return &ChatHandler{
		wsService:  wsService,
		jwtService: jwtService,
	}
}

type CreateDirectRoomRequest struct {
	UserID2 string `json:"user_id_2"`
}

type CreateGroupRoomRequest struct {
	Name    string   `json:"name"`
	UserIDs []string `json:"user_ids"`
}

type UpdateRoomRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	AvatarURL   string `json:"avatar_url,omitempty"`
}

type SendMessageRequest struct {
	Content string `json:"content"`
	Type    string `json:"type,omitempty"`
	FileURL string `json:"file_url,omitempty"`
}

func (h *ChatHandler) CreateDirectRoom(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req CreateDirectRoomRequest
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

func (h *ChatHandler) CreateGroupRoom(w http.ResponseWriter, r *http.Request) {
	var req CreateGroupRoomRequest
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

func (h *ChatHandler) ListRooms(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	rooms, err := h.wsService.ListRooms(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(rooms)
}

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

func (h *ChatHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	roomID := chi.URLParam(r, "roomId")

	if err := h.wsService.JoinRoom(roomID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ChatHandler) LeaveRoom(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	roomID := chi.URLParam(r, "roomId")

	if err := h.wsService.LeaveRoom(roomID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ChatHandler) UpdateRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")

	var req UpdateRoomRequest
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

func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	roomID := chi.URLParam(r, "roomId")

	var req SendMessageRequest
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

func (h *ChatHandler) PinMessage(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")
	messageID := chi.URLParam(r, "messageId")

	if err := h.wsService.PinMessage(roomID, messageID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ChatHandler) UnpinMessage(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")
	messageID := chi.URLParam(r, "messageId")

	if err := h.wsService.UnpinMessage(roomID, messageID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ChatHandler) ArchiveRoom(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	roomID := chi.URLParam(r, "roomId")

	if err := h.wsService.ArchiveRoom(roomID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ChatHandler) UnarchiveRoom(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	roomID := chi.URLParam(r, "roomId")

	if err := h.wsService.UnarchiveRoom(roomID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ChatHandler) MuteRoom(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	roomID := chi.URLParam(r, "roomId")

	if err := h.wsService.MuteRoom(roomID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ChatHandler) UnmuteRoom(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	roomID := chi.URLParam(r, "roomId")

	if err := h.wsService.UnmuteRoom(roomID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
