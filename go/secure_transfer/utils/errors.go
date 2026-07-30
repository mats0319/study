package utils

var (
	ErrForTest = newError(-1, "test error string")

	// support functions: 1xxxx
	// - 100xx: utils
	// - 101xx: generate key pair
	// - 102xx: initialize message file
	ErrReadDir           = newError(10001, "Read Dir Failed")
	ErrGetFileInfo       = newError(10002, "Get FileInfo Failed")
	ErrReadMessageFile   = newError(10003, "Read Message File Failed")
	ErrMarshalPrivateKey = newError(10102, "Marshal Private Key Failed")
	ErrSavePrivateKey    = newError(10103, "Save Private Key Failed")
	ErrMarshalPublicKey  = newError(10104, "Marshal Public Key Failed")
	ErrSavePublicKey     = newError(10105, "Save Public Key Failed")
	ErrInitMessageFile   = newError(10201, "Initialize Message File")

	// encrypt: 2xxxx
	// - 201xx: deserialize public key
	// - 202xx: encrypt
	ErrReadPublicKey   = newError(20101, "Read Public Key File Failed")
	ErrDecodePublicKey = newError(20102, "Decode Public Key Failed")
	ErrParsePublicKey  = newError(20103, "Parse Public Key Failed")
	ErrSaveCiphertext  = newError(20205, "Save Ciphertext Failed")

	// decrypt: 3xxxx
	// - 301xx: deserialize private key
	// - 302xx: decrypt
	ErrReadPrivateKey   = newError(30101, "Read Private Key Failed")
	ErrDecodePrivateKey = newError(30102, "Decode Private Key Failed")
	ErrParsePrivateKey  = newError(30103, "Parse Private Key Failed")
	ErrAESDecrypt       = newError(30104, "AES Decrypt Failed")
	ErrSavePlaintext    = newError(30105, "Save Plaintext Failed")

	// common: 4xxxx
	// - 400xx: key
	// - 401xx: encrypt/decrypt
	ErrInvalidPublicKey   = newError(40001, "Invalid Public Key")
	ErrInvalidPrivateKey  = newError(40002, "Invalid Private Key")
	ErrGeneratePrivateKey = newError(40003, "Generate Private Key Failed")
	ErrECDH               = newError(40101, "ECDH Failed")
	ErrKDF                = newError(40102, "KDF Failed")
	ErrAESNewCipher       = newError(40103, "AES New Cipher Failed")
	ErrAESNewGCM          = newError(40104, "AES New GCM Failed")
)

// 函数返回函数而不是实例，可以避免多处使用同一变量会继承历史数据的问题，详见测试代码
func newError(code int, detail string) func() *Error {
	return func() *Error { return NewError(code, detail) }
}
