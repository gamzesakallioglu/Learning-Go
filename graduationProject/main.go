package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/auth"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/customer"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/order"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/product"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/productCategory"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/shoppingCart"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/user"
	"github.com/gin-gonic/gin"

	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/pkg/config"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/pkg/db"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/pkg/graceful"
)

func main() {
	log.Println("Basket service is getting started")

	// Setting environments
	cfg, err := config.LoadConfig("./pkg/config/config-local")
	if err != nil {
		log.Fatalf("environment variables could not be set %v", err)
	}
	//

	// Connect to DB
	DB := db.NewPsqlDB(cfg)
	//

	// Gin - Http package

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	//

	// Server
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.ServerConfig.ReadTimeoutSecs) * time.Second,
		WriteTimeout: time.Duration(cfg.ServerConfig.WriteTimeoutSecs) * time.Second,
	}
	//

	rootRouter := r.Group(cfg.ServerConfig.RoutePrefix) // starts with route prefix

	//User Repository
	userRepo := user.NewRepository(DB)
	userRepo.Migration()

	//Customer Repository
	customerRepo := customer.NewRepository(DB)
	customerRepo.Migration()

	authRepo := auth.NewAuthRepository(DB)
	authService := auth.NewAuthService(authRepo)
	auth.NewAuthHandler(rootRouter, cfg, authService)

	productCategoryRepo := productCategory.NewRepository(DB)
	productCategoryRepo.Migration()
	productCategoryService := productCategory.NewProductCategoryService(productCategoryRepo)
	productCategory.NewProductCategoryHandler(rootRouter, cfg, productCategoryService)

	productRepo := product.NewRepository(DB)
	productRepo.Migration()
	productService := product.NewProductService(productRepo)
	product.NewProductHandler(rootRouter, cfg, productService)

	shoppingCartRepo := shoppingCart.NewRepository(DB)
	shoppingCartRepo.Migration()
	shoppingCartService := shoppingCart.NewShoppingCartService(shoppingCartRepo)
	shoppingCart.NewShoppingCartHandler(rootRouter, cfg, shoppingCartService)

	orderRepo := order.NewRepository(DB)
	orderRepo.Migration()
	orderService := order.NewOrderService(orderRepo)
	order.NewOrderHandler(rootRouter, cfg, orderService)

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Println("Book store service started")
	graceful.ShutdownGin(srv, time.Duration(cfg.ServerConfig.TimeoutSecs*int64(time.Second)))

}

func abc(c *gin.Context) {
	c.JSON(http.StatusOK, "abc")
}
