# GutenJapAlice

### Build on Ubuntu
クリーンなUbuntuでビルドするために必要なコマンド。  
不要な手順はスキップしてください。  

    sudo apt-get install golang
    sudo apt install git
    sudo apt-get install sqlite3
    export GOPATH=$HOME/go
    export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
    go get github.com/mattn/go-sqlite3
    go get github.com/mrjones/oauth
    go get golang.org/x/text/encoding
    go get gopkg.in/kyokomi/emoji.v1
    git clone https://github.com/payaneco/GutenJapAlice.git
    cd GutenJapAlice/
    go build main.go tweet.go 
