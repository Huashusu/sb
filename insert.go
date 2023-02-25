package sb

type InsertBuilder struct {
	table string
	cols  []string
	vals  []any
}

// Insert 构建insert语句
func Insert(TableName string) *InsertBuilder {
	i := new(InsertBuilder)
	i.table = TableName
	return i
}

// Column 设置插入的列名
func (i *InsertBuilder) Column(cols ...string) *InsertBuilder {
	if i.cols == nil {
		i.cols = cols
	} else if cols != nil {
		i.cols = append(i.cols, cols...)
	}
	return i
}

// Set 设置数据
func (i *InsertBuilder) Set(vals ...any) *InsertBuilder {
	if i.vals == nil {
		i.vals = make([]any, 0, len(vals))
	}
	i.vals = append(i.vals, vals...)
	return i
}

// SetSlice 批量设置数据
func (i *InsertBuilder) SetSlice(vals [][]any) *InsertBuilder {
	for j := 0; j < len(vals); j++ {
		i.Set(vals[j]...)
	}
	return i
}

// BuildWithValue 返回SQL语句和传入的参数
func (i *InsertBuilder) BuildWithValue() (string, []any) {

	if i.cols == nil {
		return "", nil
	}

	// 判断列数与值数是否匹配
	if len(i.cols) > 0 && len(i.vals)%len(i.cols) != 0 {
		return "", nil
	}

	ss := get()
	defer put(ss)
	ss.Reset()
	ss.WriteString("INSERT INTO `")
	ss.WriteString(i.table)
	ss.WriteString("` ( ")
	for j := 0; j < len(i.cols); j++ {
		ss.WriteByte('`')
		ss.WriteString(i.cols[j])
		ss.WriteByte('`')
		if j < len(i.cols)-1 {
			ss.WriteString(", ")
		}
	}
	ss.WriteString(" ) ")
	count := len(i.vals) / len(i.cols)
	if count >= 2 {
		ss.WriteString("VALUES ")
		for j := 0; j < count; j++ {
			ss.WriteString("( ")
			for g := 0; g < len(i.cols); g++ {
				ss.WriteString("? ")
				if g < len(i.cols)-1 {
					ss.WriteString(", ")
				}
			}
			ss.WriteString(")")
			if j < count-1 {
				ss.WriteString(",")
			}
		}
	} else {
		ss.WriteString("VALUE ( ")
		for g := 0; g < len(i.cols); g++ {
			ss.WriteString("? ")
			if g < len(i.cols)-1 {
				ss.WriteString(", ")
			}
		}
		ss.WriteString(")")
	}
	return ss.String(), i.vals
}

// Build 返回SQL语句
func (i *InsertBuilder) Build() string {
	str, _ := i.BuildWithValue()
	return str
}
