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

// Duration of file storing
var duration string

// Download limit of a file
var getTimes int

// Length of file ID
var idLength int

// File is can delete or not
var deletable bool

// Password for deletion
var deleteKey string

func init() {
  RootCmd.AddCommand(sendCmd)

  // Flags
  sendCmd.Flags().StringVar(&duration, "duration", "1h", "store duration (e.g. 10s, 5m, 12h, 3d)")
  sendCmd.Flags().IntVarP(&getTimes, "get-times",  "t", 100, "download limit (e.g. 1, 10)")
  sendCmd.Flags().IntVarP(&idLength, "id-length", "l",3, "length of ID (e.g. 3, 10)")
  sendCmd.Flags().BoolVar(&deletable, "deletable", true, "whether file is deletable or not")
  sendCmd.Flags().StringVar(&deleteKey, "delete-key", "", "key for delete")
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

    // (from: https://qiita.com/yyoshiki41/items/a0354d9ad70c1b8225b6)
    q := serverUrl.Query()
    q.Set("duration", duration)
    q.Set("get-times", fmt.Sprintf("%d", getTimes))
    q.Set("id-length", fmt.Sprintf("%d", idLength))
    q.Set("deletable", fmt.Sprintf("%t", deletable))
    if deleteKey != "" {
      q.Set("delete-key", deleteKey)
    }
    serverUrl.RawQuery = q.Encode()

    resp, err := http.Post(serverUrl.String(), "application/octet-stream", file)
    fileIdBytes, _ := ioutil.ReadAll(resp.Body)
    fileId := strings.TrimRight(string(fileIdBytes), "\n")
    fmt.Println(fileId)

  },
}
