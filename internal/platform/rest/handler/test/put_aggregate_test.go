package test

import (
	"bytes"
	"github.com/bloock/bloock-managed-api/internal/config"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/connection"
	"github.com/bloock/bloock-managed-api/internal/platform/rest"
	"github.com/bloock/bloock-managed-api/pkg/request"
	"github.com/go-pdf/fpdf"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"net/http"
	"os"
	"testing"
)

func TestPutAggregate(t *testing.T) {
	logger := zerolog.Logger{}
	entConnector := connection.NewEntConnector(logger)
	conn, err := connection.NewEntConnection(config.Configuration.Db.ConnectionString, entConnector, logger)
	if err != nil {
		panic(err)
	}
	err = conn.Migrate()
	if err != nil {
		panic(err)
	}

	server, err := rest.NewServer(
		zerolog.Logger{},
		conn,
		1000,
	)
	require.NoError(t, err)
	go func() {
		err = server.Start()
		require.NoError(t, err)
	}()

	randomInputFile := make([]File, 0)
	filename := "random.pdf"
	for i := 0; i < 500; i++ {
		content, err := generateRandomPDF("random.pdf")
		require.NoError(t, err)
		randomInputFile = append(randomInputFile, File{
			Filename: filename,
			Content:  content,
		})
	}

	t.Run("given no aggregate mode enabled, should return an error", func(t *testing.T) {
		req := request.ProcessFormRequest{
			Integrity: request.ProcessFormIntegrityRequest{
				Enabled:   true,
				Aggregate: true,
			},
		}
		_, status, err := SendRequest(server, &req, randomInputFile[0], nil)
		require.Error(t, err)

		assert.Equal(t, http.StatusBadRequest, status)
	})

	t.Run("given no aggregate mode enabled, should return an error", func(t *testing.T) {
		viper.Set("bloock.api_key", "")

		req := request.ProcessFormRequest{
			Integrity: request.ProcessFormIntegrityRequest{
				Enabled:   true,
				Aggregate: true,
			},
		}
		_, status, err := SendRequest(server, &req, randomInputFile[0], nil)
		require.Error(t, err)

		assert.Equal(t, http.StatusBadRequest, status)
	})

	t.Run("given aggregate mode enabled, should return a valid hash", func(t *testing.T) {
		for _, inputFile := range randomInputFile {
			req := request.ProcessFormRequest{
				Integrity: request.ProcessFormIntegrityRequest{
					Enabled:   true,
					Aggregate: true,
				},
			}

			res, status, err := SendRequest(server, &req, inputFile, nil)
			require.NoError(t, err)

			assert.Equal(t, http.StatusOK, status)
			assert.NotEmpty(t, res.Hash)
			assert.Equal(t, 0, res.Integrity.AnchorId)
		}

		ok, status, err := PutAggregate(server)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.True(t, ok)
	})
}

func generateRandomPDF(filename string) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")

	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello, World!")

	buff := new(bytes.Buffer)
	if err := pdf.Output(buff); err != nil {
		return nil, err
	}

	random := byte(rand.Intn(256))
	buff.Write([]byte{random})
	// Change one byte of the content
	pdf.RawWriteBuf(buff)

	// Write the modified buffer to a file
	err := pdf.OutputFileAndClose(filename)
	if err != nil {
		return nil, err
	}

	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}
