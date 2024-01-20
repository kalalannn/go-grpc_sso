package validation

import "regexp"

func IsEmailValid(email string) bool {
	if len(email) > 64 {
		return false
	}
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

// 8 <= password <= 72
// One Uppercase
// One lowercase
// One specsymbol
func IsPasswordValid(password string) bool {
	if len(password) < 8 || len(password) > 72 {
		return false
	}

	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)

	return hasLower && hasUpper && hasSpecial
}

// 4 <= username <= 32
// first [a-zA-Z0-9]
func IsUsernameValid(username string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-_]{3,31}$`).MatchString(username)
}
