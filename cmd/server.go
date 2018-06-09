package cmd

import (
  "github.com/spf13/cobra"
  "fmt"
  "os"
  "github.com/spf13/viper"
  "github.com/nwtgck/trans-cli-go/settings"
)

// Show URL or not
var showUrl bool

func init() {
  RootCmd.AddCommand(serverCmd)
  serverCmd.Flags().BoolVar(&showUrl, "show", false, "show server URL")
}

var serverCmd = &cobra.Command{
  Use:   "server",
  Short: "Set server URL",
  Run: func(cmd *cobra.Command, args []string) {

    if showUrl {
      // If server URL is set
      if viper.IsSet(settings.ServerUrlKey) {
        // Set server URL
        serverUrl := viper.GetString(settings.ServerUrlKey)
        fmt.Println(serverUrl)
      } else {
        fmt.Fprint(os.Stderr, "Error: Server URL is not found\n")
      }

    } else {
      // Check the length of arguments
      if len(args) != 1 {
        fmt.Fprint(os.Stderr, "Error: Specify <URL>\n")
        os.Exit(1)
      }

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
    }
  },
}
