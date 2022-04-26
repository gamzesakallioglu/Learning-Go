package productCategory

import (
	"context"
	"fmt"

	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetProductCategoryByName(ctx context.Context, name *string) (*models.ProductCategory, error)
	GetAllProductCategories(ctx context.Context, pagesize int, page int, sorting string) (*models.ProductCategories, int64)
	CreateProductCategoryForBulk(ctx context.Context, category *models.ProductCategory) error
	DeleteProductCategoryByID(ctx context.Context, id *string) error
	GetProductCategoryByID(ctx context.Context, id *string) (*models.ProductCategory, error)
	Migration()
}

type productCategoryRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &productCategoryRepository{db: db}
}

// Migrations keeps db schema up to date
func (b *productCategoryRepository) Migration() {
	b.db.AutoMigrate(&models.ProductCategory{})
}

func (p *productCategoryRepository) DeleteProductCategoryByID(ctx context.Context, id *string) error {

	if err := p.db.Where("id = ?", id).Delete(&models.ProductCategory{}).Error; err != nil {
		return err
	}
	return nil
}

func (p *productCategoryRepository) GetProductCategoryByID(ctx context.Context, id *string) (*models.ProductCategory, error) {
	var category *models.ProductCategory
	if err := p.db.WithContext(ctx).Where("id = ?", id).First(&category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (p *productCategoryRepository) GetProductCategoryByName(ctx context.Context, name *string) (*models.ProductCategory, error) {
	var category *models.ProductCategory
	if err := p.db.WithContext(ctx).Where("name = ?", *name).First(&category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (p *productCategoryRepository) CreateProductCategoryForBulk(ctx context.Context, category *models.ProductCategory) error {
	// For Bulk action. Upsert
	if category.IsParent { // Parent categories has 0 as parent id
		category.ParentID = "0"
	}

	// If a category exists with same name -> update
	// If doesn't exist -> insert
	if err := p.db.Where(models.ProductCategory{Name: category.Name}).Assign(models.ProductCategory{Name: category.Name, Description: category.Description, IsParent: category.IsParent, ParentID: category.ParentID}).FirstOrCreate(category).Error; err != nil {
		return err
	}

	return nil

}

func (p *productCategoryRepository) GetAllProductCategories(ctx context.Context, pagesize int, page int, sorting string) (*models.ProductCategories, int64) {

	var productCategories *models.ProductCategories
	var totalcount int64

	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * pagesize
	//p.db.Offset(offset).Limit(pagesize).Order(sorting).Find(&productCategories).Count(&totalcount)

	sqlQuery := fmt.Sprintf(`SELECT id,name,parent_id,
	(
		case
		when parent_id ='0' then id::text
		else parent_id
		end
	) AS order_parent
FROM     product_categories where deleted_at IS NULL 
ORDER BY order_parent, parent_id offset %v limit %v`, offset, pagesize)

	p.db.Raw(sqlQuery).Scan(&productCategories).Count(&totalcount)

	return productCategories, totalcount
}
