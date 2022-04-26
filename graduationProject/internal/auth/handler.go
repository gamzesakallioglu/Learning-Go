package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/httpErrors"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/api"
	mwUserAdmin "github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/middleware/userAdmin"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/pkg/config"
	jwtHelper "github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
)

type authHandler struct {
	cfg     *config.Config
	service Service
}

func NewAuthHandler(r *gin.RouterGroup, cfg *config.Config, service Service) {
	authHandler := authHandler{cfg: cfg, service: service}

	r.POST("/sign-up", authHandler.signUpCustomer)
	r.POST("/login", authHandler.loginCustomer)
	r.POST("/user/sign-up", mwUserAdmin.AuthMiddleware(cfg.JWTConfig.SecretKey), authHandler.signUpUser) // needs to be authorized for this
	r.POST("/user/login", authHandler.loginUser)

}

func (a authHandler) signUpCustomer(c *gin.Context) {

	var customerSignUp api.CustomerSignUp
	if err := c.Bind(&customerSignUp); err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "something's wrong! check your request body", nil)))
		return
	}

	format := strfmt.Default
	err := customerSignUp.Validate(format)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil)))
		return
	}

	customerCreated, err := a.service.CreateCustomer(c.Request.Context(), &customerSignUp)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	jwtClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": customerCreated.ID,
		"email":  customerCreated.Email,
		"iat":    time.Now().Unix(),
		"iss":    os.Getenv("ENV"),
		"exp":    time.Now().Add(60 * 60 * 60 * time.Second).Unix(),
		"role":   "customer",
	})

	token := jwtHelper.GenerateToken(jwtClaims, a.cfg.JWTConfig.SecretKey)
	c.JSON(http.StatusOK, token)

}

func (a authHandler) signUpUser(c *gin.Context) {

	var userSignUp api.UserSignUp
	if err := c.Bind(&userSignUp); err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "something's wrong! check your request body", nil)))
		return
	}

	format := strfmt.Default
	err := userSignUp.Validate(format)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil)))
		return
	}

	userCreated, err := a.service.CreateUser(c.Request.Context(), &userSignUp)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	jwtClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userCreated.ID,
		"email":  userCreated.Email,
		"iat":    time.Now().Unix(),
		"iss":    os.Getenv("ENV"),
		"exp":    time.Now().Add(60 * 60 * 60 * time.Second).Unix(),
		"role":   "customer",
	})

	token := jwtHelper.GenerateToken(jwtClaims, a.cfg.JWTConfig.SecretKey)

	c.JSON(http.StatusOK, token)
}

func (a *authHandler) loginCustomer(c *gin.Context) {
	var loginCustomer api.Login
	if err := c.Bind(&loginCustomer); err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "something's wrong! check your request body", nil)))
		return
	}

	format := strfmt.Default
	err := loginCustomer.Validate(format)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil)))
		return
	}

	customer, err := a.service.GetCustomer(c.Request.Context(), loginCustomer.Email, loginCustomer.Password)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil)))
		return
	}
	if customer == nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "customer not found", nil)))
		return
	}

	jwtClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": customer.ID,
		"email":  customer.Email,
		"iat":    time.Now().Unix(),
		"iss":    os.Getenv("ENV"),
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
		"role":   "customer",
	})
	token := jwtHelper.GenerateToken(jwtClaims, a.cfg.JWTConfig.SecretKey)

	c.JSON(http.StatusOK, token)
}

func (a *authHandler) loginUser(c *gin.Context) {
	var loginUser api.Login
	if err := c.Bind(&loginUser); err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "something's wrong! check your request body", nil)))
		return
	}

	format := strfmt.Default
	err := loginUser.Validate(format)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil)))
		return
	}

	user, err := a.service.GetUser(c.Request.Context(), loginUser.Email, loginUser.Password)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil)))
		return
	}
	if user == nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "user not found", nil)))
		return
	}

	jwtClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID.String(),
		"email":  user.Email,
		"iat":    time.Now().Unix(),
		"iss":    os.Getenv("ENV"),
		"exp":    time.Now().Add(60 * 50 * time.Second).Unix(),
		"role":   "user-" + user.UserRole, // i.e. user-admin
	})
	token := jwtHelper.GenerateToken(jwtClaims, a.cfg.JWTConfig.SecretKey)

	c.JSON(http.StatusOK, token)
}

func (a *authHandler) VerifyToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	decodedClaims := jwtHelper.VerifyToken(token, a.cfg.JWTConfig.SecretKey)

	c.JSON(http.StatusOK, decodedClaims)
}
