rm -rfv ./wasm
#GOOS=js GOARCH=wasm go build -o ./wasm/runecalc.wasm -ldflags="-s -w" -gcflags="-m" -asmflags="-trimpath=$GOPATH"
fyne package -tags no_animations -tags hints -os web