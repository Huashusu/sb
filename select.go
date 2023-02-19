package sb

import (
	"bytes"
	"strconv"
	"sync"
)

var wherePool = &sync.Pool{
	New: func() any {
		return new(WhereBuf)
	},
}

func getWhere() *WhereBuf {
	w := wherePool.Get().(*WhereBuf)
	w.reset()
	return w
}

func putWhere(buf *WhereBuf) {
	wherePool.Put(buf)
}

var joinPool = &sync.Pool{
	New: func() any {
		return new(JoinBuf)
	},
}

func getJoin() *JoinBuf {
	j := joinPool.Get().(*JoinBuf)
	j.buf.Reset()
	j.val = nil
	return j
}

func putJoin(buf *JoinBuf) {
	joinPool.Put(buf)
}

// WhereBuf 只是创建where条件，加个占位符，暂时不提供附带值的方法，敬请期待。而且带值方法，还需要安全性校验、注入校验
type WhereBuf struct {
	ty  int
	buf bytes.Buffer
	val []any
}

func (w *WhereBuf) Byte() []byte {
	defer putWhere(w)
	return w.buf.Next(w.buf.Len())
}

const (
	AND = iota
	OR
)

// Set 手动设置条件
func (w *WhereBuf) Set(s string, val ...any) *WhereBuf {
	w.buf.WriteString(s)
	w.val = val
	return w
}

func (w *WhereBuf) reset() {
	w.buf.Reset()
	w.ty = AND
	w.val = nil
}

// And 设置这个条件是and
func (w *WhereBuf) And() *WhereBuf {
	w.ty = AND
	return w
}

// Or 设置这个条件是or
func (w *WhereBuf) Or() *WhereBuf {
	w.ty = OR
	return w
}

// IsNULL 空值
func IsNULL(col string) *WhereBuf {
	w := getWhere()
	w.buf.WriteString(col)
	w.buf.WriteString(" IS NULL")
	return w
}

// IsNotNULL 非空值
func IsNotNULL(col string) *WhereBuf {
	w := getWhere()
	w.buf.WriteString(col)
	w.buf.WriteString(" IS NOT NULL")
	return w
}

// Eq 等于
func Eq(col string, val ...any) *WhereBuf {
	w := getWhere()
	w.buf.WriteString(col)
	w.buf.WriteString(" = ?")
	w.val = val
	return w
}

// NotEq 不等于
func NotEq(col string, val ...any) *WhereBuf {
	w := getWhere()
	w.buf.WriteString(col)
	w.buf.WriteString(" != ?")
	w.val = val
	return w
}

// Gt 大于
func Gt(col string, val ...any) *WhereBuf {
	w := getWhere()
	w.buf.WriteString(col)
	w.buf.WriteString(" > ?")
	w.val = val
	return w
}

// GtEq 大于等于
func GtEq(col string, val ...any) *WhereBuf {
	w := getWhere()
	w.buf.WriteString(col)
	w.buf.WriteString(" >= ?")
	w.val = val
	return w
}

// Lt 小于
func Lt(col string, val ...any) *WhereBuf {
	w := getWhere()
	w.buf.WriteString(col)
	w.buf.WriteString(" < ?")
	w.val = val
	return w
}

// LtEq 小于等于
func LtEq(col string, val ...any) *WhereBuf {
	w := getWhere()
	w.buf.WriteString(col)
	w.buf.WriteString(" <= ?")
	w.val = val
	return w
}

// In 输入参数是要In几个值，这边给几个问号
func In(col string, n int, val ...any) *WhereBuf {
	w := getWhere()
	w.buf.WriteString(col)
	w.buf.WriteString(" IN (")
	for i := 0; i < n; i++ {
		w.buf.WriteString(" ?")
		if i < n-1 {
			w.buf.WriteString(",")
		}
	}
	w.buf.WriteString(")")
	w.val = val
	return w
}

