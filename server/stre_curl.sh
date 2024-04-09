grpcurl -plaintext -d '
{
  "word": "がいあ"
}
{
  "word": "い"
}
{
  "word": "みらい"
}
' \
localhost:50051 chatdemo.ChatDemoService/WordChainChat
