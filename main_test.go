package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupTestUsers(t *testing.T) {
	t.Helper()
	initial := map[string]User{
		"595ea76c-b130-4c3b-9afa-454a17d90d60": {
			ID:   "595ea76c-b130-4c3b-9afa-454a17d90d60",
			Name: "Jailson",
		},
		"765a42a9-8783-4808-ae5a-af57cf07d9e0": {
			ID:   "765a42a9-8783-4808-ae5a-af57cf07d9e0",
			Name: "Mendes",
		},
		"e5c4a6e0-c68d-4389-bcf3-c3f908c27e5e": {
			ID:   "e5c4a6e0-c68d-4389-bcf3-c3f908c27e5e",
			Name: "Alice",
		},
	}

	users = make(map[string]User, len(initial))
	for k, v := range initial {
		users[k] = v
	}

	t.Cleanup(func() {
		users = make(map[string]User, len(initial))
		for k, v := range initial {
			users[k] = v
		}
	})
}

func TestServeHTTP_GET_Users(t *testing.T) {
	setupTestUsers(t)
	s := &server{}

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()

	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", contentType)
	}

	var gotUsers []User
	if err := json.NewDecoder(rec.Body).Decode(&gotUsers); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if len(gotUsers) != len(users) {
		t.Errorf("expected %d users, got %d", len(users), len(gotUsers))
	}
}

func TestServeHTTP_GET_UserByID(t *testing.T) {
	s := &server{}

	t.Run("existing user", func(t *testing.T) {
		setupTestUsers(t)
		targetID := "595ea76c-b130-4c3b-9afa-454a17d90d60"

		req := httptest.NewRequest(http.MethodGet, "/users/"+targetID, nil)
		rec := httptest.NewRecorder()

		s.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		contentType := rec.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("expected Content-Type application/json, got %q", contentType)
		}

		var gotUser User
		if err := json.NewDecoder(rec.Body).Decode(&gotUser); err != nil {
			t.Fatalf("failed to decode user: %v", err)
		}

		if gotUser.ID != targetID || gotUser.Name != "Jailson" {
			t.Errorf("unexpected user data: %+v", gotUser)
		}
	})

	t.Run("non-existing user", func(t *testing.T) {
		setupTestUsers(t)
		nonExistentID := "non-existent-id"

		req := httptest.NewRequest(http.MethodGet, "/users/"+nonExistentID, nil)
		rec := httptest.NewRecorder()

		s.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
		}

		expectedBody := "user [non-existent-id] not found\n"
		if rec.Body.String() != expectedBody {
			t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
		}
	})
}

func TestServeHTTP_POST_Users(t *testing.T) {
	s := &server{}

	t.Run("success", func(t *testing.T) {
		setupTestUsers(t)
		payload := `{"name":"Carlos"}`

		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(payload))
		rec := httptest.NewRecorder()

		s.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
		}

		contentType := rec.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("expected Content-Type application/json, got %q", contentType)
		}

		var created User
		if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
			t.Fatalf("failed to decode created user: %v", err)
		}

		if created.ID == "" {
			t.Error("expected non-empty ID for created user")
		}
		if created.Name != "Carlos" {
			t.Errorf("expected name 'Carlos', got %q", created.Name)
		}

		savedUser, ok := users[created.ID]
		if !ok {
			t.Fatalf("expected user with id %s to be saved in map", created.ID)
		}
		if savedUser.Name != "Carlos" {
			t.Errorf("expected saved user name 'Carlos', got %q", savedUser.Name)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		setupTestUsers(t)
		payload := `{"name":`

		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(payload))
		rec := httptest.NewRecorder()

		s.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}

		expectedBody := "error unmarshalling!\n"
		if rec.Body.String() != expectedBody {
			t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
		}
	})

	t.Run("missing name", func(t *testing.T) {
		setupTestUsers(t)
		payload := `{"name":""}`

		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(payload))
		rec := httptest.NewRecorder()

		s.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}

		expectedBody := "name is mandatory\n"
		if rec.Body.String() != expectedBody {
			t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
		}
	})
}

