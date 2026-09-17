package users

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

type RequestUser struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}


type CreateUser struct {
	ID       uuid.UUID `json:"-"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

type ValidationError struct {
	Field string
	Msg   string
}

func (v *ValidationError) Error() string {
	if v.Field != "" {
		return fmt.Sprintf(`%s: %s`, v.Field, v.Msg)
	}
	return fmt.Sprintf(`%s`, v.Msg)
}

func (u *CreateUser) Validate() error {
	if !isValidEmail(u.Email) {
		return &ValidationError{Field: "email", Msg: "must be valid"}
	}

	if u.Password == "" {
		return &ValidationError{Field: "password", Msg: "must not be empty"}
	}

	if !isStrongPassword(u.Password) {
		return &ValidationError{Field: "", Msg: "invalid credentials"}
	}

	if u.Name == "" {
		return &ValidationError{Field: "name", Msg: "must not be empty"}
	}

	if len(u.Name) < 4 {
		return &ValidationError{Field: "name", Msg: "must be atleast 5 characters"}
	}

	if len(u.Name) > 30 {
		return &ValidationError{
			Field: "name",
			Msg:   "must not be greater than 30 characters",
		}
	}

	return nil
}

func isValidEmail(email string) bool {
	email = strings.TrimSpace(email)

	// Basic email regex
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`

	re := regexp.MustCompile(pattern)

	return re.MatchString(email)
}

func isStrongPassword(password string) bool {
	password = strings.TrimSpace(password)

	// Length: 8 to 30 characters
	if len(password) < 8 || len(password) > 30 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)

	return hasUpper && hasLower && hasNumber && hasSpecial
}
