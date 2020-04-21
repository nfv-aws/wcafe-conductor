# wcafe-conductor
# これなに
wcafe-apiで実行したAPIを契機にキューを介して処理をするコンダクター
現時点ではpetsのPOST APIを実行時にDBのstatusがなるものを、コンダクターを実行するとCREATEDへ変更する。

# リポジトリクローン
```
cd $GOPATH/src/github.com
mkdir nfv-aws
cd nfv-aws
git clone git@github.com:nfv-aws/wcafe-conductor.git
```

# 使い方
パッケージインストール
```
"github.com/aws/aws-sdk-go/aws"
"github.com/aws/aws-sdk-go/aws/session"
"github.com/aws/aws-sdk-go/service/sqs"
```
# 環境設定
DB設定とAWSのSQS操作用の設定を追加

bashrcとかに以下を追記

```
vi ~/.bashrc

export WCAFE_DB_PASSWORD=password
export WCAFE_DB_ENDPOINT=endpoint
export WCAFE_SQS_REGION=region
export WCAFE_SQS_QUEUE_URL=queue_url

source ~/.bashrc
```

```
vi config/config.toml
```

# 動作確認
```
go run main.go
```