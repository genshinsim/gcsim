const tsj = require("ts-json-schema-generator");
const fs = require("fs");

const config = {
  path: "./src/Pages/Viewer/DataType.ts",
  tsconfig: "./tsconfig.json",
  type: "*", // Or <type-name> if you want to generate schema for that one type only
  expose: "none",
  additionalProperties: true,
};

const output_path = "./src/Pages/Viewer/DataType.schema.json";

const schema = tsj.createGenerator(config).createSchema(config.type);
const schemaString = JSON.stringify(schema, null, 2);
fs.writeFile(output_path, schemaString, (err) => {
  if (err) throw err;
});
