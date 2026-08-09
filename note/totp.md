# 基于时间的一次性密码

Time-based One-Time Password

是一种根据预共享的密钥与当前时间，计算带有时效性的密码的算法。  
常用于2FA（Two-Factor Authenticate，二重因素验证）等场景。

[OTP算法步骤详解](https://notes.mengxin.science/2017/05/30/hotp-totp-algorithm-analysis/)  
[RFC 4226（HOTP）](https://datatracker.ietf.org/doc/html/rfc4226)  
[RFC 6238（TOTP）](https://datatracker.ietf.org/doc/html/rfc6238)

## 过程

> 括号内为通常的默认参数

- 生成密钥K，服务端与客户端共享
- 协商起始时间T0（0）和时间步长TI（30秒）
- 协商密码生成算法（HMAC-SHA-1）
- 协商密码长度（6位）

## 生成密码

> 介绍如何计算一次性密码

OTP -> HOTP -> TOTP

几种算法的区别在于hash函数的原文，OTP的hash原文是随机的、HOTP使用计数器、TOTP使用时间步数

- OTP(K,C) = Truncate(Hash(K,C))
    - Truncate: 截取算法，将hash结果转变成一串数字
    - hash algorithm: hmac-sha-1
    - C: hash原文，根据算法不同，使用计数器(HOTP)、时间步数(TOTP)

截取算法：

1. hash得到20 Bytes，取最后一个Byte的后4 bits作为起始索引（起始索引属于0～15）
2. 从起始索引开始，往后获取总计4个Bytes，将其后31 bits转换为十进制数字作为随机长密码（舍弃最高位）
3. 根据需要的位数d，取长密码的后d位，长度不足的在高位补0

## 问题

- 因为网络延迟、时钟偏移、用户延误等原因，有些服务器也接受使用上一个时间步长或下一个时间步长生成的密码
- TOTP没有限制重试次数，所以服务端程序需要处理用户多次失败的问题
- 2的31次方转换成十进制数是21开头的10位数，即标准算法只能计算不多于10位的密码，且计算10位的密码时，最高位各数字出现概率不同

## 代码示例

### go

``` go   
func totp() {
	keyBase32 := []byte("HFLEOZBUOVKXMVRY")
	key := make([]byte, 32)
	n, err := base32.StdEncoding.Decode(key, keyBase32)
	if err != nil {
		fmt.Println("decode base32 failed, error: ", n, err)
		return
	}
	key = key[:n]

	timestampSecond := time.Now().Unix()

	// hmac-sha-1
	hasher := hmac.New(sha1.New, key)
	hasher.Write(itob(timestampSecond / 30))
	hmacHash := hasher.Sum(nil)

	// get an int32 from hash
	offset := int(hmacHash[len(hmacHash)-1] & 0x0f)
	// 算法要求屏蔽最高有效位
	longPassword := int(hmacHash[offset]&0x7f)<<24 |
		int(hmacHash[offset+1])<<16 |
		int(hmacHash[offset+2])<<8 |
		int(hmacHash[offset+3])

	// get 6 digits
	password := longPassword % int(math.Pow10(6))

	fmt.Println(fmt.Sprintf("%06d", password))
}

func itob(integer int64) []byte {
	byteArr := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		byteArr[i] = byte(integer & 0xff)
		integer = integer >> 8
	}
	return byteArr
}
```

### dart

```dart
const int _timeInterval = 30;
const int _pwdLength = 6;

(String, double) generateTOTP(String keyBase32) {
  try {
    final Uint8List keyBytes = base32.decode(keyBase32.toUpperCase());

    final int timestampSecond = DateTime.now().millisecondsSinceEpoch ~/ 1000;
    final int timeRemain = _timeInterval - timestampSecond % _timeInterval;
    final int timeCount = timestampSecond ~/ _timeInterval;
    final Uint8List timeCountBytes = _int2Bytes(timeCount);

    final List<int> hash = Hmac(sha1, keyBytes).convert(timeCountBytes).bytes;

    final int offset = hash.last & 0xf; // 0b 0000 1111
    final int longPassword =
        (hash[offset] & 0x7f) << 24 | // 0b 0111 1111
        hash[offset + 1] << 16 |
        hash[offset + 2] << 8 |
        hash[offset + 3];

    final int totp = longPassword % pow(10, _pwdLength).toInt();

    return (totp.toString().padLeft(_pwdLength, "0"), timeRemain.toDouble());
  } catch (e) {
    return ("计算密码失败", _timeInterval.toDouble());
  }
}

Uint8List _int2Bytes(int long) {
  final byteArray = Uint8List(8);

  for (var index = byteArray.length - 1; index >= 0; index--) {
    byteArray[index] = long & 0xff;
    long >>= 8;
  }

  return byteArray;
}
```
