chain.sh install
chain.sh instantiate

set-chain-env.sh -i '{"Args":["addNewRepo","{\"repoName\":\"GoT\",\"author\":\"hassan\",\"timestamp\":0,\"hashes\":{\"*************\":{},\"0000000000000000000000000000000000000000\":{}},\"branches\":{\"master\":{\"branchName\":\"master\",\"author\":\"masterCreator\",\"timestamp\":1,\"logs\":{\"*************\":{\"message\":\"message\",\"author\":\"mickey\",\"committer\":\"mickeyAsCommiter\",\"CommitterTimestamp\":3,\"hash\":\"*************\",\"parenthashes\":null,\"signature\":null}}}}}"]}'
set-chain-env.sh -q '{"Args":["queryBranch", "hassan", "GoT", "master"]}'
