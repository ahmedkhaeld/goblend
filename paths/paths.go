package paths

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const STD_CREATE_DIR_MODE = 0755

// paths
type Paths struct {
	Root        string   // root path of the application
	FolderNames []string // List of folder that expected to be available
}

func InitAppPaths(root string) error {
	var appFolders []string = []string{"app", "app/modules", "tmp", "env"}

	paths := Paths{
		Root:        root,
		FolderNames: appFolders,
	}

	for _, path := range paths.FolderNames {
		err := paths.CreateDirIfNotExist(root + "/" + path)
		if err != nil {
			return err
		}
	}

	err := paths.CheckDotEnv(root)
	if err != nil {
		return err
	}

	return nil

}

func (p *Paths) CreateDirIfNotExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, STD_CREATE_DIR_MODE)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Paths) CreateFileIfNotExist(path string) error {
	var _, err = os.Stat(path)
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return err
		}
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
	}
	return nil
}

func (p *Paths) CheckDotEnv(path string) error {

	err := p.CreateFileIfNotExist(fmt.Sprintf("%s/env/dev.env", path))
	if err != nil {
		return err
	}

	err = p.CreateFileIfNotExist(fmt.Sprintf("%s/env/prod.env", path))
	if err != nil {
		return err
	}
	return nil
}

func (p *Paths) LoadEnv(root string, filename string) error {
	// load the file according to the environment type

	err := godotenv.Load(root + "/env/" + filename + ".env")
	if err != nil {
		return err
	}

	return nil
}
