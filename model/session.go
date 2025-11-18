package model

import (
	"time"

	"gorm.io/gorm"
)

// Session represents a user's refresh token session
// This enables token revocation and "logout all devices" functionality
type Session struct {
	gorm.Model
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	RefreshToken string    `gorm:"uniqueIndex;not null;size:512" json:"-"` // Never expose in JSON
	UserAgent    string    `gorm:"size:512" json:"user_agent"`
	IPAddress    string    `gorm:"size:45" json:"ip_address"` // IPv6 max length
	ExpiresAt    time.Time `gorm:"not null;index" json:"expires_at"`
	LastUsedAt   time.Time `gorm:"not null" json:"last_used_at"`
	IsRevoked    bool      `gorm:"default:false;index" json:"is_revoked"`

	// Relationships
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

// IsExpired checks if the session has expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsValid checks if the session is valid (not expired and not revoked)
func (s *Session) IsValid() bool {
	return !s.IsExpired() && !s.IsRevoked
}

// UpdateLastUsed updates the last used timestamp
func (s *Session) UpdateLastUsed(db *gorm.DB) error {
	s.LastUsedAt = time.Now()
	return db.Save(s).Error
}

// Revoke marks the session as revoked
func (s *Session) Revoke(db *gorm.DB) error {
	s.IsRevoked = true
	return db.Save(s).Error
}
