package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

func (s *server) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ck := getCookies(r)
		h := getHeaders(r)

		zap.L().Info(fmt.Sprintf("%s %s %s, cookies: %s", r.Method, r.RequestURI,
			r.RemoteAddr, ck))
		zap.L().Debug("headers: " + h)

		if r.Method == http.MethodPost && r.URL.String() == "/api/v1/client/result" {
			next.ServeHTTP(w, r)
			return
		}

		login, password, ok := r.BasicAuth()
		if ok {
			l, ok := s.LoginCache[login]
			if ok && l.Password == password {
				if !l.IsActive {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				ctx := context.WithValue(r.Context(), "manager", &l)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})
}

func (s *server) Routes(r *mux.Router) {
	r.Use(s.middleware)
	r.HandleFunc("/login", s.login.Login).Methods(http.MethodGet)
	r.HandleFunc("/manager", s.manager.GetUserData).Methods(http.MethodGet)
	r.HandleFunc("/manager", s.manager.ChangeManagerData).Methods(http.MethodPut)
	r.HandleFunc("/manager/add", s.manager.AddManager).Methods(http.MethodPost)
	r.HandleFunc("/managers", s.manager.GetAllManagers).Methods(http.MethodGet)
	r.HandleFunc("/manager/is-active/{id}", s.manager.ChangeActive).Methods(http.MethodPut)
	r.HandleFunc("/manager/full-access/{id}", s.manager.ChangeFullAccess).Methods(http.MethodPut)
	r.HandleFunc("/clients", s.client.GetClients).Methods(http.MethodGet)
	r.HandleFunc("/client/is-new", s.client.SetClientChecked).Methods(http.MethodPost)
	r.HandleFunc("/client/is-archive", s.client.SetClientArchive).Methods(http.MethodPost)
	r.HandleFunc("/client/result", s.client.AddResult).Methods(http.MethodPost)
	r.HandleFunc("/client/{id}", s.client.GetResultClient).Methods(http.MethodGet)
}

func getCookies(r *http.Request) string {
	var ck string
	for _, v := range r.Cookies() {
		ck += v.String() + " "
	}

	if ck == "" {
		ck = "NO"
	}

	return ck
}

func getHeaders(r *http.Request) string {
	var h, val string
	for name, values := range r.Header {
		for _, v := range values {
			val += v + " "
		}
		h += fmt.Sprintf("%s: %s", name, val)
		val = ""
	}

	return h
}
