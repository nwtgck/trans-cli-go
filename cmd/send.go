package cmd

import (
  "fmt"
  "os"
  "net/url"
  "net/http"
  "io"
  "io/ioutil"
  "strings"

  "github.com/spf13/cobra"
  "gopkg.in/cheggaaa/pb.v2"
  "golang.org/x/crypto/ssh/terminal"
  "github.com/mholt/archiver"
  "github.com/nwtgck/trans-cli-go/util"
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

// Whether secure character is used for File ID
var usesSecureChar bool

// Disable progress bar or not
var sendQuiet bool

func init() {
  RootCmd.AddCommand(sendCmd)

  // Flags
  sendCmd.Flags().StringVar(&duration, "duration", "1h", "store duration (e.g. 10s, 5m, 12h, 3d)")
  sendCmd.Flags().IntVarP(&getTimes, "get-times",  "t", 100, "download limit (e.g. 1, 10)")
  sendCmd.Flags().IntVarP(&idLength, "id-length", "l",3, "length of ID (e.g. 3, 10)")
  sendCmd.Flags().BoolVar(&deletable, "deletable", true, "whether file is deletable or not")
  sendCmd.Flags().StringVar(&deleteKey, "delete-key", "", "key for deletion")
  sendCmd.Flags().BoolVar(&usesSecureChar, "secure-char", false, "uses more complex characters for File ID")
  sendCmd.Flags().BoolVarP(&sendQuiet, "quiet", "q", false, "disable progress bar or not")
}

var sendCmd = &cobra.Command{
  Use:   "send",
  Short: "Send a file",
  Long:  "Send a file",
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


    // Input (file or piped stdin)
    var input io.Reader
    // Total number of bar
    var barTotal int64
    // If pipe is not used
    if terminal.IsTerminal(0) {
      // Check the length of arguments
      if len(args) != 1 {
        fmt.Fprint(os.Stderr, "Error: Specify a file\n")
        os.Exit(1)
      }

      // File or Directory path
      fileOrDirPath := args[0]

      // Get file or dir info
      fileInfo, err := os.Stat(fileOrDirPath)
      if err != nil {
        fmt.Fprintf(os.Stderr, "Error: Canot get file info\n")
        os.Exit(1)
      }

      // (from: https://stackoverflow.com/a/8824952/2885946)
      switch mode := fileInfo.Mode(); {
      // If it's a file
      case mode.IsRegular():
        // Open the first file
        file, err := os.Open(fileOrDirPath)
        defer file.Close()
        // Assign input as file
        input = file
        if err != nil {
          fmt.Fprintf(os.Stderr, "Error: Cannot open '%s'\n", fileOrDirPath)
          os.Exit(1)
        }
        // Set bar total
        barTotal = fileInfo.Size()
      // If it's a directory
      case mode.IsDir():
        // Create a pipe
        pr, pw := io.Pipe()
        go func() {
          // Write zip to pipe writer
          err := archiver.Zip.Write(pw, []string{fileOrDirPath})
          if err != nil {
            fmt.Fprintf(os.Stderr, "Error: Cannot open '%s'\n", fileOrDirPath)
            os.Exit(1)
          }
          defer pw.Close()
        }()
        // Assign input as pipe input
        input = pr
        // Set bar total as 0
        // (FIXME: Specify properly)
        barTotal = 0
      }
    } else {
      // === Input from pipe ===

      // Set input as stdin
      input = os.Stdin
      // Set bar total as 0
      // (FIXME: Specify properly)
      barTotal = 0
    }


    // (from: https://qiita.com/yyoshiki41/items/a0354d9ad70c1b8225b6)
    q := serverUrl.Query()
    q.Set("duration", duration)
    q.Set("get-times", fmt.Sprintf("%d", getTimes))
    q.Set("id-length", fmt.Sprintf("%d", idLength))
    q.Set("deletable", fmt.Sprintf("%t", deletable))
    q.Set("secure-char", fmt.Sprintf("%t", usesSecureChar))
    if deleteKey != "" {
      q.Set("delete-key", deleteKey)
    }
    serverUrl.RawQuery = q.Encode()

    var reader io.Reader
    if sendQuiet {
      reader = input
    } else {
      // Create bar
      bar := pb.New64(barTotal)
      bar.Start()
      reader = bar.NewProxyReader(input)
    }

    resp, err := http.Post(serverUrl.String(), "application/octet-stream", reader)
    defer resp.Body.Close()
    fileIdBytes, _ := ioutil.ReadAll(resp.Body)
    fileId := strings.TrimRight(string(fileIdBytes), "\n")
    fmt.Println(fileId)

  },
}
