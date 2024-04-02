require 'sinatra'
require 'sinatra/reloader'

messages = []

get '/' do
  erb :index, locals: { messages: messages }
end

post '/messages' do
  msg = params[:message]
  messages << msg unless msg.empty?

  puts "Message: #{msg}"

  redirect '/'
end
