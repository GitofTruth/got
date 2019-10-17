package datastructures

type RepoBranch struct {
	Author    string               `json:"author"`
	Timestamp int                  `json:"timestamp"`
	Logs      map[string]CommitLog `json:"logs"`
}

func CreateRepoBranch(author string, timestamp int, logs map[string]CommitLog) (RepoBranch, error) {
	var repoBranch RepoBranch
	repoBranch.Author = author
	repoBranch.Timestamp = timestamp

	if logs == nil {
		repoBranch.Logs = make(map[string]CommitLog)
	} else {
		repoBranch.Logs = logs
	}

	return repoBranch, nil
}
