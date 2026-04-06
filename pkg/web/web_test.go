package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

/**
 * ===========================================================================
 * Web and HTTP API Tests
 * ===========================================================================
 *
 * For a Java developer:
 * - `net/http/httptest` is similar to `MockMvc` or `RestAssured`.
 * - `httptest.NewRecorder()` acts as a response buffer (like a mock response).
 * - `handler.ServeHTTP(rr, req)` lets you test handlers without starting a real server.
 */

func TestTasksAPI(t *testing.T) {
	handler := NewWebServerHandler()

	t.Run("GET /tasks JSON", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/tasks", nil)
		req.Header.Set("Accept", "application/json")
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rr.Code)
		}
		var tasks []Task
		if err := json.NewDecoder(rr.Body).Decode(&tasks); err != nil {
			t.Fatal(err)
		}
		if len(tasks) < 1 {
			t.Error("expected at least one task")
		}
	})

	t.Run("GET /tasks HTML", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/tasks", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rr.Code)
		}
		if !strings.Contains(rr.Body.String(), "<h1>Task List</h1>") {
			t.Error("expected HTML content")
		}
	})

	t.Run("GET /tasks/1", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/tasks/1", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rr.Code)
		}
		var task Task
		if err := json.NewDecoder(rr.Body).Decode(&task); err != nil {
			t.Fatal(err)
		}
		if task.ID != 1 {
			t.Errorf("expected ID 1, got %d", task.ID)
		}
	})

	t.Run("POST /tasks", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/tasks", strings.NewReader("title=New+Test+Task"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusSeeOther {
			t.Errorf("expected 303, got %d", rr.Code)
		}

		// Verify it was added
		tasks := store.List()
		found := false
		for _, task := range tasks {
			if task.Title == "New Test Task" {
				found = true
				break
			}
		}
		if !found {
			t.Error("new task not found in store")
		}
	})
}
