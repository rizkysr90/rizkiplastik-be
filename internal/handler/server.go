package handler

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizkysr90/rizkiplastik-be/internal/config"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/authentication"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/category"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/onlinetransactions"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/packagingtypes"
	packagingtypesPg "github.com/rizkysr90/rizkiplastik-be/internal/handler/packagingtypes/repository/pg"
	productcategoryrules "github.com/rizkysr90/rizkiplastik-be/internal/handler/product_category_rules"
	productCategoryRulesPg "github.com/rizkysr90/rizkiplastik-be/internal/handler/product_category_rules/repository/pg"
	productsizeunitrules "github.com/rizkysr90/rizkiplastik-be/internal/handler/product_sizeunit_rules"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/products"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/sizeunits"
	sizeunitsPg "github.com/rizkysr90/rizkiplastik-be/internal/handler/sizeunits/repository/pg"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/summary"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/variantypes"
	variantypesPg "github.com/rizkysr90/rizkiplastik-be/internal/handler/variantypes/repository/pg"

	"github.com/rizkysr90/rizkiplastik-be/internal/middleware"
	"github.com/rizkysr90/rizkiplastik-be/internal/repository/pg"
)

// Server wraps gin.Engine
type Server struct {
	router *gin.Engine
	db     *pgxpool.Pool
	cfg    *config.Config
}

// NewServer initializes and configures the server
func NewServer(db *pgxpool.Pool, cfg *config.Config) *Server {
	// Set Gin mode - options: debug, release, test
	// For production, you should use release mode
	gin.SetMode(gin.DebugMode)

	// Create a new Gin router
	router := gin.New()

	// Use the recovery middleware to recover from panics
	router.Use(gin.Recovery())

	// Use custom logger middleware
	router.Use(middleware.Logger())

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	router.Use(cors.New(config))

	// Session middleware
	router.Use(middleware.Session())

	server := &Server{
		router: router,
		db:     db,
		cfg:    cfg,
	}
	// Register routes
	server.registerRoutes()

	return server
}

// Run starts the HTTP server
func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

// registerRoutes sets up all the routes for the server
func (s *Server) registerRoutes() {
	// Root route
	s.router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Gin API with CORS and logging",
		})
	})
	// Authentication routes
	authenticationHandler := authentication.NewAuthHandler(s.db, s.cfg)
	authenticationHandler.RegisterRoutes(s.router)

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(s.db, s.cfg)

	// Product routes
	productHandler := products.NewProductHandler(s.db)
	productHandler.RegisterRoutes(s.router, authMiddleware)

	// Online transaction routes

	onlineTransactionHandler := onlinetransactions.NewOnlineTransactions(s.db)
	onlineTransactionHandler.RegisterRoutes(s.router, authMiddleware)

	// Summary routes
	summaryHandler := summary.NewSummaryHandler(s.db)
	summaryHandler.RegisterRoutes(s.router, authMiddleware)

	// Category routes
	categoryRepo := pg.NewCategory(s.db)
	categoryService := category.NewService(categoryRepo)
	categoryHandler := category.NewCategoryHandler(categoryService)
	categoryHandler.RegisterRoutes(s.router, authMiddleware)

	// Packaging type routes
	packagingTypeRepo := packagingtypesPg.NewPackagingType(s.db)
	packagingTypeHandler := packagingtypes.NewHandler(packagingTypeRepo)
	packagingTypeHandler.RegisterRoutes(s.router)

	// Size unit routes
	sizeUnitRepo := sizeunitsPg.NewSizeUnits(s.db)
	sizeUnitHandler := sizeunits.NewHandler(sizeUnitRepo)
	sizeUnitHandler.RegisterRoutes(s.router)

	// Variant type routes
	variantTypeRepo := variantypesPg.NewVarianTypes(s.db)
	variantTypeHandler := variantypes.NewHandler(variantTypeRepo)
	variantTypeHandler.RegisterRoutes(s.router)

	// Product category rules routes
	productCategoryRulesRepo := productCategoryRulesPg.NewProductCategoryRules(s.db)
	productCategoryRulesHandler := productcategoryrules.NewHandler(productCategoryRulesRepo)
	productCategoryRulesHandler.RegisterRoutes(s.router)

	// Product size unit rules routes
	productCategoryRepoV2 := pg.NewCategory(s.db)
	sizeUnitRepoV2 := pg.NewSizeUnit(s.db)
	productSizeUnitRulesRepo := pg.NewProductSizeUnitRules(s.db, productCategoryRepoV2, sizeUnitRepoV2)
	productSizeUnitRulesHandler := productsizeunitrules.NewHandler(productSizeUnitRulesRepo)
	productSizeUnitRulesHandler.RegisterRoutes(s.router)
}
