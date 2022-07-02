# 『作ってわかる！ はじめてのgRPCを試す』を作ってわかってみる

[Zenn 作ってわかる！ はじめてのgRPC](https://zenn.dev/hsaki/books/golang-grpc-starting/viewer/codegenerate)

# 開発環境構築

```
docker-compose up
```

# 変更点

## goのパッケージインストール

4章の依存パッケージのインストールを go get から go install に変更

```
go install google.golang.org/grpc@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```
