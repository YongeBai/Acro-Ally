build-windows:
	export CGO_ENABLED=1
	export GOOS=windows
	export GOARCH=amd64
	export CC=x86_64-w64-mingw32-gcc
	fyne package -os windows -icon Icon.png

run:
	go run .