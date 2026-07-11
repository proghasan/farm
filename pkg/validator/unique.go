package validator

func (f *Form) Unique(field string, params ...string) *Form {
	if len(params) < 2 {
		return f
	}
	table, column := params[0], params[1]

	value, ok := f.data[field]
	if !ok {
		return f
	}

	var count int64
	tx := f.db.Table(table).Where(column+" = ?", value)
	if len(params) > 2 {
		tx = tx.Where("id != ?", params[2])
	}

	if err := tx.Count(&count).Error; err != nil {
		return f
	}

	if count > 0 {
		f.addError(field, "The "+field+" has already been taken.")
	}

	return f
}
