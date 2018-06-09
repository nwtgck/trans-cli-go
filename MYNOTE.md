# My Note

## Set Different Default Server URL

Here is an example to embed default server URL, <https://trans-akka.mybluemix.net/>, in build time.

```bash
go build -ldflags "-X github.com/nwtgck/trans-cli-go/settings.DefaultServerUrl=https://trans-akka.mybluemix.net/"
```