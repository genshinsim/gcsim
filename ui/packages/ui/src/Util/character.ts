import { Character } from "@gcsim/types";
import { StatToIndexMap } from "./util";

type CharViewableStats = {
  [key in string]: {
    name: string;
    val: {
      [key in string]: {
        flat: number;
        per: number;
      };
    };
    flatIndex: number;
    percentIndex: number;
    count: number;
    t: string;
  };
};

export type CharStatBlock = {
  key: string;
  name: string;
  t: string;
  flat: number;
  percent: number;
};

export function ConsolidateCharStats(chars: Character[]): {
  stats: { [key in string]: CharStatBlock[] };
  maxRows: number;
} {
  const totalStats: CharViewableStats = {
    hp: {
      name: "hp / hp%",
      flatIndex: StatToIndexMap["HP"],
      percentIndex: StatToIndexMap["HPP"],
      val: {},
      count: 0,
      t: "both",
    },
    atk: {
      name: "atk / atk%",
      flatIndex: StatToIndexMap["ATK"],
      percentIndex: StatToIndexMap["ATKP"],
      val: {},
      count: 0,
      t: "both",
    },
    def: {
      name: "def / def%",
      flatIndex: StatToIndexMap["DEF"],
      percentIndex: StatToIndexMap["DEFP"],
      val: {},
      count: 0,
      t: "both",
    },
    em: {
      name: "em",
      flatIndex: StatToIndexMap["EM"],
      percentIndex: -1,
      val: {},
      count: 0,
      t: "f",
    },
    er: {
      name: "er",
      flatIndex: -1,
      percentIndex: StatToIndexMap["ER"],
      val: {},
      count: 0,
      t: "%",
    },
    cr: {
      name: "cr",
      flatIndex: -1,
      percentIndex: StatToIndexMap["CR"],
      val: {},
      count: 0,
      t: "%",
    },
    cd: {
      name: "cd",
      flatIndex: -1,
      percentIndex: StatToIndexMap["CD"],
      val: {},
      count: 0,
      t: "%",
    },
    electro: {
      name: "electro%",
      flatIndex: -1,
      percentIndex: StatToIndexMap["ElectroP"],
      val: {},
      count: 0,
      t: "%",
    },
    pyro: {
      name: "pyro%",
      flatIndex: -1,
      percentIndex: StatToIndexMap["PyroP"],
      val: {},
      count: 0,
      t: "%",
    },
    cryo: {
      name: "cryo%",
      flatIndex: -1,
      percentIndex: StatToIndexMap["CryoP"],
      val: {},
      count: 0,
      t: "%",
    },
    hydro: {
      name: "hydro%",
      flatIndex: -1,
      percentIndex: StatToIndexMap["HydroP"],
      val: {},
      count: 0,
      t: "%",
    },
    geo: {
      name: "geo%",
      flatIndex: -1,
      percentIndex: StatToIndexMap["GeoP"],
      val: {},
      count: 0,
      t: "%",
    },
    anemo: {
      name: "anemo%",
      flatIndex: -1,
      percentIndex: StatToIndexMap["AnemoP"],
      val: {},
      count: 0,
      t: "%",
    },
    phys: {
      name: "phys%",
      flatIndex: -1,
      percentIndex: StatToIndexMap["PhyP"],
      val: {},
      count: 0,
      t: "%",
    },
    dendro: {
      name: "dendro%",
      flatIndex: -1,
      percentIndex: StatToIndexMap["DendroP"],
      val: {},
      count: 0,
      t: "%",
    },
    heal: {
      name: "heal",
      flatIndex: -1,
      percentIndex: StatToIndexMap["Heal"],
      val: {},
      count: 0,
      t: "%",
    },
  };

  let maxRowCount = 0;

  chars.forEach((char) => {
    let rowCount = 0;
    for (const key in totalStats) {
      const s = totalStats[key];
      if (!(char.name in totalStats[key].val)) {
        totalStats[key].val[char.name] = { flat: 0, per: 0 };
      }
      if (char.stats[s.percentIndex] > 0 || char.stats[s.flatIndex] > 0) {
        totalStats[key].count++;
        rowCount++;
      }
      switch (s.t) {
        case "both":
          totalStats[key].val[char.name].flat = char.stats[s.flatIndex];
          totalStats[key].val[char.name].per = char.stats[s.percentIndex];

          break;
        case "f":
          totalStats[key].val[char.name].flat = char.stats[s.flatIndex];
          break;
        case "%":
          totalStats[key].val[char.name].per = char.stats[s.percentIndex];
          break;
      }
    }
    if (rowCount > maxRowCount) {
      maxRowCount = rowCount;
    }
  });

  const stats: { [key in string]: CharStatBlock[] } = {};

  //make a block for all the chars first
  chars.forEach((c) => {
    stats[c.name] = [];
  });

  for (const key in totalStats) {
    if (totalStats[key].count > 0) {
      //loop through chars
      for (const char in totalStats[key].val) {
        stats[char].push({
          key: key,
          name: totalStats[key].name,
          t: totalStats[key].t,
          flat: totalStats[key].val[char].flat,
          percent: totalStats[key].val[char].per,
        });
      }
    }
  }

  // console.log(stats);

  return { stats: stats, maxRows: maxRowCount };
}
