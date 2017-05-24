$VERSION = '1.0.0'

glide up # https://glide.sh/

go test --coverprofile=cover.cov
go build  -ldflags "-X github.com/Ticketmaster/spinner/cmd.version=$VERSION" -o build/spinner$VERSION.exe

Compress-Archive -Path "build/spinner$VERSION.exe" -Destination "build/spinner$VERSION.zip"