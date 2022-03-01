package main

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"log"
	"path/filepath"
	"time"
)

// SaveTargetIAMProfiles saves a file which contains UserName, Access Key ID, Lifetime of the target IAM profiles.
func SaveTargetIAMProfiles(given []IAMProfile) error {
	// logging
	log.Print(given)

	// marshaling phase
	buf, _ := json.MarshalIndent(given, "", "    ")
	
	// write phase
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	result, err := os.Create(filepath.Join(cwd, "tmp", string([]byte(time.Now().String())[:19]) + ".json"))
	if err != nil {
		return err
	}
	defer result.Close()

	w := bufio.NewWriter(result)
	for _, b := range buf {
		err := w.WriteByte(b)
		if err == io.EOF {
			break
		}
	}
	w.Flush()

	return nil
}
