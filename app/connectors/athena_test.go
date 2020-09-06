package connectors

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log"
	"testing"
)

func sqlMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestAwsAthena_getDatabases(t *testing.T) {

	asserter := assert.New(t)

	db, mock := sqlMock()
	a := &AwsAthena{db}

	defer func() {
		a.Close()
	}()

	rows := sqlmock.NewRows([]string{"database_name"}).AddRow("foo-db").AddRow("bar-db")
	mock.ExpectQuery("SHOW SCHEMAS").WillReturnRows(rows)

	dbs, err := a.getDatabases()

	asserter.NoError(err)
	asserter.NotEmpty(dbs)
	asserter.Len(dbs, 2)
	asserter.Equal([]string{"foo-db", "bar-db"}, dbs)
}

func TestAwsAthena_getTables(t *testing.T) {

	asserter := assert.New(t)

	db, mock := sqlMock()
	a := &AwsAthena{db}

	defer func() {
		a.Close()
	}()

	rows := sqlmock.NewRows([]string{"tab_name"}).AddRow("foo-tab").AddRow("bar-tab")
	mock.ExpectQuery("SHOW TABLES IN foo").WillReturnRows(rows)

	tabs, err := a.getTables("foo")

	asserter.NoError(err)
	asserter.NotEmpty(tabs)
	asserter.Len(tabs, 2)
	asserter.Equal([]string{"foo-tab", "bar-tab"}, tabs)
}

func TestAwsAthena_describeTables(t *testing.T) {

	asserter := assert.New(t)

	db, mock := sqlMock()
	a := &AwsAthena{db}

	defer func() {
		a.Close()
	}()

	rows := sqlmock.NewRows([]string{"column", "type"}).
		AddRow("some-string-field\tvarchar", "").
		AddRow("some-bool-field\tboolean", "").
		AddRow("some-int-field\tbigint", "").
		AddRow("some-double-field\tdouble", "")
	mock.ExpectQuery("DESCRIBE foo.bar").WillReturnRows(rows)

	cols, err := a.describeTables("foo", "bar")

	expected := []Column{
		{Name: "some-string-field", Type: "varchar"},
		{Name: "some-bool-field", Type: "boolean"},
		{Name: "some-int-field", Type: "bigint"},
		{Name: "some-double-field", Type: "double"},
	}

	asserter.NoError(err)
	asserter.NotEmpty(cols)
	asserter.Len(cols, 4)
	asserter.Equal(expected, cols)
}

type MockAthena struct {
	mock.Mock
}

func (a *MockAthena) getDatabases() ([]string, error) {
	args := a.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (a *MockAthena) getTables(database string) ([]string, error) {
	args := a.Called(database)
	return args.Get(0).([]string), args.Error(1)
}

func (a *MockAthena) describeTables(database, table string) ([]Column, error) {
	args := a.Called(database, table)
	return args.Get(0).([]Column), args.Error(1)
}

func TestAwsAthena_Index(t *testing.T) {
	asserter := assert.New(t)

	am := new(MockAthena)
	am.On("getDatabases").Return([]string{"foo-db", "bar-db"}, nil)
	am.On("getTables", "foo-db").Return([]string{"foo1-tab", "bar1-tab"}, nil)
	am.On("getTables", "bar-db").Return([]string{"foo2-tab", "bar2-tab"}, nil)
	am.On("describeTables", "foo-db", "foo1-tab").
		Return([]Column{{Name: "some-string-field", Type: "varchar"}, {Name: "some-bool-field", Type: "boolean"}}, nil)
	am.On("describeTables", "foo-db", "bar1-tab").
		Return([]Column{{Name: "some-string-field", Type: "varchar"}, {Name: "some-bool2-field", Type: "boolean"}}, nil)
	am.On("describeTables", "bar-db", "foo2-tab").
		Return([]Column{{Name: "some-string2-field", Type: "varchar"}, {Name: "some-int-field", Type: "bigint"}}, nil)
	am.On("describeTables", "bar-db", "bar2-tab").
		Return([]Column{{Name: "some-string3-field", Type: "varchar"}, {Name: "some-int-field", Type: "bigint"}}, nil)

	nodes, err := index(am.getDatabases, am.getTables, am.describeTables)

	expected := []*Node{
		{Name: "foo-db", Context: "database", Children: []*Node{
			{Name: "foo1-tab", Context: "table", Children: []*Node{
				{Name: "some-string-field", Context: "field", Properties: map[string]string{"data-type": "varchar"}},
				{Name: "some-bool-field", Context: "field", Properties: map[string]string{"data-type": "boolean"}},
			}},
			{Name: "bar1-tab", Context: "table", Children: []*Node{
				{Name: "some-string-field", Context: "field", Properties: map[string]string{"data-type": "varchar"}},
				{Name: "some-bool2-field", Context: "field", Properties: map[string]string{"data-type": "boolean"}},
			}},
		}},
		{Name: "bar-db", Context: "database", Children: []*Node{
			{Name: "foo2-tab", Context: "table", Children: []*Node{
				{Name: "some-string2-field", Context: "field", Properties: map[string]string{"data-type": "varchar"}},
				{Name: "some-int-field", Context: "field", Properties: map[string]string{"data-type": "bigint"}},
			}},
			{Name: "bar2-tab", Context: "table", Children: []*Node{
				{Name: "some-string3-field", Context: "field", Properties: map[string]string{"data-type": "varchar"}},
				{Name: "some-int-field", Context: "field", Properties: map[string]string{"data-type": "bigint"}},
			}},
		}},
	}

	asserter.NoError(err)
	asserter.NotEmpty(nodes)
	asserter.Len(nodes, 2)
	asserter.Equal(expected, nodes)
}
