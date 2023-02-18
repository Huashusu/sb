package sb

import "fmt"

func ExampleSelect() {
	// 单列
	s1 := Select("id").From("tb_1").Build()
	// 多列
	s2 := Select("id", "name").From("tb_1").Build()
	// 全部
	s3 := Select().From("tb_1").Build()

	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println(s3)

	// output:
	// SELECT id FROM tb_1
	// SELECT id, name FROM tb_1
	// SELECT * FROM tb_1
}

func ExampleSelectBuilder_From() {
	// 单表
	s1 := Select().From("table_1 AS a").Build()
	// 多表
	s2 := Select().From("table_1", "table_2").Build()

	fmt.Println(s1)
	fmt.Println(s2)

	// output:
	// SELECT * FROM table_1 AS a
	// SELECT * FROM table_1, table_2
}

func ExampleSelectBuilder_Join() {
	// 内连接
	s1 := Select().From("tb_1").Join(
		InnerJoin("tb_2", "tb_1.id = tb_2.id"),
	).Build()
	// 左连接
	s2 := Select().From("tb_1").Join(
		LeftJoin("tb_2", "tb_1.id = tb_2.id"),
	).Build()
	// 右连接
	s3 := Select().From("tb_1").Join(
		RightJoin("tb_2", "tb_1.id = tb_2.id"),
	).Build()
	// 多多多表连接
	s4 := Select().From("tb_1").
		Join(InnerJoin("tb_2", "tb_2.id = tb_1.id")).
		Join(LeftJoin("tb_3", "tb_3.id = tb_1.id")).
		Join(RightJoin("tb_4", "tb_4.id = tb_1.id")).
		Build()
	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println(s3)
	fmt.Println(s4)

	// output:
	// SELECT * FROM tb_1 INNER JOIN tb_2 ON tb_1.id = tb_2.id
	// SELECT * FROM tb_1 LEFT JOIN tb_2 ON tb_1.id = tb_2.id
	// SELECT * FROM tb_1 RIGHT JOIN tb_2 ON tb_1.id = tb_2.id
	// SELECT * FROM tb_1 INNER JOIN tb_2 ON tb_2.id = tb_1.id LEFT JOIN tb_3 ON tb_3.id = tb_1.id RIGHT JOIN tb_4 ON tb_4.id = tb_1.id
}

func ExampleSelectBuilder_Where() {
	// 单条件
	s1 := Select().From("tb_1").Where(
		Eq("name1"),
	).Build()

	// 单条件的 AND 和 OR 不生效
	s2 := Select().From("tb_1").Where(
		NotEq("name1").Or(),
	).Build()

	// 多条件 AND ,第一位的 OR 设置不生效，默认是 AND
	s3 := Select().From("tb_1").Where(
		Eq("name1").Or(),
		Lt("name2").Or(),
	).Build()

	// 多条件 AND 和 OR
	s4 := Select().From("tb_1").Where(
		Eq("name1"),
		Lt("name2").Or(),
		LtEq("name3"),
		Gt("name4").Or(),
		GtEq("name5"),
		BetWeen("name6").Or(),
		NotBetWeen("name7"),
		IsNULL("name8").Or(),
		IsNotNULL("name9"),
		In("name10", 3).Or(),
		NotEq("name11"),
		NotIn("name12", 2).Or(),
		Like("name13"),
	).Build()

	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println(s3)
	fmt.Println(s4)

	// output:
	// SELECT * FROM tb_1 WHERE name1 = ?
	// SELECT * FROM tb_1 WHERE name1 != ?
	// SELECT * FROM tb_1 WHERE name1 = ? OR name2 < ?
	// SELECT * FROM tb_1 WHERE name1 = ? OR name2 < ? AND name3 <= ? OR name4 > ? AND name5 >= ? OR name6 BETWEEN ? AND ? AND name7 NOT BETWEEN ? AND ? OR name8 IS NULL AND name9 IS NOT NULL OR name10 IN ( ?, ?, ?) AND name11 != ? OR name12 NOT IN ( ?, ?) AND name13 LIKE ?
}

func ExampleSelectBuilder_GroupBy() {
	// 单列
	s1 := Select().From("tb_1").GroupBy("col1").Build()
	// 多列
	s2 := Select().From("tb_1").GroupBy("col1", "col2").Build()
	// 更多列
	s3 := Select().From("tb_1").GroupBy("col1", "col2", "col3", "col4").Build()

	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println(s3)

	// output:
	// SELECT * FROM tb_1 GROUP BY col1
	// SELECT * FROM tb_1 GROUP BY col1, col2
	// SELECT * FROM tb_1 GROUP BY col1, col2, col3, col4
}

