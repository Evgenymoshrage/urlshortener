package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"sync"
)

type URLShortener struct { // потокобезопасное хранение URL
	urls map[string]string
	mu   sync.RWMutex
}

func NewURLShortener() *URLShortener {
	return &URLShortener{
		urls: make(map[string]string),
	}
}

// создает короткий идентификатор для URL
func (us *URLShortener) Shorten(originalURL string) (string, error) { // принимаем оригинальный url, возвращаем короткий идентификатор
	if !isValidURL(originalURL) { // Валидируем URL (проверка на корректность url)
		return "", errors.New("invalid URL") // Если URL некорректный, возвращаем пустой короткий ID и ошибку
	}

	var shortID string // Переменная, в которую будет записан генерируемый короткий идентификатор
	for {              // Генерация уникального shortID
		shortID = generateShortID() // Генерируем случайный ID, проверяем уникальность
		us.mu.RLock()
		_, exists := us.urls[shortID] // Читаем карту безопасно
		us.mu.RUnlock()
		if !exists { // Проверяем, нет ли уже такого ID в us.urls
			break // если короткий ID уникален, выходим из цикла.
		}
	}

	us.mu.Lock()
	us.urls[shortID] = originalURL // Сохраняем короткий ID в карту безопасно
	us.mu.Unlock()
	return shortID, nil // Возвращаем результат
}

// возвращает оригинальный URL по короткому ID
func (us *URLShortener) GetOriginal(shortID string) (string, error) {
	us.mu.RLock() // Позволяет многим горутинам читать одновременно, но запрещает запись до RUnlock
	defer us.mu.RUnlock()
	original, ok := us.urls[shortID] // Читаем карту (original = значение по ключу shortID, ok = true, если ключ есть в map)
	if !ok {                         // Проверка на наличие ключа
		return "", fmt.Errorf("short URL not found")
	}
	return original, nil // Если ключ найден, возвращаем оригинальный URL и nil ошибки
}

// функция, которая возвращает случайную строку (6 символов)
func generateShortID() string {
	b := make([]byte, 4)                           // 4 байта = 6 символов base64 URL-safe
	rand.Read(b)                                   // Заполняем срез случайными байтами
	return base64.RawURLEncoding.EncodeToString(b) // Преобразуем байты в строку base64
}

// проверяет корректность URL
func isValidURL(str string) bool {
	u, err := url.ParseRequestURI(str) // Проверяет, что строка — это корректный URL
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" { // Проверка схемы
		return false
	}
	return true
}
