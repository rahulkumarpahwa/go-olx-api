package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argonTime    = uint32(1) // no of rounds
	argonMemory  = uint32(64 * 1024) // 64 MiB
	argonThreads = uint8(4)
	argonKeyLen  = uint32(32) // password length
)

func HashPassword(password string) (string, error) {
	// Generate random salt
	salt := make([]byte, 16)

	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Hash password
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		argonTime,
		argonMemory,
		argonThreads,
		argonKeyLen,
	)

	// Store parameters + salt + hash
	return fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		argonMemory,
		argonTime,
		argonThreads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	), nil
}

func CheckPassword(password, encodedHash string) bool {
	// Split:
	// $argon2id$v=19$m=65536,t=3,p=4$SALT$HASH
	parts := strings.Split(encodedHash, "$")

	if len(parts) != 6 {
		return false
	}

	// parts[0] = ""
	// parts[1] = "argon2id"
	// parts[2] = "v=19"
	// parts[3] = "m=65536,t=3,p=4"
	// parts[4] = salt
	// parts[5] = hash

	if parts[1] != "argon2id" {
		return false
	}

	// Parse parameters
	var memory uint32
	var time uint32
	var threads uint8

	params := strings.Split(parts[3], ",")

	for _, param := range params {
		kv := strings.SplitN(param, "=", 2)

		if len(kv) != 2 {
			return false
		}

		switch kv[0] {
		case "m":
			value, err := strconv.ParseUint(kv[1], 10, 32)
			if err != nil {
				return false
			}
			memory = uint32(value)

		case "t":
			value, err := strconv.ParseUint(kv[1], 10, 32)
			if err != nil {
				return false
			}
			time = uint32(value)

		case "p":
			value, err := strconv.ParseUint(kv[1], 10, 8)
			if err != nil {
				return false
			}
			threads = uint8(value)
		}
	}

	// Decode salt
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}

	// Decode stored hash
	storedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}

	// Hash the password entered during login
	newHash := argon2.IDKey(
		[]byte(password),
		salt,
		time,
		memory,
		threads,
		uint32(len(storedHash)),
	)

	// Constant-time comparison
	return subtle.ConstantTimeCompare(newHash, storedHash) == 1
}
