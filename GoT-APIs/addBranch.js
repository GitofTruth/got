
/**
 * Demonstrates the use of Gateway Network & Contract classes
 */



 // Needed for reading the connection profile as JS object
const fs = require('fs');
// Used for parsing the connection profile YAML file
const yaml = require('js-yaml');
// Import gateway class
const { Gateway, FileSystemWallet, DefaultEventHandlerStrategies, Transaction  } = require('fabric-network');

// Constants for profile
const CONNECTION_PROFILE_PATH = '../profiles/dev-connection.yaml'
// Path to the wallet
const FILESYSTEM_WALLET_PATH = './user-wallet'
// Identity context used
const USER_ID = 'Admin@acme.com'
// Channel name
const NETWORK_NAME = 'airlinechannel'
// Chaincode
const CONTRACT_ID = "GoT"

//Git Log Parser
const gitlog = require('gitlog');

const repoPath = '../';

const options =
{ repo: repoPath
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




// 1. Create an instance of the gatway
const gateway = new Gateway();

// Sets up the gateway | executes the invoke & query
main()

/**
 * Executes the functions for query & invoke
 */
async function main() {
    
    // 2. Setup the gateway object
    await setupGateway()

    // 3. Get the network
    let network = await gateway.getNetwork(NETWORK_NAME)
    // console.log(network)

       // 5. Get the contract
    const contract = await network.getContract(CONTRACT_ID);
    // console.log(contract)

    //parse git logs
    // Synchronous
    let commits = gitlog(options);
    // console.log(commits[0]);

    // let commObjs = new Map();
    let commObjs = {};

    var dict = []; // create an empty array

    dict.push({
        key:   "keyName",
        value: "the value"
    });
    // commObjs.length = commits.length

    for(var i = 0; i < commits.length; i++) {
        commObjs[commits[i].hash] = {
            Message : commits[i].subject,
            Author : commits[i].authorName,
            Committer  : commits[i].committerName,
            Timestamp  : toTimestamp(commits[i].authorDate), 
            Hash      : commits[i].hash,
            Parenthashes : [commits[i].parentHashes],
            Signature   : null
           }

        // commObjs = {
        //     hash: {

        //     }
        // }

    }



    let newBranch = {

        branchName: "master",
        author: "Hassan",
        timestamp: 2,
        logs: commObjs
    }



    // let pushlog = JSON.stringify(pushObj);

    let newBranchString = JSON.stringify(newBranch);

    //pushObj.logs[0]

    //console.log(commit_str);

    // 7. Execute the transaction
    console.log(newBranchString)
    await submitTxnContract(contract,newBranchString)
    // Must give delay or use await here otherwise Error=MVCC_READ_CONFLICT
    // await submitTxnContract(contract)

    // Query the chaincode
   await queryContract(contract)


}

function toTimestamp(strDate){
    var datum = Date.parse(strDate);
    return datum/1000;
   }

/**
 * Queries the chaincode
 * @param {object} contract 
 */
async function queryContract(contract){
    try{
        // Query the chaincode
        let response = await contract.evaluateTransaction('queryBranch', 'Hassan', 'testRepo', 'master')
        console.log(`Query Response=${response.toString()}`)
    } catch(e){
        console.log(e)
    }
}

/**
 * Submit the transaction
 * @param {object} contract 
 */
async function submitTxnContract(contract, pushlog){
    try{
        // Submit the transaction
        let response = await contract.submitTransaction('addNewBranch', "Hassan", "testRepo", pushlog)
        console.log("Submit Response=",response.toString())
    } catch(e){
        // fabric-network.TimeoutError
        console.log(e)
    }
}

/**
 * Function for setting up the gateway
 * It does not actually connect to any peer/orderer
 */
async function setupGateway() {
    
    // 2.1 load the connection profile into a JS object
    let connectionProfile = yaml.safeLoad(fs.readFileSync(CONNECTION_PROFILE_PATH, 'utf8'));

    // 2.2 Need to setup the user credentials from wallet
    const wallet = new FileSystemWallet(FILESYSTEM_WALLET_PATH)

    // 2.3 Set up the connection options
    let connectionOptions = {
        identity: USER_ID,
        wallet: wallet,
        discovery: { enabled: false, asLocalhost: true }
        /*** Uncomment lines below to disable commit listener on submit ****/
        // , eventHandlerOptions: {
        //     strategy: null
        // } 
    }        /*** Uncomment lines below to disable commit listener on submit ****/
    // , eventHandlerOptions: {
    //     strategy: null
    // } 

    // 2.4 Connect gateway to the network
    await gateway.connect(connectionProfile, connectionOptions)
    // console.log( gateway)
}

