package services

import (
	"context"
	"core/internal/repository"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"unicode"
)

// GenerateUsername генерирует username на основе ФИО и названия группы
func generateUsername(fullName, groupName string) string {
	parts := strings.Fields(fullName)
	if len(parts) == 0 {
		return "user"
	}
	lastName := parts[0]
	firstName := ""
	if len(parts) > 1 {
		firstName = parts[1]
	}

	// Транслитерация простейшая (можно расширить, но для демо так)
	translit := func(s string) string {
		// Простейшая замена русских букв на латиницу (только для демо)
		// В реальном проекте лучше использовать готовую библиотеку
		replacer := strings.NewReplacer(
			"а", "a", "б", "b", "в", "v", "г", "g", "д", "d", "е", "e", "ё", "e",
			"ж", "zh", "з", "z", "и", "i", "й", "y", "к", "k", "л", "l", "м", "m",
			"н", "n", "о", "o", "п", "p", "р", "r", "с", "s", "т", "t", "у", "u",
			"ф", "f", "х", "kh", "ц", "ts", "ч", "ch", "ш", "sh", "щ", "shch", "ъ", "",
			"ы", "y", "ь", "", "э", "e", "ю", "yu", "я", "ya",
		)
		lower := strings.ToLower(s)
		return replacer.Replace(lower)
	}
	lastNameLat := translit(lastName)
	firsNameLat := translit(firstName)

	clean := func(s string) string {
		var res strings.Builder
		for _, r := range s {
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				res.WriteRune(r)
			}
		}
		return res.String()
	}
	digits := func(s string) string {
		var res strings.Builder
		for _, r := range s {
			if unicode.IsDigit(r) {
				res.WriteRune(r)
			}
		}
		return res.String()
	}
	username := ""
	if groupName != "" {
		username += digits(groupName)
		username += "-"
	}
	username += clean(lastNameLat)
	if lastNameLat != "" {
		username += "-" + clean(firsNameLat)
	}
	if len(username) > 30 {
		username = username[:30]
	}
	return username
}

// generateRandomPassword генерирует случайный пароль длиной 10 символов
func generateRandomPassword() (string, error) {
	bytes := make([]byte, 5) // 5 байт -> 10 hex символов
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// ensureUniqueUsername пытается сделать username уникальным, добавляя числовой суффикс
func ensureUniqueUsername(ctx context.Context, userRepo *repository.UserRepo, baseUsername string) (string, error) {
	username := baseUsername
	counter := 1
	for {
		existing, err := userRepo.GetByUsername(ctx, username)
		if err != nil {
			return "", err
		}
		if existing == nil {
			break
		}
		username = fmt.Sprintf("%s%d", baseUsername, counter)
		counter++
	}
	return username, nil
}