func ExampleSelectBuilder_Having() {
	// 这里的逻辑其实与WHERE的一致

	// 单条件
	s1 := Select().From("tb_1").Having(
		Eq("name1"),
	).Build()

	// 单条件的 AND 和 OR 不生效
	s2 := Select().From("tb_1").Having(
		NotEq("name1").Or(),
	).Build()

	// 多条件 AND ,第一位的 OR 设置不生效，默认是 AND
	s3 := Select().From("tb_1").Having(
		Eq("name1").Or(),
		Lt("name2").Or(),
	).Build()

	// 多条件 AND 和 OR
	s4 := Select().From("tb_1").Having(
		Eq("name1"),
		Lt("name2").Or(),
		LtEq("name3"),
		Gt("name4").Or(),
		GtEq("name5"),
		BetWeen("name6").Or(),
		NotBetWeen("name7"),
		IsNULL("name8").Or(),
		IsNotNULL("name9"),
		In("name10", 3).Or(),
		NotEq("name11"),
		NotIn("name12", 2).Or(),
		Like("name13"),
	).Build()

	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println(s3)
	fmt.Println(s4)

	// output:
	// SELECT * FROM tb_1 HAVING name1 = ?
	// SELECT * FROM tb_1 HAVING name1 != ?
	// SELECT * FROM tb_1 HAVING name1 = ? OR name2 < ?
	// SELECT * FROM tb_1 HAVING name1 = ? OR name2 < ? AND name3 <= ? OR name4 > ? AND name5 >= ? OR name6 BETWEEN ? AND ? AND name7 NOT BETWEEN ? AND ? OR name8 IS NULL AND name9 IS NOT NULL OR name10 IN ( ?, ?, ?) AND name11 != ? OR name12 NOT IN ( ?, ?) AND name13 LIKE ?
}

func ExampleSelectBuilder_OrderBy() {
	// 单列
	s1 := Select().From("tb_1").OrderBy(
		OrderByAsc("col1"),
	).Build()
	// 多列
	s2 := Select().From("tb_1").OrderBy(
		OrderByAsc("col1"),
		OrderByDesc("col2"),
	).Build()
	// 更多列
	s3 := Select().From("tb_1").OrderBy(
		OrderByAsc("col1"),
		OrderByDesc("col2"),
		OrderByAsc("col3"),
		OrderByDesc("col4"),
	).Build()

	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println(s3)

	// output:
	// SELECT * FROM tb_1 ORDER BY col1 ASC
	// SELECT * FROM tb_1 ORDER BY col1 ASC, col2 DESC
	// SELECT * FROM tb_1 ORDER BY col1 ASC, col2 DESC, col3 ASC, col4 DESC
}

func ExampleSelectBuilder_Limit() {

	s1 := Select().From("tb_1").Limit(10).Build()
	fmt.Println(s1)

	// output:
	// SELECT * FROM tb_1 LIMIT 10
}

func ExampleSelectBuilder_Offset() {

	s1 := Select().From("tb_1").Offset(10).Build()
	fmt.Println(s1)

	// output:
	// SELECT * FROM tb_1 OFFSET 10
}

func ExampleSelectBuilder_Union() {
	// no values union
	s1 := Select().From("tb_1").Union(Select().From("tb_1").Build(), nil).Build()

	// values union
	uni1, val1 := Select().From("tb_1").Where(Eq("name1", "name1")).BuildWithValue()
	uni2, val2 := Select().From("tb_2").Where(Eq("name2", "name2")).Union(uni1, val1).BuildWithValue()

	fmt.Println(s1)
	fmt.Println(uni2)
	fmt.Printf("val:%v\n", val2)

	// output:
	// SELECT * FROM tb_1 UNION SELECT * FROM tb_1
	// SELECT * FROM tb_2 WHERE name2 = ? UNION SELECT * FROM tb_1 WHERE name1 = ?
	// val:[name2 name1]
}

func ExampleSelectBuilder_UnionAll() {
	// no values union
	s1 := Select().From("tb_1").UnionAll(Select().From("tb_1").Build(), nil).Build()

	// values union
	uni1, val1 := Select().From("tb_1").Where(Eq("name1", "name1")).BuildWithValue()
	uni2, val2 := Select().From("tb_2").Where(Eq("name2", "name2")).UnionAll(uni1, val1).BuildWithValue()

	fmt.Println(s1)
	fmt.Println(uni2)
	fmt.Printf("val:%v\n", val2)

	// output:
	// SELECT * FROM tb_1 UNION ALL SELECT * FROM tb_1
	// SELECT * FROM tb_2 WHERE name2 = ? UNION ALL SELECT * FROM tb_1 WHERE name1 = ?
	// val:[name2 name1]
}

func ExampleSelectBuilder_SubQuery() {
	// 设置一个子查询，并且要定义别名
	sub := Select().From("tb_1").SubQuery("s")
	fmt.Println(sub)

	s1 := Select().From("tb_2").Join(
		InnerJoin(sub, "tb_2.id = s.id"),
	).Build()

	fmt.Println(s1)

	// output:
	// (SELECT * FROM tb_1) AS s
	// SELECT * FROM tb_2 INNER JOIN (SELECT * FROM tb_1) AS s ON tb_2.id = s.id
}
