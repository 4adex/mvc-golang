package jwtutils

import "testing"

func TestJWT(t *testing.T) {
	username := "hawkeye"
	mail := "hawkeye@gmail.com"
	role := "admin"
	id := "1"
	token, _ := GenerateJWT(username,mail,role,id)
	claim, _ := ValidateJWT(token)
	if (username!=claim.Username) || (mail!=claim.Email) || (role!=claim.Role) || (id!=claim.Id) {
		t.Error("mismatched")
	}
}

func TestGenerateJWT(t *testing.T) {
	username := "hawkeye"
	email := "hawkeye@gmail.com"
	role := "admin"
	id := "1"

	token, err := GenerateJWT(username, email, role, id)
	if err != nil {
		t.Errorf("Failed to generate token: %v", err)
	}
	if token == "" {
		t.Error("Generated token is empty")
	}
}

func TestValidateJWT(t *testing.T) {
	username := "hawkeye"
	email := "hawkeye@gmail.com"
	role := "admin"
	id := "1"

	token, _ := GenerateJWT(username, email, role, id)
	claims, err := ValidateJWT(token)

	if err != nil {
		t.Errorf("Failed to validate token: %v", err)
	}

	if claims.Username != username {
		t.Errorf("Expected username %s, got %s", username, claims.Username)
	}
	if claims.Email != email {
		t.Errorf("Expected email %s, got %s", email, claims.Email)
	}
	if claims.Role != role {
		t.Errorf("Expected role %s, got %s", role, claims.Role)
	}
	if claims.Id != id {
		t.Errorf("Expected id %s, got %s", id, claims.Id)
	}
}

func TestInvalidToken(t *testing.T) {
	invalidToken := "invalid.token.string"
	_, err := ValidateJWT(invalidToken)
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
}