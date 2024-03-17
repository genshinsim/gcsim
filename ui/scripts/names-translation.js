const genshindb = require("genshin-db");
const fs = require("fs");
const path = require("path");

let names = []
const rawdata = fs.readFileSync('./packages/ui/src/Data/char_data.generated.json')
let chardata = JSON.parse(rawdata)

for (const [key, _ ] of Object.entries(chardata.data)) {
    //skip lumine and aether
    if (key.includes("lumine") || key.includes("aether")) {
        continue
    }
    names.push(key)
}
names.push("aether")
names.push("lumine")

const travelers = [
  "electro",
  "anemo",
  "geo",
  "hydro",
  "cryo",
  "pyro",
  "dendro",
];

let trans = {
  English: { artifact_names: {}, character_names: {}, enemy_names: {}, weapon_names: {} },
  Chinese: { artifact_names: {}, character_names: {}, enemy_names: {}, weapon_names: {} },
  Japanese: { artifact_names: {}, character_names: {}, enemy_names: {}, weapon_names: {} },
  Korean: { artifact_names: {}, character_names: {}, enemy_names: {}, weapon_names: {} },
  Spanish: { artifact_names: {}, character_names: {}, enemy_names: {}, weapon_names: {} },
  Russian: { artifact_names: {}, character_names: {}, enemy_names: {}, weapon_names: {} },
  German: { artifact_names: {}, character_names: {}, enemy_names: {}, weapon_names: {} },
};

const languages = {
  English: "EN",
  Chinese: "CHS",
  Japanese: "JP",
  Korean: "KR",
  Spanish: "ES",
  Russian: "RU",
  German: "DE"
}

names.forEach((e) => {
  const eng = genshindb.characters(e);
  if (!eng) return;
  //key is e

  for (const [langName, langCode] of Object.entries(languages)) {
    const lang = genshindb.characters(e, { resultLanguage: langCode });
    trans[langName]["character_names"][e] = lang.name;
  }
});

travelers.map((e) => {
  const key = `traveler${e}`;
  const mcm = genshindb.characters("aether");
  const mcf = genshindb.characters("lumine")
  
  const eng = genshindb.talents(key);
  if (!mcm || !mcf || !eng) return;
  
  for (const [langName, langCode] of Object.entries(languages)) {
    //abit hacky,  we only want what's in brackets
    const lang = genshindb.talents(key, { resultLanguage: langCode });
    
    //names have to be right language too...
    const aether = genshindb.characters("aether", { resultLanguage: langCode });
    const lumine = genshindb.characters("lumine", { resultLanguage: langCode });
    
    trans[langName]["character_names"][`aether${e}`] = lang.name.replace(/.*?(\(.*\)).*/, `${aether.name} $1`);
    trans[langName]["character_names"][`lumine${e}`] = lang.name.replace(/.*?(\(.*\)).*/, `${lumine.name} $1`);
  }
});

//download weapons, sets and enemies :(

const weapons = genshindb.weapons("names", { matchCategories: true });

let weap = {};

weapons.forEach((e) => {
  const eng = genshindb.weapons(e);
  if (!eng) return;

  let filename =
    "./static/images/weapons/" +
    eng.name.replace(/[^0-9a-z]/gi, "").toLowerCase() +
    ".png";

  const key = eng.name.replace(/[^0-9a-z]/gi, "").toLowerCase();
  weap[key] = eng.name;

  for (const [langName, langCode] of Object.entries(languages)) {
    const lang = genshindb.weapons(e, { resultLanguage: langCode });
    trans[langName]["weapon_names"][key] = lang.name;
  }
});


let setMap = {};
const sets = genshindb.artifacts("4", { matchCategories: true });

sets.forEach((e) => {
  const eng = genshindb.artifacts(e);
  if (!eng) return;

  const key = eng.name.replace(/[^0-9a-z]/gi, "").toLowerCase();
  setMap[key] = eng.name;

  for (const [langName, langCode] of Object.entries(languages)) {
    const lang = genshindb.artifacts(e, { resultLanguage: langCode });
    trans[langName]["artifact_names"][key] = lang.name;
  }
});

let enemiesMap = {};
const enemies = genshindb.enemies("names", { matchCategories: true });

enemies.forEach((e) => {
  const eng = genshindb.enemies(e);
  if (!eng) return;

  const key = eng.name.replace(/[^0-9a-z]/gi, "").toLowerCase();
  enemiesMap[key] = eng.name;

  for (const [langName, langCode] of Object.entries(languages)) {
    const lang = genshindb.enemies(e, { resultLanguage: langCode });
    trans[langName]["enemy_names"][key] = lang.name;
  }
});


fs.writeFileSync(
  "./IngameNames.json",
  JSON.stringify(trans, null, 2),
  "utf-8"
);
