package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/elghazx/go-shorturl/internal/core/services"
)

type HTTPHandler struct {
	urlService *services.URLService
	tmpl       *template.Template
}

func NewHTTPHandler(urlService *services.URLService) *HTTPHandler {
	tmpl := template.Must(template.ParseGlob("templates/*html"))
	return &HTTPHandler{
		urlService: urlService,
		tmpl:       tmpl,
	}
}

func (h *HTTPHandler) HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.HandleRedirect(w, r)
		return
	}
	h.tmpl.ExecuteTemplate(w, "index.html", nil)
}

func (h *HTTPHandler) HandleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	longURL := r.FormValue("url")
	if longURL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	url, err := h.urlService.ShortenURL(r.Context(), longURL)
	if err != nil {
		http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
		return
	}

	shortURL := fmt.Sprintf("%s/%s", getBaseURL(r), url.ShortCode)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<div class="bg-green-50 border border-green-200 rounded-lg p-6 text-center">
		<p class="text-lg mb-4">Short URL: <a href="%s" target="_blank" class="text-blue-500 hover:text-blue-700 font-mono">%s</a></p>
		<button onclick="navigator.clipboard.writeText('%s')" 
		        class="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600 transition-colors">
		    Copy URL
		</button>
	</div>`, shortURL, shortURL, shortURL)
}

func (h *HTTPHandler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	shortCode := strings.TrimPrefix(r.URL.Path, "/")
	if shortCode == "" || shortCode == "favicon.ico" {
		return
	}

	longURL, err := h.urlService.GetOriginalURL(r.Context(), shortCode)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)
}

func (h *HTTPHandler) HandleStats(w http.ResponseWriter, r *http.Request) {
	urls, err := h.urlService.GetStats(r.Context())
	if err != nil {
		http.Error(w, "Failed to get stats", http.StatusInternalServerError)
		return
	}

	h.tmpl.ExecuteTemplate(w, "stats.html", urls)
}

func getBaseURL(r *http.Request) string {
	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		return strings.TrimSuffix(baseURL, "/")
	}

	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s", scheme, r.Host)
}
