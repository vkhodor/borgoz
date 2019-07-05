package main

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"os"
	"github.com/labstack/echo"
	"errors"
	"strconv"
)

var Version string = "v0.0.0"

type Application struct {
	Config *Configuration
	Echo   *echo.Echo
	Logger *log.Logger
}

type Configuration struct {
	BorgBin        string
	Host           string
	Port           int
	ReposDirectory string
	DefaultRepoKey string
	LogLevel       log.Lvl
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

	borgBin, ok := os.LookupEnv(EnvVariableBorgBin)
	if !ok {
		borgBin = DefaultBorgBin
	}

	return Configuration{
		Host:           host,
		Port:           port,
		ReposDirectory: reposDirectory,
		DefaultRepoKey: defaultRepoKey,
		LogLevel:       String2LogLevel(logLevel),
		BorgBin:        borgBin,
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

	app := Application{Config: &cfg, Echo: e, Logger: log.New(LogPrefix + " " + Version)}
	app.Logger.SetLevel(app.Config.LogLevel)
	app.Logger.Debugf("%v", app.Config)

	e.GET("/check/lastBackupTime/:repo/:time/:key", app.handlerBackupNotOlderThen)
	e.GET("/check/lastBackupTime/:repo/:time", app.handlerBackupNotOlderThen)

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
	borgRepo, err := NewBorgRepo(fmt.Sprintf("%v/%v", a.Config.ReposDirectory, repo), a.Config.BorgBin, key, a.Logger)
	if err != nil {
		a.Logger.Errorf("NewBorgRepo returned: %v", err)
		return echo.NewHTTPError(404, fmt.Sprintf("%v is not valid borg repo", repo))
	}

	a.Logger.Debugf("%v is valid borg repo", repo)

	intTime, err := strconv.Atoi(time)
	if err != nil {
		return echo.NewHTTPError(400, fmt.Sprintf("Can't parse string %v to int: %v", time, err))
	}

	if ok, err := borgRepo.IsLastBackupEarlierThen(intTime); !ok {
		msg := fmt.Sprintf("%v has too old last backup", repo)
		if err != nil {
			msg = fmt.Sprintf("%v: IsLastBackupEarlierThen returned error: %v", repo, err)
		}
		return echo.NewHTTPError(404, msg)
	}
	a.Logger.Debugf("%v: IsLastBackupEarlierThen returned OK", repo)

	return echo.NewHTTPError(200, fmt.Sprintf("%v LastBackupTime is OK", repo))
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
