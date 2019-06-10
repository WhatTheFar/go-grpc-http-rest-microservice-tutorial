package cmd

import (
	"context"
	"database/sql"
	"flag"
	"fmt"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"

	"github.com/amsokol/go-grpc-http-rest-microservice-tutorial/config"
	"github.com/amsokol/go-grpc-http-rest-microservice-tutorial/pkg/logger"
	"github.com/amsokol/go-grpc-http-rest-microservice-tutorial/pkg/protocol/grpc"
	"github.com/amsokol/go-grpc-http-rest-microservice-tutorial/pkg/protocol/rest"
	v1 "github.com/amsokol/go-grpc-http-rest-microservice-tutorial/pkg/service/v1"
)

type Config struct {
	Env string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	// get configuration
	var cfg Config
	flag.StringVar(&cfg.Env, "env", "", "environment to build")
	flag.Parse()
	config.InitViper("../../config", cfg.Env)
	v := config.GetViper()

	if len(v.Protocol.Grpc.Port) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", v.Protocol.Grpc.Port)
	}

	if len(v.Protocol.Http.Port) == 0 {
		return fmt.Errorf("invalid TCP port for HTTP gateway: '%s'", v.Protocol.Http.Port)
	}

	// initialize logger
	if err := logger.Init(v.Logging.LogLevel, v.Logging.LogTimeFormat); err != nil {
		return fmt.Errorf("failed to initialize logger: %v", err)
	}

	// add MySQL driver specific parameter to parse date/time
	// Drop it for another database
	param := "parseTime=true"

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		v.Database.MySqlDB.Username,
		v.Database.MySqlDB.Password,
		v.Database.MySqlDB.Host,
		v.Database.MySqlDB.DBName,
		param)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	v1API := v1.NewToDoServiceServer(db)
	v1HealthAPI := v1.NewHealthcheckServiceServer()

	// run HTTP gateway
	go func() {
		_ = rest.RunServer(ctx, v.Protocol.Grpc.Port, v.Protocol.Http.Port)
	}()

	return grpc.RunServer(ctx, v1API, v1HealthAPI, v.Protocol.Grpc.Port)
}
