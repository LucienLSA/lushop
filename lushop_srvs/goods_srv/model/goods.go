package model

type Category struct {
	BaseModel
	Name             string `gorm:"type:varchar(20);not null;comment:'商品分类名称'" json:"name"`
	Level            int32  `gorm:"type:int;not null;default:1;comment:'1表示商品分类的等级'" json:"level"`
	ParentCategoryID int32  `json:"-"`
	ParentCategory   *Category
	SubCategory      []*Category `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_category"`
	IsTab            bool        `gorm:"default:false;not null;comment:'是否Tap栏显示'" json:"is_tab"`
}

type Brands struct {
	BaseModel
	Name string `gorm:"type:varchar(50);not null;comment:'品牌名称'"`
	Logo string `gorm:"type:varchar(200);default:'';not null;comment:'品牌Logo图片'"`
}

type GoodsCategroyBrand struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Category   Category
	BrandsID   int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Brands     Brands
}

func (GoodsCategroyBrand) TableName() string {
	return "goodscategorybrand"
}

type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null"`
	Url   string `gorm:"type:varchar(200);not null;comment:'图片链接'"`
	Index int32  `gorm:"type:int;default:1;not null;comment:'轮播图的索引'"`
}

type Goods struct {
	BaseModel
	CategoryID      int32 `gorm:"type:int;not null"`
	Category        Category
	BrandsID        int32 `gorm:"type:int;not null"`
	Brands          Brands
	OnSale          bool     `gorm:"default:false;not null;comment:'是否特价'"`
	ShipFree        bool     `gorm:"default:false;not null;comment:'是否免运费'"`
	IsNew           bool     `gorm:"default:false;not null;comment:'是否新品'"`
	IsHot           bool     `gorm:"default:false;not null;comment:'是否热卖'"`
	Name            string   `gorm:"type:varchar(100);not null;comment:'商品名称'"`
	GoodsSn         string   `gorm:"type:varchar(50);not null;comment:'商品编号'"`
	ClickNum        int32    `gorm:"type:int;default:0;not null;comment:'商品点击数'"`
	SoldNum         int32    `gorm:"type:int;default:0;not null;comment:'商品销量'"`
	FavNum          int32    `gorm:"type:int;default:0;not null;comment:'商品收藏数'"`
	MarketPrice     float32  `gorm:"not null;comment:'商品市场价'"`
	ShopPrice       float32  `gorm:"not null;comment:'商品实际价'"`
	GoodsBrief      string   `gorm:"type:varchar(100);not null;comment:'商品简介'"`
	Images          GormList `gorm:"type:varchar(1000);not null;comment:'商品图片'"`
	DescImages      GormList `gorm:"type:varchar(5000);not null;comment:'商品详情图片'"`
	GoodsFrontImage string   `gorm:"type:varchar(200);not null;comment:'商品封面图'"`
}
