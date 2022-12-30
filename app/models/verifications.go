package models

import "time"

type VerificationDataType int

const (
	MailConfirmation VerificationDataType = iota + 1
	PassReset
)

// VerificationData represents the type for the data stored for verification.
type VerificationData struct {
	Email     string               `json:"email" validate:"required" db:"email"`
	Code      string               `json:"code" validate:"required" db:"code"`
	ExpiresAt time.Time            `json:"expiresat" db:"expiresat"`
	Type      VerificationDataType `json:"type" db:"type"`
}
