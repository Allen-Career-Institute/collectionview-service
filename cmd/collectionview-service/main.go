package main

import (
	"collectionview-service/internal/utils"
	"flag"
	"fmt"
	"os"

	"collectionview-service/internal/conf"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

const (
	SandboxEnv    = "sandbox"
	DevEnv        = "dev"
	StageEnv      = "stage"
	ProdEnv       = "prod"
	DockerConfDir = "/app/data/conf/"
	LocalConfDir  = "collectionview-service/configs/"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = "collectionview-service"
	// Version is the version of the compiled software.
	Version = os.Getenv("tag")
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	fmt.Println("init called ::")
	flag.StringVar(&flagconf, "conf", utils.JoinStrings(LocalConfDir, "config_local.yaml"), "config path, eg: -conf config_local.yaml")
	fmt.Println("init successful ::")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	c := config.New(
		config.WithSource(
			file.NewSource(getConfigPath()),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bc.Server, bc.Data, logger, bc.Redis)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func getConfigPath() string {
	env := os.Getenv("ENV")
	log.Infof("using env  : " + env)

	configFile := flagconf
	switch env {
	case SandboxEnv:
		configFile = utils.JoinStrings(LocalConfDir, "config_sandbox.yaml")
		log.Infof("using sandbox config")
	case DevEnv:
		configFile = utils.JoinStrings(DockerConfDir, "config_dev.yaml")
		log.Infof("using dev config")
	case StageEnv:
		configFile = utils.JoinStrings(DockerConfDir, "config_stage.yaml")
		log.Infof("using stage config")
	case ProdEnv:
		configFile = utils.JoinStrings(DockerConfDir, "config_prod.yaml")
		log.Infof("using stage config")
	default:
		log.Infof("using local config : %s", configFile)
	}
	return configFile
}
