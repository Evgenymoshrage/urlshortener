package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	shortener := NewURLShortener() // Создаём экземпляр URLShortener

	http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) { // Обработчик POST
		if r.Method != http.MethodPost { // Проверяем, что метод запроса POST
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct { // Декодирование JSON запроса
			URL string `json:"url"` // тег структуры
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil { // читает тело запроса и декодирует JSON в структуру
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		shortID, err := shortener.Shorten(req.URL) // Создание короткого URL
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp := map[string]string{ // Формирование ответа JSON
			"short_url":    shortID,
			"original_url": req.URL,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { // Обработчик GET
		shortID := r.URL.Path[1:]                       // Берём часть URL после /, чтобы получить короткий идентификатор.
		original, err := shortener.GetOriginal(shortID) // ищем оригинальный URL
		if err != nil {                                 // Если короткого ID нет возвращаем 404
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, original, http.StatusFound) // Редирект на оригинальный URL
	})

	log.Println("Server running at :8080")
	http.ListenAndServe(":8080", nil) // Запуск сервера на порту 8080
}