// NotIn 输入参数是列名和值的数量，解释同 In 方法
func NotIn(col string, n int, val ...any) *WhereBuf {
	w := getWhere()
	w.buf.WriteString(col)
	w.buf.WriteString(" NOT IN (")
	for i := 0; i < n; i++ {
		w.buf.WriteString(" ?")
		if i < n-1 {
			w.buf.WriteString(",")
		}
	}
	w.buf.WriteString(")")
	w.val = val
	return w
}

// Like 对列进行like匹配
func Like(col string, val ...any) *WhereBuf {
	w := getWhere()
	w.buf.WriteString(col)
	w.buf.WriteString(" LIKE ?")
	w.val = val
	return w
}

// BetWeen 范围查询
func BetWeen(col string, val ...any) *WhereBuf {
	w := getWhere()
	w.buf.WriteString(col)
	w.buf.WriteString(" BETWEEN ? AND ?")
	w.val = val
	return w
}

// NotBetWeen 范围排除查询
func NotBetWeen(col string, val ...any) *WhereBuf {
	w := getWhere()
	w.buf.WriteString(col)
	w.buf.WriteString(" NOT BETWEEN ? AND ?")
	w.val = val
	return w
}

const (
	desc = "DESC"
	asc  = "ASC"
)

type OrderBuf struct {
	buf bytes.Buffer
}

func (o *OrderBuf) Byte() []byte {
	return o.buf.Bytes()
}

// OrderByAsc 对col升序排列
func OrderByAsc(col string) *OrderBuf {
	w := new(OrderBuf)
	w.buf.WriteString(col)
	w.buf.WriteByte(' ')
	w.buf.WriteString(asc)
	return w
}

// OrderByDesc 对col降序排列
func OrderByDesc(col string) *OrderBuf {
	w := new(OrderBuf)
	w.buf.WriteString(col)
	w.buf.WriteByte(' ')
	w.buf.WriteString(desc)
	return w
}

type SelectBuilder struct {
	limitStatus    bool
	offsetStatus   bool
	subQueryStatus bool
	unionStatus    bool
	unionALlStatus bool
	limit          int
	offset         int
	alias          string
	unionSQL       string
	join           *bytes.Buffer
	group          []string
	from           []string
	cols           []string
	where          []*WhereBuf
	having         []*WhereBuf
	order          []*OrderBuf
	values         []any
	unionVal       []any
	joinVal        []any
}

const (
	AllCol  = "*"
	growNum = 16
)

// Select 用于构建查询所需的Builder
func Select(Cols ...string) *SelectBuilder {
	return &SelectBuilder{
		cols:        Cols,
		join:        get(),
		limitStatus: false,
	}
}

// From 设置从哪张表查询，多表情况用Join添加
func (s *SelectBuilder) From(TableName ...string) *SelectBuilder {
	s.from = TableName
	return s
}

type JoinBuf struct {
	val []any
	buf bytes.Buffer
}

func (j *JoinBuf) Byte() []byte {
	defer putJoin(j)
	return j.buf.Bytes()
}

// InnerJoin 构建内连接buf
func InnerJoin(InnerTable, Condition string, where ...*WhereBuf) *JoinBuf {
	j := getJoin()
	j.buf.WriteString("INNER JOIN ")
	j.buf.WriteString(InnerTable)
	j.buf.WriteString(" ON ")
	j.buf.WriteString(Condition)
	if where != nil && len(where) > 0 {
		j.buf.WriteString(" ")
		for i := 0; i < len(where); i++ {
			if where[i].ty == AND {
				j.buf.WriteString(" AND ")
			} else {
				j.buf.WriteString(" OR ")
			}
			j.buf.Write(where[i].Byte())
			j.val = append(j.val, where[i].val...)
		}
	}
	return j
}

