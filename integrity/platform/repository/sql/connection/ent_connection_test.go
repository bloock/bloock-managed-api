package connection

import (
	"bloock-managed-api/ent"
	mock_connection "bloock-managed-api/integrity/platform/repository/sql/connection/mocks"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConnection_NewConnection(t *testing.T) {
	ctrl := gomock.NewController(t)
	storageConnector := mock_connection.NewMockSQLConnector(ctrl)

	tests := []struct {
		name   string
		url    string
		driver string
	}{
		{name: "given mysql url it should be detected", url: "mysql://username:password@localhost:3306/mydatabase", driver: Mysql},
		{name: "given postgres url it should be detected", url: "postgresql://username:password@localhost:5432/mydatabase", driver: Postgres},
		{name: "given sqlite memory url it should be detected", url: "file:ent?mode=memory&cache=shared&_fk=1", driver: Sqlite},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := ent.NewClient()
			storageConnector.EXPECT().Connect(test.driver, test.url).Return(client, nil)
			_, err := NewEntConnection(test.url, storageConnector, zerolog.Logger{})
			assert.NoError(t, err)
		})
	}

	t.Run("given unsupported database error should be returned", func(t *testing.T) {
		_, err := NewEntConnection("somedb://username:password@localhost:3306/mydatabase", storageConnector, zerolog.Logger{})
		assert.Error(t, err)
	})
}
