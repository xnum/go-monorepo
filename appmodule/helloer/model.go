package helloer

import (
	"go-monorepo/models"
)

// RequestLog defines database model of one request log.
type RequestLog struct {
	models.UUIDModel
	Type     string `gorm:"not null"`
	Request  string `gorm:"not null"`
	Response string `gorm:"not null"`
}
