package main

import (
	"bufio"
	"errors"
	"fmt"
	"go/types"
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
//	if b.lastBackupTime == nil {
//		//Get from repo
//	}

	err := os.Setenv("BORG_PASSPHRASE", "iewei1ahdeij9ni8geChieKee3ohm3")
	if err != nil {
		return time.Time{}, err
	}

	borgOut, err := exec.Command(b.borgBin, "list", "--json", b.path).Output()
	if err != nil {
		return time.Time{}, err
	}

	return b.lastBackupTime, nil
}

func NewBorgRepo(path string, borgBin string, key string) (*BorgRepo, error) {
	if ok, err := IsValidBorgRepo(path); !ok {
		return &BorgRepo{}, err
	}

	return &BorgRepo{path: path, borgBin: borgBin, key: key}, nil
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
