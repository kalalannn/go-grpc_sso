package validation_test

import (
	"grpc-sso/internal/validation"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmailValid(t *testing.T) {
	assert.True(t, validation.IsEmailValid("aa@bb.cc"))
	assert.True(t, validation.IsEmailValid("aa@bb-ccccccc.dddd"))
	assert.False(t, validation.IsEmailValid("aa@bb-ccccccc.ddddd"))
	assert.False(t, validation.IsEmailValid("aa@bb-cc.dd-c"))
	assert.False(t, validation.IsEmailValid("aa@bb/fc.dd"))
	assert.False(t, validation.IsEmailValid("aabb"))
	assert.False(t, validation.IsEmailValid("aabb@cc"))
	assert.False(t, validation.IsEmailValid("01234567890123456789012345678901234567890123456789012345678@cc.co"))
}

func TestIsPasswordValid(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{name: "Valid password",
			password: "Password123!",
			expected: true,
		}, {name: "No uppercase letter",
			password: "password123!",
			expected: false,
		}, {name: "No lowercase letter",
			password: "PASSWORD123!",
			expected: false,
		}, {name: "No special character",
			password: "Password123",
			expected: false,
		}, {name: "Too short",
			password: "Pass!",
			expected: false,
		}, {name: "Not Exceeds maximum length",
			password: "ThisIsAVeryLongPasswordThatExceedsTheMaximumAllowedLengthOf72Characters!",
			expected: true,
		}, {name: "Exceeds maximum length",
			password: "ThisIsAVeryLongPasswordThatExceedsTheMaximumAllowedLengthOf72Characters!!",
			expected: false,
		}, {name: "Empty password",
			password: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, validation.IsPasswordValid(tt.password))
		})
	}
}

func TestIsUsernameValid(t *testing.T) {
	tests := []struct {
		name     string
		username string
		expected bool
	}{
		{name: "Valid username with letters and numbers",
			username: "User123",
			expected: true,
		}, {name: "Valid username with letters, numbers, and special characters",
			username: "User_123",
			expected: true,
		}, {name: "Valid username with hyphen",
			username: "User-123",
			expected: true,
		}, {name: "Username starts with a number",
			username: "123User",
			expected: true,
		}, {name: "Too short username",
			username: "Usr",
			expected: false,
		}, {name: "Too long username",
			username: "ThisUsernameIsWayTooLong1234567890",
			expected: false,
		}, {name: "Username starts with special character",
			username: "_User123",
			expected: false,
		}, {name: "Username with invalid special character",
			username: "user@123",
			expected: false,
		}, {name: "Empty username",
			username: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, validation.IsUsernameValid(tt.username))
		})
	}
}
