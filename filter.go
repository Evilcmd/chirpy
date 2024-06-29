package main

import "strings"

func checkProfane(message string) string {
	splitMessage := strings.Split(message, " ")
	for i, msg := range splitMessage {
		if strings.ToLower(msg) == "kerfuffle" || strings.ToLower(msg) == "sharbert" || strings.ToLower(msg) == "fornax" {
			splitMessage[i] = "****"
		}
	}
	return strings.Join(splitMessage, " ")
}
