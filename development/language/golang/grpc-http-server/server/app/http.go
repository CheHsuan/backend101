package app

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	commpb "grpc-http-server/pb/common"
	pb "grpc-http-server/pb/phonebook"
)

func (app *Application) initHTTPServer(ctx context.Context) error {
	app.HTTPServer = http.NewServeMux()
	gwmux := runtime.NewServeMux(
		runtime.WithErrorHandler(customErrorHandler),
	)
	if err := pb.RegisterPhonebookHandlerFromEndpoint(ctx, gwmux, app.Config.GetAddr(), []grpc.DialOption{grpc.WithInsecure()}); err != nil {
		return err
	}
	app.HTTPServer.Handle("/", gwmux)
	app.HTTPServer.Handle("/openapi/", http.StripPrefix("/openapi/", http.FileServer(http.Dir("./openapi"))))

	return nil
}

func customErrorHandler(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	const fallback = `{"error": "failed to marshal error message"}`

	w.WriteHeader(runtime.HTTPStatusFromCode(grpc.Code(err)))
	resp := &commpb.ErrorResponse{
		Status: &commpb.Status{
			Code:    grpc.Code(err).String(),
			Message: grpc.ErrorDesc(err),
		},
	}

	jErr := json.NewEncoder(w).Encode(resp)
	if jErr != nil {
		w.Write([]byte(fallback))
	}
}
