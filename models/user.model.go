package models

import (
	"github.com/google/uuid"
)

// Database models

type Segment struct {
	SegmentID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	SegmentSlug string    `gorm:"type:varchar(255);not null"`
}

type User struct {
	UserUUID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID   int64     `gorm:"type:int;not null"`
}

type UserSegment struct {
	UserUUID  uuid.UUID `gorm:"type:uuid;not null"`
	SegmentID uuid.UUID `gorm:"type:uuid;not null"`
}

// Payload struct

type UserPayload struct {
	UserID         int64    `json:"userId"`
	SegmentsAdd    []string `json:"segmentsAdd"`
	SegmentsDelete []string `json:"segmentsDelete"`
}

type UserResponse struct {
	UserID   int64    `json:"userID"`
	Segments []string `json:"segments"`
}

type SegmentPayload struct {
	SegmentID      uuid.UUID `json:"segmentID"`
	SegmentSlug    string    `json:"segmentSlug"`
	AutoAdd        bool      `json:"autoAdd"`
	UserPercentage float64   `json:"userPercentage"`
}
