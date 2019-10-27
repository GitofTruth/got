const gitlog = require('gitlog');

__dirname = '/home/hkandil/Desktop/Git of Truth/HLF-Dev-Chaincode-V1.4-1.3/gocc/src/github.com/GitofTruth';


const options =
    { repo: __dirname + '/GoT'
    , number: 2
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
  console.log(commits)
});

// // Synchronous
// let commits = gitlog(options);
// console.log(commits);