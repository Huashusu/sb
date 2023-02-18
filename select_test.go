package sb

import (
	"reflect"
	"testing"
)

func TestSelectSimple(t *testing.T) {
	want := `SELECT * FROM tb_user WHERE username = ?`
	got, _ := Select().
		From("tb_user").
		Where(Eq("username")).
		BuildWithValue()
	if !(want == got) {
		t.Errorf("simple sql error:\n\tWant:%v\n\t Got:%v\n", want, got)
	}
}

func TestSelectSimpleWithValue(t *testing.T) {
	wantSql := `SELECT * FROM tb_user WHERE username = ?`
	wantVal := []any{"alice"}
	gotSql, gotVal := Select().
		From("tb_user").
		Where(Eq("username", "alice")).
		BuildWithValue()
	if !(wantSql == gotSql) || !reflect.DeepEqual(wantVal, gotVal) {
		t.Errorf("simple sql error:\n\tWant:%v\n\t Got:%v\n", wantSql, gotSql)
		t.Errorf("simple val error:\n\tWant:%v\n\t Got:%v\n", wantVal, gotVal)
	}
}

func BenchmarkSelectSimpleWithValue(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Select().
			From("tb_user").
			Where(Eq("username", "alice")).
			BuildWithValue()
	}
}

func TestSelectComplex(t *testing.T) {
	wantSql := "SELECT * FROM tb_user AS u1 " +
		"LEFT JOIN (SELECT id FROM tb_user WHERE name = ?) AS u2 ON u1.id = u2.id AND u1.id > ? " +
		"WHERE u2.id > ? AND u1.create_time < ? AND u1.nickname IS NOT NULL OR u1.email IS NULL GROUP BY u1.id HAVING u1.id >= ? ORDER BY u1.id ASC, u2.create_time DESC LIMIT 100 OFFSET 10"

	subSql, _ := Select("id").
		From("tb_user").
		Where(Eq("name", "name1")).
		SubQueryWithValue("u2")
	gotSql, _ := Select().
		From("tb_user AS u1").
		Join(LeftJoin(subSql, "u1.id = u2.id", Gt("u1.id", "id1"))).
		Where(
			Gt("u2.id", "id2"),
			Lt("u1.create_time", "time1"),
			IsNotNULL("u1.nickname"),
			IsNULL("u1.email").Or(),
		).
		GroupBy("u1.id").
		OrderBy(
			OrderByAsc("u1.id"),
			OrderByDesc("u2.create_time"),
		).
		Having(
			GtEq("u1.id", "id3"),
		).
		Limit(100).
		Offset(10).
		BuildWithValue()
	if !(wantSql == gotSql) {
		t.Errorf("select complex sql:\nWant:%v\n Got:%v\n", wantSql, gotSql)
	}
}

func TestSelectComplexWithValue(t *testing.T) {
	wantSql := "SELECT * FROM tb_user AS u1 " +
		"LEFT JOIN (" +
		"SELECT id FROM tb_user WHERE name = ?" +
		") AS u2 ON u1.id = u2.id AND u1.id > ? " +
		"WHERE u2.id > ? AND u1.create_time < ? AND u1.nickname IS NOT NULL OR u1.email IS NULL " +
		"GROUP BY u1.id " +
		"HAVING u1.id >= ? " +
		"ORDER BY u1.id ASC, u2.create_time DESC " +
		"LIMIT 100 OFFSET 10"
	wantVal := []any{"name1", "id1", "id2", "time1", "id3"}

	subSql, subVal := Select("id").
		From("tb_user").
		Where(Eq("name", "name1")).
		SubQueryWithValue("u2")
	gotSql, gotVal := Select().
		From("tb_user AS u1").
		Join(
			LeftJoin(subSql, "u1.id = u2.id", Gt("u1.id", "id1")),
			subVal...,
		).
		Where(
			Gt("u2.id", "id2"),
			Lt("u1.create_time", "time1"),
			IsNotNULL("u1.nickname"),
			IsNULL("u1.email").Or(),
		).
		GroupBy("u1.id").
		OrderBy(
			OrderByAsc("u1.id"),
			OrderByDesc("u2.create_time"),
		).
		Having(
			GtEq("u1.id", "id3"),
		).
		Limit(100).
		Offset(10).
		BuildWithValue()
	if !(wantSql == gotSql) || !(reflect.DeepEqual(wantVal, gotVal)) {
		t.Errorf("select complex sql:\nWant:%v\n Got:%v\n", wantSql, gotSql)
		t.Errorf("select complex val:\nWant:%v\n Got:%v\n", wantVal, gotVal)
	}
}

func BenchmarkSelectComplexWithValue(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		subSql, subVal := Select("id").
			From("tb_user").
			Where(Eq("name", "name1")).
			SubQueryWithValue("u2")
		_, _ = Select().
			From("tb_user AS u1").
			Join(LeftJoin(subSql, "u1.id = u2.id", Gt("u1.id", "id1")), subVal...).
			Where(
				Gt("u2.id", "id2"),
				Lt("u1.create_time", "time1"),
				IsNotNULL("u1.nickname"),
				IsNULL("u1.email").Or(),
			).
			GroupBy("u1.id").
			OrderBy(
				OrderByAsc("u1.id"),
				OrderByDesc("u2.create_time"),
			).
			Having(
				GtEq("u1.id", "id3"),
			).
			Limit(100).
			Offset(10).
			BuildWithValue()
	}
}
