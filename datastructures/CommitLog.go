package datastructures

type CommitLog struct {
	// Status       []string `json:"status"`
	// Files        []string `json:"files"`

	Message      string   `json:"message"`
	Author       string   `json:"author"`
	Committer    string   `json:"committer"`
	Timestamp    int      `json:"timestamp"`
	Hash         string   `json:"hash"`
	Parenthashes []string `json:"parenthash"`
	//Signature    []string `json:"signature"`
}

func CreateNewCommitLog(message string, author string, commiter string, timestamp int, hash string, parenthashes []string) (CommitLog, error) {
	var log CommitLog
	log.Message = message
	log.Author = author
	log.Committer = commiter
	log.Timestamp = timestamp
	log.Hash = hash
	log.Parenthashes = parenthashes
	//log.Signature = signature

	return log, nil
}
