package log

import (
	"os"
	"strings"
)

//openFile ログを出力するファイルを設定する。
//ファイルが存在する場合、ファイルにログを追記。
//ファイルが存在しない場合、ファイルを作成し、ログを出力。
func openFile(fileName string) (*os.File, error) {
	if exists(fileName) {
		f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0666)
		return f, err
	}
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	return f, err
}

//formatFilePath ログに記載するファイル名の抽出
func formatFilePath(path string) string {
	arr := strings.Split(path, "/")
	return arr[len(arr)-1]

}

//exists　ファイルが存在するか確認する。
func exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
