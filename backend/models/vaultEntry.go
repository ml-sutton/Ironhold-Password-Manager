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

type CreateVaultEntryRequest struct {
}
type CreateVaultEntryHeaders struct {
}
