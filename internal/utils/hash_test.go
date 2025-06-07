package utils

import "testing"

func TestHashPassword(t *testing.T) {
	password := "secure_password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if hash == "" {
		t.Fatal("HashPassword() returned an empty hash")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "secure_password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if !CheckPasswordHash(password, hash) {
		t.Fatal("CheckPasswordHash() failed for the correct password")
	}

	incorrectPassword := "wrong_password"
	if CheckPasswordHash(incorrectPassword, hash) {
		t.Fatal("CheckPasswordHash() succeeded for an incorrect password")
	}
}