func TestServeHTTP_PUT_Users(t *testing.T) {
	s := &server{}

	t.Run("success", func(t *testing.T) {
		setupTestUsers(t)
		targetID := "595ea76c-b130-4c3b-9afa-454a17d90d60"
		payload := `{"name":"Jailson Updated"}`

		req := httptest.NewRequest(http.MethodPut, "/users/"+targetID, strings.NewReader(payload))
		rec := httptest.NewRecorder()

		s.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var message string
		if err := json.NewDecoder(rec.Body).Decode(&message); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		expectedMsg := "User with id [" + targetID + "] updated!"
		if message != expectedMsg {
			t.Errorf("expected message %q, got %q", expectedMsg, message)
		}

		if users[targetID].Name != "Jailson Updated" {
			t.Errorf("expected user name to be updated in map, got %q", users[targetID].Name)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		setupTestUsers(t)
		nonExistentID := "non-existent-id"
		payload := `{"name":"New Name"}`

		req := httptest.NewRequest(http.MethodPut, "/users/"+nonExistentID, strings.NewReader(payload))
		rec := httptest.NewRecorder()

		s.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
		}

		expectedBody := "user [non-existent-id] not found\n"
		if rec.Body.String() != expectedBody {
			t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
		}
	})

	t.Run("empty body", func(t *testing.T) {
		setupTestUsers(t)
		targetID := "595ea76c-b130-4c3b-9afa-454a17d90d60"

		req := httptest.NewRequest(http.MethodPut, "/users/"+targetID, bytes.NewReader([]byte{}))
		rec := httptest.NewRecorder()

		s.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}

		expectedBody := "name is mandatory\n"
		if rec.Body.String() != expectedBody {
			t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
		}
	})

	t.Run("invalid json body", func(t *testing.T) {
		setupTestUsers(t)
		targetID := "595ea76c-b130-4c3b-9afa-454a17d90d60"
		payload := `{"name":`

		req := httptest.NewRequest(http.MethodPut, "/users/"+targetID, strings.NewReader(payload))
		rec := httptest.NewRecorder()

		s.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
		}

		expectedBody := "error unmarshalling!\n"
		if rec.Body.String() != expectedBody {
			t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
		}
	})
}

func TestServeHTTP_PATCH_Users(t *testing.T) {
	s := &server{}

	t.Run("success", func(t *testing.T) {
		setupTestUsers(t)
		targetID := "765a42a9-8783-4808-ae5a-af57cf07d9e0"
		payload := `{"name":"Mendes Patched"}`

		req := httptest.NewRequest(http.MethodPatch, "/users/"+targetID, strings.NewReader(payload))
		rec := httptest.NewRecorder()

		s.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var message string
		if err := json.NewDecoder(rec.Body).Decode(&message); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		expectedMsg := "User with id [" + targetID + "] updated!"
		if message != expectedMsg {
			t.Errorf("expected message %q, got %q", expectedMsg, message)
		}

		if users[targetID].Name != "Mendes Patched" {
			t.Errorf("expected user name to be updated in map, got %q", users[targetID].Name)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		setupTestUsers(t)
		nonExistentID := "non-existent-id"
		payload := `{"name":"New Name"}`

		req := httptest.NewRequest(http.MethodPatch, "/users/"+nonExistentID, strings.NewReader(payload))
		rec := httptest.NewRecorder()

		s.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
		}

		expectedBody := "user [non-existent-id] not found\n"
		if rec.Body.String() != expectedBody {
			t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
		}
	})

	t.Run("empty body", func(t *testing.T) {
		setupTestUsers(t)
		targetID := "765a42a9-8783-4808-ae5a-af57cf07d9e0"

		req := httptest.NewRequest(http.MethodPatch, "/users/"+targetID, bytes.NewReader([]byte{}))
		rec := httptest.NewRecorder()

		s.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}

		expectedBody := "name is mandatory\n"
		if rec.Body.String() != expectedBody {
			t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
		}
	})

	t.Run("invalid json body", func(t *testing.T) {
		setupTestUsers(t)
		targetID := "765a42a9-8783-4808-ae5a-af57cf07d9e0"
		payload := `{"name":`

		req := httptest.NewRequest(http.MethodPatch, "/users/"+targetID, strings.NewReader(payload))
		rec := httptest.NewRecorder()

		s.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
		}

		expectedBody := "error unmarshalling!\n"
		if rec.Body.String() != expectedBody {
			t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
		}
	})
}

func TestServeHTTP_DELETE_Users(t *testing.T) {
	s := &server{}

	t.Run("success", func(t *testing.T) {
		setupTestUsers(t)
		targetID := "e5c4a6e0-c68d-4389-bcf3-c3f908c27e5e"

		req := httptest.NewRequest(http.MethodDelete, "/users/"+targetID, nil)
		rec := httptest.NewRecorder()

		s.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var message string
		if err := json.NewDecoder(rec.Body).Decode(&message); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		expectedMsg := "User with id [" + targetID + "] deleted!"
		if message != expectedMsg {
			t.Errorf("expected message %q, got %q", expectedMsg, message)
		}

		if _, exists := users[targetID]; exists {
			t.Errorf("expected user with id %s to be deleted from map", targetID)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		setupTestUsers(t)
		nonExistentID := "non-existent-id"

		req := httptest.NewRequest(http.MethodDelete, "/users/"+nonExistentID, nil)
		rec := httptest.NewRecorder()

		s.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
		}

		expectedBody := "user [non-existent-id] not found\n"
		if rec.Body.String() != expectedBody {
			t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
		}
	})
}

func TestServeHTTP_Default_MethodNotFound(t *testing.T) {
	setupTestUsers(t)
	s := &server{}

	req := httptest.NewRequest(http.MethodOptions, "/users", nil)
	rec := httptest.NewRecorder()

	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}

	expectedBody := "not found\n"
	if rec.Body.String() != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
	}
}
