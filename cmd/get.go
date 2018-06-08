package cmd

import (
	"fmt"
	"os"
	"net/url"
	"net/http"
	"path"
	"io"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Download a file",
	Long:  "Download a file",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Extract command parts (Almost part is the same as send command)

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
		outFile, err := os.Create(fileFileName)
		if err != nil {
			fmt.Fprint(os.Stderr, "Error: Cannot open '%s'\n", fileFileName)
		}

		// Save the file
		io.Copy(outFile, resp.Body)
	},
}
