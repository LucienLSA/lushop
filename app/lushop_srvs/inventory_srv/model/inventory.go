package model

// type Stock struct {
// 	BaseModel
// 	name    string
// 	Address string
// }

type Inventory struct {
	BaseModel
	Goods   int32 `gorm:"type:int;index;comment:商品id"`
	Stocks  int32 `gorm:"type:int;comment:仓库"`
	Version int32 `gorm:"type:int;comment:分布式锁-乐观锁"` // 分布式锁的乐观锁
}
type InventoryNew struct {
	BaseModel
	Goods   int32 `gorm:"type:int;index"`
	Stocks  int32 `gorm:"type:int"`
	Version int32 `gorm:"type:int"` //分布式锁的乐观锁
	Freeze  int32 `gorm:"type:int"` //冻结库存
}

type Delivery struct {
	Goods   int32  `gorm:"type:int;index"`
	Nums    int32  `gorm:"type:int"`
	OrderSn string `gorm:"type:varchar(200)"`
	Status  string `gorm:"type:varchar(200)"` // 1.代表等待支付，2.代表支付成功，3.支付失败
}

type StockSellDetail struct {
	BaseModel
	OrderSn string         `gorm:"type:varchar(200);index:idx_order_sn,unique;comment:订单编号"`
	Status  int32          `gorm:"type:varchar(200);comment:1.已扣减,2.已归还"` // 1.代表已扣减，2.代表已归还，3.失败
	Detail  GormDetailList `gorm:"type:varchar(200);comment:详细商品"`
}

func (StockSellDetail) TableName() string {
	return "stockselldetail"
}
