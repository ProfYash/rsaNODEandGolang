const axios = require('axios')
const NodeRSA = require('node-rsa');
const key = new NodeRSA();
async function getpublickey() {
    let resp = await axios.get('http://localhost:9002/api/getpublicKey')
    if (resp.data == null) {
        console.log("no publickey")
    } else {
        console.log(resp.data)
        let keydata = Buffer.from(resp.data)
        key.importKey(keydata, 'pkcs8-public-pem')
        // console.log(key)
        await encrypt()
    }
}
getpublickey()

async function encrypt(){
    let encryptedtxt = key.encrypt('hello AMAR','hex')
    console.log(encryptedtxt)
}