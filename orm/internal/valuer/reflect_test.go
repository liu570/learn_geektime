package valuer

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"learn_geektime/orm/model"
	"testing"
)

func Test_ReflectValue_SetColumn(t *testing.T) {
	testSetColumn(t, NewReflectValue)
}

func testSetColumn(t *testing.T, creator Creator) {
	testCases := []struct {
		name   string
		entity any
		rows   func() *sqlmock.Rows

		wantErr    error
		wantEntity any
	}{
		{
			name:   "simple entity",
			entity: &TestModel{},
			rows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
				rows.AddRow("1", "Tom", 18, "Jerry")
				return rows
			},
			wantEntity: &TestModel{
				Id:        1,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "Jerry"},
			},
		},
		{
			name:   "partial columns",
			entity: &TestModel{},
			rows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "last_name"})
				rows.AddRow("1", "Jerry")
				return rows
			},
			wantEntity: &TestModel{
				Id:       1,
				LastName: &sql.NullString{Valid: true, String: "Jerry"},
			},
		},
		{
			name:   "order",
			entity: &TestModel{},
			rows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "last_name", "first_name", "age"})
				rows.AddRow("1", "Jerry", "Tom", 18)
				return rows
			},
			wantEntity: &TestModel{
				Id:        1,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "Jerry"},
			},
		},
	}
	r := model.NewRegistry()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			m, err := r.Get(tt.entity)
			require.NoError(t, err)
			val := creator(tt.entity, m)
			mock.ExpectQuery("SELECT XXX").WillReturnRows(tt.rows())
			rows, err := mockDB.Query("SELECT XXX")
			require.NoError(t, err)
			err = val.SetColumn(rows)
			assert.Equal(t, tt.wantErr, err)
			if err == nil {
				assert.Equal(t, tt.wantEntity, tt.entity)
			}

		})
	}
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}
