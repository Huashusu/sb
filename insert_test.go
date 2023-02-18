package sb

import (
	"reflect"
	"testing"
)

func TestInsertSimple(t *testing.T) {
	want := "INSERT INTO `tb_1` ( `id`, `name`, `password` ) VALUE ( ? , ? , ? )"

	got := Insert("tb_1").
		Column("id", "name", "password").
		Set(1, "name1", "password").
		Build()
	if want != got {
		t.Errorf("insert simple sql:\nwant:%s\n got:%s\n", want, got)
	}
}

func TestInsertSimpleWithValue(t *testing.T) {
	wantSql := "INSERT INTO `tb_1` ( `id`, `name`, `password` ) VALUE ( ? , ? , ? )"
	wantVal := []any{1, "name1", "password"}

	gotSql, gotVal := Insert("tb_1").
		Column("id", "name", "password").
		Set(1, "name1", "password").
		BuildWithValue()
	if wantSql != gotSql || !reflect.DeepEqual(wantVal, gotVal) {
		t.Errorf("insert simple sql:\nwant:%s\n got:%s\n", wantSql, gotSql)
		t.Errorf("insert simple val:\nwant:%v\n got:%v\n", wantVal, gotVal)
	}
}

func BenchmarkInsertSimpleWithValue(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Insert("tb_1").
			Column("id", "name", "password").
			Set(1, "name1", "password").
			BuildWithValue()
	}
}

func TestInsertComplex(t *testing.T) {
	want := "INSERT INTO `tb_1` ( `id`, `name`, `password` ) VALUES ( ? , ? , ? ),( ? , ? , ? ),( ? , ? , ? )"
	got := Insert("tb_1").
		Column("id", "name", "password").
		SetSlice([][]any{
			{1, "name1", "pass1"},
			{2, "name2", "pass2"},
			{3, "name3", "pass3"},
		}).
		Build()
	if want != got {
		t.Errorf("insert complex sql:\nwant:%s\n got:%s\n", want, got)
	}
}

func TestInsertComplexWithValue(t *testing.T) {
	wantSql := "INSERT INTO `tb_1` ( `id`, `name`, `password` ) VALUES ( ? , ? , ? ),( ? , ? , ? ),( ? , ? , ? )"
	wantVal := []any{
		1, "name1", "pass1",
		2, "name2", "pass2",
		3, "name3", "pass3",
	}
	gotSql, gotVal := Insert("tb_1").
		Column("id", "name", "password").
		SetSlice([][]any{
			{1, "name1", "pass1"},
			{2, "name2", "pass2"},
			{3, "name3", "pass3"},
		}).
		BuildWithValue()
	if wantSql != gotSql || !reflect.DeepEqual(wantVal, gotVal) {
		t.Errorf("insert complex sql:\nwant:%s\n got:%s\n", wantSql, gotSql)
		t.Errorf("insert complex val:\nwant:%v\n got:%v\n", wantVal, gotVal)
	}
}

func BenchmarkInsertComplexWithValue(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Insert("tb_1").
			Column("id", "name", "password").
			SetSlice([][]any{
				{1, "name1", "pass1"},
				{2, "name2", "pass2"},
				{3, "name3", "pass3"},
			}).
			BuildWithValue()
	}
}
