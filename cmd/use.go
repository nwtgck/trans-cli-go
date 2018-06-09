package cmd

import (
  "fmt"
  "os"

  "github.com/spf13/cobra"
  "github.com/spf13/viper"
  "github.com/nwtgck/trans-cli-go/constants"
)


func init() {
  RootCmd.AddCommand(useCmd)
}

var useCmd = &cobra.Command{
  Use:   "use",
  Short: "Set server URL",
  Run: func(cmd *cobra.Command, args []string) {

    // Check the length of arguments
    if len(args) != 1 {
      fmt.Fprint(os.Stderr, "Error: Specify <URL>\n")
      os.Exit(1)
    }

    // Get server URL
    serverUrlStr   := args[0]

    // Set server URL
    viper.Set(constants.ServerUrlKey, serverUrlStr)
    // Save config
    err := viper.WriteConfig()
    if err != nil {
      fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
      os.Exit(1)
    }

    // Print message
    fmt.Printf("Set '%s'\n", serverUrlStr)

  },
}
