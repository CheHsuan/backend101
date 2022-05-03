package phonebook

import (
	"context"
	commpb "grpc-http-server/pb/common"
	pb "grpc-http-server/pb/phonebook"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serviceServerImpl struct {
	pb.UnimplementedPhonebookServer
	phonebook map[string]string
}

func NewPhonebookServer() pb.PhonebookServer {
	return &serviceServerImpl{
		phonebook: map[string]string{},
	}
}

func (s *serviceServerImpl) CreatePhoneRecord(ctx context.Context, req *pb.CreatePhoneRecordRequest) (*pb.CreatePhoneRecordResponse, error) {
	s.phonebook[req.GetName()] = req.GetPhone()
	return &pb.CreatePhoneRecordResponse{
		Status: &commpb.Status{
			Code:    "200",
			Message: "success",
		},
	}, nil
}

func (s *serviceServerImpl) QueryPhoneRecord(ctx context.Context, req *pb.QueryPhoneRecordRequest) (*pb.QueryPhoneRecordResponse, error) {
	phone, ok := s.phonebook[req.GetName()]
	if !ok {
		return nil, status.Error(codes.NotFound, "invalid name")
	}
	return &pb.QueryPhoneRecordResponse{
		Status: &commpb.Status{
			Code:    "200",
			Message: "success",
		},
		Name:  req.GetName(),
		Phone: phone,
	}, nil
}
