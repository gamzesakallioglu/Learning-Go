package productCategory

import (
	"bufio"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/httpErrors"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/api"
	mwUserAdmin "github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/middleware/userAdmin"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/pkg/config"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/pkg/pagination"
	"github.com/gin-gonic/gin"
)

type Form struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type productCategoryHandler struct {
	cfg     *config.Config
	service Service
}

func NewProductCategoryHandler(r *gin.RouterGroup, cfg *config.Config, service Service) {
	productCategoryHandler := productCategoryHandler{cfg: cfg, service: service}

	r.POST("/productCategories", mwUserAdmin.AuthMiddleware(cfg.JWTConfig.SecretKey), productCategoryHandler.createBulkProductCategories)
	r.GET("/productCategories", productCategoryHandler.GetAllProductCategories)
	r.DELETE("/productCategories/:id", mwUserAdmin.AuthMiddleware(cfg.JWTConfig.SecretKey), productCategoryHandler.deleteProductCategoryByID)
}

func (p *productCategoryHandler) deleteProductCategoryByID(c *gin.Context) {

	id := c.Param("id")
	if len(id) <= 0 {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "id cannot be null", nil)))
		return
	}

	p.service.DeleteProductCategoryByID(c.Request.Context(), &id)

}

func (p *productCategoryHandler) createBulkProductCategories(c *gin.Context) {
	// File is going to readed and upserted into db
	// If a record exists in db with same name -> update
	// If not -> insert

	var form Form
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "something's wrong! check your request body", nil)))
		return
	}

	formFile, _ := c.FormFile("file")

	openedFile, _ := formFile.Open()
	defer openedFile.Close()

	scanner := bufio.NewScanner(openedFile)
	lineCount := 0
	//lineCountP := &lineCount

	for scanner.Scan() {
		lineCount++
		if lineCount == 1 {
			continue
		}

		lineSplit := strings.Split(scanner.Text(), ",")
		name := lineSplit[0]
		description := lineSplit[1]
		isParentString := lineSplit[2]
		parentID := lineSplit[3]

		isParent, _ := strconv.ParseBool(isParentString)
		productCategory := api.ProductCategory{Name: &name, Description: description, IsParent: isParent, ParentID: parentID}

		err = p.service.CreateProductCategoryForBulk(c.Request.Context(), &productCategory)
		if err != nil {
			c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil)))
			return
		}
	}

	c.String(http.StatusOK, "File's readed and inserted into database")
}

func (p *productCategoryHandler) GetAllProductCategories(c *gin.Context) {

	page, pageSize := pagination.GetPaginationParametersFromRequest(c)

	pagination, _ := p.service.GetAllProductCategories(c.Request.Context(), page, pageSize)

	c.JSON(http.StatusOK, pagination)

}
