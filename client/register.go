package client

// RegisterRequest struct
type RegisterRequest struct {
	Name        string   `json:"name"`
	Path        string   `json:"path"`
	Endpoints   []string `json:"end_points"`
	Host        string   `json:"host"`
	Description string   `json:"description"`
}

// RegisterResponse struct
type RegisterResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// NewRegisterResponse returns a new registration response
func NewRegisterResponse(id, name string) *RegisterResponse {
	return &RegisterResponse{
		ID:   id,
		Name: name,
	}
}
