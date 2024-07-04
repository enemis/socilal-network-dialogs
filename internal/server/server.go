package server

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"social-network-dialogs/internal/config"
	"social-network-dialogs/internal/dialog"
	"social-network-dialogs/internal/logger"
	pb "social-network-dialogs/internal/server/proto"
)

type Server struct {
	server  *grpc.Server
	config  *config.Config
	logger  logger.LoggerInterface
	dialogs *dialog.DialogService
	pb.UnimplementedDialogServiceServer
}

func New(logger logger.LoggerInterface, config *config.Config, dialog *dialog.DialogService) *Server {
	server := Server{
		server:                           grpc.NewServer(),
		config:                           config,
		logger:                           logger,
		UnimplementedDialogServiceServer: pb.UnimplementedDialogServiceServer{},
		dialogs:                          dialog,
	}

	return &server
}

func InitHooks(lc fx.Lifecycle, server *Server, logger logger.LoggerInterface) {
	var listener net.Listener
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			listener, err := net.Listen("tcp", server.config.GRPCAddress)
			if err != nil {
				logger.Fatal(fmt.Sprintf("failed to listen %s:", server.config.GRPCAddress), err, nil)
			}

			pb.RegisterDialogServiceServer(server.server, server)
			log.Printf("gRPC server listening at %v", listener.Addr())
			go func() {
				if err := server.server.Serve(listener); err != nil {
					logger.Fatal("failed to serve grpc:", err, nil)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			server.server.GracefulStop()
			err := listener.Close()
			return err
		},
	})
}

func (s *Server) SendMessage(_ context.Context, in *pb.SendMessageDialogRequest) (*pb.SendMessageDialogResponse, error) {
	message, err := s.dialogs.SendDirectMessage(in.Sender.Id.Value, in.Reciever.Id.Value, in.Message)
	if err != nil {
		return nil, err
	}

	return &pb.SendMessageDialogResponse{Message: &pb.Message{
		Id:        &pb.UUID{Value: message.Id.String()},
		UserId:    &pb.UUID{Value: message.UserId.String()},
		DialogId:  &pb.UUID{Value: message.DialogId.String()},
		CreateAt:  timestamppb.New(message.CreatedAt),
		DeletedAt: nil,
		Message:   message.Message,
	}}, nil
}
