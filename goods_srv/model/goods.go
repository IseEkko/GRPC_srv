package model

//分类信息表
//字段是不是能为null这个是很重要的
//开发的过程中尽量设置为不为null
//这些类型我们使用int32还是int，我们尽量的使用int类型
type Category struct {
	BaseModel
	Name             string `gorm:"type:varchar(20);not null"`
	ParentCategoryID int32
	ParentCategory   *Category //指向自己的时候是需要定义为指针的
	Level            int32     `gorm:"type;int;not null"`
	IsTag            bool      `gorm:"default:false;not null"`
}

/**
多对多的关系，这个时候我们需要另外加一个表进行多对多关系的建立
这里品牌和分类之间就是多对多的关系
*/

type Brands struct {
	BaseModel
	Name string `gorm:"type:varchar(20);not null"`
	Logo string `gorm:"type:varchar(200);default:'';not null"`
}

type GoodsCategoryBrand struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Category   Category

	BrandsID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Brands   Brands
}

func (GoodsCategoryBrand) TableName() string {
	return "goodscategorybrand"
}

type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null"`
	Url   string `gorm:"type:varchar(200);not null"`
	Index int32  `gorm:"type:int;default:1;not null"`
}
