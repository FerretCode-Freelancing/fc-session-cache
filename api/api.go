package api

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ferretcode-hosting/fc-session-cache/cache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Api struct {
	Cache cache.Cache 
}

func (a *Api) NewApi() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)

	r.Post("/put", func (w http.ResponseWriter, r *http.Request) {
		a.Put(w, r)
	})

	fmt.Println("API started.")

	http.ListenAndServe(":3000", r)
}

func (a *Api) Auth(realm string, creds map[string]string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()

			if !ok {
				basicAuthFailed(w, realm)

				return
			}

			credPass, credUserOk := creds[user]
			if !credUserOk || subtle.ConstantTimeCompare([]byte(pass), []byte(credPass)) != 1 {
				basicAuthFailed(w, realm)	

				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

func basicAuthFailed(w http.ResponseWriter, realm string) {
	w.Header().Add("WWW-Authenticate", fmt.Sprintf("Basic realm=%s", realm))
	w.WriteHeader(http.StatusUnauthorized)
}

type PutRequest struct {
	Cookie string `json:"cookie"`
	Session SessionObject `json:"session"`
}
type SessionObject struct {
	C string
	S map[string]string
}
func (a *Api) Put(w http.ResponseWriter, r *http.Request) error {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	pr := &PutRequest{}
	if err := json.Unmarshal(body, pr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	res := a.Cache.Put(pr.Cookie, pr.Session)

	if res != nil {
		http.Error(w, res.Error(), http.StatusInternalServerError)
		return res
	}

	w.WriteHeader(200)

	w.Write([]byte("Cached successfully."))

	return nil
}
