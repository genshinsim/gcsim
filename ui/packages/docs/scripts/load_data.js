var Papa = require("papaparse");
var fs = require("fs");

var fields = Papa.parse(fs.readFileSync("./data/fields.csv").toString());
let fieldsJson = {};
fields.data.forEach((e, i) => {
  if (i == 0) {
    return; //skip first row
  }
  if (!(e[0] in fieldsJson)) {
    fieldsJson[e[0]] = [];
  }
  //split fields
  const f = e[1].split(" | ");
  fieldsJson[e[0]].push({
    fields: f,
    desc: e[2],
  });
});
fs.writeFileSync(
  "./src/components/Fields/data.json",
  JSON.stringify(fieldsJson, null, 2)
);
console.log("Done loading fields");

var frames = Papa.parse(fs.readFileSync("./data/frames.csv").toString());
let framesJson = {};
frames.data.forEach((e, i) => {
  if (i == 0) {
    return; //skip first row
  }
  if (!(e[0] in framesJson)) {
    framesJson[e[0]] = [];
  }
  //split fields
  framesJson[e[0]].push({
    vid_credit: e[1],
    count_credit: e[2],
    vid: e[3],
    count: e[4],
  });
});
fs.writeFileSync(
  "./src/components/Frames/data.json",
  JSON.stringify(framesJson, null, 2)
);
console.log("Done loading frames");

var iss = Papa.parse(fs.readFileSync("./data/issues.csv").toString());
let issObj = {};
iss.data.forEach((e, i) => {
  if (i == 0) {
    return; //skip first row
  }
  if (!(e[0] in issObj)) {
    issObj[e[0]] = [];
  }
  //split fields
  issObj[e[0]].push(e[1]);
});
fs.writeFileSync(
  "./src/components/Issues/data.json",
  JSON.stringify(issObj, null, 2)
);
console.log("Done loading issues");

var params = Papa.parse(fs.readFileSync("./data/params.csv").toString());
let paramsObj = {};
params.data.forEach((e, i) => {
  if (i == 0) {
    return; //skip first row
  }
  if (!(e[0] in paramsObj)) {
    paramsObj[e[0]] = [];
  }
  //split fields
  paramsObj[e[0]].push({
    ability: e[1],
    param: e[2],
    desc: e[3],
  });
});
fs.writeFileSync(
  "./src/components/Params/data.json",
  JSON.stringify(paramsObj, null, 2)
);
console.log("Done loading params");

var hl = Papa.parse(fs.readFileSync("./data/hitlag.csv").toString());
let hlObj = {};
hl.data.forEach((e, i) => {
  if (i == 0) {
    return; //skip first row
  }
  if (!(e[0] in hlObj)) {
    hlObj[e[0]] = {};
  }
  if (!(e[1] in hlObj[e[0]])) {
    hlObj[e[0]][e[1]] = [];
  }
  //split fields
  hlObj[e[0]][e[1]].push({
    ability: e[2],
    hitHaltTime: parseFloat(e[3]),
    hitHaltTimeScale: parseFloat(e[4]),
    canBeDefenseHalt: e[5].toLowerCase() === "true" ? true : false,
    deployable: e[6].toLowerCase() === "true" ? true : false,
  });
});
fs.writeFileSync(
  "./src/components/Hitlag/data.json",
  JSON.stringify(hlObj, null, 2)
);
console.log("Done loading hitlag");
