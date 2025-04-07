package crypto

func VerifyMasterPassword(hash string, masterPassword string) (bool, error) {
	return hash == masterPassword, nil
}