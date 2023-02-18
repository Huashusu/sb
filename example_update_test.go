package sb

import "fmt"

func ExampleUpdateBuilder_Build() {
	// 简简单单没有where条件
	const (
		id   = "id"
		name = "name"
	)

	s1 := Update("tb_1").Column(id, name).Build()
	fmt.Println(s1)

	// 简单上加一点where
	s2 := Update("tb_1").
		Column(id, name).
		Where(
			Eq(id).Or(),
			LtEq(name).Or(),
		).
		Build()
	fmt.Println(s2)

	// output:
	// UPDATE tb_1 SET `id` = ? , `name` = ?
	// UPDATE tb_1 SET `id` = ? , `name` = ? WHERE id = ? OR name <= ?
}

func ExampleUpdateBuilder_BuildWithValue() {
	const (
		id   = "id"
		name = "name"
	)

	// 简单带值sql
	s1, v1 := Update("tb_1").
		Column(id, name).
		Set("idVal", "nameVal").
		BuildWithValue()
	fmt.Printf("s1:%s\nv1:%v\n", s1, v1)

	// 整点复杂的，例如加上where的条件和值
	s2, v2 := Update("tb_1").
		Column(id, name).
		Set("idVal", "nameVal").
		Where(
			Eq("password", "passVal"),
			Eq("nickname", "nickVal"),
		).
		BuildWithValue()
	fmt.Printf("s2:%s\nv2:%v\n", s2, v2)

	// output:
	// s1:UPDATE tb_1 SET `id` = ? , `name` = ?
	// v1:[idVal nameVal]
	// s2:UPDATE tb_1 SET `id` = ? , `name` = ? WHERE password = ? AND nickname = ?
	// v2:[idVal nameVal passVal nickVal]
}
