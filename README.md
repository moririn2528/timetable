# timetable
 時間割を作りたい

# 公開用
サイト: http://timetable.torimari.site:8080/
説明: https://qiita.com/moririn2528/items/e10070d47275fd10f169
説明(雑、技術寄り): https://qiita.com/moririn2528/items/994d0185d2b55b3c1a97
非公開部分を匿名化しています。そのため、少し見づらくなっています。ご了承ください。

## 説明資料


## 補足
### クラス
クラスの包含関係をグラフ化したものです。どちらかというと開発、説明のためのものです。

### 時間割
"クラス:"の右のプラスボタンからクラス20,クラス47,クラス50を選択すると一般的な時間割の表示となります。

### 先生用時間割
時間割変更はここから可能です。授業できない時間帯をすべて選択、計算ボタンを押すと変更プレビューとなります。描画を押すと、授業がなくなる枠に赤色、追加される枠に青色、授業が変更される枠(クラスなり科目なりが変わる)が黄色で表示されます。

# 開発者向け
## 環境構築
```
git config --local core.hooksPath .githooks
```

install git secrets
https://github.com/awslabs/git-secrets#windows


## VM 構築
git clone して、
.ssh/config を適当にして
`bash start-scp.sh (Host)`

ssh 接続した後 startup.sh を実行