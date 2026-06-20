package cmd

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var workspace string

var rootCmd = &cobra.Command{
	Use: "esol",
	Short: "A CLI tool to analyze and generate statistics for your " +
		"Exercism solutions.",
}

func Execute(version string) error {
	rootCmd.Version = version
	return rootCmd.Execute()
}

func init() {
	viper.AddConfigPath(configDir())
	viper.SetConfigName("user")
	viper.SetConfigType("json")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	workspace = viper.GetString("workspace")
	if workspace == "" {
		panic("workspace not configured")
	}
}

func executableDir() string {
	executable, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(executable)
}

func configDir() string {
	dir := ""
	if runtime.GOOS == "windows" {
		dir = os.Getenv("APPDATA")
		if dir != "" {
			return filepath.Join(dir, "exercism")
		}
	} else {
		dir := os.Getenv("EXERCISM_CONFIG_HOME")
		if dir != "" {
			return dir
		}
		dir = os.Getenv("XDG_CONFIG_HOME")
		if dir == "" {
			dir = filepath.Join(os.Getenv("HOME"), ".config")
		}
		if dir != "" {
			return filepath.Join(dir, "exercism")
		}
	}
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return dir
}
