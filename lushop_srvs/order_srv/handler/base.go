package handler

import "gorm.io/gorm"

func Paginate(pageNum, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// db = global.DB
		if pageNum < 1 {
			pageNum = 1
		}
		switch {
		case pageSize > 10:
			pageSize = 10
		case pageSize < 1:
			pageSize = 5
		}
		offset := (pageNum - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
