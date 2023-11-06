package src

import (
	"os"
	"runtime/pprof"
)

func StartCPUProfile(profileFilePath string) (*os.File, error) {
	profileFile, err := os.OpenFile(profileFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	err = pprof.StartCPUProfile(profileFile)
	if err != nil {
		profileFile.Close()
		return nil, err
	}

	return profileFile, nil
}

func StopCPUProfile(profileFile *os.File) error {
	pprof.StopCPUProfile()
	return profileFile.Close()
}
