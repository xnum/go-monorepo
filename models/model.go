package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Model defines model. We need lower case ID.
type Model struct {
	ID        uint           `gorm:"primary_key" json:"id"`
	CreatedAt time.Time      `                   json:"created_at"`
	UpdatedAt time.Time      `                   json:"updated_at"`
	DeletedAt gorm.DeletedAt `                   json:"-"          sql:"index"`
}

// UUIDModel defines model with uuid as primary key.
type UUIDModel struct {
	ID        uuid.UUID      `gorm:"primary_key;type:uuid;not null;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time      `                                                                 json:"created_at"`
	UpdatedAt time.Time      `                                                                 json:"updated_at"`
	DeletedAt gorm.DeletedAt `                                                                 json:"-"          sql:"index"`
}

var models map[any]struct{}

func init() {
	models = make(map[any]struct{})
}

// RegisterModel saves the model to auto migrate to database.
// x must be a pointer.
func RegisterModel(x ...any) {
	for _, t := range x {
		models[t] = struct{}{}
	}
}

// Models returns registered models.
func Models() []any {
	arr := make([]any, 0, len(models))
	for k := range models {
		arr = append(arr, k)
	}

	return arr
}
