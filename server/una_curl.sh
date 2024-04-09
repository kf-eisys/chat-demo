grpcurl -plaintext -d '
{
  "name": "fuga"
}' \
localhost:50051 chatdemo.ChatDemoService/SayHello
