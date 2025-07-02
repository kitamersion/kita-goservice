package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kitamersion/go-goservice/graph"
	"github.com/kitamersion/go-goservice/internal/config"
	"github.com/kitamersion/go-goservice/internal/database"
	"github.com/kitamersion/go-goservice/internal/domain/repositories"
	"github.com/kitamersion/go-goservice/internal/domain/services"
	"github.com/kitamersion/go-goservice/internal/events"
	"github.com/kitamersion/go-goservice/internal/events/producer"
	"github.com/sirupsen/logrus"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8000"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Load configuration
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Setup logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Database connection
	db, err := database.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize Kafka topics (runs every startup, safe if already exists)
	if err := events.InitKafkaTopics(&cfg.Kafka, logger); err != nil {
		logger.WithError(err).Fatal("Failed to initialize Kafka topics")
	}

	// Initialize producer
	eventProducer := producer.NewProducer(&cfg.Kafka, logger)
	defer eventProducer.Close()

	// Initialize repositories and services
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo, eventProducer)

	// Initialize GraphQL resolver
	gqlResolver := &graph.Resolver{
		UserService: userService,
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: gqlResolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))
	http.Handle("/graphql", srv)

	log.Printf("connect to http://localhost:%s/playground for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
