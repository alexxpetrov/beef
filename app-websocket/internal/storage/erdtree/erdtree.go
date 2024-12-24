package erdtree

import (
	"app-websocket/gen/erdtree/v1/erdtreev1connect"
	"net/http"

	"connectrpc.com/connect"
)

func New() (erdtreev1connect.ErdtreeStoreClient, error) {
	// erdTreeUrl := "https://erdtree.fly.dev"
	erdTreeUrl := ""

	if erdTreeUrl == "" {
		erdTreeUrl = "http://host.docker.internal:50051"
	}

	return erdtreev1connect.NewErdtreeStoreClient(
		http.DefaultClient,
		erdTreeUrl, // Server URL
		connect.WithGRPCWeb(),
	), nil
}
