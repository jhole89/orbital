package connectors

import (
	"bou.ke/monkey"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
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

func TestAwsAthenaConnector_getDatabases(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := []string{"database_name"}

	mock.ExpectQuery("SHOW SCHEMAS").WillReturnRows(sqlmock.NewRows(columns).AddRow("foo-db").AddRow("bar-db"))

	var a AwsAthenaConnector
	a.Connection = db
	dbs := a.getDatabases()
	assert.Equal(t, []string{"foo-db", "bar-db"}, dbs)
}

func TestAwsAthenaConnector_getTables(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := []string{"tab_name"}

	mock.ExpectQuery("SHOW TABLES IN foo").WillReturnRows(sqlmock.NewRows(columns).AddRow("foo-tab").AddRow("bar-tab"))

	var a AwsAthenaConnector
	a.Connection = db
	tabs := a.getTables("foo")

	assert.Equal(t, []string{"foo-tab", "bar-tab"}, tabs)
}

func TestAwsAthenaConnector_describeTables(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := []string{"column", "type"}

	mock.ExpectQuery("DESCRIBE foo.bar").WillReturnRows(
		sqlmock.NewRows(
			columns,
		).AddRow(
			"some-string-field\tvarchar", "",
		).AddRow(
			"some-bool-field\tboolean", "",
		).AddRow(
			"some-int-field\tbigint", "",
		).AddRow(
			"some-double-field\tdouble", "",
		),
	)

	var a AwsAthenaConnector
	a.Connection = db
	cols := a.describeTables("foo", "bar")

	expected := []Column{
		{Name: "some-string-field", Type: "varchar"},
		{Name: "some-bool-field", Type: "boolean"},
		{Name: "some-int-field", Type: "bigint"},
		{Name: "some-double-field", Type: "double"},
	}

	assert.Equal(t, expected, cols)
}

func TestAwsAthenaConnector_Index(t *testing.T) {

}
