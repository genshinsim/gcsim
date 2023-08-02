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
  English: { artifact_names: {}, character_names: {}, weapon_names: {} },
  Chinese: { artifact_names: {}, character_names: {}, weapon_names: {} },
  Japanese: { artifact_names: {}, character_names: {}, weapon_names: {} },
  Spanish: { artifact_names: {}, character_names: {}, weapon_names: {} },
  Russian: { artifact_names: {}, character_names: {}, weapon_names: {} },
  German: { artifact_names: {}, character_names: {}, weapon_names: {} },
};

names.forEach((e) => {
  const eng = genshindb.characters(e);
  if (!eng) return;
  //key is e 
  let key = e

  const cn = genshindb.characters(e, { resultLanguage: "CHS" });
  const jp = genshindb.characters(e, { resultLanguage: "JP" });
  const es = genshindb.characters(e, { resultLanguage: "ES" });
  const ru = genshindb.characters(e, { resultLanguage: "RU" });
  const de = genshindb.characters(e, { resultLanguage: "DE" });

  trans["English"]["character_names"][e] = eng.name;
  trans["Chinese"]["character_names"][e] = cn.name;
  trans["Japanese"]["character_names"][e] = jp.name;
  trans["Spanish"]["character_names"][e] = es.name;
  trans["Russian"]["character_names"][e] = ru.name;
  trans["German"]["character_names"][e] = de.name;
});

travelers.map((e) => {
  const key = `traveler${e}`;
  const mcm = genshindb.characters("aether");
  const mcf = genshindb.characters("lumine")
  
  const eng = genshindb.talents(key);
  if (!mcm || !mcf || !eng) return;

  //abit hacky,  we only want what's in brackets
  const cn = genshindb.talents(key, { resultLanguage: "CHS" });
  const jp = genshindb.talents(key, { resultLanguage: "JP" });
  const es = genshindb.talents(key, { resultLanguage: "ES" });
  const ru = genshindb.talents(key, { resultLanguage: "RU" });
  const de = genshindb.talents(key, { resultLanguage: "DE" });

  //names have to be right language too...
  const mcn = genshindb.characters("aether", { resultLanguage: "CHS" });
  const mjp = genshindb.characters("aether", { resultLanguage: "JP" });
  const mes = genshindb.characters("aether", { resultLanguage: "ES" });
  const mru = genshindb.characters("aether", { resultLanguage: "RU" });
  const mde = genshindb.characters("aether", { resultLanguage: "DE" });

  const fcn = genshindb.characters("lumine", { resultLanguage: "CHS" });
  const fjp = genshindb.characters("lumine", { resultLanguage: "JP" });
  const fes = genshindb.characters("lumine", { resultLanguage: "ES" });
  const fru = genshindb.characters("lumine", { resultLanguage: "RU" });
  const fde = genshindb.characters("lumine", { resultLanguage: "DE" });

  trans["English"]["character_names"][`aether${e}`] = eng.name.replace(/.*?(\(.*\)).*/, `${mcm.name} $1`);
  trans["Chinese"]["character_names"][`aether${e}`] = cn.name.replace(/.*?(\(.*\)).*/, `${mcn.name} $1`);
  trans["Japanese"]["character_names"][`aether${e}`] = jp.name.replace(/.*?(\(.*\)).*/, `${mjp.name} $1`);
  trans["Spanish"]["character_names"][`aether${e}`] = es.name.replace(/.*?(\(.*\)).*/, `${mes.name} $1`);
  trans["Russian"]["character_names"][`aether${e}`] = ru.name.replace(/.*?(\(.*\)).*/, `${mru.name} $1`);
  trans["German"]["character_names"][`aether${e}`] = de.name.replace(/.*?(\(.*\)).*/, `${mde.name} $1`);

  trans["English"]["character_names"][`lumine${e}`] = eng.name.replace(/.*?(\(.*\)).*/, `${mcf.name} $1`);
  trans["Chinese"]["character_names"][`lumine${e}`] = cn.name.replace(/.*?(\(.*\)).*/, `${fcn.name} $1`);
  trans["Japanese"]["character_names"][`lumine${e}`] = jp.name.replace(/.*?(\(.*\)).*/, `${fjp.name} $1`);
  trans["Spanish"]["character_names"][`lumine${e}`] = es.name.replace(/.*?(\(.*\)).*/, `${fes.name} $1`);
  trans["Russian"]["character_names"][`lumine${e}`] = ru.name.replace(/.*?(\(.*\)).*/, `${fru.name} $1`);
  trans["German"]["character_names"][`lumine${e}`] = de.name.replace(/.*?(\(.*\)).*/, `${fde.name} $1`);
});

//download weapons and sets :(

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

  const cn = genshindb.weapons(e, { resultLanguage: "CHS" });
  const jp = genshindb.weapons(e, { resultLanguage: "JP" });
  const es = genshindb.weapons(e, { resultLanguage: "ES" });
  const ru = genshindb.weapons(e, { resultLanguage: "RU" });
  const de = genshindb.weapons(e, { resultLanguage: "DE" });

  trans["English"]["weapon_names"][key] = eng.name;
  trans["Chinese"]["weapon_names"][key] = cn.name;
  trans["Japanese"]["weapon_names"][key] = jp.name;
  trans["Spanish"]["weapon_names"][key] = es.name;
  trans["Russian"]["weapon_names"][key] = ru.name;
  trans["German"]["weapon_names"][key] = de.name;

});


let setMap = {};
const sets = genshindb.artifacts("4", { matchCategories: true });

sets.forEach((e) => {
  const eng = genshindb.artifacts(e);
  if (!eng) return;

  let art = eng.name.replace(/[^0-9a-z]/gi, "").toLowerCase();
  setMap[art] = eng.name;

  const cn = genshindb.artifacts(e, { resultLanguage: "CHS" });
  const jp = genshindb.artifacts(e, { resultLanguage: "JP" });
  const es = genshindb.artifacts(e, { resultLanguage: "ES" });
  const ru = genshindb.artifacts(e, { resultLanguage: "RU" });
  const de = genshindb.artifacts(e, { resultLanguage: "DE" });

  trans["English"]["artifact_names"][art] = eng.name;
  trans["Chinese"]["artifact_names"][art] = cn.name;
  trans["Japanese"]["artifact_names"][art] = jp.name;
  trans["Spanish"]["artifact_names"][art] = es.name;
  trans["Russian"]["artifact_names"][art] = ru.name;
  trans["German"]["artifact_names"][art] = de.name;

});

fs.writeFileSync(
  "./IngameNames.json",
  JSON.stringify(trans, null, 2),
  "utf-8"
);
