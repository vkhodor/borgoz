package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const BorgRepoDetectionString = "This is a Borg Backup repository."

type BorgRepo struct {
	borgBin        string
	path           string
	key            string
	lastBackupTime time.Time
	lastBackupId string
	logger 			*log.Logger
}

type BorgBackup struct {
	Archive string `json:"archive"`
	BArchive string `json:"barchive"`
	Id string `json:"id"`
	Name string `json:"name"`
	Start string `json:"start"`
	Time string `json:"time"`
}

type BorgEncryption struct {
	Mode string `json:"mode"`
}

type BorgRepository struct {
	Id string `json:"id"`
	LastModified string `json:"last_modified"`
	Location string `json:"location"`
}

type BorgBackupList struct {
	Archives []BorgBackup `json:"archives"`
	Encryption BorgEncryption `json:"encryption"`
	Repository BorgRepository `json:"repository"`
}

func (b *BorgRepo) GetLastBorgBackupTime() (time.Time, error) {
	err := os.Setenv(EnvVariableBorgPassphrase, b.key)
	defer os.Unsetenv(EnvVariableBorgPassphrase)


	if err != nil {
		b.logger.Errorf("%v", err)
		return time.Time{}, err
	}

	borgOut, err := exec.Command(b.borgBin, "list", "--json", b.path).Output()
	b.logger.Debugf("%v list --json %v", b.borgBin, b.path)
	
	if err != nil {
		b.logger.Errorf("%v", err)
		b.logger.Debugf("%v", borgOut)
		return time.Time{}, err
	}

	backupOut := BorgBackupList{}
	err = json.Unmarshal(borgOut, &backupOut)
	if err != nil {
		b.logger.Errorf("%v", err)
		return time.Time{}, err
	}

	b.lastBackupTime, err = ParseTimeInCurrentLocation(backupOut.Repository.LastModified)
	if err != nil {
		b.logger.Errorf("%v", err)
		return time.Time{}, err
	}

	return b.lastBackupTime, nil
}

func (b *BorgRepo) IsLastBackupEarlierThen(seconds int) (bool, error) {
	lastBackupTime, err := b.GetLastBorgBackupTime()
	if err != nil {
		return false, err
	}
	now := time.Now()
	passedDuration := time.Duration(now.Unix() - lastBackupTime.Unix()) * time.Second

	if passedDuration > time.Duration(seconds) * time.Second {
		return false, nil
	}

	return true, nil
}

func NewBorgRepo(path string, borgBin string, key string, logger *log.Logger) (*BorgRepo, error) {
	if ok, err := IsValidBorgRepo(path); !ok {
		return &BorgRepo{}, err
	}

	return &BorgRepo{path: path, borgBin: borgBin, key: key, logger: logger}, nil
}

func IsValidBorgRepo(path string) (bool, error) {
	readmefile, err := os.Open(fmt.Sprintf("%v/README", path))
	if err != nil {
		return false, err
	}

	scanner := bufio.NewScanner(readmefile)

	scanner.Scan()
	line := scanner.Text()

	if !strings.Contains(line, BorgRepoDetectionString) {
		return false, errors.New(fmt.Sprintf("%v not contains %v", line, BorgRepoDetectionString))
	}

	return true, nil
}
