package datastructures

// This structure is modeling the necessary data to be stored for each commit
// including data nomrally stored through git and data required to get the git
// objects through the GoT IPFS Cluster.
// The EncryptionKey stored inside this object represents the hash of the symmetric
// encryption key used to encrypt the git objects. However, when a query is done
// in the block chain, this field is filled with the encrypted form of the EncryptionKey
// instead of its hash.
type CommitLog struct {
	Message            string            `json:"message"`
	Author             string            `json:"author"`
	Committer          string            `json:"committer"`
	CommitterTimestamp int               `json:"CommitterTimestamp"`
	Hash               string            `json:"hash"`
	Parenthashes       []string          `json:"parentHashes"`
	Signature          []byte            `json:"signature"`
	EncryptionKey      string            `json:"encryptionKey"`
	StorageHashes      map[string]string `json:"storageHashes"`
}

// this is a helper function to initialize a new CommitLog object instance
func CreateNewCommitLog(message string, author string, commiter string, timestamp int, hash string, parenthashes []string, signature []byte, encryptionKey string, storageHashes map[string]string) (CommitLog, error) {
	var log CommitLog

	log.Message = message
	log.Author = author
	log.Committer = commiter
	log.CommitterTimestamp = timestamp
	log.Hash = hash
	log.Parenthashes = parenthashes
	log.Signature = signature
	log.EncryptionKey = encryptionKey
	log.StorageHashes = storageHashes

	return log, nil
}
