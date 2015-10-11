# [WIP] Batch Request Server

複数のリクエストを1つにまとめて返すサーバーです。

APIの仕様は、基本的に[FacebookのBatch Requests](https://developers.facebook.com/docs/graph-api/making-multiple-requests)に従っています。

## 使い方

### Requestの例
```
curl \
-d end_point=http://localhost:3000 \
-d include_headers=true \
-d 'batch=[
  {
    "method":"POST",
    "name":"create-ad",
    "relative_url":"11077200629332/ads",
    "body":"ads=%5B%7B%22name%22%3A%22test_ad%22%2C%22billing_entity_id%22%3A111200774273%7D%5D"
  },
  {
    "method":"GET",
    "relative_url":"?ids={result=create-ad:$.data.*.id}"
  }
]' \
http://localhost:8080
```

| パラメータ名 | 必須 | 説明 |
|:-:|:-:|:-:|
| end_point | o | バッチリクエストするEndpointのホストを指定します |
| batch | o | 複数のリクエストの中身をJSON形式の配列で指定します |
| include_headers |  | 結果に`headers`を含めるかを指定します (デフォルト: true) |

#### batchパラメータの形式

| 項目名 | 必須 | 説明 |
|:-:|:-:|:-:|
| method | o | HTTPリクエストメソッドを指定します (GET/POST/PUT/PATCH/DELETE) |
| relative_url | o | リクエストの`end_point`からのpathを指定します |
| body |  | POST/PUT/PATCHリクエストの中身を指定します |
| name |  | 別のリクエストのパラメータとして使えるように結果に名前をつけます |


### Responseの例

```
[
    { "code": 200,
      "headers": [
          { "name":"Content-Type",
            "value":"text/javascript; charset=UTF-8"}
       ],
      "body":"{\"id\":\"…\"}"
    },
    { "code": 200,
      "headers": [
          { "name":"Content-Type",
            "value":"text/javascript; charset=UTF-8"
          },
          { "name":"ETag",
            "value": "…"
          }
      ],
      "body": "{\"data\": [{…}]}
    }
]
```

以下の項目がリクエストした順にJSON形式の配列になって返ります。

| 項目名 |  説明 |
|:-:|:-:|
| code | HTTPのステータスコードを返します |
| headers | レスポンスのHEADERをハッシュ形式で返します (`includes_headers=false`だとこの項目は省略されます) |
| body | レスポンスのBODYを返します |

※ なお、各リクエストがタイムアウトした場合`null`が返ります。
