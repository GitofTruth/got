git --bare init .got
git remote set-url origin .got
mkdir .gotconfig

#adding wallet info
cp -r $GOPATH/src/github.com/GitofTruth/GoT/profiles .gotconfig
cp $GOPATH/src/github.com/GitofTruth/GoT/client/client.js .gotconfig
cp $GOPATH/src/github.com/GitofTruth/GoT/client/enrollAdmin.js .gotconfig
cp $GOPATH/src/github.com/GitofTruth/GoT/client/wallet.js .gotconfig
cp $GOPATH/src/github.com/GitofTruth/GoT/client/package.json .gotconfig

cd .gotconfig

npm install
npm install js-yaml fabric-network fabric-ca-client

cd ..

sudo node .gotconfig/wallet.js add acme Admin
node .gotconfig/client.js

#cp $GOPATH/src/github.com/GitofTruth/GoT/hooks/repo/pre-push .git/hooks
# cp $GOPATH/src/github.com/GitofTruth/GoT/hooks/bare/post-receive .got/hooks
# cp $GOPATH/src/github.com/GitofTruth/GoT/hooks/bare/pre-receive .got/hooks
# mkdir .gotconfig
# git push
# git pull
