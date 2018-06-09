package cmd

import (
  "fmt"
  "os"
  "net/url"
  "net/http"
  "path"
  "io"

  "github.com/spf13/cobra"
  "gopkg.in/cheggaaa/pb.v2"
  "github.com/spf13/viper"
  "github.com/nwtgck/trans-cli-go/settings"
)


// Disable progress bar or not
var getQuiet bool
// Outputs a file to stdout or not
var outputsToStdout bool

func init() {
  RootCmd.AddCommand(getCmd)
  getCmd.Flags().BoolVarP(&getQuiet, "quiet", "q", false, "disable progress bar or not")
  getCmd.Flags().BoolVar(&outputsToStdout, "stdout", false, "outputs a file to stdout")
}

var getCmd = &cobra.Command{
  Use:   "get",
  Short: "Download a file",
  Long:  "Download a file",
  Run: func(cmd *cobra.Command, args []string) {
    // TODO: Extract command parts

    // If server URL is not set
    if !viper.IsSet(settings.ServerUrlKey) {
      fmt.Fprint(os.Stderr, "Error: Server URL is not found\n")
      os.Exit(1)
    }

    // Get server URL
    serverUrlStr := viper.GetString(settings.ServerUrlKey)

    // Check validity of server URL
    serverUrl, err := url.Parse(serverUrlStr)
    if err != nil {
      // Exit if URL is not valid
      fmt.Fprintf(os.Stderr, "Error: Server URL '%s' is not valid\n", serverUrlStr)
      os.Exit(1)
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

    // Download the file
    resp, err := http.Get(serverUrl.String())

    // Create a file
    fileFileName := fileId

    // Output file
    var outFile io.WriteCloser

    // If outputs to stdout
    if outputsToStdout {
      outFile = os.Stdout
    } else {
      outFile, err = os.Create(fileFileName)
      defer outFile.Close()
      if err != nil {
        fmt.Fprintf(os.Stderr, "Error: Cannot open '%s'\n", fileFileName)
      }
    }

    var reader io.Reader
    var bar *pb.ProgressBar = nil
    if getQuiet {
      reader = resp.Body
    } else {
      // Create a bar
      bar = pb.New64(resp.ContentLength)
      bar.Start()
      reader = bar.NewProxyReader(resp.Body)
    }

    // Save the file
    io.Copy(outFile, reader)

    if bar != nil {
      // Finish progress bar
      bar.Finish()
    }
  },
}
