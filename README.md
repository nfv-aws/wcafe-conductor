# wcafe-conductor
## これなに
wcafe-apiで実行したAPIを契機にキューを介して処理をするコンダクター
現時点では以下の処理が走る。
- petsのPOST APIを実行時にDBのstatusが、コンダクターを実行するとCREATEDへ変更する。
- storesのPOST APIを実行時にDBのstrong_pointが、コンダクターを実行すると"sqs_test"に変更する。
- usersのPOST APIを実行時にDBのadressが、コンダクターを実行すると"Tokyo"に変更する。

## リポジトリクローン
```
cd $GOPATH/src/github.com
mkdir nfv-aws
cd nfv-aws
git clone git@github.com:nfv-aws/wcafe-conductor.git
```

## 使い方
### 事前準備  
#### パッケージインストール
```
"github.com/aws/aws-sdk-go/aws"
"github.com/aws/aws-sdk-go/aws/session"
"github.com/aws/aws-sdk-go/service/sqs"
"github.com/sirupsen/logrus"
```
#### 環境設定
DB設定とAWSのSQS操作用の設定を追加

bashrcとかに以下を追記

```
vi ~/.bashrc

export WCAFE_DB_PASSWORD=password
export WCAFE_DB_ENDPOINT=endpoint
export WCAFE_SQS_REGION=region
export WCAFE_SQS_PETS_QUEUE_URL=queue_url_1
export WCAFE_SQS_Stores_QUEUE_URL=queue_url_2
export WCAFE_SQS_Users_QUEUE_URL=queue_url_3

source ~/.bashrc
```

```
vi config/config.toml
```

## 動作確認
```
go run main.go
```
# ログの設定方法
ログは以下の3パターンを用意しており、DefaultではInfoモードとなっている。
- Debugモード
- Infoモード
- Errorモード

切り替え方は以下のように環境変数を設定して、プログラムを実行すればよい。
```
export LOG_LVE="Debug"
```

## UnitTest  
### Conductor

```
go test -vet=off -v ./conductor/... 

PASS
ok      github.com/nfv-aws/wcafe-conductor/conductor    0.014s
```

`-vet=off` オプションをつけないとテスト対象コードに含まれるフォーマット文の影響で以下のエラーが出るので注意

```
$ go test -v ./conductor/...
# github.com/nfv-aws/wcafe-conductor/conductor
conductor/supplies_conductor.go:50:3: Fatal call has possible formatting directive %v
FAIL    github.com/nfv-aws/wcafe-conductor/conductor [build failed]
FAIL
```
