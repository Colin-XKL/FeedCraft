package dao

import (
	"database/sql"
	"time"
)

type BaseModelWithoutPK struct {
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `gorm:"index" json:"deleted_at,omitempty"`
}
