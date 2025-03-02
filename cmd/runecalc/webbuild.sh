fyne package -os web
rm -fv ./wasm/runecalc.wasm
GOOS=js GOARCH=wasm go build -o ./wasm/runecalc.wasm -ldflags="-s -w" -gcflags="-m" -asmflags="-trimpath=$GOPATH"