// LeftJoin 构建左连接buf
func LeftJoin(LeftTable, Condition string, where ...*WhereBuf) *JoinBuf {
	j := getJoin()
	j.buf.WriteString("LEFT JOIN ")
	j.buf.WriteString(LeftTable)
	j.buf.WriteString(" ON ")
	j.buf.WriteString(Condition)
	if where != nil && len(where) > 0 {
		j.buf.WriteString(" ")
		for i := 0; i < len(where); i++ {
			if where[i].ty == AND {
				j.buf.WriteString("AND ")
			} else {
				j.buf.WriteString("OR ")
			}
			j.buf.Write(where[i].Byte())
			j.val = append(j.val, where[i].val...)
		}
	}
	return j
}

// RightJoin 构建右连接buf
func RightJoin(RightTable, Condition string, where ...*WhereBuf) *JoinBuf {
	j := getJoin()
	j.buf.WriteString("RIGHT JOIN ")
	j.buf.WriteString(RightTable)
	j.buf.WriteString(" ON ")
	j.buf.WriteString(Condition)
	if where != nil && len(where) > 0 {
		j.buf.WriteString(" ")
		for i := 0; i < len(where); i++ {
			if where[i].ty == AND {
				j.buf.WriteString(" AND ")
			} else {
				j.buf.WriteString(" OR ")
			}
			j.buf.Write(where[i].Byte())
			j.val = append(j.val, where[i].val...)
		}
	}
	return j
}

// Join 添加连接，注意Join顺序
func (s *SelectBuilder) Join(join *JoinBuf, val ...any) *SelectBuilder {
	s.join.Grow(growNum)
	if s.join.Len() > 0 {
		s.join.WriteByte(' ')
	}
	s.join.Write(join.Byte())
	s.joinVal = append(s.joinVal, val...)
	s.joinVal = append(s.joinVal, join.val...)
	return s
}

// Where 添加where条件，提供了几个常用的where条件可供使用。
// 默认是AND连接条件，如果是OR请不要放第一个，原因：不能对条件排序，避免破坏顺序造成的潜在sql优化问题，
func (s *SelectBuilder) Where(cons ...*WhereBuf) *SelectBuilder {
	s.where = cons
	return s
}

// GroupBy 暂时只对单列分组，多列分组情况还需考虑
func (s *SelectBuilder) GroupBy(GroupCol ...string) *SelectBuilder {
	s.group = GroupCol
	return s
}

// Having 是对聚合函数后的值进行where判断，所以这里就直接拿WhereBuf那套方法了
func (s *SelectBuilder) Having(cons ...*WhereBuf) *SelectBuilder {
	s.having = cons
	return s
}

// OrderBy 多多个字段进行排序处理
func (s *SelectBuilder) OrderBy(orders ...*OrderBuf) *SelectBuilder {
	s.order = orders
	return s
}

// Limit 查询数量
func (s *SelectBuilder) Limit(limit int) *SelectBuilder {
	s.limitStatus = true
	s.limit = limit
	return s
}

// Offset 偏移量
func (s *SelectBuilder) Offset(offset int) *SelectBuilder {
	s.offsetStatus = true
	s.offset = offset
	return s
}

// SubQueryWithValue 设置自身为子查询，同时设置别名并返回sql字符串和value
func (s *SelectBuilder) SubQueryWithValue(as ...string) (string, []any) {
	s.subQueryStatus = true
	if as != nil && len(as) > 0 {
		s.alias = as[0]
	}
	return s.BuildWithValue()
}

// SubQuery 设置自身为子查询，同时设置别名并返回sql字符串
func (s *SelectBuilder) SubQuery(as ...string) string {
	s.subQueryStatus = true
	if as != nil && len(as) > 0 {
		s.alias = as[0]
	}
	str, _ := s.BuildWithValue()
	return str
}

// Union union SQL
func (s *SelectBuilder) Union(sql string, val []any) *SelectBuilder {
	s.unionStatus = true
	s.unionSQL = sql
	s.unionVal = val
	return s
}

// UnionAll union all SQL
func (s *SelectBuilder) UnionAll(sql string, val []any) *SelectBuilder {
	s.unionALlStatus = true
	return s.Union(sql, val)
}

