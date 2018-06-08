package main

import (
  "fmt"
  "os"

  "github.com/nwtgck/trans-cli-go/cmd"
)

func main() {
  if err := cmd.RootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(-1)
  }
}
