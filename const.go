package main

const LogPrefix = "BorGoZ"

const (
	EnvVariableHost           = "BORGOZ_HOST"
	EnvVariablePort           = "BORGOZ_PORT"
	EnvVariableReposDirectory = "BORGOZ_REPOS_DIR"
	EnvVariableDefaultRepoKey = "BORGOZ_DEFAULT_REPO_KEY"
	EnvVariableLogLevel       = "BORGOZ_LOG_LEVEL"
	EnvVariableBorgBin        = "BORGOZ_BORG_BIN"
)

const (
	DefaultHost           = "127.0.0.1"
	DefaultPort           = "8080"
	DefaultReposDirectory = "./"
	DefaultLogLevel       = "INFO"
	DefaultBorgBin        = "borg"
)

const ErrorExitNewApplication int = 1 // Exit when Fatal error during NewApplication
