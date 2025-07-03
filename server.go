package main

import (
	"blog-fanchiikawa-service/db"
	"blog-fanchiikawa-service/graph"
	"blog-fanchiikawa-service/repository"
	"blog-fanchiikawa-service/resolver"
	"blog-fanchiikawa-service/scheduler"
	"blog-fanchiikawa-service/sdk"
	"blog-fanchiikawa-service/service"
	"blog-fanchiikawa-service/websocket"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Initialize infrastructure
	db.InitMySQL()
	sdk.InitAWS() // Unified AWS initialization for both SDK v1 and v2
	sdk.InitS3()
	sdk.InitComprehend()
	sdk.InitTranslate()
	sdk.InitPolly()
	sdk.InitRekognition()
	sdk.InitTextract()

	// Initialize repositories
	userRepo := repository.NewUserRepository()
	deviceRepo := repository.NewUserDeviceRepository()
	imageRepo := repository.NewImageReposity()
	labelRepo := repository.NewLabelRepository()
	imageLabelRepo := repository.NewImageLabelRepository()
	textKeywordRepo := repository.NewTextKeywordRepository()
	imageTextKeywordRepo := repository.NewImageTextKeywordRepository()
	transactionMgr := repository.NewTransactionManager()
	chatRepo := repository.NewChatRepository(db.GetEngine())
	chatMessageRepo := repository.NewChatMessageRepository(db.GetEngine())

	// Initialize services
	languageService := service.NewLanguageService()
	translateService := service.NewTranslateService()
	speechService := service.NewSpeechService(languageService)
	storageService := service.NewStorageService()
	userService := service.NewUserService(userRepo, deviceRepo, transactionMgr)
	mediaService := service.NewMediaService(imageRepo, labelRepo, imageLabelRepo, textKeywordRepo, imageTextKeywordRepo, transactionMgr)
	lexService := sdk.NewLexService()
	chatService := service.NewChatService(chatRepo, chatMessageRepo, lexService)
	configService := service.NewConfigService()
	customLabelsService := service.NewCustomLabelsService()

	// Initialize WebSocket hub
	hub := websocket.NewHub(chatService)
	go hub.Run()

	// Initialize resolver
	resolverInstance := resolver.NewResolver(
		userService,
		languageService,
		translateService,
		speechService,
		storageService,
		chatService,
		configService,
		customLabelsService,
	)

	// Initialize Scheduler
	scheduler := scheduler.NewScheduler(mediaService)
	defer scheduler.Shutdown()
	scheduler.ImageSync()
	scheduler.ImageLabelDetect()
	scheduler.ImageTextDetect()

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
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	// Serve static files from web directory
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("./web/"))))
	
	// Serve individual pages
	http.HandleFunc("/chat/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/chat.html")
	})
	http.HandleFunc("/custom-labels/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/custom-labels.html")
	})
	http.HandleFunc("/test/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/test.html")
	})
	http.HandleFunc("/home/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})
	
	http.Handle("/playground/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	http.HandleFunc("/ws", hub.ServeWS)
	
	// Default route to navigation page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "./web/index.html")
		} else {
			http.NotFound(w, r)
		}
	})

	log.Printf("Navigation page available at http://localhost:%s/", port)
	log.Printf("GraphQL playground available at http://localhost:%s/playground/", port)
	log.Printf("Chat interface available at http://localhost:%s/chat/", port)
	log.Printf("Custom Labels interface available at http://localhost:%s/custom-labels/", port)
	log.Printf("API Test page available at http://localhost:%s/test/", port)
	log.Printf("WebSocket endpoint available at ws://localhost:%s/ws", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
