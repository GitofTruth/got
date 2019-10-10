package GoT

import (
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
)

type LogsExplorer struct {
	Repo *git.Repository
}


func CreateNewLogsExplorer(path string) LogsExplorer, err{

}


func (le *LogsExplorer) OpenRepo(path string) (bool, error) {
	fs := memfs.New()
	repo, err := git.PlainOpenWithOptions(path, fs)
	if err != nil {
		panic(err)
	} else {
		le.Repo = repo
		return true, nil
	}
}
