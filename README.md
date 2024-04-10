双方向ストリーミングを利用したしりとりチャットアプリです
※入力は全てひらがなでお願いします

## ざっくり構成図
```
root/
  ├ chatdemo/
  │  └ chatdemo.proto
  ├ client/
  │  └ app.rb
  └ server/
     └ main.go
```

## 起動方法
- server
```
$ cd server && ./pegasus
```

- client
```
$ cd client && bundle && ruby app.rb
```

## protobuf生成
- server
```
$ buf generate --template buf.gen.server.yaml
```

- client
```
$ buf generate --template buf.gen.client.yaml
```
