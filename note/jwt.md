# JWT学习笔记

JWT：json web token[RFC 7519](https://datatracker.ietf.org/doc/html/rfc7519)， 是一种客户端与服务端之间确认身份的方案，此方案中服务端不保存状态信息

jwt分为两种：JWS (JW sign，RFC 7515)和JWE (JW encrypt，RFC 7516)，它们在格式上相似，JWE更复杂一些，所以我们用JWS举例

服务端在通过验证后（例如验证密码），构造并向客户端发送一个token， 客户端每次向服务端发送请求时携带该token，服务端验证token（验签或解密）并认为token携带的数据是可信的。

## JWT的结构

e.g.（换行符仅用作展示） eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ. SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c

上方JWT前两段的明文为：
`{"alg":"HS256","typ":"JWT"}.{"sub":"1234567890","name":"John Doe","iat":1516239022}`

- Header：base64url编码，通常包含算法（S部分的算法）、token类型等信息
- Payload：base64url编码，包含想要传输的实际数据
- Signature：使用H中的算法得到的签名，服务端信任payload的基础

## JWT的使用

- 服务端构造P，根据H、P计算S，整体拼好传给客户端
- 客户端每次请求时带上该token（放在HTTP请求的头信息Authorization字段）
- 服务端拿到token并验证，通过验证后认为P部分的内容是可信的

## base64url编码

与base64编码相似，只是将其中几个在url中有特殊含义的字符替换成其他字符：

- '+' -> '-'
- '/' -> '_'
- 不使用'='在末尾填充
