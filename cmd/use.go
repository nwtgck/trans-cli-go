package cmd

import (
  "github.com/spf13/cobra"
  "fmt"
  "os"
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
      fmt.Fprint(os.Stderr, "Error: Specify <URL|Alias>\n")
      os.Exit(1)
    }

    // Get server URL or alias
    serverUrlOrAlias   := args[0]

    // Set server
    setServer(serverUrlOrAlias)
  },
}
