<img align="right" src="spinner.png" alt="Spinner" />

# Spinner

[![Build status](https://ci.appveyor.com/api/projects/status/1us41ajgvlwu9dcb?svg=true)](https://ci.appveyor.com/project/cdhunt/spinner)

Spinner, A Service Monitor and file tail executable for Windows Containers.

## Usage

```text
spinner [command]

Available Commands:
  help        Help about any command
  service     Watch a Windows Service
  site        Watch a Site
  version     All software has versions. This is Spinner's.

Flags:
      --config string   config file (default is $HOME/.spinner.yaml)
  -d, --debug           Print debug logging
  -h, --help            help for spinner
  -t, --tail string     Path to file to tail and pipe to STDOUT
  -o, --out             Path to file to write response body if site response status >= 300

Service Usage:
  spinner service [name] [flags]

Site Usage:
  spinner site [url] [counter limit] [time to exit] [flags]
```


## Examples

```powershell
  spinner.exe service W3SVC -t c:\\iislog\\W3SVC\\u_extend1.log
  spinner.exe site http://localhost -t c:\\iislog\\W3SVC\\u_extend1.log
```

## Note

Spinner will exit if the file does not already exist. Make sure to generate an event
and flush the buffer immediately or wait for the file to generate before starting
Spinner.

```dockerfile
  CMD Start-Service W3SVC; `
    Invoke-WebRequest http://localhost -UseBasicParsing | Out-Null; `
    netsh http flush logbuffer | Out-Null; `
    spinner.exe service W3SVC -t c:\\iislog\\W3SVC\\u_extend1.log
```

The Linux build is experimental and does not include the `service` command.
