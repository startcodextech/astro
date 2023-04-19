package rest

type (
	CreateClientRequest struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password"`
	}
)
