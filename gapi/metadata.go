package gapi

import (
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	grpcGatewayClientIP        = "x-forwarded-for"
	userAgentHeader            = "user-agent"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		} else if userAgents = md.Get(userAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		if clientIPs := md.Get(grpcGatewayClientIP); len(clientIPs) > 0 {
			mtdt.ClientIP = clientIPs[0]
		}
	}
	if len(mtdt.ClientIP) == 0 {
		if p, ok := peer.FromContext(ctx); ok {
			mtdt.ClientIP = p.Addr.String()
		}
	}
	return mtdt
}
