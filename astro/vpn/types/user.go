package types

import "time"

type User struct {
	Username     string     `json:"username"`
	UsedPassword bool       `json:"used_password"`
	CreatedAt    time.Time  `json:"created_at"`
	ExpiredAt    time.Time  `json:"expired_at,omitempty,null"`
	RevokedAt    *time.Time `json:"revoked_at,omitempty,null"`
	Serial       string     `json:"serial,omitempty"`
	Status       string     `json:"status,omitempty"`
}
