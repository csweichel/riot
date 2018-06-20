[![Gitpod.io](https://img.shields.io/badge/gitpod.io-master-green.svg)](https://gitpod.io/github.com/32leaves/master)
[![Build Status](https://travis-ci.org/32leaves/riot.svg?branch=master)](https://travis-ci.org/32leaves/riot)

# riot
As simple-as-they-come docker orchestrator targeting IoT and the Raspberry Pi.
It supports building and deploying _applications_ to _nodes_. That's it.
It requires a registry that the nodes can push and pull from.
It does not require any agents/processes/deployments on the nodes.

## Usage
```
$ riot
As simple-as-they-come docker orchestrator targeting IoT and the Raspberry Pi.

Usage:
  riot [command]

Available Commands:
  build       Builds all applications of this project
  deploy      Deploys all applications of this project
  help        Help about any command
  init        Initializes this directory as a riot project
  status      Displays the status of all applications and their deployment
  version     Prints the version of riot
  vet         Validates a riot project

Flags:
      --alsologtostderr                  log to standard error as well as files
  -h, --help                             help for riot
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --logtostderr                      log to standard error instead of files
      --project string                   riot project directory (default is the current working directory)
      --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
  -v, --v Level                          log level for V logs
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging

Use "riot [command] --help" for more information about a command.
```