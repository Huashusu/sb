package sb

import "fmt"

func ExampleInsertBuilder_BuildWithValue() {
	// 带入参的SQL生成
	s1, v1 := Insert("tb_1").
		Column("id", "name").
		Set("1", "2").
		BuildWithValue()
	fmt.Println(s1, v1)

	s2, v2 := Insert("tb_1").
		Column("id", "name").
		Set("1", "2", "3", "4", "5").
		BuildWithValue()
	// 列数与传入的参数不匹配的情况
	fmt.Printf("sql:%#v val is nil:%#v\n", s2, v2 == nil)

	// output:
	// INSERT INTO `tb_1` ( `id`, `name` ) VALUE ( ? , ? ) [1 2]
	// sql:"" val is nil:true
}

func ExampleInsertBuilder_Build() {
	// 仅生成SQL，不设置值
	s1 := Insert("tb_1").
		Column("id", "name", "salt").
		Build()

	fmt.Println(s1)

	// 非常逆天的情况下，就是不设置插入列名，会返回空串
	s2 := Insert("tb_1").
		Build()

	fmt.Printf("no column insert sql:%#v\n", s2)

	// output:
	// INSERT INTO `tb_1` ( `id`, `name`, `salt` ) VALUE ( ? , ? , ? )
	// no column insert sql:""
}
