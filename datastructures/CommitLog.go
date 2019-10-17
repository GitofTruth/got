package datastructures

type CommitLog struct {
	Message      string   `json:"message"`
	Author       string   `json:"author"`
	Committer    string   `json:"committer"`
	Timestamp    int      `json:"timestamp"`
	Hash         []byte   `json:"hash"`
	Parenthashes [][]byte `json:"parenthash"`
	Signature    []byte   `json:"signature"`
}

func CreateCommitLog(message string, author string, commiter string, timestamp int, hash []byte, parenthashes [][]byte, signature []byte) (CommitLog, error) {
	var log CommitLog
	log.Message = message
	log.Author = author
	log.Committer = commiter
	log.Timestamp = timestamp
	log.Hash = hash
	log.Parenthashes = parenthashes
	log.Signature = signature

	return log, nil
}
