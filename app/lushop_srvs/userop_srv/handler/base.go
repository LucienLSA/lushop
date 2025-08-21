package handler

import (
	proto "useropsrv/proto"

	"gorm.io/gorm"
)

type UserOpServer struct {
	// proto_address.UnimplementedAddressServer
	// proto_message.UnimplementedMessageServer
	// proto_userfav.UnimplementedUserFavServer
	proto.UnimplementedUserOpServer
}

func Paginate(pageNum, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// db = global.DB
		if pageNum < 1 {
			pageNum = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize < 1:
			pageSize = 10
		}
		offset := (pageNum - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
