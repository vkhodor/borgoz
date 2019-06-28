package main

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"os"
	"github.com/labstack/echo"
	"errors"
	"strconv"
)

type Application struct {
	Config *Configuration
	Echo   *echo.Echo
	Logger *log.Logger
}

type Configuration struct {
	Host string
	Port int
	ReposDirectory string
	DefaultRepoKey string
	LogLevel log.Lvl
}

func String2LogLevel(logLevel string) log.Lvl {
	switch logLevel {
		case "DEBUG":
			return log.DEBUG
		case "INFO":
			return log.INFO
		case "WARN":
			return log.WARN
		case "ERROR":
			return log.ERROR
		default:
			return log.OFF
	}
}

func NewConfiguration() (Configuration, error){

	logLevel, ok := os.LookupEnv(EnvVariableLogLevel)
	if !ok {
		logLevel = DefaultLogLevel
	}

	host, ok := os.LookupEnv(EnvVariableHost)
	if !ok {
		host = DefaultHost
	}

	stringPort, ok := os.LookupEnv(EnvVariablePort)
	if !ok {
		stringPort = DefaultPort
	}

	port, err := strconv.Atoi(stringPort)
	if err != nil {
		return Configuration{}, errors.New(
			fmt.Sprintf(
				"Can't %v=%v cant convert to Integer",
						EnvVariablePort, stringPort,
				),
			)
	}

	reposDirectory, ok := os.LookupEnv(EnvVariableReposDirectory)
	if !ok {
		reposDirectory = DefaultReposDirectory
	}

	defaultRepoKey, ok := os.LookupEnv(EnvVariableDefaultRepoKey)
	if !ok {
		return Configuration{}, errors.New(
			fmt.Sprintf("Env variable %v is not set", EnvVariableDefaultRepoKey),
			)
	}

	return Configuration{
		Host: host,
		Port: port,
		ReposDirectory: reposDirectory,
		DefaultRepoKey: defaultRepoKey,
		LogLevel: String2LogLevel(logLevel),
	}, nil
}

func NewApplication() (Application, error) {

	cfg, err := NewConfiguration()
	if err != nil {
		return Application{}, err
	}

	e := echo.New()
	if cfg.LogLevel != log.DEBUG {
		e.HideBanner = true
	}
	
	app := Application{Config: &cfg, Echo: e, Logger: log.New(LogPrefix)}
	app.Logger.SetLevel(app.Config.LogLevel)
	app.Logger.Debugf("%v", app.Config)

	e.GET("/check/backup/:repo/not_older_then/:time/:key", app.handlerBackupNotOlderThen)
	e.GET("/check/backup/:repo/not_older_then/:time", app.handlerBackupNotOlderThen)

	return app, nil
}

func (a *Application)handlerBackupNotOlderThen(c echo.Context) error {
	repo := c.Param("repo")
	time := c.Param("time")
	key := c.Param("key")
	if key == "" {
		key = a.Config.DefaultRepoKey
	}

	a.Logger.Infof("repo=%v time=%v key=%v", repo, time, key)
	return nil
}

func (a *Application) Start() error {
	err := a.Echo.Start(fmt.Sprintf("%v:%v", a.Config.Host, a.Config.Port))
	return err
}

func main() {
	app, err := NewApplication()
	if err != nil {
		log.Fatal(err)
		os.Exit(ErrorExitNewApplication)
	}

	log.Fatal(app.Start())

}
