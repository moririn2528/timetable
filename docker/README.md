# docker
- 起動:  `docker comopose up -d` -d でバックグラウンド実行
- 終了: `docker compose down`

# mysql のチェック
`docker compose exec db /bin/bash` db はサービス名

## 削除
docker system prune --volumes

