package valuer

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"learn_geektime/orm/model"
	"testing"
)

func BenchmarkSetColumn(b *testing.B) {

	fn := func(b *testing.B, creator Creator) {
		mockDB, mock, err := sqlmock.New()
		require.NoError(b, err)
		defer mockDB.Close()
		mockRows := sqlmock.NewRows([]string{"id", "last_name", "first_name", "age"})
		row := []driver.Value{"1", "Jerry", "Tom", 18}
		for i := 0; i < b.N; i++ {
			mockRows.AddRow(row...)
		}

		mock.ExpectQuery("SELECT XXX").WillReturnRows(mockRows)
		rows, err := mockDB.Query("SELECT XXX")
		r := model.NewRegistry()
		m, err := r.Get(&TestModel{})
		require.NoError(b, err)
		b.ResetTimer() // 重置时间
		for i := 0; i < b.N; i++ {
			rows.Next()
			val := creator(&TestModel{}, m)
			err = val.SetColumn(rows)
		}
	}
	b.Run("reflect", func(b *testing.B) {
		fn(b, NewReflectValue)
	})
	b.Run("unsafe", func(b *testing.B) {
		fn(b, NewUnsafeValue)
	})
}
