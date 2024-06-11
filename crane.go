package main

import (
	"fmt"
	"log"
	"os"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/logs"
)

func init() {
	logs.Warn.SetOutput(os.Stderr)
	logs.Progress.SetOutput(os.Stderr)
}

func copy_image(src string, dst string) error {
	log.Printf("COPY EVENT: Copy %s to %s", src, dst)
	if err := crane.Copy(src, dst, crane.WithAuthFromKeychain(authn.DefaultKeychain)); err != nil {
		fmt.Printf("log.Logger: %s", err.Error())
		return err
	}
	return nil
}
