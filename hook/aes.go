package hook

import (
	"TFService/model"

	"github.com/zhibingzhou/go_public/common"
)

var HookAES *common.AES

func getPrivateKey(mer_code string) string {
	//第一步骤，读取密钥
	a_list := model.MerInfoRedis(mer_code)
	//设置密钥
	return a_list["private_key"]
}

func HookAesDecrypt(mer_code, decode_str string) string {
	private_key := getPrivateKey(mer_code)
	mer_aes := common.SetAES(private_key, "", "", 16)
	aes_res := mer_aes.AesDecryptString(decode_str)
	return aes_res
}

func HookAesEncrypt(mer_code, encode_str string) string {
	private_key := getPrivateKey(mer_code)
	mer_aes := common.SetAES(private_key, "", "", 16)
	aes_res := mer_aes.AesEncryptString(encode_str)
	return aes_res
}

func AesEncrypt(private_key, encode_str string) string {
	mer_aes := common.SetAES(private_key, "", "", 16)
	aes_res := mer_aes.AesEncryptString(encode_str)
	return aes_res
}
