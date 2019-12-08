git --bare init .got
git remote set-url origin .got
cp $GOPATH/src/github.com/GitofTruth/GoT/hooks/pre-push .git/hooks
cp $GOPATH/src/github.com/GitofTruth/GoT/hooks/pre-push .got/hooks
