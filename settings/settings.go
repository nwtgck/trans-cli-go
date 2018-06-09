package settings

// (NOTE: This variable can be replaced in BUILD TIME)
var DefaultServerUrl = "https://trans-akka.herokuapp.com"

const (
  ConfigName       = "config"
  ConfigExt        = "json"
  ServerUrlKey     = "server_url"
  ServerUrlEnvName = "TRANS_SERVER_URL"
)
