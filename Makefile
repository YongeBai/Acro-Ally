build-windows:
	export CC=x86_64-w64-mingw32-gcc
	export CGO_ENABLED=1
	export GOOS=windows
	export GOARCH=amd64
	fyne package -os windows -icon Icon.png

build-linux:
	fyne install .
run:
	go run .