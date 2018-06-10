package cmd

import (
  "fmt"
  "os"

  "github.com/spf13/cobra"
  "github.com/spf13/viper"
  "github.com/nwtgck/trans-cli-go/settings"
)


// Unset value or not
//var unSet bool

func init() {
  RootCmd.AddCommand(configCmd)
  //configCmd.Flags().BoolVarP(&unSet, "unset", "u",false, "unset setting")
  configCmd.AddCommand(configServerCmd)
}

var configCmd = &cobra.Command{
  Use:   "config",
  Short: "Configure settings",
  Run: func(cmd *cobra.Command, args []string) {
    // Print help
    cmd.Help()
  },
}

// `config server` command
var configServerCmd =  &cobra.Command{
  Use:   "server",
  Short: "Configure/Show Server URL",
  Run: func(cmd *cobra.Command, args []string) {
    if len(args) == 1 {
      // Get server URL
      serverUrlStr   := args[0]

      // Set server URL
      viper.Set(settings.ServerUrlKey, serverUrlStr)
      // Save config
      err := viper.WriteConfig()
      if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
        os.Exit(1)
      }

      // Print message
      fmt.Printf("Set '%s'\n", serverUrlStr)


    } else {
      // If server URL is set
      if viper.IsSet(settings.ServerUrlKey) {
        // Set server URL
        serverUrl := viper.GetString(settings.ServerUrlKey)
        fmt.Println(serverUrl)
      } else {
        fmt.Fprint(os.Stderr, "Error: Server URL is not found\n")
      }
    }
  },
}