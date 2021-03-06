package cmd

import (
  "fmt"
  "os"
  "net/url"
  "net/http"
  "path"
  "io/ioutil"

  "github.com/spf13/cobra"
  "github.com/nwtgck/trans-cli-go/util"
)


// Password for deletion
// TODO: Rename
var key string

func init() {
  RootCmd.AddCommand(deleteCmd)
  deleteCmd.Flags().StringVar(&key, "delete-key", "", "key for deletion")
}

var deleteCmd = &cobra.Command{
  Use:   "delete",
  Short: "Delete a file",
  Run: func(cmd *cobra.Command, args []string) {
    // Check validity of server URL
    serverUrlStr, err := util.GetUrlAndValidate()
    if err != nil {
      // Exit if URL is not valid
      fmt.Fprintf(os.Stderr, "Error: Server URL '%s' is not valid\n", serverUrlStr)
      os.Exit(1)
    }
    // Convert URL string to server URL
    serverUrl, err := url.Parse(serverUrlStr)
    if err != nil {
      // NOTE: This will never happen because of the validation
      panic(err)
    }

    // Check the length of arguments
    if len(args) != 1 {
      fmt.Fprint(os.Stderr, "Error: Specify a file ID\n")
      os.Exit(1)
    }

    // Get file ID
    fileId := args[0]

    // Join server url with file ID
    serverUrl.Path = path.Join(serverUrl.Path, fileId)

    // (from: https://qiita.com/yyoshiki41/items/a0354d9ad70c1b8225b6)
    q := serverUrl.Query()
    if key != "" {
      q.Set("delete-key", key)
    }
    serverUrl.RawQuery = q.Encode()

    // Create deletion of the file request
    req, err := http.NewRequest("DELETE", serverUrl.String(), nil)
    if err != nil {
      fmt.Fprintf(os.Stderr, err.Error() + "\n")
      os.Exit(1)
    }

    // Delete the file
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
      fmt.Fprintf(os.Stderr, err.Error() + "\n")
      os.Exit(1)
    }

    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      fmt.Println(err.Error())
      os.Exit(1)
    }

    // Print message from server
    fmt.Print(string(respBody))
  },
}
