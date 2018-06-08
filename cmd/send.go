package cmd

import (
  "fmt"
  "os"
  "net/url"
  "net/http"
  "io/ioutil"
  "strings"

  "github.com/spf13/cobra"
)

func init() {
  RootCmd.AddCommand(sendCmd)
}

var sendCmd = &cobra.Command{
  Use:   "send",
  Short: "Send a file",
  Long:  "Send a file",
  Run: func(cmd *cobra.Command, args []string) {

    // TODO: Hard code: ENV NAME
    const TRANS_SERVER_URL_NAME = "TRANS_SERVER_URL"

    // Get server URL
    serverUrlStr, exist := os.LookupEnv(TRANS_SERVER_URL_NAME)

    // If $TRANS_SERVER_URL does not exist
    if !exist {
      // Emit an error and exit
      fmt.Fprintf(os.Stderr, "Error: Set $%s properly\n", TRANS_SERVER_URL_NAME)
      os.Exit(1)
    }

    // Check validity of server URL
    serverUrl, err := url.Parse(serverUrlStr)
    if err != nil {
      // Exit if URL is not valid
      fmt.Fprint(os.Stderr, "Error: Server URL '%s' is not valid\n", serverUrlStr)
      os.Exit(1)
    }

    // Check the length of arguments
    if len(args) != 1 {
      fmt.Fprint(os.Stderr, "Error: Specify a file\n")
      os.Exit(1)
    }

    // Open the first file
    file, err := os.Open(args[0])
    if err != nil {
      fmt.Fprintf(os.Stderr, "Error: Cannot open '%s'\n", args[0])
      os.Exit(1)
    }

    resp, err := http.Post(serverUrl.String(), "application/octet-stream", file)
    fileIdBytes, _ := ioutil.ReadAll(resp.Body)
    fileId := strings.TrimRight(string(fileIdBytes), "\n")
    fmt.Println(fileId)

  },
}
