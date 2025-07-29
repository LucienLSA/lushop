package model

import (
	"context"
	"goodssrv/global"
	"strconv"

	"gorm.io/gorm"
)

type Category struct {
	BaseModel
	Name             string      `gorm:"type:varchar(20);not null;comment:'商品分类名称'" json:"name"`
	ParentCategoryID int32       `json:"parent_category_id"`
	Level            int32       `gorm:"type:int;not null;default:1;comment:'1表示商品分类的等级'" json:"level"`
	IsTab            bool        `gorm:"default:false;not null;comment:'是否Tap栏显示'" json:"is_tab"`
	ParentCategory   *Category   `json:"-"`
	SubCategory      []*Category `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_category"`
}

// 强制指定表名为 "category"
func (Category) TableName() string {
	return "category"
}

type Brand struct {
	BaseModel
	Name string `gorm:"type:varchar(50);not null;comment:'品牌名称'"`
	Logo string `gorm:"type:varchar(200);default:'';not null;comment:'品牌Logo图片'"`
}

type GoodsCategoryBrand struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,unique;comment:'分类ID'"`
	Category   Category
	BrandID    int32 `gorm:"type:int;index:idx_category_brand,unique;comment:'品牌ID'"`
	Brand      Brand
}

func (GoodsCategoryBrand) TableName() string {
	return "goodscategorybrand"
}

type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null;comment:轮播图"`
	Url   string `gorm:"type:varchar(200);not null;comment:'图片链接'"`
	Index int32  `gorm:"type:int;default:1;not null;comment:'轮播图的索引'"`
}

type Goods struct {
	BaseModel
	CategoryID      int32 `gorm:"type:int;not null;comment:'商品分类ID'"`
	Category        Category
	BrandID         int32 `gorm:"type:int;not null"`
	Brand           Brand
	OnSale          bool     `gorm:"default:false;not null;comment:'是否特价'"`
	GoodsSn         string   `gorm:"type:varchar(50);not null;comment:'商品编号'"`
	Name            string   `gorm:"type:varchar(100);not null;comment:'商品名称'"`
	ClickNum        int32    `gorm:"type:int;default:0;not null;comment:'商品点击数'"`
	SoldNum         int32    `gorm:"type:int;default:0;not null;comment:'商品销量'"`
	FavNum          int32    `gorm:"type:int;default:0;not null;comment:'商品收藏数'"`
	MarketPrice     float32  `gorm:"not null;comment:'商品市场价'"`
	ShopPrice       float32  `gorm:"not null;comment:'商品实际价'"`
	GoodsBrief      string   `gorm:"type:varchar(100);not null;comment:'商品简介'"`
	ShipFree        bool     `gorm:"default:false;not null;comment:'是否免运费'"`
	Images          GormList `gorm:"type:varchar(1000);not null;comment:'商品图片'"`
	DescImages      GormList `gorm:"type:varchar(5000);not null;comment:'商品详情图片'"`
	GoodsFrontImage string   `gorm:"type:varchar(200);not null;comment:'商品封面图'"`
	IsNew           bool     `gorm:"default:false;not null;comment:'是否新品'"`
	IsHot           bool     `gorm:"default:false;not null;comment:'是否热卖'"`
}

// 钩子函数 创建完数据库之后进行同步到ES
func (g *Goods) AfterCreate(tx *gorm.DB) (err error) {
	esModel := EsGoods{
		ID:          g.ID,
		CategoryID:  g.CategoryID,
		BrandID:     g.BrandID,
		OnSale:      g.OnSale,
		ShipFree:    g.ShipFree,
		IsNew:       g.IsNew,
		IsHot:       g.IsHot,
		Name:        g.Name,
		ClickNum:    g.ClickNum,
		SoldNum:     g.SoldNum,
		FavNum:      g.FavNum,
		MarketPrice: g.MarketPrice,
		GoodsBrief:  g.GoodsBrief,
		ShopPrice:   g.ShopPrice,
	}
	_, err = global.EsClient.Index().Index(esModel.GetIndexName()).BodyJson(esModel).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// 更新数据库操作之后，同步更新到ES
func (g *Goods) AfterUpdate(tx *gorm.DB) (err error) {
	esModel := EsGoods{
		ID:          g.ID,
		CategoryID:  g.CategoryID,
		BrandID:     g.BrandID,
		OnSale:      g.OnSale,
		ShipFree:    g.ShipFree,
		IsNew:       g.IsNew,
		IsHot:       g.IsHot,
		Name:        g.Name,
		ClickNum:    g.ClickNum,
		SoldNum:     g.SoldNum,
		FavNum:      g.FavNum,
		MarketPrice: g.MarketPrice,
		GoodsBrief:  g.GoodsBrief,
		ShopPrice:   g.ShopPrice,
	}
	_, err = global.EsClient.Update().Index(esModel.GetIndexName()).Doc(esModel).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// 删除同步
func (g *Goods) AfterDelete(tx *gorm.DB) (err error) {
	_, err = global.EsClient.Delete().Index(EsGoods{}.GetIndexName()).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
