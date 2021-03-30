# salesforce-api-kube

## 概要
salesforce-api-kubeは、SalesforceとのAPI連携を担当するサービスです。
Status-Kanbanからメッセージを受信すると、Status Kanban, redis-clusterに書き込みをして、Salesforceに読み込み・書き込みを行います。


## 動作環境
salesforce-api-kubeはAIONのプラットフォーム上での動作を前提としています。
使用する際は、事前にAIONの動作環境を用意してください。

OS: Linux
CPU: Intel64/AMD64/ARM64
Kubernetes
AION

## sand box での実行
下記の設定を行う事で接続先を salesforce sand box に切り替える事ができます。
1. `config.test.json` に クレデンシャル情報を記載する。
2. 環境変数 DEV=true で起動する。


## セットアップ
このリポジトリをクローンし、makeコマンドを用いてDocker container imageのビルドを行ってください。
```
$ cd salesforce-api-kube
$ make docker-build
```

## kanban との通信
### kanban から受信するデータ
kanban から受信する metadata に下記の情報を含む必要があります。

| key | type | description |
| --- | --- | --- |
| method | string | 使用する HTTP メソッド |
| object | string | 操作対象の Salesforce オブジェクト |
| path_param | string | パスパラメータ (必要な場合のみ)|
| query_params | map[string]strnig | クエリパラメータ (必要な場合のみ)|

path_param, query_params は必要な場合のみ含まれます。  

具体例: 
```example
# metadata (map[string]interface{}) の中身

"method": "get"
"object": "Account"
"path_param": "15"
```

### kanban に送信するデータ
kanban に送信する metadata に下記の情報を含める必要があります。

| key | value |
| --- | --- |
| key | 送信するデータの Object 名 |
| content | 送信するデータの中身 |

具体例:
```example
# metadata (map[string]interface{}) の中身

"key": "Account"
"content": `[{
    "attributes": {
        "type": "Account",
        "url": "/services/data/v51.0/sobjects/Account/0010I000028S9KmQAK"
    },
    "Id": "xxxxxxx",
    "Name": "サンプル太郎",
    "LastName": "サンプル",
    "FirstName": "太郎",
    "CustomerNameKana__c": "サンプル　タロウ",
    "PersonEmail": "sample25a@gmail.com", 
    "KokyakuStatus__c": "契約中（継続）",
    "Phone_nohyphen__c": "041234567",
    "MailingAddress__c": "000-0000 神奈川県ほげほげ市",
    "todoufukenSikugun__c": "神奈川県ほげほげ市"
}]`
```
