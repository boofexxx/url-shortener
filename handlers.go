package main

import (
	"crypto/md5"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/boofexxx/url-shortener/store"
	"github.com/gorilla/mux"
)

type ServerMux struct {
	store *store.Store
	log   *log.Logger
	*mux.Router
}

// NewServerMux returns a new router instance
// with redis store and standard logger
func NewServerMux() (*ServerMux, error) {
	store := store.NewStore()
	err := store.Ping()
	if err != nil {
		return nil, err
	}
	return &ServerMux{
		store:  store,
		log:    log.New(os.Stdout, "shortener: ", log.Flags()),
		Router: mux.NewRouter(),
	}, nil
}

// CreateShortURLHandler generates and maps short url to
// url specified in query parameters with key "url"
func (sm *ServerMux) CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	origURL := r.URL.Query().Get("url")
	if origURL == "" {
		http.Error(w, "Redirect to nothing? O_O", http.StatusBadRequest)
		return
	}

	shortURL := GenerateShortURL(origURL)
	err := sm.store.Add(shortURL, origURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(shortURL))
}

// GetShortURLHandler redirects to original url mapped to
// short url that comes as a part of path
func (sm *ServerMux) GetShortURLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := vars["shortURL"]

	origURL, err := sm.store.Get(shortURL)
	if err != nil {
		code := http.StatusNotFound
		http.Error(w, http.StatusText(code), code)
		return
	}
	http.Redirect(w, r, httpPrepend(origURL), http.StatusFound)
}

// httpPrenend checks if there is any protocol specified
// if not it prepends https:// to url
func httpPrepend(url string) string {
	if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		return "https://" + url
	}
	return url
}

// GenerateShortURL generates short url with 8 characters
func GenerateShortURL(url string) string {
	urlSum := md5.Sum([]byte(url))
	shortURL := base64.StdEncoding.EncodeToString(urlSum[:])

	return shortURL[:8]
}

func (sm *ServerMux) LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sm.log.Printf("handle %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
