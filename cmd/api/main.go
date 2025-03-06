package main

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/vlkhvnn/DocCollab/internal/auth"
	"github.com/vlkhvnn/DocCollab/internal/db"
	"github.com/vlkhvnn/DocCollab/internal/env"
	"github.com/vlkhvnn/DocCollab/internal/store"
	ws "github.com/vlkhvnn/DocCollab/internal/websocket"
	"go.uber.org/zap"
)

// upgrader is used to upgrade HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()
	err := godotenv.Load("../../.env")
	if err != nil {
		logger.Fatal("Error loading .env file")
	}
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbconfig{
			addr:         env.GetString("DB_ADDR", "postgres://postgres:1234@localhost/doccollab?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		auth: authConfig{
			basic: basicConfig{
				user: env.GetString("AUTH_BASIC_USER", ""),
				pass: env.GetString("AUTH_BASIC_PASS", ""),
			},
			token: tokenConfig{
				secret: env.GetString("AUTH_TOKEN_SECRET", ""),
				exp:    time.Hour * 24 * 3,
				iss:    "doccollab",
			},
		},
	}

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("DB connection pool established")

	store := store.NewStorage(db)

	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.iss, cfg.auth.token.iss)

	app := &application{
		config:        cfg,
		store:         store,
		authenticator: jwtAuthenticator,
		logger:        logger,
		hub:           ws.NewHub(&store),
	}
	mux := app.mount()
	logger.Fatal(app.run(mux))
}
