<img align="right" src="spinner.png" alt="Spinner" />

# Spinner

Spinner, A Service Monitor and file tail executable for Windows Containers.

## Usage

```text
  spinner.exe flags
The flags are:
  --service <name>    Service name to Watch
  --path <path>       Path to log file to tail
  --usage             Print usage
  --debug             Print debug logging
  --version           Print version of this binary
```

## Examples

```powershell
  # Watch IIS and tail the access log:
  .\spinner.exe -service W3SVC -path c:\iislog\W3SVC\u_extend1.log
```

## Note

Spinner will exit if the file does not already exist. Make sure to generate an event
and flush the buffer immediately or wait for the file to generate before starting
Spinner.

```dockerfile
  CMD Start-Service W3SVC; `
    Invoke-WebRequest http://localhost -UseBasicParsing | Out-Null; `
    netsh http flush logbuffer | Out-Null; `
    spinner.exe -service W3SVC -path c:\iislog\W3SVC\u_extend1.log
```