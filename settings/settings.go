package settings

// (NOTE: This variable can be replaced in BUILD TIME)
var DefaultServerUrl = "https://trans-akka.herokuapp.com"

const (
  ConfigName       = "config"
  ConfigExt        = "json"
  ServerUrlEnvName = "TRANS_SERVER_URL"

  ServerUrlKey         = "server_url"
  ServerAliasesKey     = "server_aliases"
  ServerAliasesNameKey = "name"
  ServerAliasesUrlKey  = "url"
)
