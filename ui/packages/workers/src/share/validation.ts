import { Validator } from "@cfworker/json-schema";

export const validator = new Validator({
  type: "object",
  required: ["sim_version"],
  properties: {
    sim_version: { type: "string" }, // sim hash
  },
});
