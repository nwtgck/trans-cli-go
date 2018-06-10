package cmd

import (
  "fmt"
  "os"
  "net/url"

  "github.com/spf13/cobra"
  "github.com/spf13/viper"
  "github.com/nwtgck/trans-cli-go/settings"
  "github.com/k0kubun/pp"
)

// Set server
func setServer(serverUrlOrAlias string ){
  // Server URL
  var serverUrlStr string
  // If valid url
  if _, err := url.ParseRequestURI(serverUrlOrAlias); err == nil {
    serverUrlStr = serverUrlOrAlias
  } else {
    // Set server aliases
    serverAliases := toArrayMapStringString(viper.Get(settings.ServerAliasesKey))

    // If alias not found
    if u := findServerUrl(serverAliases, serverUrlOrAlias); u == nil {
      fmt.Fprintf(os.Stderr, "Error: Alias '%s' not found or invalid URL\n", serverUrlOrAlias)
      os.Exit(1)
    } else {
      // Set server URL
      serverUrlStr = *u
    }
  }

  // Set server URL
  viper.Set(settings.ServerUrlKey, serverUrlStr)
  // Save config
  err := viper.WriteConfig()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
    os.Exit(1)
  }

  // Print message
  fmt.Printf("Set '%s'\n", serverUrlStr)
}

func init() {
  RootCmd.AddCommand(configCmd)
  configCmd.AddCommand(configServerCmd, configAliasCmd)
}

var configCmd = &cobra.Command{
  Use:   "config",
  Short: "Configure settings",
  Run: func(cmd *cobra.Command, args []string) {
    // Print help
    cmd.Help()
  },
}

// Find server URL from aliases
func findServerUrl(serverAliases []map[string]string, alias string) *string {
  for _, a := range serverAliases {
    if a[settings.ServerAliasesNameKey] == alias {
      serverUrl := a[settings.ServerAliasesUrlKey]
      return &serverUrl
    }
  }
  return nil
}

// `config server` command
var configServerCmd =  &cobra.Command{
  Use:   "server",
  Short: "Set/Show Server URL",
  Run: func(cmd *cobra.Command, args []string) {
    if len(args) == 1 {
      // Get server URL or alias
      serverUrlOrAlias   := args[0]

      setServer(serverUrlOrAlias)
    } else {
      // If server URL is set
      if viper.IsSet(settings.ServerUrlKey) {
        // Set server URL
        serverUrl := viper.GetString(settings.ServerUrlKey)
        fmt.Println(serverUrl)
      } else {
        fmt.Fprint(os.Stderr, "Error: Server URL is not found\n")
      }
    }
  },
}

// Replace new alias if exist
// Add an alias otherwise
func updateAliases(serverAliases []map[string]string, alias map[string]string) []map[string]string {
  result := []map[string]string{}
  for _, a := range serverAliases {
    if a[settings.ServerAliasesNameKey] != alias[settings.ServerAliasesNameKey] {
      result = append(result, a)
    }
  }
  return append(result, alias)
}

// Convert interface{} => []map[string]string
// (deep casting)
func toArrayMapStringString(a interface{}) []map[string]string {
  result := []map[string]string{}
  for _, e := range a.([]interface{}) {
    e2 := e.(map[string]interface{})
    m  := map[string]string{}
    for k, v := range e2 {
      m[k] = v.(string)
    }
    result = append(result, m)
  }
  return result
}

// `config alias` command
var configAliasCmd =  &cobra.Command{
  Use:   "alias",
  Short: "Set/Show server alias",
  Run: func(cmd *cobra.Command, args []string) {
    if len(args) == 2 {
      serverAlias := args[0]
      serverUrl   := args[1]

      // Create an alias
      alias := map[string]string{
        settings.ServerAliasesNameKey: serverAlias,
        settings.ServerAliasesUrlKey: serverUrl,
      }

      var serverAliases []map[string]string
      // If server_aliases is set
      if viper.IsSet(settings.ServerAliasesKey) {
        serverAliases = toArrayMapStringString(viper.Get(settings.ServerAliasesKey))
      } else {
        serverAliases = []map[string]string{}
      }

      // Update server aliases
      serverAliases = updateAliases(serverAliases, alias)

      // Set the aliases
      viper.Set(settings.ServerAliasesKey, serverAliases)
      // Save config
      err := viper.WriteConfig()
      if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
        os.Exit(1)
      }

      // Print message
      fmt.Printf("Set '%s' = '%s'\n", serverAlias, serverUrl)
    } else {
      // If server_aliases is set
      if viper.IsSet(settings.ServerAliasesKey) {
        // Set server aliases
        serverAliases := toArrayMapStringString(viper.Get(settings.ServerAliasesKey))
        // Print server aliases
        pp.Println(serverAliases)
      } else {
        fmt.Println("No server aliases")
      }
    }
  },
}