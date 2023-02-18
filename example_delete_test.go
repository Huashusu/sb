package sb

import "fmt"

func ExampleDeleteBuilder_Build() {
	//简简单单的SQL，但是delete语句不加where会删全表
	s1 := Delete("tb_1").Build()
	fmt.Println(s1)

	// 附加上where条件的
	s2 := Delete("tb_1").Where(
		Eq("id"),
	).Build()
	fmt.Println(s2)

	// output:
	// DELETE FROM `tb_1`
	// DELETE FROM `tb_1` WHERE id = ?
}

func ExampleDeleteBuilder_BuildWithValue() {
	// 简单带值sql
	s1, v1 := Delete("tb_1").Where(Eq("id", "idVal")).BuildWithValue()
	fmt.Printf("s1:%s v1:%v\n", s1, v1)

	// 多来几个WHERE？
	s2, v2 := Delete("tb_1").
		Where(
			Eq("id", "idVal"),
			Gt("time", "timeVal"),
		).
		BuildWithValue()
	fmt.Printf("s2:%s v2:%v\n", s2, v2)

	// 用子查询来作为删除条件
	subSql, subVal := Select("id").
		From("tb_1").
		Where(
			Gt("create_time", "2006-01-02 15:04:05"),
		).SubQueryWithValue()
	subStr := new(WhereBuf).Set("id IN "+subSql, subVal...)

	s3, v3 := Delete("tb_1").Where(subStr).BuildWithValue()
	fmt.Printf("s3:%s v3:%v\n", s3, v3)

	// output:
	// s1:DELETE FROM `tb_1` WHERE id = ? v1:[idVal]
	// s2:DELETE FROM `tb_1` WHERE id = ? AND time > ? v2:[idVal timeVal]
	// s3:DELETE FROM `tb_1` WHERE id IN (SELECT id FROM tb_1 WHERE create_time > ?) v3:[2006-01-02 15:04:05]
}
