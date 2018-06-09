package cmd

import (
  "os"
  "github.com/spf13/cobra"
  "fmt"
  "path"
  "io/ioutil"

  "github.com/nwtgck/trans-cli-go/version"
  "github.com/spf13/viper"
  "github.com/nwtgck/trans-cli-go/settings"
)

// Server URL string
var ServerUrlStr string

// Path of config directory
var ConfigDirPath    = path.Join(os.Getenv("HOME"), ".config", "trans-cli")

var RootCmd = &cobra.Command{
  Use:     os.Args[0],
  Short:   "Trans CLI",
  Long:    "Trans CLI",
  Version: version.Version,
  Run: func(cmd *cobra.Command, args []string) {
    cmd.Help()
  },
}

// (from: http://yusukeiwaki.hatenablog.com/entry/2018/05/25/cobra/viper%E3%81%A7%E3%82%B3%E3%83%9E%E3%83%B3%E3%83%89%E3%83%A9%E3%82%A4%E3%83%B3%E3%83%84%E3%83%BC%E3%83%AB%E3%82%92%E4%BD%9C%E3%81%A3%E3%81%A6%E3%81%84%E3%81%A6%E8%89%AF%E3%81%8B%E3%81%A3)
func initConfig(){
  // Add config path
  viper.AddConfigPath(ConfigDirPath)
  // Set config file name
  viper.SetConfigName(settings.ConfigName)
  // Set config type
  viper.SetConfigType(settings.ConfigExt)
  // Set default server url
  viper.SetDefault(settings.ServerUrlKey, settings.DefaultServerUrl)

  if _, err := os.Stat(ConfigDirPath); os.IsNotExist(err) {
    // Make nested directories for config
    os.MkdirAll(ConfigDirPath, os.ModePerm)
  }

  // Path of config file
  defaultConfigFilePath := path.Join(ConfigDirPath, fmt.Sprintf("%s.%s", settings.ConfigName, settings.ConfigExt))
  // If config file doesn't exist
  if _, err := os.Stat(defaultConfigFilePath); os.IsNotExist(err) {
    // Write empty JSON file
    ioutil.WriteFile(defaultConfigFilePath, []byte("{}"), 0644)
  }
  // Bind env
  viper.BindEnv(settings.ServerUrlKey, settings.ServerUrlEnvName)

  // Read config file
  err := viper.ReadInConfig()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
    os.Exit(1)
  }

}

func init() {
  cobra.OnInitialize(initConfig)

  RootCmd.AddCommand(versionCmd)

  const serverFlag = "server"
  RootCmd.PersistentFlags().StringVarP(&ServerUrlStr, serverFlag, "s", "", "Trans Server URL")
  viper.BindPFlag(settings.ServerUrlKey, RootCmd.PersistentFlags().Lookup(serverFlag))
}

var versionCmd = &cobra.Command{
  Use:   "version",
  Short: "Show version",
  Long:  `Show version`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println(version.Version)
  },
}
