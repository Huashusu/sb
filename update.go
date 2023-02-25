package sb

import "strings"

type UpdateBuilder struct {
	table  string
	cols   []string
	values []any
	where  []*WhereBuf
}

func Update(TableName string) *UpdateBuilder {
	return &UpdateBuilder{
		table: TableName,
	}
}

func (u *UpdateBuilder) Column(cols ...string) *UpdateBuilder {
	if u.cols == nil {
		u.cols = cols
	} else if cols != nil {
		u.cols = append(u.cols, cols...)
	}
	return u
}

func (u *UpdateBuilder) Set(vals ...any) *UpdateBuilder {
	if u.values == nil {
		u.values = vals
	} else if vals != nil {
		u.values = append(u.values, vals...)
	}
	return u
}

func (u *UpdateBuilder) Where(cons ...*WhereBuf) *UpdateBuilder {
	if u.where == nil {
		u.where = cons
	} else if cons != nil {
		u.where = append(u.where, cons...)
	}
	return u
}

func (u *UpdateBuilder) BuildWithValue() (string, []any) {
	count := 0
	ss := get()
	defer put(ss)
	ss.Reset()
	ss.WriteString("UPDATE ")
	ss.WriteString(u.table)
	ss.WriteString(" SET ")
	for i := 0; i < len(u.cols); i++ {
		// 判断案例来源：col = col + 1
		if strings.Index(u.cols[i], "=") > 0 {
			ss.WriteString(u.cols[i])
			count++
		} else {
			ss.WriteByte('`')
			ss.WriteString(u.cols[i])
			ss.WriteByte('`')
			ss.WriteString(" = ?")
		}
		if i < len(u.cols)-1 {
			ss.WriteString(" ")
			ss.WriteString(", ")
		}
	}
	if u.where != nil && len(u.where) > 0 {
		ss.WriteString(" WHERE ")
		for i := 0; i < len(u.where); i++ {
			if i > 0 && i < len(u.where) {
				if u.where[i].ty == AND {
					ss.WriteString(" AND ")
				} else {
					ss.WriteString(" OR ")
				}
			}
			u.values = append(u.values, u.where[i].val...)
			ss.Write(u.where[i].Byte())
		}
	}

	return ss.String(), u.values
}

func (u *UpdateBuilder) Build() string {
	s, _ := u.BuildWithValue()
	return s
}
