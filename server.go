package main

import (
	"duit-pasutri-be/graph"
	"duit-pasutri-be/graph/resolvers"
	"duit-pasutri-be/internal/database"
	"duit-pasutri-be/internal/directives"
	"duit-pasutri-be/internal/middleware"
	"duit-pasutri-be/internal/repositories"
	"duit-pasutri-be/internal/services"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8000"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repositories.NewUserRepository(db)
	accountRepo := repositories.NewAccountRepository(db)
	categoryRepo := repositories.NewCategoryRepositoryImpl(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	userService := services.NewUserService(userRepo, db)
	accountService := services.NewAccountService(db, accountRepo, userRepo)
	categoryService := services.NewCategoryService(db, categoryRepo)
	transactionService := services.NewTransactionService(db, transactionRepo, accountRepo, categoryRepo, userRepo)

	config := graph.Config{Resolvers: &resolvers.Resolver{
		UserService:        userService,
		AccountService:     accountService,
		CategoryService:    categoryService,
		TransactionService: transactionService,
	}}

	config.Directives.Auth = directives.Auth

	srv := handler.New(graph.NewExecutableSchema(config))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", middleware.AuthMiddleware(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
