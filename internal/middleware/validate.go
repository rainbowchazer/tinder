package middleware

import (
	"encoding/json"
	"net/http"
	"regexp"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Gender   string `json:"gender"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func ValidateRegister(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Некорректные данные", http.StatusBadRequest)
			return
		}

		if !isValidEmail(req.Email) {
			http.Error(w, "Некорректный email", http.StatusBadRequest)
			return
		}
		if len(req.Password) < 6 {
			http.Error(w, "Пароль должен содержать минимум 6 символов", http.StatusBadRequest)
			return
		}
		if req.Age < 18 {
			http.Error(w, "Возраст должен быть не меньше 18 лет", http.StatusBadRequest)
			return
		}

		// Возвращаем тело запроса для следующего обработчика
		r.Body.Close()
		r.Body = reencodeRequest(req)

		next.ServeHTTP(w, r)
	})
}

func ValidateLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Некорректные данные", http.StatusBadRequest)
			return
		}

		if !isValidEmail(req.Email) {
			http.Error(w, "Некорректный email", http.StatusBadRequest)
			return
		}
		if len(req.Password) < 6 {
			http.Error(w, "Пароль должен содержать минимум 6 символов", http.StatusBadRequest)
			return
		}

		r.Body.Close()
		r.Body = reencodeRequest(req)

		next.ServeHTTP(w, r)
	})
}

func isValidEmail(email string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return regex.MatchString(email)
}

func reencodeRequest(body interface{}) *FakeReadCloser {
	data, _ := json.Marshal(body)
	return &FakeReadCloser{data: data}
}

type FakeReadCloser struct {
	data []byte
	pos  int
}

func (r *FakeReadCloser) Read(p []byte) (int, error) {
	if r.pos >= len(r.data) {
		return 0, http.ErrBodyReadAfterClose
	}
	n := copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

func (r *FakeReadCloser) Close() error {
	return nil
}
