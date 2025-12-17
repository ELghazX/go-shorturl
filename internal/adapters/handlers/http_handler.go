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
		<div class="border-2 border-black p-5 bg-gray-100">
			<div class="text-xl mb-4 break-all">%s</div>
			<button onclick="navigator.clipboard.writeText('%s')" 
					class="w-full px-4 py-4 border-2 border-black bg-black text-white font-bold uppercase tracking-wider hover:bg-white hover:text-black transition-colors">
				COPY URL
			</button>
		</div>
	`, shortURL, shortURL)
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

func (h *HTTPHandler) HandleAPIStats(w http.ResponseWriter, r *http.Request) {
	urls, err := h.urlService.GetStats(r.Context())
	if err != nil {
		http.Error(w, "Failed to get stats", http.StatusInternalServerError)
		return
	}

	if len(urls) == 0 {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<div class="py-16 text-center text-xl uppercase">No Data</div>`)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `<table class="w-full border-collapse mb-10"><thead><tr class="bg-black text-white"><th class="px-4 py-4 text-left font-bold uppercase tracking-wider">Code</th><th class="px-4 py-4 text-left font-bold uppercase tracking-wider">URL</th><th class="px-4 py-4 text-left font-bold uppercase tracking-wider">Clicks</th><th class="px-4 py-4 text-left font-bold uppercase tracking-wider">Created</th></tr></thead><tbody>`)
	for _, url := range urls {
		fmt.Fprintf(w, `<tr class="border-b border-black hover:bg-gray-100"><td class="px-4 py-4 font-bold"><a href="/%s" target="_blank" class="underline hover:bg-black hover:text-white">%s</a></td><td class="px-4 py-4 max-w-md truncate">%s</td><td class="px-4 py-4 font-bold text-xl">%d</td><td class="px-4 py-4">%s</td></tr>`,
			url.ShortCode, url.ShortCode, url.LongURL, url.Clicks, url.CreatedAt.Format("2006-01-02 15:04"))
	}
	fmt.Fprint(w, `</tbody></table>`)
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
