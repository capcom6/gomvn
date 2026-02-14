package entity

import (
	"time"
)

type Path struct {
	UserID    uint      `gorm:"primary_key;auto_increment:false;not null;"`
	Path      string    `gorm:"primary_key;size:255;not null"`
	Deploy    bool      `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func NewPath(userID uint, path string, deploy bool) Path {
	now := time.Now()

	return Path{
		UserID:    userID,
		Path:      path,
		Deploy:    deploy,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewPathID(userID uint, path string) Path {
	//nolint:exhaustruct // partial constructor
	return Path{
		UserID: userID,
		Path:   path,
	}
}
