package app

import (
	"context"
	"net/http"
	"strings"

	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"

	"grpc-http-server/config"
)

type Application struct {
	config.Config
	GRPCServer *grpc.Server
	HTTPServer *http.ServeMux
}

func NewApplication(config config.Config) *Application {
	return &Application{
		Config: config,
	}
}

func (app *Application) Start(ctx context.Context) error {
	if err := app.initGRPCServer(ctx); err != nil {
		return err
	}
	if err := app.initHTTPServer(ctx); err != nil {
		return err
	}

	c := cors.New(cors.Options{
		// NOTE: DO NOT USE THIS IN PRODUCTION
		AllowedOrigins: []string{"*"},
	})

	return http.ListenAndServe(app.Config.GetAddr(), c.Handler(app.handler()))
}

// Handler serves conn with GRPC or HTTP server
func (app *Application) handler() http.Handler {
	return h2c.NewHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
				app.GRPCServer.ServeHTTP(w, r)
			} else {
				app.HTTPServer.ServeHTTP(w, r)
			}
		}),
		&http2.Server{},
	)
}
