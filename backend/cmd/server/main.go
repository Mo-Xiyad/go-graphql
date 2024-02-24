// package main

// import (
// 	"log"
// 	generated "server/graph"

// 	"github.com/99designs/gqlgen/graphql/handler"
// 	"github.com/99designs/gqlgen/graphql/playground"
// 	fiber "github.com/gofiber/fiber/v2"
// 	"github.com/valyala/fasthttp/fasthttpadaptor"
// )

// func main() {
// 	schema := generated.NewExecutableSchema(generated.Config{Resolvers: &generated.Resolver{}})
// 	graphqlHandler := handler.New(schema)

// 	app := fiber.New()

// 	app.Post("/graphql", func(c *fiber.Ctx) error {
// 		fasthttpadaptor.NewFastHTTPHandler(graphqlHandler)(c.Context())
// 		return nil
// 	})

// 	// Set up GraphQL Playground
// 	playgroundHandler := playground.Handler("GraphQL Playground", "/graphql")
// 	app.Get("/playground", func(c *fiber.Ctx) error {
// 		fasthttpadaptor.NewFastHTTPHandler(playgroundHandler)(c.Context())
// 		return nil
// 	})
// 	app.Get("/playground", func(c *fiber.Ctx) error {
// 		playground.Handler("GraphQL Playground", "/graphql")
// 		return nil
// 	})

//		// Start server
//		log.Fatal(app.Listen(":4000"))
//	}
package main

import (
	"context"
	"log"
	"server"
	"server/cmd/resolvers"
	generated "server/graph"
	"server/pkg/db"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func graphqlHandler(ctx *server.Context) gin.HandlerFunc {
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers.Resolver{}}))

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

func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "ServerContextKey", c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func main() {
	ctx, err := server.NewContext()
	if err != nil {
		log.Fatalf("failed to initialize context: %v", err)
	}

	router := gin.Default()

	db.InitializeDB()

	//TODO: Fix CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	router.Use(cors.New(config))

	router.Use(GinContextToContextMiddleware())

	router.POST("/query", graphqlHandler(ctx))
	router.GET("/", playgroundHandler())

	if err := router.Run(":4000"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
