const genshindb = require("genshin-db");
const axios = require("axios");
const fs = require("fs");
const path = require("path");

const download_image = (url, image_path) =>
  axios({
    url,
    responseType: "stream",
  })
    .then(
      (response) =>
        new Promise((resolve, reject) => {
          response.data
            .pipe(fs.createWriteStream(image_path))
            .on("finish", () => resolve())
            .on("error", (e) => reject(e));
        })
    )
    .catch((e) => console.log(e));

const names = [
  "aether",
  "lumine",
  "albedo",
  "aloy",
  "amber",
  "barbara",
  "beidou",
  "bennett",
  "chongyun",
  "diluc",
  "diona",
  "eula",
  "fischl",
  "ganyu",
  "hutao",
  "jean",
  "kazuha",
  "kaeya",
  "ayaka",
  "keqing",
  "klee",
  "sara",
  "lisa",
  "mona",
  "ningguang",
  "noelle",
  "qiqi",
  "raiden",
  "razor",
  "rosaria",
  "kokomi",
  "sayu",
  "sucrose",
  "tartaglia",
  "thoma",
  "venti",
  "xiangling",
  "xiao",
  "xingqiu",
  "xinyan",
  "yanfei",
  "yoimiya",
  "zhongli",
  "gorou",
  "itto",
  "shenhe",
  "yunjin",
];

let chars = {};
let properKeyToChar = {};

names.forEach((e) => {
  const x = genshindb.characters(e);
  let key = x.name.replace(/[^0-9a-z]/gi, "").toLowerCase();

  chars[e] = {
    key: e,
    name: x.name,
    element: x.element.toLowerCase(),
    weapon_type: x.weapontype.toLowerCase(),
  };
  properKeyToChar[key] = e;

  let filename = "./static/images/avatar/" + e + ".png";

  if (!fs.existsSync(filename)) {
    console.log(e + ": " + x.images.icon);

    download_image(x.images.icon, filename)
      .then((msg) => {
        console.log("done downloading to file: ", filename);
      })
      .catch((e) => {
        console.log(e);
      });
  }
});

fs.writeFileSync(
  "./src/Components/data/charNames.json",
  JSON.stringify(chars),
  "utf-8"
);

fs.writeFileSync(
  "./src/Components/data/charKeyToShort.json",
  JSON.stringify(properKeyToChar),
  "utf-8"
);

//download weapons and sets :(

const weapons = genshindb.weapons("names", { matchCategories: true });

let weap = {};

weapons.forEach((e) => {
  const x = genshindb.weapons(e);

  let filename =
    "./static/images/weapons/" +
    x.name.replace(/[^0-9a-z]/gi, "").toLowerCase() +
    ".png";

  weap[x.name.replace(/[^0-9a-z]/gi, "").toLowerCase()] = x.name;

  if (!fs.existsSync(filename)) {
    download_image(x.images.icon, filename)
      .then((msg) => {
        console.log("done downloading to file: ", filename);
      })
      .catch((e) => {
        console.log(e);
      });
  }
});

fs.writeFileSync(
  "./src/Components/data/weaponNames.json",
  JSON.stringify(weap),
  "utf-8"
);

let setMap = {};

const sets = genshindb.artifacts("4", { matchCategories: true });

sets.forEach((e) => {
  const x = genshindb.artifacts(e);

  let art = x.name.replace(/[^0-9a-z]/gi, "").toLowerCase();
  setMap[art] = x.name;

  let filename;
  for (const [key, value] of Object.entries(x.images)) {
    filename = `./static/images/artifacts/${art}_${key}.png`;

    if (!fs.existsSync(filename)) {
      console.log(`${key}: ${value}`);
      download_image(value, filename)
        .then(() => {
          console.log("done downloading to file: ", filename);
        })
        .catch((e) => {
          console.log(e);
        });
    }
  }
});

fs.writeFileSync(
  "./src/Components/data/artifactNames.json",
  JSON.stringify(setMap),
  "utf-8"
);
