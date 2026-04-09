package dto

import (
	"time"
)

type NomineeRequest struct {
	FirstName    string    `json:"first_name" binding:"required"`
	LastName     string    `json:"last_name" binding:"required"`
	Mobile       string    `json:"mobile" binding:"required"`
	Email        string    `json:"email" binding:"required,email"`
	DOB          time.Time `json:"dob" binding:"required"` // should be date only
	GuardianName string    `json:"guardian_name,omitempty"`
}
