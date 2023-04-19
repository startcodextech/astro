package types

import "time"

type (
	User struct {
		ID        string    `json:"id"`
		Username  string    `json:"username"`
		Password  string    `json:"password"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Email     string    `json:"email"`
		Enable    bool      `json:"enable"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)
