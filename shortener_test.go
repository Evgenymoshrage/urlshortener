package main

import "testing"

func TestURLShortener_Shorten(t *testing.T) {
	tests := []struct { //  Объявляем таблицу тестов
		name    string // название кейс
		url     string // входной URL
		wantErr bool   // ожидаемое наличие ошибки
	}{
		{"Валидный HTTP URL", "http://example.com", false},
		{"Валидный HTTPS URL", "https://google.com/search?q=test", false},
		{"Невалидный URL", "not-a-url", true},
		{"Пустая строка", "", true},
	}

	shortener := NewURLShortener() // Создаём экземпляр URL-shortener, который будет использоваться для всех тестов
	for _, tt := range tests {     // Проходим по всем тест-кейсам циклом
		t.Run(tt.name, func(t *testing.T) { // создаём под-тест (позволяет запускать каждый кейс отдельно, go test -v видно название кейса и результат)
			shortID, err := shortener.Shorten(tt.url) // Вызываем Shorten для текущего URL и проверяем ошибку
			if (err != nil) != tt.wantErr {           // есть ли фактическая ошибка и равна ли она ожидаемой
				t.Errorf("Ошибка = %v, ожидали ошибку = %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(shortID) < 6 { // Проверка длины shortID
				t.Errorf("ShortID слишком короткий: %s", shortID)
			}
		})
	}
}

func TestURLShortener_GetOriginal(t *testing.T) { // Тест для GetOriginal
	shortener := NewURLShortener() // создаем тестируемый объект

	shortID, err := shortener.Shorten("http://example.com")
	if err != nil {
		t.Fatalf("Не удалось создать shortID: %v", err)
	}

	tests := []struct { // Таблица тестов
		name    string // имя теста
		id      string // shortID
		wantErr bool   // ожидаем ошибку
	}{
		{"existing shortID", shortID, false},      // Первый кейс - мы передаём существующий shortID, ошибки нет
		{"non-existing shortID", "unknown", true}, // Второй кейс - мы передаём НЕсуществующий shortID, ошибка есть
	}

	for _, tt := range tests { // Цикл по тестам
		t.Run(tt.name, func(t *testing.T) {
			_, err := shortener.GetOriginal(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ошибка = %v, ожидали ошибку = %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsValidURL(t *testing.T) {
	tests := []struct { // Таблица тестов
		url  string // входное значение
		want bool   // входное значение
	}{ // Конкретные тест-кейсы
		{"http://example.com", true},
		{"https://example.com", true},
		{"ftp://example.com", false},
		{"invalid", false},
	}

	for _, tt := range tests { // Цикл по тестам
		if got := isValidURL(tt.url); got != tt.want {
			t.Errorf("isValidURL(%q) = %v, want %v", tt.url, got, tt.want)
		}
	}
}