// BuildWithValue 返回SQL字符串和传入的参数
func (s *SelectBuilder) BuildWithValue() (string, []any) {

	//global func param
	sep := " "
	i := 0

	// select string
	ss := get()
	defer put(ss)
	ss.Reset()
	ss.Grow(1 << 10)
	if s.subQueryStatus {
		ss.WriteByte('(')
	}
	ss.WriteString("SELECT ")

	// set column
	if s.cols != nil && len(s.cols) > 0 {
		for i = 0; i < len(s.cols); i++ {
			ss.WriteString(s.cols[i])
			if i < len(s.cols)-1 {
				ss.WriteString(", ")
			}
		}
	} else {
		ss.WriteString("*")
	}

	// set from table
	ss.WriteString(sep)
	ss.WriteString("FROM ")
	if s.from != nil && len(s.from) > 0 {
		for i = 0; i < len(s.from); i++ {
			ss.WriteString(s.from[i])
			if i < len(s.from)-1 {
				ss.WriteString(", ")
			}
		}
	}

	// set join table
	if s.join.Len() > 0 {
		ss.WriteString(sep)
		ss.WriteString(s.join.String())
		s.values = append(s.values, s.joinVal...)
	}

	// set where condition
	if s.where != nil && len(s.where) > 0 {
		ss.WriteString(sep)
		ss.WriteString("WHERE ")
		for i = 0; i < len(s.where); i++ {
			if i > 0 && i < len(s.where) {
				if s.where[i].ty == AND {
					ss.WriteString(" AND ")
				} else {
					ss.WriteString(" OR ")
				}
			}
			s.values = append(s.values, s.where[i].val...)
			ss.Write(s.where[i].Byte())
		}
	}

	// set group by
	if s.group != nil && len(s.group) > 0 {
		ss.WriteString(sep)
		ss.WriteString("GROUP BY ")
		for i = 0; i < len(s.group); i++ {
			if i > 0 && i < len(s.group) {
				ss.WriteString(", ")
			}
			ss.WriteString(s.group[i])
		}
	}

	// set having
	if s.having != nil && len(s.having) > 0 {
		ss.WriteString(sep)
		ss.WriteString("HAVING ")
		for i = 0; i < len(s.having); i++ {
			if i > 0 && i < len(s.having) {
				if s.having[i].ty == AND {
					ss.WriteString(" AND ")
				} else {
					ss.WriteString(" OR ")
				}
			}
			s.values = append(s.values, s.having[i].val...)
			ss.Write(s.having[i].Byte())
		}
	}

	// set order by
	if s.order != nil && len(s.order) > 0 {
		ss.WriteString(sep)
		ss.WriteString("ORDER BY ")
		for i = 0; i < len(s.order); i++ {
			ss.Write(s.order[i].Byte())
			if i < len(s.order)-1 {
				ss.WriteString(", ")
			}
		}
	}

	// set limit N,M
	if s.limitStatus || s.offsetStatus {
		if s.limitStatus {
			ss.WriteString(sep)
			ss.WriteString("LIMIT ")
			ss.WriteString(strconv.Itoa(s.limit))
		}
		if s.offsetStatus {
			ss.WriteString(sep)
			ss.WriteString("OFFSET ")
			ss.WriteString(strconv.Itoa(s.offset))
		}
	}

	if s.subQueryStatus {
		ss.WriteByte(')')
		if s.alias != "" {
			ss.WriteString(" AS ")
			ss.WriteString(s.alias)
		}
	}

	// union
	if s.unionStatus {
		// union all
		if s.unionALlStatus {
			ss.WriteString(" UNION ALL ")
		} else {
			ss.WriteString(" UNION ")
		}
		s.values = append(s.values, s.unionVal...)
		ss.WriteString(s.unionSQL)
	}

	return ss.String(), s.values
}

// Build 返回SQL字符串，其实是调用了BuildWithValue方法，忽略了val
func (s *SelectBuilder) Build() string {
	str, _ := s.BuildWithValue()
	return str
}
