package forms

type GoodsForm struct {
	Name        string   `form:"name" json:"name" binding:"required,min=2,max=100"`
	GoodsSn     string   `form:"goods_sn" json:"goods_sn" binding:"required,min=2,lt=20"`
	Stocks      int32    `form:"stocks" json:"stocks" binding:"required,min=1"`
	CategoryId  int32    `form:"category" json:"category" binding:"required"`
	MarketPrice float32  `form:"market_price" json:"market_price" binding:"required,min=0"`
	ShopPrice   float32  `form:"shop_price" json:"shop_price" binding:"required,min=0"`
	GoodsBrief  string   `form:"goods_brief" json:"goods_brief" binding:"required,min=3"`
	Images      []string `form:"images" json:"images" binding:"required,min=1"`
	DescImages  []string `form:"desc_images" json:"desc_images" binding:"required,min=1"`
	//表单认证bool类型需要用指针 版本问题  以后可能会修复
	ShipFree   *bool  `form:"ship_free" json:"ship_free" binding:"required"`
	FrontImage string `form:"front_image" json:"front_image" binding:"required,url"`
	Brand      int32  `form:"brand" json:"brand" binding:"required"`
}

type GoodsStatusForm struct {
	//表单认证bool类型需要用指针 版本问题  以后可能会修复
	IsNew  *bool `form:"new" json:"new" binding:"required"`
	IsHot  *bool `form:"hot" json:"hot" binding:"required"`
	OnSale *bool `form:"sale" json:"sale" binding:"required"`
}
