const ipfsClient = require('ipfs-http-client');
const express = require('express');
const bodyParser = require('body-parser');
const fileupload = require('express-fileupload');
const fs = require('fs');


var ipfs = new ipfsClient(
    {
        host: 'localhost',
        port: '5001',
        protocol: 'http'
    }
);
var app = express();


app.use(function(req, res, next) {
    res.setHeader("Access-Control-Allow-Origin", "*");
    res.setHeader("Access-Control-Allow-Credentials", "true");
    res.setHeader("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT");
    res.setHeader("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers");
    next();
  });

app.set('view engine', 'ejs');
app.use(bodyParser.urlencoded({extended: true}));

app.use(fileupload());

var port = process.env.PORT || 3000;

app.listen(port, function() {
    console.log(`Server listening on 'http://localhost:${port}'.`);
  });

app.get('/', (req,res) => {
    console.log('rendering home page');
    res.render('home');
});

app.get('/hw', (req, res)=>{
    res.send('Hello World!');
})

app.post('/upload', (req,res) => {

    const file = req.files.file;
    const fileName = req.body.fileName;
    const filePath = 'files/' + fileName;



    file.mv(filePath, async (err) => {
            if(err) {
                console.log("Error: failed to download the file");
                return res.status(500).send(err);
            }

            const fileHash = await addFile(fileName, filePath);

            fs.unlink(filePath, (err) =>{
                if(err) console.log(err);
            }
            );

            res.render('upload', {
                fileName,
                fileHash
            });
    });


});

const addFile = async (fileName, filePath) => {
    const file =  fs.readFileSync(filePath);

    for await (const result of ipfs.add(file)) {
        console.log(result)
        return result.path;
      }

}


module.exports = app;