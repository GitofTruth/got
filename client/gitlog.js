const gitlog = require('gitlog');

const repoPath = '/home/hkandil/Desktop/Git of Truth/HLF-Dev-Chaincode-V1.4-1.3/gocc/src/github.com/GitofTruth/GoT';

const options =
    { repo: repoPath
    , number: 1
    , fields:
      [ 'subject'
      , 'authorName'
      , 'committerName'
      , 'hash'
      , 'parentHashes'
      ]
    , execOptions:
      { maxBuffer: 1000 * 1024,
        status: false
      }
    };
// Asynchronous (with Callback)
gitlog(options, function(error, commits) {
  // Commits is an array of commits in the repo
  //console.log(commits)
});

// Synchronous
let commits = gitlog(options);
console.log(commits[0]);

let commit_str = JSON.stringify(commits[0]);

console.log(commit_str);