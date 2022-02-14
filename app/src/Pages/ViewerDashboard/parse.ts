import pako from "pako";
import { ResultsSummary } from "~src/types";
import Ajv from "ajv";

import schema from "./DataType.schema.json";

const ajv = new Ajv();

export function Uint8ArrayFromBase64(base64: string) {
  return Uint8Array.from(window.atob(base64), (v) => v.charCodeAt(0));
}

export function extractJSONStringFromBinary(binaryStr: Uint8Array): {
  err: string;
  data: string;
} {
  try {
    const restored = pako.inflate(binaryStr, { to: "string" });
    return {
      err: "",
      data: restored,
    };
  } catch {
    return {
      err: "Not a valid gzipped JSON file",
      data: "",
    };
  }
}

export function parseAndValidate(jsonStr: string): ResultsSummary | string {
  let data: ResultsSummary = JSON.parse(jsonStr);
  const validate = ajv.compile(schema.definitions["*"]);
  const valid = validate(data);

  if (valid) {
    return data;
  }
  return JSON.stringify(validate.errors, null, 2);
}
