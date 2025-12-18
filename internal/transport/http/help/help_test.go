package help_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alonsoF100/authorization-service/internal/transport/http/help"
	"github.com/stretchr/testify/require"
)

func TestWriteJSON(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		data       any
		wantBody   string
	}{
		{
			name:       "some data",
			statusCode: 200,
			data:       map[string]string{"ok": "true"},
			wantBody:   `{"ok":"true"}`,
		},
		{
			name:       "empty data",
			statusCode: 204,
			data:       nil,
			wantBody:   "null",
		},
		{
			name:       "slice",
			statusCode: 200,
			data:       []string{"a", "b", "c"},
			wantBody:   `["a","b","c"]`,
		},
		{
			name:       "struct with tags",
			statusCode: 200,
			data: struct {
				ID   string `json:"id"`
				Name string `json:"name,omitempty"`
			}{
				ID:   "123",
				Name: "John",
			},
			wantBody: `{"id":"123","name":"John"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			help.WriteJSON(w, tt.statusCode, tt.data)

			require.Equal(t, tt.statusCode, w.Code)
			require.Equal(t, "application/json", w.Header().Get("Content-Type"))
			require.JSONEq(t, tt.wantBody, w.Body.String())
		})
	}

	t.Run("error handling", func(t *testing.T) {
		w := httptest.NewRecorder()
		unserializable := make(chan int)

		help.WriteJSON(w, 200, unserializable)

		require.Equal(t, 200, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	})
}
