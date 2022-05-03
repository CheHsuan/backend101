package app

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "grpc-http-server/pb/phonebook"
	"grpc-http-server/service/phonebook"
)

func (app *Application) initGRPCServer(ctx context.Context) error {
	app.GRPCServer = grpc.NewServer()
	pb.RegisterPhonebookServer(app.GRPCServer, phonebook.NewPhonebookServer())
	reflection.Register(app.GRPCServer)

	return nil
}
