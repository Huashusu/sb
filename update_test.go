package sb

import (
	"reflect"
	"testing"
)

func TestUpdateSimple(t *testing.T) {
	wantSql := "UPDATE tb_user SET `name` = ? , count = count + 1 , `password` = ?"
	gotSql := Update("tb_user").
		Column("name", "count = count + 1", "password").
		Build()
	if wantSql != gotSql {
		t.Errorf("update simple sql:\nwant:%s\n got:%s\n", wantSql, gotSql)
	}
}

func TestUpdateSimpleWithValue(t *testing.T) {
	wantSql := "UPDATE tb_user SET `name` = ? , count = count + 1 , `password` = ? WHERE id1 = ? AND id2 = ?"
	wantVal := []any{"nameVal", "passwordVal", "idVal1", "idVal2"}
	gotSql, gotVal := Update("tb_user").
		Column("name", "count = count + 1", "password").
		Where(
			Eq("id1", "idVal1"),
			Eq("id2", "idVal2"),
		).
		Set("nameVal", "passwordVal").
		BuildWithValue()
	if wantSql != gotSql || !reflect.DeepEqual(wantVal, gotVal) {
		t.Errorf("update simple sql:\nwant:%s\n got:%s\n", wantSql, gotSql)
		t.Errorf("update simple val:\nwant:%v\n got:%v\n", wantVal, gotVal)
	}
}

func BenchmarkUpdateSimpleWithValue(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Update("tb_user").
			Column("name", "count = count + 1", "password").
			Where(
				Eq("id1", "idVal1"),
				Eq("id2", "idVal2"),
			).
			Set("nameVal", "passwordVal").
			BuildWithValue()
	}
}

func TestUpdateComplex(t *testing.T) {
	wantSql := "UPDATE tb_user SET `name` = ? , count = count + 1 , `password` = ? WHERE id = ? OR name = ?"
	gotSql := Update("tb_user").
		Column("name", "count = count + 1", "password").
		Where(Eq("id"), Eq("name").Or()).
		Build()
	if wantSql != gotSql {
		t.Errorf("update simple sql:\nwant:%s\n got:%s\n", wantSql, gotSql)
	}
}

func TestUpdateComplexWithValue(t *testing.T) {
	wantSql := "UPDATE tb_user SET `name` = ? , count = count + 1 , `password` = ? WHERE id = ? OR name = ?"
	wantVal := []any{"nameVal1", "passwordVal", "idVal", "nameVal2"}
	gotSql, gotVal := Update("tb_user").
		Column("name", "count = count + 1", "password").
		Set("nameVal1", "passwordVal").
		Where(
			Eq("id", "idVal"),
			Eq("name", "nameVal2").Or(),
		).
		BuildWithValue()
	if wantSql != gotSql || !reflect.DeepEqual(wantVal, gotVal) {
		t.Errorf("update simple sql:\nwant:%s\n got:%s\n", wantSql, gotSql)
	}
}

func BenchmarkUpdateComplexWithValue(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Update("tb_user").
			Column("name", "count = count + 1", "password").
			Set("nameVal1", "passwordVal").
			Where(
				Eq("id", "idVal"),
				Eq("name", "nameVal2").Or(),
			).
			BuildWithValue()
	}
}
