package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShortenHandler(t *testing.T) {
	shortener := NewURLShortener() // Создаём экземпляр бизнес-логики, который хранит карту shortID -> originalURL

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // Создаём HTTP-обработчик
		if r.Method != http.MethodPost { // Проверка метода
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed) // Только POST разрешён
			return                                                           // Если метод другой → возвращаем 405
		}
		var req struct { // cоздаём временную структуру для запроса
			URL string `json:"url"` // JSON должен быть вида: {"url":"http://example.com"}
		}
		json.NewDecoder(r.Body).Decode(&req)       // Decode превращает тело запроса в Go-структуру req
		shortID, err := shortener.Shorten(req.URL) // Создание короткого URL
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"short_url": shortID, "original_url": req.URL}) // Возврат ответа JSON
	})

	tests := []struct { // Таблица тестов
		name     string // название теста
		inputURL string // URL для запроса
		wantErr  bool   // ожидаем наличие ошибки
	}{
		{"valid URL", "http://example.com", false},
		{"invalid URL", "not-a-url", true},
	}

	for _, tt := range tests { // Цикл по кейсам
		t.Run(tt.name, func(t *testing.T) { // Создаём под-тест, чтобы каждый кейс запускался отдельно и видно в отчёте
			body, _ := json.Marshal(map[string]string{"url": tt.inputURL})                 // Создаём тело запроса, кодируем входной URL в JSON
			req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(body)) // Создаём тестовый HTTP-запрос
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req) // Вызов обработчика

			if tt.wantErr { // Проверка результата
				if rec.Code == http.StatusOK {
					t.Errorf("ожидали ошибку, получили 200")
				}
			} else {
				if rec.Code != http.StatusOK {
					t.Errorf("ожидали 200, получили %d", rec.Code)
				}
			}
		})
	}
}
