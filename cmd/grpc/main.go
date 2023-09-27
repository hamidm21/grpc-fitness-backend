/*
 *
 * Copyright 2018 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// The server demonstrates how to consume and validate OAuth2 tokens provided by
// clients for each RPC.
package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/sirupsen/logrus"

	"gitlab.com/mefit/mefit-server/controller"

	"gitlab.com/mefit/mefit-server/utils/config"
	"gitlab.com/mefit/mefit-server/utils/log"

	pb "gitlab.com/mefit/mefit-api/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	// "google.golang.org/grpc/credentials"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"gitlab.com/mefit/mefit-server/utils"
	"gitlab.com/mefit/mefit-server/utils/initializer"
	"gitlab.com/mefit/mefit-server/utils/signal"
)

var customFunc grpc_recovery.RecoveryHandlerFunc = func(p interface{}) error {
	log.Logger().Errorf("Recovered from panic: %v", p)
	return status.Errorf(codes.Internal, "something unexpected happens")
}

func runHealthCheck() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		log.Logger().Info("/healthz returns ok")
		fmt.Fprintf(w, "ok!")
	})
	log.Logger().Info("Starting healthz daemon on port 10256 ...")
	go http.ListenAndServe(":10256", nil)
}
func main() {
	defer initializer.Initialize()()
	debugMode := config.Config().GetBool(utils.KeyDebug)
	runHealthCheck()
	listenPort := fmt.Sprintf(":%s", config.Config().GetString(utils.KeyPort))
	log.Logger().Infof("server starting on %s...", listenPort)
	// creds, err := credentials.NewServerTLSFromFile(utils.GRPCKeyPair())
	// if err != nil {
	// 	log.Logger().Panic("can't load grpc cert/key")
	// }

	//////////// Interceptor options///////////////////////////////////////////////////////////
	// Logrus entry is used, allowing pre-definition of certain fields by the user.
	logrusEntry := logrus.NewEntry(log.Logger())
	// Shared options for the logger, with a custom duration to log field function.
	logOpts := []grpc_logrus.Option{
		grpc_logrus.WithDurationField(func(duration time.Duration) (key string, value interface{}) {
			return "grpc.time_ns", duration.Nanoseconds()
		}),
	}
	// Shared options for the logger, with a custom gRPC code to log level function.
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(customFunc),
	}
	//////////// Interceptor options///////////////////////////////////////////////////////////
	var s *grpc.Server
	if debugMode {
		s = grpc.NewServer(
			// grpc.Creds(creds),
			grpc_middleware.WithUnaryServerChain(
				showReq,
				ensureValidToken,
			),
		)
	} else {
		s = grpc.NewServer(
			// grpc.Creds(creds),
			grpc_middleware.WithUnaryServerChain(
				ensureValidToken,
				grpc_prometheus.UnaryServerInterceptor,
				// grpc_ctxtags.UnaryServerInterceptor(),
				grpc_logrus.UnaryServerInterceptor(logrusEntry, logOpts...),
				grpc_recovery.UnaryServerInterceptor(opts...),
			),
			grpc_middleware.WithStreamServerChain(
				grpc_prometheus.StreamServerInterceptor,
				// grpc_ctxtags.StreamServerInterceptor(),
				grpc_logrus.StreamServerInterceptor(logrusEntry, logOpts...),
				grpc_recovery.StreamServerInterceptor(opts...),
			),
		)
	}

	pb.RegisterMefitServer(s, &controller.Controller{})
	// pb_admin.RegisterAdminServer(s, &controller.Controller{})

	lis, err := net.Listen("tcp", listenPort)
	if err != nil {
		log.Logger().Fatalf("failed to listen: %v", err)
	}
	if err := s.Serve(lis); err != nil {
		log.Logger().Fatalf("failed to serve: %v", err)
	}

	<-signal.OnExit()
	print("app closed")
}

//Debug purpose
func showReq(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Logger().Debugf("info: %v, req: %v", spew.Sdump(info), spew.Sdump(req))
	res, err := handler(ctx, req)
	log.Logger().Debugf("res: %v", spew.Sdump(res))
	return res, err
}

// ensureValidToken ensures a valid token exists within a request's metadata. If
// the token is missing or invalid, the interceptor blocks execution of the
// handler and returns an error. Otherwise, the interceptor invokes the unary
// handler.
func ensureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if strings.HasSuffix(info.FullMethod, "SignIn") || strings.HasSuffix(info.FullMethod, "SignUp") || strings.HasSuffix(info.FullMethod, "AnonySignUp") || strings.HasSuffix(info.FullMethod, "ResetPassword") {
		return handler(ctx, req)
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, utils.ErrMissingMetadata
	}
	// The keys within metadata.MD are normalized to lowercase.
	// See: https://godoc.org/google.golang.org/grpc/metadata#New
	if newCtx, ok := utils.AuthValid(ctx, md); ok {
		// Continue execution of handler after ensuring a valid token.
		return handler(newCtx, req)
	}
	//Invalid creds
	return nil, utils.ErrInvalidToken
}
