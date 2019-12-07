package datastructures

type CommitLog struct {
	Message            string   `json:"message"`
	Author             string   `json:"author"`
	Committer          string   `json:"committer"`
	CommitterTimestamp int      `json:"CommitterTimestamp"`
	Hash               string   `json:"hash"`
	Parenthashes       []string `json:"parenthashes"`
	Signature          []byte   `json:"signature"`
}

func CreateNewCommitLog(message string, author string, commiter string, timestamp int, hash string, parenthashes []string, signature []byte) (CommitLog, error) {
	var log CommitLog
	log.Message = message
	log.Author = author
	log.Committer = commiter
	log.CommitterTimestamp = timestamp
	log.Hash = hash
	log.Parenthashes = parenthashes
	log.Signature = signature

	return log, nil
}
