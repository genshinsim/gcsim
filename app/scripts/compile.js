// import fs from "fs";
const fs = require("fs");
const toml = require("toml");
const Fuse = require("fuse.js");
const genshindb = require("genshin-db");

const dir = "./gsimactions";

const getFiles = (path) => {
  const files = [];
  for (const file of fs.readdirSync(path)) {
    const fullPath = path + "/" + file;
    if (fs.lstatSync(fullPath).isDirectory())
      getFiles(fullPath).forEach((x) => files.push(file + "/" + x));
    else files.push(file);
  }
  return files;
};

let result = [];

const files = getFiles(dir);

// console.log(files);

files.forEach((file) => {
  //skip none toml files
  if (!file.endsWith(".toml")) {
    return;
  }
  console.log(`reading ${dir + "/" + file}`);

  const data = toml.parse(fs.readFileSync(dir + "/" + file, "utf-8"));

  data.path = dir + "/" + file;
  //   console.log(data);

  //add to result
  result.push(data);
});

// Create the Fuse index
const myIndex = Fuse.createIndex(
  ["title", "author", "description", "characters"],
  result
);
// Serialize and save it
fs.writeFileSync(
  "./src/data/fuse-index.json",
  JSON.stringify(myIndex.toJSON())
);
fs.writeFileSync("./src/data/configs.json", JSON.stringify(result));

//grap character names

let charImages = {};

const chars = genshindb.characters("names", { matchCategories: true });

chars.forEach((char) => {
  const x = genshindb.characters(char);

  let key = x.name.replace(/[^0-9a-z]/gi, "").toLowerCase();

  charImages[key] = x.images.icon;
});

fs.writeFileSync(
  "./src/data/character_images.json",
  JSON.stringify(charImages)
);
