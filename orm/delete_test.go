package orm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleter_Build(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name    string
		d       QueryBuilder
		want    *Query
		wantErr error
	}{
		{
			name: "single delete",
			d:    NewDeleter[TestModel](db),
			want: &Query{
				SQL: "DELETE FROM `test_model`",
			},
			wantErr: nil,
		},
		{
			name: "Order delete",
			d:    NewDeleter[TestModel](db),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.d.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.want, query)
		})
	}
}
