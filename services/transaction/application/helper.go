package application

import "evolve/services"

func encryptPassword(encryptor services.EncryptorManager) {
	byteValueOfPassword, _ := encryptor.GenerateFromPasscode(universalPassword)
	universalEncryptedPassword = string(byteValueOfPassword)
}

func (t transactionAppHandler) validatePin(pin string) error {
	return t.encryptor.ComparePasscode(pin, universalEncryptedPassword)
}
