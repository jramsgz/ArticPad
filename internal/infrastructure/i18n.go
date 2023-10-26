package infrastructure

import (
	"fmt"
	"os"

	"github.com/jramsgz/articpad/pkg/i18n"
)

// startI18n starts the i18n service and loads the locales from the specified path
func startI18n(localesPath string) (*i18n.I18n, error) {
	i := i18n.New()

	if _, err := os.Stat(localesPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("i18n locales directory does not exist: %s", localesPath)
	}
	files, err := os.ReadDir(localesPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		// Read the file contents and load them into the i18n instance.
		b, err := os.ReadFile(localesPath + "/" + file.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to read i18n file: %s - %s", file.Name(), err.Error())
		}

		if err := i.Load(b, file.Name() == "en.json"); err != nil {
			return nil, fmt.Errorf("failed to load i18n file: %s - %s", file.Name(), err.Error())
		}
	}

	return i, nil
}
