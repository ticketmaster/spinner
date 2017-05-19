$VERSION = '0.2.1'

glide up # https://glide.sh/

go test --coverprofile=cover.cov
go build  -ldflags "-X main.version=$VERSION" -o build/spinner$VERSION.exe