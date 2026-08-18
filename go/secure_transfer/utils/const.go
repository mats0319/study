package utils

const (
	FileHeaderStart = "ST"
	DeriveKeyInfo   = "secure transfer derive key info"

	PrivateKeyFileName = "priv.key"
	PublicKeyFileName  = "PUB.KEY"

	OriginFileName    = "message"
	DefaultExtension  = "txt"
	EncryptedFileName = "CIPHER"
	DecryptedFileName = "message_decrypted"
)

const (
	DeriveKeySaltLength = 16
	AESBaseNonceLength  = 7 // aes-gcm nonce: 12 = 7 + 1 + 4
	PublicKeyLength     = 32
	FileHeaderAADLength = 32 // sha256 length

	DeriveKeyLength = 32 // 32 Bytes, match AES-256

	OnceEncMaxSize = 10 * 1024 * 1024 // 10 MB
	FrameSize      = 1024 * 1024      // 1 MB
)

type EncryptMethod = byte

const (
	EncryptMethod_Once   EncryptMethod = 1 // 一次性读取全部文件到内存中加密
	EncryptMethod_Stream               = 2 // 流式加密
)

func IsEncryptMethod(method byte) bool {
	return method == EncryptMethod_Once || method == EncryptMethod_Stream
}
