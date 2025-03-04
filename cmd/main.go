package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/dotcreep/go-automate-deploy/docs"
	"github.com/dotcreep/go-automate-deploy/internal/api/cloudflare_api"
	"github.com/dotcreep/go-automate-deploy/internal/api/deploy_api"
	"github.com/dotcreep/go-automate-deploy/internal/api/jenkins_api"
	"github.com/dotcreep/go-automate-deploy/internal/api/portainer_api"
	"github.com/dotcreep/go-automate-deploy/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func KeyMiddlewareOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Json := utils.Json{}
		ymlconf, err := utils.Open()
		if err != nil {
			Json.NewResponse(false, nil, nil, "config.yml not found", http.StatusInternalServerError, nil)
			return
		}
		apiKey := ymlconf.Config.XToken
		if apiKey == "" {
			Json.NewResponse(false, w, nil, "token not found", http.StatusUnauthorized, nil)
			return
		}
		if apiKey != r.Header.Get("X-Token") {
			Json.NewResponse(false, w, nil, "unauthorized", http.StatusUnauthorized, nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func KeyMiddlewareTwo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Json := utils.Json{}
		ymlconf, err := utils.Open()
		if err != nil {
			Json.NewResponse(false, nil, nil, "config.yml not found", http.StatusInternalServerError, err.Error())
			return
		}
		apiKey := ymlconf.Config.XTokenX
		if apiKey == "" {
			Json.NewResponse(false, w, nil, "token not found", http.StatusUnauthorized, nil)
			return
		}
		if apiKey != r.Header.Get("X-Token") {
			Json.NewResponse(false, w, nil, "unauthorized", http.StatusUnauthorized, nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// @title						Automate Deployment API
// @version					1.0
// @description				Documentation for Automate Deployment Restful API
// @license.name				MIT
// @license.url				https://opensource.org/licenses/MIT
// @BasePath					/
//
// @SecurityDefinitions.apikey	X-Token
// @Name						X-Token
// @In							header
// @Description				Input your token authorized
func main() {
	Json := utils.Json{}
	ymlconf, err := utils.Open()
	if err != nil {
		Json.NewResponse(false, nil, nil, "config.yml not found", http.StatusInternalServerError, err.Error())
		return
	}
	port := ymlconf.Config.PORT
	r := chi.NewRouter()
	cors := cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Token"},
		AllowCredentials: true,
	})
	r.Use(cors)
	r.Route("/api/v1/deploy", func(r chi.Router) {
		r.Use(KeyMiddlewareTwo)
		r.Post("/start", deploy_api.Deploy)      // Docs
		r.Delete("/remove", deploy_api.Undeploy) // Docs
	})
	r.Route("/api/v1/domain", func(r chi.Router) {
		r.Use(KeyMiddlewareTwo)
		r.Post("/add", cloudflare_api.AddDomain)
		r.Post("/nameserver", cloudflare_api.GetNameserver)             // Docs
		r.Post("/status", cloudflare_api.StatusDomain)                  // Docs
		r.Post("/register-status", cloudflare_api.StatusRegisterDomain) // Docs
		r.Post("/check", cloudflare_api.GetBasedomainRegisteredStatus)  // Docs
		r.Post("/is-not-exists", cloudflare_api.GetDomainIsNotExists)   // Docs
	})
	r.Route("/api/v1/system", func(r chi.Router) {
		r.Use(KeyMiddlewareTwo)
		r.Get("/stack", portainer_api.GetStack)           // Docs
		r.Post("/status", portainer_api.GetStatusOfStack) // Docs
		r.Post("/is-not-exists", portainer_api.GetStackIsNotExists)
		r.Post("/update", portainer_api.UpdateStackByName)
	})
	r.Route("/api/v1/mobile", func(r chi.Router) {
		r.Use(KeyMiddlewareTwo)
		r.Post("/is-not-exists", jenkins_api.GetBuilderIsNotExists)
		r.Post("/status", jenkins_api.GetStatusOfItem)
	})
	// r.Route("/api/v1/let-me-in", func(r chi.Router) {
	// 	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 		Json.NewResponse(true, w, "fine", "validated", http.StatusOK, nil)
	// 	})
	// })
	r.Get("/docs/*", httpSwagger.WrapHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		Json.NewResponse(false, w, nil, "page not found", http.StatusNotFound, nil)
	})
	// mux := http.NewServeMux()
	// mux.HandleFunc("PUT /add", api.AddDomain)
	// mux.HandleFunc("PUT /delete", api.DeleteDomain)
	// mux.HandleFunc("POST /register", api.RegisterDomain)
	// mux.HandleFunc("GET /status", api.StatusDomain)
	// mux.HandleFunc("GET /stack", api.GetStack)
	// mux.HandleFunc("POST /stack", api.AddStack)
	// mux.HandleFunc("POST /deploy", api.Deploy)
	fmt.Printf("Server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
