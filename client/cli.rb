# chatdemo配下のファイル読み込み
this_dir = File.expand_path(File.dirname(__FILE__))
proto_dir = File.join(this_dir, 'chatdemo')
$LOAD_PATH.unshift(proto_dir)

require 'chatdemo_pb'
require 'chatdemo_services_pb'
require 'grpc'

class ChatClient
  def initialize
    @stub = Chatdemo::ChatService::Stub.new('localhost:50051', :this_channel_is_insecure)
  end

  def start_chat
    # 送受信スレッド
    send_thread = Thread.new do
      loop do
        print 'You: '
        message = gets.chomp
        next if message.empty?
        break if message == 'exit' || message == 'quit'

        req = Chatdemo::SendMessageRequest.new(message: message)
        res = @stub.send_message([req])
        res.each do |msg|
          puts "Server: #{msg.message}"
        end
      end
    end

    send_thread.join
  end
end

client = ChatClient.new
client.start_chat
