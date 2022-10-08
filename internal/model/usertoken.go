package model

import uuid "github.com/jackc/pgtype/ext/gofrs-uuid"

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserToken struct {
	UserID uuid.UUID `gorm:"user_id"`
}
