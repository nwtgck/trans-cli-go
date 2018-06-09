# trans-cli

Command Line Interface for [Trans](https://github.com/nwtgck/trans-server-akka) Server

| branch | Travis status|
| --- | --- |
| [`master`](https://github.com/nwtgck/trans-cli-go/tree/master) |[![Build Status](https://travis-ci.com/nwtgck/trans-cli-go.svg?token=TuxNpqznwwyy7hyJwBVm&branch=master)](https://travis-ci.com/nwtgck/trans-cli-go) |
| [`develop`](https://github.com/nwtgck/trans-cli-go/tree/develop) | [![Build Status](https://travis-ci.com/nwtgck/trans-cli-go.svg?token=TuxNpqznwwyy7hyJwBVm&branch=develop)](https://travis-ci.com/nwtgck/trans-cli-go) |

## Usages

```bash
# Send a file
trans send ~/Documents/hello.txt
```

```bash
# Get a file
trans get d84
```

```bash
# Delete the file
trans delete d84
```

(File ID, `d84` is save as `d84` file in `$PWD`)
