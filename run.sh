GOPATH="$(pwd)"
export GOPATH=$GOPATH

# Clean binary output
rm -rf $GOPATH/bin/

if [[ ! -d "$GOPATH/bin/" ]]; then
    mkdir -p "$GOPATH/bin/"
fi

go build -o $GOPATH/bin/freakingmath freakingmath
cd ./bin/
./freakingmath
