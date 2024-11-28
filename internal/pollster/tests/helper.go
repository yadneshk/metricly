package tests

import (
	"fmt"
	"os"
)

func SetupCollectorSources(fileName, fileContent string) error {

	colltrFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create mounts file: %s", err)
	}

	if _, err = colltrFile.Write([]byte(fileContent)); err != nil {
		return fmt.Errorf("failed to write mock content: %s", err)
	}
	colltrFile.Close()

	return nil
}
