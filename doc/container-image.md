# container image

## binary

預設go編譯會產生出靜態連結檔案 (CGO_ENABLED=0)

所以我們在本地編譯後再複製進container image內提供使用。

## image

預設編譯使用的dockerfile從`docker/dockerfile.tmpl`模板轉換成實際的dockerfile

如果`docker`資料夾內有與binary檔案同名的`XXX.dockerfile`則會使用該特化版本

image內需要特別的靜態資源檔則是放在`static/`資料夾中

在go1.16以後也可以使用embed的方式。

## deploy infra code

放在`deploy/`資料夾中，存放`docker-compose.yaml` `helm chart` `kubernetes yaml`。
