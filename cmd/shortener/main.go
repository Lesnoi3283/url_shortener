package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/gRPC/grpchandlers"
	"github.com/Lesnoi3283/url_shortener/internal/app/gRPC/interceptors"
	"github.com/Lesnoi3283/url_shortener/internal/app/gRPC/proto"
	"github.com/Lesnoi3283/url_shortener/internal/app/handlers"
	"github.com/Lesnoi3283/url_shortener/internal/app/logic"
	"github.com/Lesnoi3283/url_shortener/pkg/databases"
	"github.com/Lesnoi3283/url_shortener/pkg/secure"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// naOrValue returns "N/A" if v contains a default value. Returns v if not.
func naOrValue(v string) string {
	if v == "" {
		return "N/A"
	} else {
		return v
	}
}

// gracefulShutdown listens for os signals syscall.SIGTERM, syscall.SIGINT and syscall.SIGQUIT.
// Calls HTTPServer.Shutdown and gRPCServer.GracefulStop if signal received.
// This func have to be called in different goroutine, because it has an endless loop.
func gracefulShutdown(HTTPServer *http.Server, gRPCServer *grpc.Server, log zap.SugaredLogger, wg *sync.WaitGroup) {
	defer wg.Done()
	shutDownCh := make(chan os.Signal, 1)
	signal.Notify(shutDownCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for v := range shutDownCh {
		log.Infof("Received an os signal '%s', graceful shutting down...", v.String())
		err := HTTPServer.Shutdown(context.Background())
		if err != nil {
			log.Error("failed to shutdown gracefully", zap.Error(err))
		}
		gRPCServer.GracefulStop()
		close(shutDownCh)
	}
}

func main() {

	//print version data
	fmt.Printf("Build version: %s\n", naOrValue(buildVersion))
	fmt.Printf("Build date: %s\n", naOrValue(buildDate))
	fmt.Printf("Build commit: %s\n", naOrValue(buildCommit))

	//conf
	conf := config.Config{}
	err := conf.Configure()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	//storages set
	var URLStore logic.URLStorageInterface
	if conf.DBConnString != "" {
		var err error
		URLStore, err = databases.NewPostgresql(conf.DBConnString)
		if err != nil {
			log.Fatalf("Problem with starting postgresql: %v", err.Error())
		}
	} else if conf.FileStoragePath != "" {
		URLStore = databases.NewJSONFileStorage(conf.FileStoragePath)
	} else {
		URLStore = databases.NewJustAMap()
	}

	//logger set
	logLevel, err := zap.ParseAtomicLevel(conf.LogLevel)
	if err != nil {
		log.Fatalf("logger was not started, err: %v", err)
	}
	zCfg := zap.NewProductionConfig()
	zCfg.Level = logLevel
	zCfg.DisableStacktrace = true
	zapLogger, err := zCfg.Build()
	if err != nil {
		log.Fatalf("logger was not started, err: %v", err)
	}
	sugar := zapLogger.Sugar()

	//JWTHelper set
	JWTHelper := secure.NewJWTHelper(conf.JWTSecret, conf.JWTTimeoutHours)

	//HTTP server building
	r, err := handlers.NewRouter(conf, URLStore, *sugar, JWTHelper)
	if err != nil {
		sugar.Fatalf("Error creating new router: %v", err)
	}
	var server *http.Server
	wg := &sync.WaitGroup{}

	if conf.EnableHTTPS {
		//HTTPS server starting
		sugar.Info("Starting HTTPS server")

		manager := autocert.Manager{
			Cache:      autocert.DirCache("cache-dir"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("urlshortener.ru"),
		}
		server = &http.Server{
			Addr:      conf.ServerAddress,
			Handler:   r,
			TLSConfig: manager.TLSConfig(),
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			err := server.ListenAndServeTLS("", "")
			if err != nil {
				sugar.Errorf("server error: %v", err)
			}
		}()

	} else {
		//HTTP server starting
		sugar.Info("Starting HTTP server")
		server = &http.Server{
			Addr:    conf.ServerAddress,
			Handler: r,
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			err := server.ListenAndServe()
			if err != nil {
				sugar.Errorf("server error: %v", err)
			}
		}()
	}

	//run gRPC
	gRPCServer, err := runGRPCServer(&conf, URLStore, *sugar, JWTHelper)
	if err != nil {
		sugar.Fatalf("Error starting gRPC server: %v", err)
	}

	//graceful shutdown
	wg.Add(1)
	go gracefulShutdown(server, gRPCServer, *sugar, wg)
	wg.Wait()
}

// runGRPCServer creates and runs a new gRPC server. Calls logger.Fatal if starting gRPC is not possible.
func runGRPCServer(conf *config.Config, storage logic.URLStorageInterface, logger zap.SugaredLogger, jh *secure.JWTHelper) (*grpc.Server, error) {
	listen, err := net.Listen("tcp", conf.GRPCAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to listen gRPC: %v", err)
	}

	//prepare interceptors
	trustedSubnet := &net.IPNet{}
	if conf.TrustedSubnet != "" {
		var err error
		_, trustedSubnet, err = net.ParseCIDR(conf.TrustedSubnet)
		if err != nil {
			return nil, fmt.Errorf("error parsing trusted subnet: %w", err)
		}
	}

	//prepare gRPC server
	grpcServer := &grpc.Server{}
	if conf.EnableHTTPS {
		logger.Info("Preparing TLS gRPC server")

		manager := autocert.Manager{
			Cache:      autocert.DirCache("cache-dir"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("urlshortener.ru"),
		}

		grpcServer = grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				interceptors.NewIPInterceptor(trustedSubnet),
				interceptors.NewUnaryAuthInterceptor(jh),
			),
			grpc.Creds(credentials.NewTLS(&tls.Config{
				GetCertificate: manager.GetCertificate,
			})),
		)
	} else {
		logger.Info("Preparing gRPC server (no TLS)")
		grpcServer = grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				interceptors.NewIPInterceptor(trustedSubnet),
				interceptors.NewUnaryAuthInterceptor(jh),
			),
		)
	}

	proto.RegisterURLShortenerServiceServer(grpcServer, &grpchandlers.ShortenerServer{
		Storage: storage,
		Logger:  logger,
		Conf:    conf,
	})

	//start gRPC server
	logger.Info("Starting gRPC server...")
	go func() {
		err := grpcServer.Serve(listen)
		if err != nil {
			logger.Fatalf("failed to start gRPC server: %v", err)
		}
	}()
	return grpcServer, nil
}
