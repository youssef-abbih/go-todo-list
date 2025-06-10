package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/youssef-abbih/go-todo-list/models"
)

// contextKey type for context keys
type contextKey string

// setParam adds a key/value param to request context (to mock URL params)
func setParam(ctx context.Context, key, value string) context.Context {
	return context.WithValue(ctx, contextKey(key), value)
}

// Test DefaultResponse handler
func TestDefaultResponse(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	DefaultResponse(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	expected := "Welcome to my Todo List API"
	if strings.TrimSpace(string(body)) != expected {
		t.Errorf("expected body %q, got %q", expected, string(body))
	}
}

func setup() {
	models.InitDB()
}

// Test POST /tasks
func TestPostTask(t *testing.T) {
	setup()

	validTask := models.Task{Title: "Test", Description: "Test desc", Completed: false}
	body, _ := json.Marshal(validTask)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	PostTask(res, req)

	if res.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d", res.Code)
	}

	malformedReq := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader("invalid json"))
	malformedReq.Header.Set("Content-Type", "application/json")
	malformedRes := httptest.NewRecorder()
	PostTask(malformedRes, malformedReq)
	if malformedRes.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request for malformed JSON, got %d", malformedRes.Code)
	}
}

// Test GET /tasks
func TestGetTasks(t *testing.T) {
	setup()
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	res := httptest.NewRecorder()
	GetTasks(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", res.Code)
	}
}

// Test GET /tasks/{id}
func TestGetTask(t *testing.T) {
	setup()

	req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	req = req.WithContext(setParam(req.Context(), "id", "1"))
	res := httptest.NewRecorder()
	GetTask(res, req)
	if res.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", res.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/tasks/9999", nil)
	req = req.WithContext(setParam(req.Context(), "id", "9999"))
	res = httptest.NewRecorder()
	GetTask(res, req)
	if res.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found, got %d", res.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/tasks/abc", nil)
	req = req.WithContext(setParam(req.Context(), "id", "abc"))
	res = httptest.NewRecorder()
	GetTask(res, req)
	if res.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", res.Code)
	}
}

// Test PUT /tasks/{id}
func TestPutTask(t *testing.T) {
	setup()

	updated := models.Task{Title: "Updated", Description: "Updated desc", Completed: true}
	body, _ := json.Marshal(updated)
	req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewReader(body))
	req = req.WithContext(setParam(req.Context(), "id", "1"))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	PutTask(res, req)
	if res.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", res.Code)
	}

	nonexistent := httptest.NewRequest(http.MethodPut, "/tasks/9999", bytes.NewReader(body))
	nonexistent = nonexistent.WithContext(setParam(nonexistent.Context(), "id", "9999"))
	nonexistent.Header.Set("Content-Type", "application/json")
	res = httptest.NewRecorder()
	PutTask(res, nonexistent)
	if res.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found, got %d", res.Code)
	}

	malformed := httptest.NewRequest(http.MethodPut, "/tasks/1", strings.NewReader("bad json"))
	malformed = malformed.WithContext(setParam(malformed.Context(), "id", "1"))
	malformed.Header.Set("Content-Type", "application/json")
	res = httptest.NewRecorder()
	PutTask(res, malformed)
	if res.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", res.Code)
	}
}

// Test DELETE /tasks/{id}
func TestDeleteTask(t *testing.T) {
	setup()

	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	req = req.WithContext(setParam(req.Context(), "id", "1"))
	res := httptest.NewRecorder()
	DeleteTask(res, req)
	if res.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", res.Code)
	}

	nonexistent := httptest.NewRequest(http.MethodDelete, "/tasks/9999", nil)
	nonexistent = nonexistent.WithContext(setParam(nonexistent.Context(), "id", "9999"))
	res = httptest.NewRecorder()
	DeleteTask(res, nonexistent)
	if res.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found, got %d", res.Code)
	}

	invalid := httptest.NewRequest(http.MethodDelete, "/tasks/abc", nil)
	invalid = invalid.WithContext(setParam(invalid.Context(), "id", "abc"))
	res = httptest.NewRecorder()
	DeleteTask(res, invalid)
	if res.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", res.Code)
	}
}
