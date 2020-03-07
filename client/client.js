// internal imports
// var Constants = require("constants");
const CONNECTION_PROFILE_PATH = '../profiles/dev-connection.yaml';
const FILESYSTEM_WALLET_PATH = './user-wallet';


const USER_ID = 'Admin@acme.com';
const NETWORK_NAME = 'airlinechannel';
const CONTRACT_ID = "GoT";

const repoPath = '.';
var repo = {};
var branchesNames = [];
var branchesContent = [];

var options =
{ repo: repoPath
, branch: 'master'
, number: 10
, fields:
  [ 'subject'
  , 'authorName'
  , 'committerName'
  , 'authorDate'
  , 'hash'
  , 'parentHashes'
  ]
, execOptions:
  { maxBuffer: 1000 * 1024,
    status: false
  }
};

// imports
const path = require('path');
const fs = require('fs');
const yaml = require('js-yaml');

const { Gateway, FileSystemWallet, DefaultEventHandlerStrategies, Transaction  } = require('fabric-network');

const gitlog = require('gitlog');

const gitlog = require('gitlog');
const simpleGit = require('simple-git')(repoPath);

class client {
  constructor (){
    await this.setupGateway();
    let this.network = await gateway.getNetwork(NETWORK_NAME);
    const this.contract = await network.getContract(CONTRACT_ID);
    loadCurrentRepo();
  }

  loadCurrentRepo(){

    simpleGit.branchLocal(function(e,d){
      branchesNames = d['all']
      branchesContent = d['branches']

      var branchObjs = {};
      var hashObjs = {};
      for (var branchInd = 0; branchInd < branchesNames.length; branchInd++){
        options['branch'] = branchesNames[branchInd]
        console.log(options)
        let commits = gitlog(options);
        var commObjs = {};
        console.log(commits)

        for(var i = 0; i < commits.length; i++) {
            hashObjs[commits[i].hash]={}
            commObjs[commits[i].hash] = {
                Message : commits[i].subject,
                Author : commits[i].authorName,
                Committer  : commits[i].committerName,
                Timestamp  : toTimestamp(commits[i].authorDate),
                Hash      : commits[i].hash,
                Parenthashes : [commits[i].parentHashes],
                Signature   : null
               }
        }

        var newBranch = {
            branchName: branchesNames[branchInd],
            author: USER_ID,
            timestamp: 1,
            logs: commObjs
        }
        branchObjs[newBranch.branchName] = newBranch
      }

      repo = {
        repoName: path.basename(path.resolve(process.cwd())),
        author: USER_ID,
        timeStamp: 0,
        hashes: hashObjs,
        branches: branchObjs
      }

      console.log(USER_ID)
      console.log(repo)
    });
  }

  addRepo(){
    try{
        let response = await this.contract.submitTransaction('addNewRepo', "Hassan", "testRepo", pushlog)
        console.log("Submit Response=",response.toString())
    } catch(e){
        console.log(e)
    }
  }

  queryRepo(){

  }

  cloneRepo(){

  }

  addBranch(){

  }

  queryBranches(){

  }

  queryBranch(){

  }

  addCommits(){

  }

  queryBranchCommits(){

  }

  queryLastBranchCommit(){

  }

  async static setupGateway(){
    // 2.1 load the connection profile into a JS object
    let connectionProfile = yaml.safeLoad(fs.readFileSync(CONNECTION_PROFILE_PATH, 'utf8'));

    // 2.2 Need to setup the user credentials from wallet
    const wallet = new FileSystemWallet(FILESYSTEM_WALLET_PATH);

    // 2.3 Set up the connection options
    let connectionOptions = {
        identity: USER_ID,
        wallet: wallet,
        discovery: { enabled: false, asLocalhost: true }
    };

    const this.gateway = new Gateway();
    await this.gateway.connect(connectionProfile, connectionOptions);
  }

  async function submitTxnContract(contract, pushlog){
      try{
          let response = await this.contract.submitTransaction('addNewBranch', "Hassan", "testRepo", pushlog)
          console.log("Submit Response=",response.toString())
      } catch(e){
          console.log(e)
      }
  }

  async function queryContract(contract){
      try{
          let response = await this.contract.evaluateTransaction('queryBranch', 'Hassan', 'testRepo', 'master')
          console.log(`Query Response=${response.toString()}`)
      } catch(e){
          console.log(e)
      }
  }
}


function toTimestamp(strDate){
    var datum = Date.parse(strDate);
    return datum/1000;
}


var client = new Client()
client.loadCurrentRepo();
