package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

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
	
	username := strings.Trim(os.Getenv("FC_SESSION_CACHE_USERNAME"), "\n")
	password := strings.Trim(os.Getenv("FC_SESSION_CACHE_PASSWORD"), "\n")

	if username != "" && password != "" {
		r.Use(basicauth.New("fc-hosting", map[string][]string{
			username: { password },
		}))
	}
	
	r.Post("/put", func (w http.ResponseWriter, r *http.Request) {
		a.Put(w, r)
	})

	r.Get("/get", func (w http.ResponseWriter, r *http.Request) {
		a.Get(w, r)
	})

	r.Delete("/remove", func(w http.ResponseWriter, r *http.Request) {
		a.Remove(w, r)
	})

	r.Delete("/flush", func(w http.ResponseWriter, r *http.Request) {
		a.Flush(w, r)
	})

	fmt.Println("API started.")

	http.ListenAndServe(":3000", r)
}

func (a *Api) ProcessBody(w http.ResponseWriter, r *http.Request, s interface{}) error {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		return err
	}

	if jsonErr := json.Unmarshal(body, s); jsonErr != nil {
		return jsonErr
	}

	return nil
}

type SessionObject struct {
	C string `json:"cookie"`
	S map[string]interface{} `json:"session"`
}
func (a *Api) Put(w http.ResponseWriter, r *http.Request) error {
	pr := &SessionObject{}
	err := a.ProcessBody(w, r, pr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return err
	}

	fmt.Printf("putting:\n %s\n %s", pr.C, pr.S)

	res := a.Cache.Put(pr.C, pr.S)

	if res != nil {
		http.Error(w, res.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return res
	}

	w.WriteHeader(200)

	w.Write([]byte("Cached successfully."))

	return nil
}

type GetRequest struct {
	Cookie string `json:"cookie"`
}
func (a *Api) Get(w http.ResponseWriter, r *http.Request) error {
	cookie := r.URL.Query().Get("sid")

	unescaped, err := url.QueryUnescape(cookie)

	if err != nil {
		http.Error(w, "The provided session id might be invalid!", http.StatusBadRequest)
		return errors.New("you need to add `sid` to your query parameters")
	}

	cookie = strings.Trim(unescaped, "\n")

	fmt.Println(cookie)

	if cookie == "" {
		http.Error(w, "You need to add `sid` to your query parameters!", http.StatusBadRequest)
		return errors.New("you need to add `sid` to your query parameters")
	}

	res, err := a.Cache.Get(cookie)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	bytes, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	if string(bytes) == "" {
		http.Error(w, "The given session was not found.", http.StatusNotFound)
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	fmt.Printf("session was fetched; returning:\n %s", string(bytes))

	w.Write(bytes)

	return nil
}

type RemoveRequest struct {
	Cookie string `json:"cookie"`
}
func (a *Api) Remove(w http.ResponseWriter, r *http.Request) {
	rr := &RemoveRequest{}
	err := a.ProcessBody(w, r, rr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	isFound, err := a.Cache.Remove(rr.Cookie)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if !isFound {
		http.Error(w, "The given cookie was not found.", http.StatusNotFound)
		return
	}

	w.WriteHeader(200)
	
	w.Write([]byte("The session was removed successfully."))

	return
}

func (a *Api) Flush(w http.ResponseWriter, r *http.Request) {
	if err := a.Cache.Flush(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(200)
	
	w.Write([]byte("The cache was flushed successfully."))

	return
}
