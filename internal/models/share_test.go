package models

import (
	"database/sql"
	"testing"
	"time"
)

func TestShare_IsValid(t *testing.T) {
	future := time.Now().Add(1 * time.Hour)
	past := time.Now().Add(-1 * time.Hour)

	tests := []struct {
		name     string
		share    Share
		expected bool
	}{
		{
			name: "valid share",
			share: Share{
				ExpiresAt: future,
			},
			expected: true,
		},
		{
			name: "expired share",
			share: Share{
				ExpiresAt: past,
			},
			expected: false,
		},
		{
			name: "revoked share",
			share: Share{
				ExpiresAt: future,
				RevokedAt: sql.NullTime{Time: time.Now(), Valid: true},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.share.IsValid(); got != tt.expected {
				t.Errorf("IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestShare_CanStartNewPlay(t *testing.T) {
	future := time.Now().Add(1 * time.Hour)
	past := time.Now().Add(-1 * time.Hour)

	tests := []struct {
		name     string
		share    Share
		expected bool
	}{
		{
			name: "no limits",
			share: Share{
				ExpiresAt: future,
			},
			expected: true,
		},
		{
			name: "expired",
			share: Share{
				ExpiresAt: past,
			},
			expected: false,
		},
		{
			name: "max total plays reached",
			share: Share{
				ExpiresAt:   future,
				TotalPlays:  10,
				MaxTotalPlays: sql.NullInt64{Int64: 10, Valid: true},
			},
			expected: false,
		},
		{
			name: "max total plays not reached",
			share: Share{
				ExpiresAt:   future,
				TotalPlays:  9,
				MaxTotalPlays: sql.NullInt64{Int64: 10, Valid: true},
			},
			expected: true,
		},
		{
			name: "max concurrent reached",
			share: Share{
				ExpiresAt:                future,
				CurrentConcurrentViewers: 5,
				MaxConcurrentViewers:     sql.NullInt64{Int64: 5, Valid: true},
			},
			expected: false,
		},
		{
			name: "max concurrent not reached",
			share: Share{
				ExpiresAt:                future,
				CurrentConcurrentViewers: 4,
				MaxConcurrentViewers:     sql.NullInt64{Int64: 5, Valid: true},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.share.CanStartNewPlay(); got != tt.expected {
				t.Errorf("CanStartNewPlay() = %v, want %v", got, tt.expected)
			}
		})
	}
}
