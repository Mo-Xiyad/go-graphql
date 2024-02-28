package main

import (
	"log"
	"server"
	"server/cmd/resolvers"
	auth "server/cmd/services/auth"
	user "server/cmd/services/user"
	"server/config"
	generated "server/graph"
	"server/pkg/db"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// define types for all services
type Services struct {
	AuthService auth.IAuthService
	UserService user.IUserService
}

func graphqlHandler(ctx *server.Context, services Services) gin.HandlerFunc {
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &resolvers.Resolver{
			AuthService: services.AuthService,
			UserService: services.UserService,
		},
	}))

	return func(c *gin.Context) {

		ctxWithDB := server.WithDB(c.Request.Context(), ctx.DB)
		c.Request = c.Request.WithContext(ctxWithDB)

		h.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

type Initializer struct {
	ctx  *server.Context
	db   *gorm.DB
	conf *config.Config
}

func initializer() (*Initializer, error) {
	conf := config.New()

	database, err := db.InitializeDB(conf)
	if err != nil {
		return nil, err
	}

	ctx, err := server.NewContext(database)
	if err != nil {
		return nil, err
	}

	return &Initializer{
		ctx:  ctx,
		db:   database,
		conf: conf,
	}, nil
}

func main() {
	initializer, err := initializer()

	if err != nil {
		log.Fatalf("failed to initialize server: %v", err)
	}

	router := gin.Default()

	//TODO: Fix CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	router.Use(cors.New(config))

	router.Use(server.GinContextToContextMiddleware())

	// Repository layer to communicate with the database
	userRepo := user.NewUserRepo(initializer.db)

	// Service layer to handle business logic
	authTokenService := auth.NewTokenService(initializer.conf)
	authService := auth.NewAuthService(userRepo, authTokenService)
	userService := user.NewUserService(userRepo)

	allowedOperations := map[string]bool{
		"Login": true,
	}
	router.Use(authMiddleware(authTokenService, allowedOperations))

	router.POST("/query", graphqlHandler(initializer.ctx, Services{
		AuthService: authService,
		UserService: userService,
	}))
	router.GET("/", playgroundHandler())

	if err := router.Run(":4000"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
