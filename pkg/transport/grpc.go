package transport

import (
	"context"
	"watermark-service/api/v1/pb"
	"watermark-service/internal"
	"watermark-service/pkg/endpoints"

	grpctransport "github.com/go-kit/kit/transport/grpc"
)

// grpcServer is a struct responsible for creating gRPC servers and handling endpoint mappings.
type grpcServer struct {
	get           grpctransport.Handler
	status        grpctransport.Handler
	addDocument   grpctransport.Handler
	watermark     grpctransport.Handler
	serviceStatus grpctransport.Handler
}

// NewGRPCServer initializes a gRPC server for the provided endpoints.
func NewGRPCServer(ep endpoints.Set) pb.WatermarkServer {
	return &grpcServer{
		get: grpctransport.NewServer(
			ep.GetEndpoint,
			decodeGRPCGetRequest,
			decodeGRPCGetResponse,
		),
		status: grpctransport.NewServer(
			ep.StatusEndpoint,
			decodeGRPCStatusRequest,
			decodeGRPCStatusResponse,
		),
		addDocument: grpctransport.NewServer(
			ep.AddDocumentEndpoint,
			decodeGRPCAddDocumentRequest,
			decodeGRPCAddDocumentResponse,
		),
		watermark: grpctransport.NewServer(
			ep.WatermarkEndpoint,
			decodeGRPCWatermarkRequest,
			decodeGRPCWatermarkResponse,
		),
		serviceStatus: grpctransport.NewServer(
			ep.ServiceStatusEndpoint,
			decodeGRPCServiceStatusRequest,
			decodeGRPCServiceStatusResponse,
		),
	}
}

// Functions for implementing gRPC service endpoints.

func (g *grpcServer) Get(ctx context.Context, r *pb.GetRequest) (*pb.GetReply, error) {
	_, rep, err := g.get.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetReply), nil
}

func (g *grpcServer) ServiceStatus(ctx context.Context, r *pb.ServiceStatusRequest) (*pb.ServiceStatusReply, error) {
	_, rep, err := g.get.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ServiceStatusReply), nil
}

func (g *grpcServer) AddDocument(ctx context.Context, r *pb.AddDocumentRequest) (*pb.AddDocumentReply, error) {
	_, rep, err := g.addDocument.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.AddDocumentReply), nil
}

func (g *grpcServer) Status(ctx context.Context, r *pb.StatusRequest) (*pb.StatusReply, error) {
	_, rep, err := g.status.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.StatusReply), nil
}

func (g *grpcServer) Watermark(ctx context.Context, r *pb.WatermarkRequest) (*pb.WatermarkReply, error) {
	_, rep, err := g.watermark.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.WatermarkReply), nil
}

func decodeGRPCGetRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetRequest)
	var filters []internal.Filter
	for _, f := range req.Filters {
		filters = append(filters, internal.Filter{Key: f.Key, Value: f.Value})
	}
	return endpoints.GetRequest{Filters: filters}, nil
}

func decodeGRPCStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.StatusRequest)
	return endpoints.StatusRequest{TicketID: req.TicketID}, nil
}

func decodeGRPCWatermarkRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.WatermarkRequest)
	return endpoints.WatermarkRequest{TicketID: req.TicketID, Mark: req.Mark}, nil
}

func decodeGRPCAddDocumentRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.AddDocumentRequest)
	doc := &internal.Document{
		Content:   req.Document.Content,
		Title:     req.Document.Title,
		Author:    req.Document.Author,
		Topic:     req.Document.Topic,
		Watermark: req.Document.Watermark,
	}
	return endpoints.AddDocumentRequest{Document: doc}, nil
}

func decodeGRPCServiceStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return endpoints.ServiceStatusRequest{}, nil
}

func decodeGRPCGetResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.GetReply)
	var docs []internal.Document
	for _, d := range reply.Documents {
		doc := internal.Document{
			Content:   d.Content,
			Title:     d.Title,
			Author:    d.Author,
			Topic:     d.Topic,
			Watermark: d.Watermark,
		}
		docs = append(docs, doc)
	}
	return endpoints.GetResponse{Documents: docs, Err: reply.Err}, nil
}

func decodeGRPCStatusResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.StatusReply)
	return endpoints.StatusResponse{Status: internal.Status(reply.Status), Err: reply.Err}, nil
}

func decodeGRPCWatermarkResponse(ctx context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.WatermarkReply)
	return endpoints.WatermarkResponse{Code: int(reply.Code), Err: reply.Err}, nil
}

func decodeGRPCAddDocumentResponse(ctx context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.AddDocumentReply)
	return endpoints.AddDocumentResponse{TicketID: reply.TicketID, Err: reply.Err}, nil
}

func decodeGRPCServiceStatusResponse(ctx context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.ServiceStatusReply)
	return endpoints.ServiceStatusResponse{Code: int(reply.Code), Err: reply.Err}, nil
}
