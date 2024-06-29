package main

import "testing"

func TestCheckProfane(t *testing.T) {
	checks := []struct {
		inputMessage    string
		expectedMessage string
	}{
		{"Hi", "Hi"},
		{"kerfuffle", "****"},
		{"sharbert", "****"},
		{"fornax", "****"},
		{"This is a kerfuffle opinion I need to sharbert with the world", "This is a **** opinion I need to **** with the world"},
		{"This is a kerfuffle! opinion I need to fornax with the world", "This is a kerfuffle! opinion I need to **** with the world"},
	}
	for _, testCheck := range checks {
		receivedMessage := checkProfane(testCheck.inputMessage)
		if receivedMessage != testCheck.expectedMessage {
			t.Errorf("received message '%v' does not match the expected message '%v'", receivedMessage, testCheck.expectedMessage)
		}
	}
}
