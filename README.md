# grpc-practice

### 環境構築

```
$ brew install protobuf
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

```

https://grpc.io/docs/languages/go/quickstart/

### 認証処理

認証にまつわるインターセプターの処理がまとめられているので DL する

```
go get github.com/grpc-ecosystem/go-grpc-middleware
```

### gRPC のエラーハンドリング

|       code        | http |              内容              |
| :---------------: | :--: | :----------------------------: |
| DEADLINE_EXCEEDED | 504  |          タイムアウト          |
|   UNIMPLEMENTED   | 501  |      処理が実装されてない      |
|    UNAVAILABLE    | 503  | サービスが一時的に利用できない |
|      UNKNOWN      | 500  |          不明なエラー          |
|  UNAUTHENTICATED  | 401  |         認証情報がない         |
| PERMISSION_DENIED | 403  |         実行権限がない         |

https://grpc.io/docs/guides/error/
