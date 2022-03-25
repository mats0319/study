# server-sent events

## 介绍

一种服务端主动向客户端推送消息的技术，使用http协议  
服务端在响应http请求时，表示返回值类型是流式的(```text/event-stream```)，参考视频传输

本demo包含：

1. 基础sse功能
2. https
3. go读写cookie
5. 服务端设置cookie
6. 通知多人

## 使用手册

### 修改参数

修改```./main.go```文件中的参数  
若修改了监听地址，请同步修改```./html/index.html```中的源地址

### 生成https证书

> 若上一步中修改了监听地址，请同步修改下方第二条命令，```/CN=127.0.0.1```中的```127.0.0.1```，改成上一步设置的监听地址即可

```powershell 
openssl genrsa -out server.key 2048
openssl req -nodes -new -key server.key -subj "/CN=127.0.0.1" -out server.csr
openssl x509 -req -sha256 -days 365 -in server.csr -signkey server.key -out server.crt
```

### 启动

> make sure you're at 'study/points/server_sent_events/' dir

```powershell 
go mod tidy
go run .
```

程序会自动启动一个网页，在windows系统上

### FAQ

- https证书不可信——信任证书并刷新页面
- 网页自启动失败——请检查当前位置，我们默认启动当前路径下的```/html/index.html```文件
