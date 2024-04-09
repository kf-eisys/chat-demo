# chatdemo配下のファイル読み込み
this_dir = File.expand_path(File.dirname(__FILE__))
proto_dir = File.join(this_dir, 'chatdemo')
$LOAD_PATH.unshift(proto_dir)

require 'chatdemo_pb'
require 'chatdemo_services_pb'
require 'grpc'

class ChatClient
  def initialize
    @stub = Chatdemo::ChatDemoService::Stub.new('localhost:50051', :this_channel_is_insecure)
  end

  def start_chat
    puts 'しりとりを始めます。終える時はexitかquitを入力してください。'
    print 'You: '

    requests = Enumerator.new do |y|
      loop do
        message = gets.chomp
        next if message.empty?
        break if message == 'exit' || message == 'quit'

        y << Chatdemo::WordChain.new(word: message)
      end
    end

    @stub.word_chain_chat(requests).each do |msg|
      case msg.result
      when :RESULT_WIN
        puts "\nServer: #{msg.message} (You win)"
        break
      when :RESULT_LOSE
        puts "\nServer: #{msg.message} (You lose)"
        break
      else
        puts "\nServer: #{msg.word}"
      end

      print "You: " # タイミングずれるのでここでprint
    end

    puts 'しりとりを終了します。'
  end
end

client = ChatClient.new
client.start_chat
