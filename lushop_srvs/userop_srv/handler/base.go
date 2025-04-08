package handler

import (
	proto_address "useropsrv/proto/gen/address"
	proto_message "useropsrv/proto/gen/message"
	proto_userfav "useropsrv/proto/gen/userfav"

	"gorm.io/gorm"
)

type UserOpServer struct {
	proto_address.UnimplementedAddressServer
	proto_message.UnimplementedMessageServer
	proto_userfav.UnimplementedUserFavServer
}

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
