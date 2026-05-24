package models

type VaultEntry struct {
	GUID string
	Blob []byte
}

func CreateVaultEntry(guid string, blob []byte) *VaultEntry {
	var VaultEntryObject VaultEntry = VaultEntry{
		GUID: guid,
		Blob: blob,
	}
	return &VaultEntryObject
}
func CreateVaultEntryFromRequest(guid string, req CreateVaultEntryRequest) *VaultEntry {
	var VaultEntryObject VaultEntry = VaultEntry{
		GUID: guid,
		Blob: req.Blob,
	}
	return &VaultEntryObject
}

type CreateVaultEntryRequest struct {
	Blob []byte
}
type CreateVaultEntryHeaders struct {
}
