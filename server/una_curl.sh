grpcurl -plaintext -d '
{
  "name": "hoge"
}' \
localhost:50051 chatdemo.ChatService/SayHello
