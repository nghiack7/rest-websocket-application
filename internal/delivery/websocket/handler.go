package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/personal/task-management/internal/usecase"
	"github.com/personal/task-management/pkg/utils/jwt"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, implement proper origin checking
	},
}

type Handler struct {
	wsService  usecase.WebSocketService
	jwtService jwt.JWTTokenServicer
}

func NewHandler(wsService usecase.WebSocketService, jwtService jwt.JWTTokenServicer) *Handler {
	return &Handler{
		wsService:  wsService,
		jwtService: jwtService,
	}
}

func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {

	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "missing token", http.StatusBadRequest)
		return
	}
	// decode token
	claims, err := h.jwtService.ValidateToken(token)
	if err != nil {
		http.Error(w, "invalid token", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "could not upgrade connection", http.StatusInternalServerError)
		return
	}

	h.wsService.HandleConnection(conn, claims.UserID.String())
}
