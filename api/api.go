package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/99designs/basicauth-go"
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
	
	username := os.Getenv("FC_SESSION_CACHE_USERNAME")
	password := os.Getenv("FC_SESSION_CACHE_PASSWORD")

	r.Use(basicauth.New("fc-hosting", map[string][]string{
		username: { password },
	}))

	r.Post("/put", func (w http.ResponseWriter, r *http.Request) {
		a.Put(w, r)
	})

	fmt.Println("API started.")

	http.ListenAndServe(":3000", r)
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
