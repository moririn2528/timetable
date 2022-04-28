# 結局用いる
tsc でコンパイルすると、サーバーが go 言語だとエラー吐かれる。これは、tsc のコンパイルによって生まれる javascript は意図的にそのまま動かせないようにしているかららしい。webpack を用いることで解決。

# webpack 仕様時の実行方法

`yarn start` で実行、localhost:8080 で確認可能

`yarn run build` で [name].js が出力される。

# typescript to javascript
`tsc` 

`tsc -w` でセーブ毎にコンパイル

# ポートを閉じる
`netstat -ano | find ":8080"` からプロセス ID を検索、タスクマネージャで削除


# 参考
環境構築
https://qiita.com/sansaisoba/items/921438a19cbf5a31ec53


