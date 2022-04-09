# timetable
 時間割を作りたい


## 詰まったところ

### フロントサイドについて
フロントサイドは html+js, それから api を呼び出す構造。フロント側とサーバー側は分かれているという考えだったため、CORS の設定、http-server で別の port で起動させてようとしたが、Go 言語では

http.Handle("/", http.FileServer(http.Dir("static")))

で同じオリジンで起動でき、CORS もいらなくなった。
html などのフロント側のファイルは (main,go が入っているフォルダ)/static に入れる。"static" は変更可能、"../html" でも動いたが、"/" を "/top" などに変えると動かなかった。