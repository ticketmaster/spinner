Get-ChildItem .\build\ -Directory | ForEach-Object {Push-Location $_.FullName -StackName build}

while ((Get-Location -Stack -StackName build -ErrorAction SilentlyContinue).count -gt 0) {
    $osarch = $pwd | Split-Path -leaf
    Compress-Archive .\spinner* "..\spinner_$osarch-v$env:APPVEYOR_BUILD_VERSION.zip"
    Pop-Location -StackName build
}