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

### How to run
動かすためには、以下の設定が必要です。

#### データ取得
不思議の国のアリスのデータを取得して`sqlite`のDBファイルである`alice.db`にデータを突っ込むために必要な設定

+ sqliteを使えるようにする
  - Windowsなら、必要なdll等をGutenJapAliceと同じフォルダに配置する
  - Linuxならapt-getなり何なりで適当に取得する
+ 必要なGoパッケージを取得する
  - mattn/go-sqlite3
  - oauth
  - text/encoding
  - kyokomi/emoji.v1
+ `go build main.go twitter.go`でビルドできる
+ ネットにつながれば`go run main.go twitter.go`で実行できる
### ツイート

データ取得のための設定のほかに、ツイートに必要な追加設定

+ `oauth.json`を作成する
  - `oauth_sample.json`をリネームとかして内容を書き換える
+ ネットにつながれば`go run main.go twitter.go -b bookmark.json`で実行できる
  - 1回実行するのに最低でも3分はかかる(15分くらいはザラにある)ので注意すること
