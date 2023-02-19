# Sql Builder 

用来生成SQL语句，非ORM项目可以考虑用下试试，搭配[sqlx](https://github.com/jmoiron/sqlx#sqlx)一些库，快速将结果集`Scan`到结构体也挺不错的。

支持的数据库类型：

- [x] MySQL

- [ ] PostgreSQL

  ...

## Install

```bash
go get github.com/huashusu/sb@v0.1.1
```

## Benchmark Test

平台：windows 64位

CPU：AMD Ryzen 5 4600H with Radeon Graphics

| Benchmark Name |  |      | ||
| -------------- | :--: | :--: |:--:|:--:|
|BenchmarkDeleteSimple-12|10306125|113.9 ns/op|88 B/op|2 allocs/op|
|BenchmarkDeleteComplexWithValue-12|2509845|481.8 ns/op|208 B/op|7 allocs/op|
|BenchmarkInsertSimpleWithValue-12|3310828|354.2 ns/op|240 B/op|4 allocs/op|
|BenchmarkInsertComplexWithValue-12|1541556|754.8 ns/op|544 B/op|6 allocs/op|
|BenchmarkSelectSimpleWithValue-12|1341458|880.4 ns/op|1465 B/op|8 allocs/op|
|BenchmarkSelectComplexWithValue-12|331906|3771 ns/op|3929 B/op|34 allocs/op|
|BenchmarkUpdateSimpleWithValue-12|1731949|684.1 ns/op|384 B/op|8 allocs/op|
|BenchmarkUpdateComplexWithValue-12|1762831|675.7 ns/op|384 B/op|8 allocs/op|

## Examples

更多样例可以查看test文件和example文件：

- Insert： [example_insert_test](https://github.com/huashusu/sb/blob/master/example_insert_test.go)和[insert_test](https://github.com/huashusu/sb/blob/master/insert_test.go)
- Update： [example_update_test](https://github.com/huashusu/sb/blob/master/example_update_test.go)和[update_test](https://github.com/huashusu/sb/blob/master/update_test.go)
- Delete： [example_delete_test](https://github.com/huashusu/sb/blob/master/example_delete_test.go)和[delete_test](https://github.com/huashusu/sb/blob/master/delete_test.go)
- Select： [example_select_test](https://github.com/huashusu/sb/blob/master/example_select_test.go)和[select_test](https://github.com/huashusu/sb/blob/master/select_test.go)

```go
package main

import (
	"fmt"

	"github.com/huashusu/sb"
)

func main() {
	// 只生成sql的情况，传入的值什么的会抛弃掉,Insert()、Update()、Delete()、Select()都是一样的操作
	sql := sb.Insert("tb_user").
		Column("id", "name", "pass").
		// 会检测传入列数和值数量是否一致还有是否倍数关系
		Set("idVal", "nameVal", "passVal").
		Build()
	fmt.Printf("sql:%s\n", sql)

	// 生成带值的插入语句
	i1, vi1 := sb.Insert("tb_user").
		Column("id", "name", "pass").
		// 会检测传入列数和值数量是否一致还有是否倍数关系
		Set("idVal", "nameVal", "passVal").
		BuildWithValue()
	fmt.Printf("i1:%s vi1:%v\n", i1, vi1)
	// INSERT INTO `tb_user` ( `id`, `name`, `pass` ) VALUE ( ? , ? , ? ) vi1:[idVal nameVal passVal]

	// 批量插入的实现
	i2, vi2 := sb.Insert("tb_user").
		Column("id", "name", "pass").
		// 传入二维数组，会自动变成values和增加相应的占位符
		// 这里的值与返回的值呢，就那什么...展开了
		SetSlice([][]any{
			{"id1", "name1", "pass1"},
			{"id2", "name2", "pass2"},
			{"id3", "name3", "pass3"},
		}).
		BuildWithValue()
	fmt.Printf("i2:%s vi2:%v\n", i2, vi2)
	// i2:INSERT INTO `tb_user` ( `id`, `name`, `pass` ) VALUES ( ? , ? , ? ),( ? , ? , ? ),( ? , ? , ? ) vi2:[id1 name1 pass1 id2 name2 pass2 id3 name3 pass3]

	// 简单带值sql
	u1, vu1 := sb.Update("tb_1").
		Column("id", "name").
		Set("idVal", "nameVal").
		BuildWithValue()
	fmt.Printf("u1:%s vu1:%v\n", u1, vu1)
	// u1:UPDATE tb_1 SET `id` = ? , `name` = ? vu1:[idVal nameVal]

	// 整点复杂的，例如加上where的条件和值
	u2, vu2 := sb.Update("tb_1").
		Column("id", "name").
		Set("idVal", "nameVal").
		Where(
			sb.Eq("password", "passVal"),
			sb.Eq("nickname", "nickVal"),
		).
		BuildWithValue()
	fmt.Printf("u2:%s vu2:%v\n", u2, vu2)
	// u2:UPDATE tb_1 SET `id` = ? , `name` = ? WHERE password = ? AND nickname = ? vu2:[idVal nameVal passVal nickVal]

	// 简单带值删除sql
	d1, vd1 := sb.Delete("tb_1").Where(sb.Eq("id", "idVal")).BuildWithValue()
	fmt.Printf("d1:%s vd1:%v\n", d1, vd1)
	// d1:DELETE FROM `tb_1` WHERE id = ? vd1:[idVal]

	// 删除多来几个WHERE？
	d2, vd2 := sb.Delete("tb_1").
		Where(
			sb.Eq("id", "idVal"),
			sb.Gt("time", "timeVal"),
		).
		BuildWithValue()
	fmt.Printf("d2:%s vd2:%v\n", d2, vd2)
	// d2:DELETE FROM `tb_1` WHERE id = ? AND time > ? vd2:[idVal timeVal]

	// 用子查询来作为删除条件
	subSql, subVal := sb.Select("id").
		From("tb_1").
		Where(
			sb.Gt("create_time", "2006-01-02 15:04:05"),
		).SubQueryWithValue()
	subStr := new(sb.WhereBuf).Set("id IN "+subSql, subVal...)
	d3, vd3 := sb.Delete("tb_1").Where(subStr).BuildWithValue()
	fmt.Printf("d3:%s vd3:%v\n", d3, vd3)
	// d3:DELETE FROM `tb_1` WHERE id IN (SELECT id FROM tb_1 WHERE create_time > ?) vd3:[2006-01-02 15:04:05]

	// 生成普通的查询语句带like条件和值的
	s1, vs1 := sb.
		Select("id", "phone", "nick_name").
		From("tb_user").
		Where(sb.Like("nick_name", "小%")).
		Limit(10).
		Offset(0).
		// Build 和 BuildWithValue 的区别就是需不需要val的问题
		// Build 只返回string，不会有额外传入的值。
		// BuildWithValue 返回string和额外传入的值
		BuildWithValue()
	fmt.Printf("s1:%s vs1:%v\n", s1, vs1)
	// s1:SELECT id, phone, nick_name FROM tb_user WHERE nick_name LIKE ? LIMIT 10 OFFSET 0 v1:[小%]

	// 好了，现在来尝试一下复杂Select包含Join和子查询的
	// 先定义子查询的语句和值
	subSql, subVal = sb.Select("id").
		From("tb_user").
		Where(sb.Eq("name", "name1")).
		SubQueryWithValue("u2")
	s2, vs2 := sb.Select().
		From("tb_user AS u1").
		Join(
			sb.LeftJoin(subSql, "u1.id = u2.id", sb.Gt("u1.id", "id1")), // 这里是join的语句
			subVal..., // 这里是join的带值问题，记得展开数组
		).
		// 多Join语句要多.Join()
		Join(
			// 不需要带值也可以这样写
			sb.RightJoin(
				sb.Select("id").
					From("tb_user").
					SubQuery("u3"),
				"u3.id = u2.id"),
		).
		Where(
			sb.Gt("u2.id", "id2"),
			sb.Lt("u1.create_time", "time1"),
			sb.IsNotNULL("u1.nickname"),
			sb.IsNULL("u1.email").Or(),
		).
		GroupBy("u1.id").
		OrderBy(
			sb.OrderByAsc("u1.id"),
			sb.OrderByDesc("u2.create_time"),
		).
		Having(
			sb.GtEq("u1.id", "id3"),
		).
		Limit(100).
		Offset(10).
		BuildWithValue()
	// Tips:没有实际案例参考，只是往复杂了写
	fmt.Printf("s2:%s vs2:%v", s2, vs2)
	// SELECT * FROM tb_user AS u1
	// LEFT JOIN (SELECT id FROM tb_user WHERE name = ?) AS u2 ON u1.id = u2.id AND u1.id > ?
	// RIGHT JOIN (SELECT id FROM tb_user) AS u3 ON u3.id = u2.id
	// WHERE u2.id > ? AND u1.create_time < ? AND u1.nickname IS NOT NULL OR u1.email IS NULL
	// GROUP BY u1.id HAVING u1.id >= ?
	// ORDER BY u1.id ASC, u2.create_time DESC
	// LIMIT 100
	// OFFSET 10
	// v2:[name1 id1 id2 time1 id3]
	// Tips:生成样例，注意SubQuery与SubQueryWithValue的区别和Build与BuildWithValue的区别一样
}

```



### 后续计划：

1. 优化，感觉可以上Pool优化`BuildWithValue`方法
2. ...
