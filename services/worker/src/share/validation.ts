import { Validator } from "@cfworker/json-schema";
import { SummaryStats, Character } from "../dataType";

export interface uploadData {
  data: string;
  meta: {
    char_names: string[];
    dps: SummaryStats;
    sim_duration: SummaryStats;
    dps_by_target: { [key: number]: SummaryStats };
    iter: number;
    runtime: number;
    num_targets: number;
    char_details: Character[];
  };
  path?: string; //for organization purposes
  perm?: boolean;
}

export const validator = new Validator({
  type: "object",
  required: ["data", "meta"],
  properties: {
    data: { type: "string" }, //base64 string of zlib data
    meta: {
      type: "object",
      required: [
        "char_names",
        "dps",
        "dps_by_target",
        "sim_duration",
        "iter",
        "runtime",
        "char_details",
        "num_targets",
      ],
      description: { type: "string" },
      char_names: {
        type: "array",
        items: {
          type: "string",
        },
      },
      dps: {
        type: "object",
        properties: {
          mean: {
            type: "number",
          },
          min: {
            type: "number",
          },
          max: {
            type: "number",
          },
          sd: {
            type: "number",
          },
        },
        required: ["mean", "min", "max"],
      },
      dps_by_target: {
        type: "object",
        additionalProperties: {
          type: "object",
          properties: {
            mean: {
              type: "number",
            },
            min: {
              type: "number",
            },
            max: {
              type: "number",
            },
            sd: {
              type: "number",
            },
          },
          required: ["mean", "min", "max"],
        },
      },
      sim_duration: {
        type: "object",
        properties: {
          mean: {
            type: "number",
          },
          min: {
            type: "number",
          },
          max: {
            type: "number",
          },
          sd: {
            type: "number",
          },
        },
        required: ["mean", "min", "max"],
      },
      iter: {
        type: "number",
      },
      runtime: {
        type: "number",
      },
      num_targets: {
        type: "number",
      },
      char_details: {
        type: "array",
        items: {
          type: "object",
          properties: {
            name: {
              type: "string",
            },
            level: {
              type: "number",
            },
            element: {
              type: "string",
            },
            max_level: {
              type: "number",
            },
            cons: {
              type: "number",
            },
            weapon: {
              type: "object",
              properties: {
                name: {
                  type: "string",
                },
                refine: {
                  type: "number",
                },
                level: {
                  type: "number",
                },
                max_level: {
                  type: "number",
                },
              },
              required: ["name", "refine", "level", "max_level"],
            },
            talents: {
              type: "object",
              properties: {
                attack: {
                  type: "number",
                },
                skill: {
                  type: "number",
                },
                burst: {
                  type: "number",
                },
              },
              required: ["attack", "skill", "burst"],
            },
            sets: {
              type: "object",
              additionalProperties: {
                type: "number",
              },
            },
          },
          required: [
            "name",
            "level",
            "element",
            "max_level",
            "cons",
            "weapon",
            "talents",
            "sets",
          ],
        },
      },
    },
    path: { type: "string" }, //organization purpose (fake folder names)
    perm: { type: "boolean" }, //if link should be permanent
  },
});
