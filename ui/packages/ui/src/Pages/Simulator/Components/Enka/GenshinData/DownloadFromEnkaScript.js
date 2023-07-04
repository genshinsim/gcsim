const fs = require('fs');
const https = require('https');
let url = 'https://raw.githubusercontent.com/EnkaNetwork/API-docs/master/store/characters.json';
https.get(url,(res) => {
  const path = `./src/Components/Enka/GenshinData/EnkaCharacterMap.json`; 
  const filePath = fs.createWriteStream(path);
  res.pipe(filePath);
  filePath.on('finish',() => {
      filePath.close();
      console.log('Download Completed'); 
  })
})

url = 'https://raw.githubusercontent.com/EnkaNetwork/API-docs/master/store/loc.json'
https.get(url,(res) => {
  // use key "en" to extract english json
  var str = ''
  res.on('data', (data) => {
      str += data
  })
  res.on('end', async function() {
    var json = JSON.parse(str)
    var en = json['en']
    fs.writeFileSync(
          './src/Components/Enka/GenshinData/EnkaTextMapEN.json',
          JSON.stringify({"en":en})
        );
  })
})