package grpcServer

import (
	"app-websocket/gen/erdtree/v1/erdtreev1connect"
	"app-websocket/gen/user/v1/userInfov1connect"
	"app-websocket/internal/user"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func NewAuthInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {

			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

type GrpcServer struct {
	mux             *http.ServeMux
	path            string
	corsHandler     http.Handler
	srv             *http.Server
	shutDownTimeout time.Duration
	logger          *slog.Logger
}

func New(logger *slog.Logger, erdtreeClient erdtreev1connect.ErdtreeStoreClient) (*GrpcServer, error) {
	// port := os.Getenv("PORT")

	user := user.NewExternalClient(erdtreeClient)
	mux := http.NewServeMux()
	interceptors := connect.WithInterceptors(NewAuthInterceptor())
	path, handler := userInfov1connect.NewUserInfoServiceHandler(user, interceptors)

	corsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Set-Cookie, connect-protocol-version")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			return
		}

		ctx := context.Background()
		// Call the next handler with the updated context
		handler.ServeHTTP(w, r.WithContext(ctx))
	})

	srv := &http.Server{
		Addr:    "0.0.0.0:4600",
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}
	shutDownTimeout := 10 * time.Second

	return &GrpcServer{
		mux,
		path,
		corsHandler,
		srv,
		shutDownTimeout,
		logger,
	}, nil
}

func (server *GrpcServer) Run(ctx context.Context) error {
	errResult := make(chan error)
	go func() {
		server.mux.Handle(server.path, server.corsHandler)

		server.logger.Info(fmt.Sprintf("starting listening: %s", server.srv.Addr))

		// if server.certFilePath != "" && server.keyFilePath != "" {
		// 	errResult <- server.srv.ListenAndServeTLS(server.certFilePath, server.keyFilePath)
		// }
		fmt.Println("SERVING AT: ", "0.0.0.0:")

		server.srv.ListenAndServe()
	}()

	var err error
	select {
	case <-ctx.Done():
		return ctx.Err()

	case err = <-errResult:
	}
	return err
}

func (server *GrpcServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), server.shutDownTimeout)
	defer cancel()
	err := server.srv.Shutdown(ctx)
	if err != nil {
		server.logger.Error("failed to shutdown HTTP Server", slog.String("error", err.Error()))
	}
}
