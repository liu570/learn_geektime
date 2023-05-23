package orm

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"learn_geektime/orm/internal/errs"
	"testing"
)

func TestRawQuerier_Get(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	testCases := []struct {
		name     string
		query    string
		mockErr  error
		mockRows *sqlmock.Rows
		wantErr  error
		wantVal  *TestModel
	}{
		{
			name:    "raw select row",
			query:   "SELECT `id`,`first_name`,`age`,`last_name` FROM `test_model`",
			mockErr: nil,
			mockRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
				rows.AddRow([]byte("123"), []byte("Ming"), []byte("18"), []byte("Deng"))
				return rows
			}(),
			wantVal: &TestModel{
				Id:        123,
				FirstName: "Ming",
				Age:       18,
				LastName: &sql.NullString{
					Valid:  true,
					String: "Deng",
				},
			},
		},
		{
			name:    "invalid columns",
			query:   "SELECT `id`,`first_name`,`age`,`last_name` FROM `test_model`",
			mockErr: nil,
			mockRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "first_name", "gender"})
				rows.AddRow([]byte("123"), []byte("Ming"), []byte("male"))
				return rows
			}(),
			wantErr: errs.NewErrUnknownColumn("gender"),
		},
		{
			name:    "less columns",
			query:   "SELECT `id`,`first_name` FROM `test_model`",
			mockErr: nil,
			mockRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "first_name"})
				rows.AddRow([]byte("123"), []byte("Ming"))
				return rows
			}(),
			wantVal: &TestModel{
				Id:        123,
				FirstName: "Ming",
			},
		},
	}

	for _, tc := range testCases {
		if tc.mockErr != nil {
			mock.ExpectQuery(tc.query).WillReturnError(tc.mockErr)
		} else {
			mock.ExpectQuery(tc.query).WillReturnRows(tc.mockRows)
		}
	}
	db, err := OpenDB(mockDB, DBUseReflectValuer())
	require.NoError(t, err)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			args := []any{}
			res, err := RawQuery[TestModel](db, tt.query, args...).Get(context.Background())
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantVal, res)
		})
	}
}

func TestRawQuerier_GetMulti(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	testCases := []struct {
		name     string
		query    string
		mockErr  error
		mockRows *sqlmock.Rows
		wantErr  error
		wantVal  []*TestModel
	}{
		{
			name:    "raw select rows",
			query:   "SELECT `id`,`first_name`,`age`,`last_name` FROM `test_model`",
			mockErr: nil,
			mockRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
				rows.AddRow([]byte("123"), []byte("Ming"), []byte("18"), []byte("Deng"))
				rows.AddRow([]byte("124"), []byte("WenJie"), []byte("22"), []byte("Liu"))
				return rows
			}(),
			wantVal: func() []*TestModel {
				res := make([]*TestModel, 0, 2)
				res = append(res, &TestModel{
					Id:        123,
					FirstName: "Ming",
					Age:       18,
					LastName: &sql.NullString{
						Valid:  true,
						String: "Deng",
					},
				})
				res = append(res, &TestModel{
					Id:        124,
					FirstName: "WenJie",
					Age:       22,
					LastName: &sql.NullString{
						Valid:  true,
						String: "Liu",
					},
				})
				return res
			}(),
		},
		{
			name:    "invalid columns",
			query:   "SELECT `id`,`first_name`,`age`,`last_name` FROM `test_model`",
			mockErr: nil,
			mockRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "first_name", "gender"})
				rows.AddRow([]byte("123"), []byte("Ming"), []byte("male"))
				rows.AddRow([]byte("124"), []byte("Ming"), []byte("male"))
				return rows
			}(),
			wantErr: errs.NewErrUnknownColumn("gender"),
		},
		{
			name:    "less columns",
			query:   "SELECT `id`,`first_name` FROM `test_model`",
			mockErr: nil,
			mockRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "first_name"})
				rows.AddRow([]byte("123"), []byte("Ming"))
				rows.AddRow([]byte("124"), []byte("WenJie"))
				return rows
			}(),
			wantVal: func() []*TestModel {
				res := make([]*TestModel, 0, 2)
				res = append(res, &TestModel{
					Id:        123,
					FirstName: "Ming",
				})
				res = append(res, &TestModel{
					Id:        124,
					FirstName: "WenJie",
				})
				return res
			}(),
		},
	}

	for _, tc := range testCases {
		if tc.mockErr != nil {
			mock.ExpectQuery(tc.query).WillReturnError(tc.mockErr)
		} else {
			mock.ExpectQuery(tc.query).WillReturnRows(tc.mockRows)
		}
	}
	db, err := OpenDB(mockDB, DBUseReflectValuer())
	require.NoError(t, err)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			args := []any{}
			res, err := RawQuery[TestModel](db, tt.query, args...).GetMulti(context.Background())
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantVal, res)
		})
	}
}
