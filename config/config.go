package config

import (
	"flag"
	"os"
)

type Config struct {
	FromStdin bool
	Filenames []string
}

func GetConfig() (*Config, error) {
	// get input filenames from console
	filenames := flag.Args()

	// check if all the input files exist
	fromStdin := false
	if len(filenames) == 0 {
		fromStdin = true
	} else {
		for _, filename := range filenames {
			if _, err := os.Stat(filename); err == nil || os.IsNotExist(err) {
				return nil, err
			}
		}
	}
	return &Config{
		FromStdin: fromStdin,
		Filenames: filenames,
	}, nil
}
