package handler

import (
	"context"

	connect_go "github.com/bufbuild/connect-go"
	"github.com/ras0q/connect-web-playground/pkg/bufgen/api/proto"
	"github.com/ras0q/connect-web-playground/pkg/bufgen/api/proto/protoconnect"
	"google.golang.org/protobuf/types/known/emptypb"
)

type readyHandler struct{}

func NewReadyHandler() protoconnect.ReadyServiceHandler {
	return &readyHandler{}
}

func (h *readyHandler) Ready(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[proto.ReadyResponse], error) {
	return &connect_go.Response[proto.ReadyResponse]{
		Msg: &proto.ReadyResponse{
			Ready: true,
		},
	}, nil
}

func (h *readyHandler) Unready(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[proto.ReadyResponse], error) {
	return &connect_go.Response[proto.ReadyResponse]{
		Msg: &proto.ReadyResponse{
			Ready: false,
		},
	}, nil
}
