golang monorepo example
===

# 說明

一個golang的monolithic repository結構範例。

現用於天鏡科技公司內部產品開發，使大部分go開發的程式能夠共用同一套codebase。

package命名遵守官方指引：名稱是一個有意義的單字，而不是common、util等空泛的命名。

建議package組織方式如下：

service `app/quote-recv/{client,server,model}.go`

api `app/rest-api/quote/controller.go` (only rest-api imports quote)

需要抽出共用模組時

共用模組 `appmodule/candlestick/{client,helper}.go` `lib/calendar/calculator.go`

多模組共用資料庫model時

`models/trader/setting.go`

# 使用方式

git clone 後，確認開啟環境變數 `GO111MODULE=on`，進行一次 `go mod download` 下載所有依賴。

使用 `./build_docker.sh` 可以將程式碼編譯成執行檔或 docker images。

預設情況下 docker image tag 會從 git commit hash 自動抓取，例如： `develop-1347acf1`，如果想要自訂，在編譯前先執行 `export CUSTOM_TAG=latest`。

`./build_docker.sh` 提供了一些選項
- `--bin` 只產生執行檔，不產生 docker image
- `--dep` 依據原始碼的修改，重新編譯相關的執行檔
- `--push` 產生出 docker image 後順帶推送到 docker registry

# Config

所有config藉由`cliflag.Register`在啟動時進行自動註冊。

singleton參考`database`的func init，

其他參考`internal/rpc/hello/cli-flags.go`和`cmd/hellocli/main.go`

# Folder Structure

- app: app implementation
- appmodule: shared module across the apps
- build: artifacts
- cliflag: config helper
- cmd
- cache, database, logging, env: singleton package
- deploy: helm charts
- doc
- docker: dockerfiles
- health: to support prometheus and k8s alive/ready probe
- internal/rpc: generated code
- models: gorm model definitions
- pkg: forked and modified packages from Internet
- protos: protobuf definition files
