git --bare init .got
git remote set-url origin .got
#cp $GOPATH/src/github.com/GitofTruth/GoT/hooks/repo/pre-push .git/hooks
cp $GOPATH/src/github.com/GitofTruth/GoT/hooks/bare/post-receive .got/hooks
cp $GOPATH/src/github.com/GitofTruth/GoT/hooks/bare/pre-receive .got/hooks
mkdir .gotconfig
git push
git pull
