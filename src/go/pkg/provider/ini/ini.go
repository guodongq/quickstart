package ini

import (
	"gopkg.in/ini.v1"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/guodongq/quickstart/pkg/provider"
)

type Ini struct {
	provider.AbstractProvider

	options IniOptions
	keys    []string
}

func New(optionFuncs ...func(*IniOptions)) *Ini {
	defaultOptions := getDefaultIniOptions()
	options := &defaultOptions

	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}

	return &Ini{options: *options}
}

func (i *Ini) Init() error {
	if !i.options.Enabled {
		return nil
	}

	iniPath := i.filePath()
	file, err := ini.Load(iniPath)
	if err != nil {
		return err
	}
	sections := file.Section(i.options.Section)
	for _, key := range sections.Keys() {
		_ = os.Setenv(key.Name(), key.Value())
		i.keys = append(i.keys, key.Name())
	}
	return nil
}

func (i *Ini) Close() error {
	for _, key := range i.keys {
		_ = os.Unsetenv(key)
	}
	return nil
}

func (i *Ini) filePath() string {
	if i.options.IniFile != "" {
		return i.options.IniFile
	}

	curDir, _ := os.Getwd()
	iniPath := filepath.Join(curDir, "config.ini")
	if fileExists(iniPath) {
		return iniPath
	}

	binaryPath, _ := exec.LookPath(os.Args[0])
	execDir := filepath.Dir(binaryPath)
	iniPath = filepath.Join(execDir, "config.ini")
	if fileExists(iniPath) {
		return iniPath
	}

	panic("Invalid ini file path")
}

func fileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}
