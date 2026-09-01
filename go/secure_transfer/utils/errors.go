package utils

var (
	// common: 9xxxx
	// - 900xx: key
	// - 901xx: encrypt/decrypt
	ErrInvalidPublicKey   = newError(90001, "Invalid Public Key")
	ErrGeneratePrivateKey = newError(90002, "Generate Private Key Failed")
	ErrECDH               = newError(90101, "ECDH Failed")
	ErrKDF                = newError(90102, "KDF Failed")
	ErrAESNewCipher       = newError(90103, "AES New Cipher Failed")
	ErrAESNewGCM          = newError(90104, "AES New GCM Failed")
	ErrOpenEncFile        = newError(90105, "Open Encrypt File Failed")

	// support functions: 1xxxx
	// - 100xx: utils
	// - 101xx: generate key pair
	// - 102xx: initialize message file
	ErrReadDir           = newError(10001, "Read Current Dir Failed")
	ErrGetFileInfo       = newError(10002, "Get FileInfo Failed")
	ErrNoMatchedFile     = newError(10003, "No Matched File")
	ErrMarshalPrivateKey = newError(10101, "Marshal Private Key Failed")
	ErrSavePrivateKey    = newError(10102, "Save Private Key Failed")
	ErrFileExist         = newError(10103, "File Already Exist")
	ErrMarshalPublicKey  = newError(10104, "Marshal Public Key Failed")
	ErrSavePublicKey     = newError(10105, "Save Public Key Failed")
	ErrInitMessageFile   = newError(10201, "Initialize Message File Failed")

	// encrypt: 2xxxx
	// - 201xx: deserialize public key
	// - 202xx: encrypt once
	ErrReadPublicKey   = newError(20101, "Read Public Key File Failed")
	ErrDecodePublicKey = newError(20102, "Decode Public Key Failed")
	ErrParsePublicKey  = newError(20103, "Parse Public Key Failed")
	ErrReadOriginFile  = newError(20201, "Read Origin File Failed")
	ErrWriteEncFile    = newError(20202, "Write Encrypt File Failed")
	ErrOpenOriginFile  = newError(20203, "Open Origin File Failed")

	// decrypt: 3xxxx
	// - 301xx: deserialize private key
	// - 302xx: decrypt
	ErrReadPrivateKey   = newError(30101, "Read Private Key Failed")
	ErrDecodePrivateKey = newError(30102, "Decode Private Key Failed")
	ErrParsePrivateKey  = newError(30103, "Parse Private Key Failed")
	ErrEncryptedFile    = newError(30201, "Invalid Encrypted File")
	ErrAESDecrypt       = newError(30202, "AES Decrypt Failed")
	ErrWriteDecFile     = newError(30203, "Write Decrypt File Failed")
	ErrOpenDecFile      = newError(30204, "Open Decrypt File Failed")

	// lib: 4xxxx
	// - 400xx: file header
	// - 401xx: stream frame
	ErrFileHeader = newError(40001, "Invalid File Header")
	ErrReadExact  = newError(40002, "Read Exact Failed")
	ErrFrame      = newError(40101, "Invalid Frame")
)

// 函数返回函数而不是实例，可以避免多处使用同一变量会继承历史数据的问题，详见测试代码
func newError(code int, detail string) func() *Error {
	return func() *Error { return NewError(code, detail) }
}
