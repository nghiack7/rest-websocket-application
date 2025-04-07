package dtos

// CreateDirectRoomRequest represents the request body for creating a direct chat room
type CreateDirectRoomRequest struct {
	UserID2 string `json:"user_id_2" example:"user-123"`
}

// CreateGroupRoomRequest represents the request body for creating a group chat room
type CreateGroupRoomRequest struct {
	Name    string   `json:"name" example:"Team Chat"`
	UserIDs []string `json:"user_ids" example:"[\"user-123\", \"user-456\"]"`
}

// UpdateRoomRequest represents the request body for updating a chat room
type UpdateRoomRequest struct {
	Name        string `json:"name,omitempty" example:"New Room Name"`
	Description string `json:"description,omitempty" example:"Updated room description"`
	AvatarURL   string `json:"avatar_url,omitempty" example:"https://example.com/avatar.jpg"`
}

// SendMessageRequest represents the request body for sending a message
type SendMessageRequest struct {
	Content string `json:"content" example:"Hello, world!"`
	Type    string `json:"type,omitempty" example:"text" enums:"text,file,image,video,audio"`
	FileURL string `json:"file_url,omitempty" example:"https://example.com/file.pdf"`
}
