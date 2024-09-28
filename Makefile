build-windows:
	sudo usermod -aG docker $USER
	fyne-cross windows -output Acro-Ally.exe -icon Icon.png

build-linux:
	fyne install .
run:
	go run .