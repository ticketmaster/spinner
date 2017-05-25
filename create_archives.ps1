Get-ChildItem .\build\ -Directory | ForEach-Object {Push-Location $_.FullName -StackName build}

while ((Get-Location -Stack -StackName build).count -gt 1) {
    $pwd
    $osarch = $pwd | Split-Path -leaf
    Compress-Archive .\spinner_v1.0.0* "..\spinner_$osarch-v$env:APPVEYOR_BUILD_VERSION.zip"
    Pop-Location -StackName build
}

Pop-Location -StackName build

ls .\build -rec | select FullName