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
  "yaemiko",
];

let chars = {};
let properKeyToChar = {};

let trans = {
  English: {},
  Chinese: {},
  Japanese: {},
  Spanish: {},
};

names.forEach((e) => {
  const eng = genshindb.characters(e);
  let key = eng.name.replace(/[^0-9a-z]/gi, "").toLowerCase();

  chars[e] = {
    key: e,
    name: eng.name,
    element: eng.element.toLowerCase(),
    weapon_type: eng.weapontype.toLowerCase(),
  };
  properKeyToChar[key] = e;

  const cn = genshindb.characters(e, { resultLanguage: "CHS" });
  const jp = genshindb.characters(e, { resultLanguage: "JP" });
  const es = genshindb.characters(e, { resultLanguage: "ES" });

  trans["English"][e] = eng.name;
  trans["Chinese"][e] = cn.name;
  trans["Japanese"][e] = jp.name;
  trans["Spanish"][e] = es.name;

  let filename = "./static/images/avatar/" + e + ".png";

  if (!fs.existsSync(filename)) {
    console.log(e + ": " + eng.images.icon);

    download_image(eng.images.icon, filename)
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

fs.writeFileSync(
  "./public/locales/characters.json",
  JSON.stringify(trans),
  "utf-8"
);

//download weapons and sets :(

const weapons = genshindb.weapons("names", { matchCategories: true });

let weap = {};
let weapTrans = {
  English: {},
  Chinese: {},
  Japanese: {},
  Spanish: {},
};

weapons.forEach((e) => {
  const eng = genshindb.weapons(e);

  let filename =
    "./static/images/weapons/" +
    eng.name.replace(/[^0-9a-z]/gi, "").toLowerCase() +
    ".png";

  const key = eng.name.replace(/[^0-9a-z]/gi, "").toLowerCase();
  weap[key] = eng.name;

  const cn = genshindb.weapons(e, { resultLanguage: "CHS" });
  const jp = genshindb.weapons(e, { resultLanguage: "JP" });
  const es = genshindb.weapons(e, { resultLanguage: "ES" });

  weapTrans["English"][key] = eng.name;
  weapTrans["Chinese"][key] = cn.name;
  weapTrans["Japanese"][key] = jp.name;
  weapTrans["Spanish"][key] = es.name;

  if (!fs.existsSync(filename)) {
    download_image(eng.images.icon, filename)
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

fs.writeFileSync(
  "./public/locales/weapons.json",
  JSON.stringify(weapTrans),
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
