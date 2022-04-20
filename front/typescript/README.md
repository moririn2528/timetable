# 使用しない
```
class.js:11 Uncaught ReferenceError: exports is not defined
```
というエラーが出て、この対処法は tsconfig.json の "module": "commonjs" を消すというものばかり。これを消すと "typescript-json-decorder" が使えない。
これが使えないのであれば、型はアノテーションレベルに下がる。

また、js にコンパイルしているので、実行時のエラーでどこでエラーが出ているのかわからない。これらの不便さは型がつくメリットを超えるため、使用しないことにした。

# webpack 仕様時の実行方法

`npm start` で実行、localhost:8080 で確認可能

`npm run build` で [name].js が出力される。

# typescript to javascript
`tsc` 

`tsc -w` でセーブ毎にコンパイル

# ポートを閉じる
`netstat -ano | find ":8080"` からプロセス ID を検索、タスクマネージャで削除