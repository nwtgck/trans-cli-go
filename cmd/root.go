package cmd

import (
  "os"
  "github.com/spf13/cobra"
  "fmt"

  "github.com/nwtgck/trans-cli-go/version"
)

var RootCmd = &cobra.Command{
  Use:     os.Args[0],
  Short:   "Trans CLI",
  Long:    "Trans CLI",
  Version: version.Version,
  Run: func(cmd *cobra.Command, args []string) {
    cmd.Help()
  },
}

func init() {
  cobra.OnInitialize()
  RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
  Use:   "version",
  Short: "Show version",
  Long:  `Show version`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println(version.Version)
  },
}
