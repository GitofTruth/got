// internal imports
// var Constants = require("constants");
const CONNECTION_PROFILE_PATH = '.gotconfig/profiles/dev-connection.yaml';
const FILESYSTEM_WALLET_PATH = '.gotconfig/user-wallet';

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

const simpleGit = require('simple-git')(repoPath);


var gateway, network, contract;

class Client {
  constructor (){
  }

  async loadCurrentRepo(){
    await simpleGit.branchLocal(function(e,d){
      branchesNames = d['all']
      branchesContent = d['branches']

      var branchObjs = {};
      var hashObjs = {};
      for (var branchInd = 0; branchInd < branchesNames.length; branchInd++){
        options['branch'] = branchesNames[branchInd]
        // console.log(options)
        let commits = gitlog(options);
        var commObjs = {};
        // console.log(commits)

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
      };

      console.log(USER_ID)
      console.log(repo)
    });
  }

  async addRepo(){
    try{
        console.log("sending this repo:", repo)
        let newRepoString = JSON.stringify(repo);
        let response = await contract.submitTransaction('addNewRepo', newRepoString)
        console.log("Submit Response=",response.toString())
    } catch(e){
        console.log(e)
    }
  }

  async queryRepo(){
    try{
        let response = await contract.evaluateTransaction('queryRepo', repo['author'], repo['repoName'])
        console.log("Submit Response=",response.toString())
    } catch(e){
        console.log(e)
    }
  }

  async cloneRepo(){
    try{
        let response = await contract.evaluateTransaction('clone', repo['author'], repo['repoName'])
        console.log("Submit Response=",response.toString())
    } catch(e){
        console.log(e)
    }
  }

  async addBranch(branchName){
    try{
        let newBranchString = JSON.stringify(repo['branches'][branchName]);
        let response = await contract.submitTransaction('addNewBranch', repo['author'], repo['repoName'], newBranchString)
        console.log("Submit Response=",response.toString())
    } catch(e){
        console.log(e)
    }
  }

  async queryBranches(){
    try{
        let response = await contract.evaluateTransaction('queryBranches', repo['author'], repo['repoName'])
        console.log("Submit Response=",response.toString())
    } catch(e){
        console.log(e)
    }
  }

  async queryBranch(branchName){
    try{
        let response = await contract.evaluateTransaction('queryBranch', repo['author'], repo['repoName'], branchName)
        console.log("Submit Response=",response.toString())
    } catch(e){
        console.log(e)
    }
  }

  async addCommits(){

  }

  async queryBranchCommits(branchName, lastCommit){
    try{
        let response = await contract.evaluateTransaction('queryBranchCommits', repo['author'], repo['repoName'], branchName, lastCommit)
        console.log("Submit Response=",response.toString())
    } catch(e){
        console.log(e)
    }
  }

  async queryLastBranchCommit(branchName){
    try{
        let response = await contract.evaluateTransaction('queryLastBranchCommit', repo['author'], repo['repoName'], branchName)
        console.log("Submit Response=",response.toString())
    } catch(e){
        console.log(e)
    }
  }
  //
  // async function submitTxnContract(contract, pushlog){
  //     try{
  //         let response = await this.contract.submitTransaction('addNewBranch', "Hassan", "testRepo", pushlog)
  //         console.log("Submit Response=",response.toString())
  //     } catch(e){
  //         console.log(e)
  //     }
  // }
  //
  // async function queryContract(contract){
  //     try{
  //         let response = await this.contract.evaluateTransaction('queryBranch', 'Hassan', 'testRepo', 'master')
  //         console.log(`Query Response=${response.toString()}`)
  //     } catch(e){
  //         console.log(e)
  //     }
  // }
}

async function setupGateway(){
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

  await gateway.connect(connectionProfile, connectionOptions);
}


function toTimestamp(strDate){
    var datum = Date.parse(strDate);
    return datum/1000;
}

async function main(){
  gateway = new Gateway();
  await setupGateway();
  network = await gateway.getNetwork(NETWORK_NAME);
  contract = await network.getContract(CONTRACT_ID);


  var client = new Client();
  await client.loadCurrentRepo();
  await client.addRepo();
  await client.queryRepo();
  await client.cloneRepo();
  await client.queryBranches();
  await client.queryLastBranchCommit("master");

  return 0;
}

main()
