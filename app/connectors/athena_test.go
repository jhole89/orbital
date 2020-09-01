package connectors

import (
	"bou.ke/monkey"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestAwsAthenaConnector_Connect(t *testing.T) {
	monkey.Patch(sql.Open, func(driverName string, dataSourceName string) (*sql.DB, error) {
		return (*sql.DB)(nil), nil
	})
	var a AwsAthenaConnector
	a.Connect("some-address")
	assert.IsType(t, (*sql.DB)(nil), a.Connection, "Connect should establish an sql connection.")

	monkey.Unpatch(sql.Open)
}

func TestAwsAthenaConnector_Index(t *testing.T) {

}

func TestAwsAthenaConnector_Query(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO product_viewers").WithArgs(2, 3).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()


	var s *sql.DB
	monkey.PatchInstanceMethod(reflect.TypeOf(s), "Query", func(query string, args ...interface{}) (*sql.Rows, error) {
		return
	}

return (*sql.DB)(nil), nil
	})

}

func TestAwsAthenaConnector_describeTables(t *testing.T) {

}

func TestAwsAthenaConnector_getDatabases(t *testing.T) {

}

func TestAwsAthenaConnector_getTableMeta(t *testing.T) {

}

func TestAwsAthenaConnector_getTables(t *testing.T) {

}
