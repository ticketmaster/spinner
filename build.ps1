$VERSION = '1.0.0'

Remove-Item .\build\*

glide up # https://glide.sh/

go test .\cmd\ --coverprofile=cover.cov
go build  -ldflags "-X cmd.version=$VERSION" -o build/spinner$VERSION.exe


Compress-Archive -Path "build/spinner$VERSION.exe" -Destination "build/spinner$VERSION.zip"