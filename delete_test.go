package sb

import (
	"reflect"
	"testing"
)

func TestDeleteSimple(t *testing.T) {
	wantSql := "DELETE FROM `tb_user`"
	gotSql := Delete("tb_user").Build()

	if wantSql != gotSql {
		t.Errorf("delete simple sql:\nwant:%s\n got:%s\n", wantSql, gotSql)
	}
}
func BenchmarkDeleteSimple(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Delete("tb_user").Build()
	}
}

func TestDeleteComplex(t *testing.T) {
	wantSql := "DELETE FROM `tb_user` WHERE id = ? OR name = ?"
	gotSql := Delete("tb_user").
		Where(Eq("id"), Eq("name").Or()).
		Build()

	if wantSql != gotSql {
		t.Errorf("delete complex sql:\nwant:%s\n got:%s\n", wantSql, gotSql)
	}
}

func TestDeleteComplexWithValue(t *testing.T) {
	wantSql := "DELETE FROM `tb_user` WHERE id = ? OR name = ?"
	wantVal := []any{"idVal", "nameVal"}
	gotSql, gotVal := Delete("tb_user").
		Where(
			Eq("id", "idVal"),
			Eq("name", "nameVal").Or(),
		).
		BuildWithValue()

	if wantSql != gotSql || !reflect.DeepEqual(wantVal, gotVal) {
		t.Errorf("delete complex sql:\nwant:%s\n got:%s\n", wantSql, gotSql)
		t.Errorf("delete complex val:\nwant:%v\n got:%v\n", wantVal, gotVal)
	}
}

func BenchmarkDeleteComplexWithValue(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Delete("tb_user").
			Where(
				Eq("id", "idVal"),
				Eq("name", "nameVal").Or(),
			).
			BuildWithValue()
	}
}
