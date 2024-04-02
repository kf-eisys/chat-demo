grpcurl -plaintext -d '
{
  "message": "hoge"
}
{
  "message": "fuga"
}
{
  "message": "piyo"
}
{
  "message": "さようなら"
}
' \
localhost:50051 chatdemo.ChatService/SendMessage
