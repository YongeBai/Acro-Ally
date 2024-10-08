build-windows:
	export CC=x86_64-w64-mingw32-gcc
	export CXX=x86_64-w64-mingw32-g++
	export CGO_ENABLED=1
	export GOOS=windows
	export GOARCH=amd64
	fyne-cross windows -arch=amd64


build-linux:
	fyne install .
run:
	go run .