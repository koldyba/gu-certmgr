package certificates

import (
	"context"
	"errors"

	"golang-united-certificates/pkg/api"
	"golang-united-certificates/pkg/db"
	"golang-united-certificates/pkg/helpers"

	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct {
	api.UnimplementedCertificatesServer
	Database db.DB
}

func (srv *GRPCServer) Get(ctx context.Context, req *api.GetRequest) (*api.GetResponse, error) {
	cert, err := srv.Database.GetCertById(req.CertId)
	return &api.GetResponse{Certificate: helpers.WriteApiCert(cert)}, err
}

func (srv *GRPCServer) Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error) {
	if srv.Database.IsCertExistsByUserAndCourse(req.UserId, req.CourseId) {
		return &api.CreateResponse{}, errors.New("Cert for this user for this course was already Created")
	}
	cert, err := srv.Database.Create(req.UserId, req.CourseId)
	if err != nil {
		return nil, err
	}
	return &api.CreateResponse{Certificate: helpers.WriteApiCert(cert)}, nil
}

func (srv *GRPCServer) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	data, npt, err := srv.Database.List(int(req.PageSize), req.PageToken)
	resp := []*api.Cert{}
	for _, cert := range data {
		resp = append(resp, helpers.WriteApiCert(cert))
	}
	return &api.ListResponse{Certificates: resp, NextPageToken: npt}, err
}

func (srv *GRPCServer) ListForUser(ctx context.Context, req *api.ListForUserRequest) (*api.ListResponse, error) {
	data, npt, err := srv.Database.ListForUser(int(req.PageSize), req.PageToken, req.UserId)
	resp := []*api.Cert{}
	for _, cert := range data {
		resp = append(resp, helpers.WriteApiCert(cert))
	}
	return &api.ListResponse{Certificates: resp, NextPageToken: npt}, err
}
func (srv *GRPCServer) ListForCourse(ctx context.Context, req *api.ListForCourseRequest) (*api.ListResponse, error) {
	data, npt, err := srv.Database.ListForCourse(int(req.PageSize), req.PageToken, req.CourseId)
	resp := []*api.Cert{}
	for _, cert := range data {
		resp = append(resp, helpers.WriteApiCert(cert))
	}
	return &api.ListResponse{Certificates: resp, NextPageToken: npt}, err
}
func (srv *GRPCServer) Delete(ctx context.Context, req *api.DeleteRequest) (*emptypb.Empty, error) {
	srv.Database.Delete(req.CertId)
	return &emptypb.Empty{}, nil
}
