# 转发程序

Listen: 监听的端口号
Type: ws是websocket的连接， tcp是tcp的连接
Target: 目标的连接地址，比如 127.0.0.1:3306

发送顺序
client -> netchannel -> target
接受顺序
target -> netchannel -> client

程序功能：
实现转发功能，程序和target使用tcp长连接。

ws协议的服务器可以通过此服务转换成tcp的方式。