package erdtree

import (
	"app-websocket/gen/erdtree/v1/erdtreev1connect"
	"net/http"
)

func New() (erdtreev1connect.ErdtreeStoreClient, error) {
	return erdtreev1connect.NewErdtreeStoreClient(
		http.DefaultClient,
		"http://host.docker.internal:50051", // Server URL
	), nil
}
