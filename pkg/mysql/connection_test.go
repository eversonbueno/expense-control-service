package mysql

import (
"testing"

"github.com/DATA-DOG/go-sqlmock"
"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	t.Run("Quando falha na conexão, deve retornar erro", func(t *testing.T) {
		mysqlInstance := New("localhost", "user", "password", "test_db", 0, 0, 0)
		conn, err := mysqlInstance.Connect()
		assert.Error(t, err)
		assert.Nil(t, conn)
	})

}

func TestClose(t *testing.T) {
	t.Run("Quando a conexão é fechada com sucesso, não deve retornar erro", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer mockDB.Close()

		mock.ExpectClose().WillReturnError(nil)

		mysqlInstance := &MySQL{}
		mysqlInstance.conn = mockDB

		err = mysqlInstance.Close()
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Quando a conexão é nil, deve não retornar erro", func(t *testing.T) {
		mysqlInstance := &MySQL{}
		err := mysqlInstance.Close()
		assert.Error(t, err)
	})
}