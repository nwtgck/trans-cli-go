package util

import (
  "github.com/spf13/viper"
  "fmt"
  "os"
  "github.com/nwtgck/trans-cli-go/settings"
  "github.com/asaskevich/govalidator"
  "github.com/pkg/errors"
)

// Get URL from config and validate it
func GetUrlAndValidate() (string, error) {
  // If server URL is not set
  if !viper.IsSet(settings.ServerUrlKey) {
    fmt.Fprint(os.Stderr, "Error: Server URL is not found\n")
    os.Exit(1)
  }
  // Get server URL
  serverUrlStr := viper.GetString(settings.ServerUrlKey)

  // Check whether the string is URL or not
  if !govalidator.IsRequestURL(serverUrlStr) {
    return serverUrlStr, errors.Errorf("'%s' is not URL\n", serverUrlStr)
  }

  return serverUrlStr, nil
}