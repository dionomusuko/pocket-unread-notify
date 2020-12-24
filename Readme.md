# unread pocket slack notify

まだ読んでいないpocketの記事をslackに通知します。

## 環境変数
| キー | 値 |
|:----|:----|
|SLACK_API_KEY|slackのwebhookURL|
|POCKET_CONSUMER_KEY|pocketのコンシューマーキー|
|POCKET_ACCESS_TOKEN|pocketのアクセストークン|

## 使ったサービス等
- Cloud functions
    - entrypoint Function
    - runtime go113
    - trigger-topic
- Cloud pub/sub

- Cloud scheduler

