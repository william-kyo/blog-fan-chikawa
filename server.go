package main

import (
	"blog-fanchiikawa-service/db"
	"blog-fanchiikawa-service/graph"
	"blog-fanchiikawa-service/repository"
	"blog-fanchiikawa-service/resolver"
	"blog-fanchiikawa-service/scheduler"
	"blog-fanchiikawa-service/sdk"
	"blog-fanchiikawa-service/service"
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

const defaultPort = "8080"

func main() {
	// Initialize infrastructure
	db.InitMySQL()
	sdk.InitAWSSession()
	sdk.InitS3()
	sdk.InitComprehend()
	sdk.InitTranslate()
	sdk.InitPolly()

	// Initialize repositories
	userRepo := repository.NewUserRepository()
	deviceRepo := repository.NewUserDeviceRepository()
	imageRepo := repository.NewImageReposity()
	transactionMgr := repository.NewTransactionManager()

	// Initialize services
	languageService := service.NewLanguageService()
	translateService := service.NewTranslateService()
	speechService := service.NewSpeechService(languageService)
	storageService := service.NewStorageService()
	userService := service.NewUserService(userRepo, deviceRepo, transactionMgr)
	mediaService := service.NewMediaService(imageRepo)

	// Initialize resolver
	resolverInstance := resolver.NewResolver(
		userService,
		languageService,
		translateService,
		speechService,
		storageService,
	)

	// Initialize Scheduler
	scheduler := scheduler.NewScheduler(mediaService)
	defer scheduler.Shutdown()
	scheduler.DataSync()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{Resolver: resolverInstance},
	}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
