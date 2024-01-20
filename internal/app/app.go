package app

import (
	"fmt"
	"grpc-sso/internal/db"
	"grpc-sso/internal/jwt"
	"grpc-sso/internal/sso"
	"grpc-sso/pkg/utils"
	"log"
	"net"

	"google.golang.org/grpc"
)

func Run() {
	config := utils.MustLoadConfig()
	serverPort := fmt.Sprintf("%d", config.Server.Port)

	listener, err := net.Listen("tcp", ":"+serverPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Init UserService
	us, err := db.NewUserService(config.Database.ConnectionString)
	if err != nil {
		log.Fatalf("unable to init UserService")
	}

	// Init JWTService
	jwts := jwt.NewJWTService(
		config.JWT.SecretKey,
		config.JWT.AccessTokenExpiry,
		config.JWT.RefreshTokenExpiry,
	)

	// Create service
	ssoService := sso.NewSSOService(jwts, us)

	// Register service to Server
	sso.RegisterSSOServiceServer(grpcServer, ssoService)

	log.Println("Server is running on port: " + serverPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
