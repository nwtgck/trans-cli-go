package cmd

import (
  "fmt"
  "os"
  "net/http"
  "path"
  "io"

  "github.com/spf13/cobra"
  "gopkg.in/cheggaaa/pb.v2"
  "github.com/Code-Hex/pget"
  "github.com/nwtgck/trans-cli-go/util"
  "net/url"
)


// Disable progress bar or not
var getQuiet bool
// Outputs a file to stdout or not
var outputsToStdout bool
// Parallel download  or not
var parallel bool

func init() {
  RootCmd.AddCommand(getCmd)
  getCmd.Flags().BoolVarP(&getQuiet, "quiet", "q", false, "disable progress bar or not")
  getCmd.Flags().BoolVar(&outputsToStdout, "stdout", false, "outputs a file to stdout")
  getCmd.Flags().BoolVarP(&parallel, "parallel", "p", false, "enable parallel download")
}

var getCmd = &cobra.Command{
  Use:   "get",
  Short: "Download a file",
  Long:  "Download a file",
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

    // FIXME: Allow users to use both --quiet and --parallel
    // If --quiet and --parallel
    if getQuiet && parallel {
      // Print warning message
      fmt.Fprint(os.Stderr, "Warning: Disable '--quiet'\n")
    }

    // FIXME: Allow users to use both --stdout and --parallel
    // If --stdout and --parallel
    if outputsToStdout && parallel {
      // Print warning message
      fmt.Fprint(os.Stderr, "Warning: Disable '--stdout'\n")
    }

    // Get file ID
    fileId := args[0]

    // Join server url with file ID
    serverUrl.Path = path.Join(serverUrl.Path, fileId)

    // Create a file
    fileFileName := fileId


    // If parallel download is enable
    if parallel {
      // pget setting
      pg := pget.New()
      pg.URLs = []string{serverUrl.String()}
      pg.TargetDir = ""
      pg.Utils.SetFileName(fileFileName)


      if err := pg.Checking(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %s\n", err)
        os.Exit(1)
      }

      if err := pg.Download(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %s\n", err)
        os.Exit(1)
      }

      if err := pg.Utils.BindwithFiles(pg.Procs); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %s\n", err)
        os.Exit(1)
      }
    } else {
      // Download the file
      resp, err := http.Get(serverUrl.String())

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
    }
  },
}
