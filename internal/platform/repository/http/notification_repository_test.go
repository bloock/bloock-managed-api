package http

import (
	"bloock-managed-api/internal/platform/test_utils/fixtures"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpNotificationRepository_NotifyCertification(t *testing.T) {
	expectedWhResp := `{"hash":"testHash","File":null,"webhook_response":"test"}`
	pdfContent := fixtures.PDFContent

	server := MockServer{}.Expect(t, expectedWhResp, pdfContent)

	err := NewHttpNotificationRepository(http.Client{}, server.URL, zerolog.Logger{}).
		NotifyCertification("testHash", "test", pdfContent)

	assert.NoError(t, err)
}

type MockServer struct {
}

func (s MockServer) Expect(t *testing.T, expectedWhResp string, expectedFile []byte) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		method := r.Method
		assert.Equal(t, "POST", method)
		m, err := r.MultipartReader()
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		whPart, err := m.NextPart()
		if err != nil {
			http.Error(w, "wh_response not found", http.StatusBadRequest)
			return
		}
		whRespContent, err := io.ReadAll(whPart)
		if err != nil {
			http.Error(w, "wh_response not found", http.StatusBadRequest)
			return
		}
		assert.Equal(t, expectedWhResp, string(whRespContent))

		file, err := m.NextPart()
		if err != nil {
			http.Error(w, "file not found", http.StatusBadRequest)
			return
		}
		fileContent, _ := io.ReadAll(file)
		assert.Equal(t, expectedFile, fileContent)

		w.WriteHeader(http.StatusOK)
	}))
}
