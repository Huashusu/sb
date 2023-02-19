package sb

type DeleteBuilder struct {
	table  string
	where  []*WhereBuf
	values []any
}

func Delete(TableName string) *DeleteBuilder {
	return &DeleteBuilder{
		table: TableName,
	}
}

func (d *DeleteBuilder) Where(cons ...*WhereBuf) *DeleteBuilder {
	d.where = cons
	return d
}

func (d *DeleteBuilder) BuildWithValue() (string, []any) {
	ss := get()
	defer put(ss)
	ss.Reset()
	ss.WriteString("DELETE FROM ")
	ss.WriteByte('`')
	ss.WriteString(d.table)
	ss.WriteByte('`')
	if d.where != nil && len(d.where) > 0 {
		ss.WriteString(" WHERE ")
		for i := 0; i < len(d.where); i++ {
			if i > 0 && i < len(d.where) {
				if d.where[i].ty == AND {
					ss.WriteString(" AND ")
				} else {
					ss.WriteString(" OR ")
				}
			}
			d.values = append(d.values, d.where[i].val...)
			ss.Write(d.where[i].Byte())
		}
	}
	return ss.String(), d.values
}

func (d *DeleteBuilder) Build() string {
	s, _ := d.BuildWithValue()
	return s
}
