package app

import (
	"github.com/nishant01/procard-go-api/routes"
	"github.com/rs/cors"

	"net/http"
)

func StartApplication() {
	routes.Router.Use(jwtAuthentication) //attach JWT auth middleware
	routes.MapUrls()

	//handler := cors.Default().Handler(routes.Router)
	srv := http.Server{
		Handler: cors.Default().Handler(routes.Router),
		Addr:    ":8080",
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
