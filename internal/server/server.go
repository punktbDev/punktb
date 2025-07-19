package server

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gitlab.com/freelance/punkt-b/backend/config"
	"gitlab.com/freelance/punkt-b/backend/internal/controller"
	"gitlab.com/freelance/punkt-b/backend/internal/domain"
	"gitlab.com/freelance/punkt-b/backend/internal/dto"
	"gitlab.com/freelance/punkt-b/backend/internal/service"
	"gitlab.com/freelance/punkt-b/backend/pkg/database"
	"go.uber.org/zap"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	DefaultTimeout = 3 * time.Second
)

type (
	server struct {
		httpCfg     config.HTTP
		login       controller.Login
		manager     controller.Manager
		client      controller.Client
		srvManagers service.Manager
		cache
	}
	cache struct {
		m      sync.Mutex
		logins map[string]dto.Manager
	}
	Server interface {
		Run() error
	}
)

func NewServer(cfg *config.Config) (Server, error) {
	db, err := database.NewDatabase(fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable&connect_timeout=%d",
		cfg.PostgreSQL.Username, cfg.PostgreSQL.Password, net.JoinHostPort(cfg.PostgreSQL.Host, cfg.PostgreSQL.Port),
		cfg.PostgreSQL.Database, DefaultTimeout))
	if err != nil {
		return nil, err
	}

	zap.L().Info("connecting to the database is successfully")

	srv := service.NewManager(domain.NewManager(db))

	return &server{
		httpCfg:     cfg.HTTP,
		login:       controller.NewLogin(service.NewLogin(domain.NewLogin(db))),
		manager:     controller.NewManager(srv),
		client:      controller.NewClient(service.NewClient(domain.NewClient(db))),
		srvManagers: srv,
	}, nil
}

func (s *server) Run() error {
	go s.fillLoginCache()

	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()
	s.Routes(api)

	handler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"access-control-allow-headers", "access-control-allow-origin", "Accept",
			"X-Requested-With", "If-Modified-Since", "User-Agent", "Keep-Alive", "X-CustomHeader", "Accept-Language",
			"Content-Type", "Content-Language", "Origin", "api_key", "Authorization", "DNT", "Cache-Control"}),
	)(router)
	srv := http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadTimeout:       5 * time.Second,   // Таймаут на чтение данных (запроса) от клиента
		WriteTimeout:      30 * time.Second,  // Таймаут на запись данных (ответа)
		IdleTimeout:       120 * time.Second, // Таймаут для ожидания новых запросов на соединении
		ReadHeaderTimeout: 10 * time.Second,
	}

	zap.L().Info("server started")

	if s.httpCfg.IsHTTPS {
		zap.L().Info("https mode is on")
		srv.Addr = ":8443"
		return srv.ListenAndServeTLS(s.httpCfg.ServerCertPath, s.httpCfg.PrivateKeyPath)
	} else {
		return srv.ListenAndServe()
	}
}

func (s *server) fillLoginCache() {
	s.cache.logins = make(map[string]dto.Manager, 10)
	ticker := time.NewTicker(5 * time.Second)
loop:
	for {
		select {
		case <-ticker.C:
			ms, err := s.srvManagers.GetAllManagers()
			if err != nil {
				zap.L().Debug("fillLoginCache", zap.Error(err))
				continue loop
			}
			s.cache.logins = make(map[string]dto.Manager, 10)

			for _, v := range ms {
				s.cache.logins[v.Login] = v
			}
		}
	}
}
