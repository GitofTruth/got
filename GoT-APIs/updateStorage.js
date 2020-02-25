// const { readdir, stat } = require("fs").promises
const fs = require('fs');
const { join } = require("path")

let testpath = "../.git/objects/"
let defStoragePath = "../../storage/"
let excptList = ["info", "pack"]
//
// async function ListDirectories(path){
//   console.log("getting dirs")
//    path => {
//     let dirs = []
//     for (const file of await readdir(path)) {
//       if ((await stat(join(path, file))).isDirectory()) {
//         dirs = [...dirs, file]
//         console.log(file)
//       }
//     }
//
//     console.log("finished getting dirs")
//     return dirs
//   }
// }


function ListFiles(testFolder){
  let files = []
  fs.readdirSync(testFolder).forEach(file => {
    if( excptList.indexOf(file) > -1){
    }else {
      files = [... files, file]
      // console.log(file);
    }
  });

  return files
}


function GetObjectsName(objHomeDirectory) {
  let names = []
  let dirs = ListFiles(objHomeDirectory)
  dirs.forEach(d => {
    files = ListFiles(objHomeDirectory + d + "/")
    files.forEach(f => {
      names = [...names, d + f]
    })
  })

  return names;
}

function CopyObjects(namesList = [], storagePath = defStoragePath){
  let localObjsPath = ".git/objects/"
  let localObjs = GetObjectsName(localObjsPath)
  let missing = []

  if (namesList.length > 0){
    namesList.forEach(n => {
      if(localObjs.indexOf(n) <= -1){
        missing = [...missing, n]
      }
    })
  }{
    missing = localObjs
  }

  missing.forEach( m => {
    fs.copyFile(storagePath + m[0:1] + "/" + m[2:], localObjsPath + m[0:1] + "/" + m[2:], (err) => {
    if (err) throw err;
    console.log(m + " was added to bare repo");
    });
  })
}

// GetObjectsName(testpath)

console.log(GetObjectsName(testpath))


// var ncp = require('ncp').ncp;
//
// ncp.limit = 16;
//
// ncp(source, destination, function (err) {
//  if (err) {
//    return console.error(err);
//  }
//  console.log('done!');
// });




// const fs = require('fs');
//
// // destination.txt will be created or overwritten by default.
