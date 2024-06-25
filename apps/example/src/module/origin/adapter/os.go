package adapter

import (
	"log"

	"github.com/spf13/afero"
)

func NewOSFs(basePath string) (afero.Fs, error) {
	osFs := afero.NewOsFs()

	log.Printf("init basePath %+v", basePath)

	fs := afero.NewBasePathFs(osFs, basePath)

	exists, err := afero.DirExists(osFs, basePath)
	if err != nil {
		return nil, err
	}
	if !exists {
		err := osFs.MkdirAll(basePath, 0755)
		if err != nil {
			return nil, err
		}
	}

	return fs, nil
}
