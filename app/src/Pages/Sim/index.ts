export * from "./Simple";

import { CharDetail } from "/src/Components/Character";

export const charTestConfig: CharDetail[] = [
  {
    name: "bennett",
    element: "pyro",
    level: 70,
    max_level: 80,
    cons: 6,
    weapon: {
      name: "aquilafavonia",
      refine: 1,
      level: 90,
      max_level: 90,
    },
    talents: { attack: 9, skill: 9, burst: 10 },
    sets: { noblesseoblige: 4 },
    stats: [
      0, 0, 0, 5437, 0.146, 391, 0.105, 1.114, 46, 0.318, 0.98, 0, 0.466, 0, 0,
      0, 0, 0, 0, 0, 0, 0,
    ],
    snapshot: [
      0, 0, 0, 5437, 0.146, 391, 0.5549999999999999, 1.3806999995708467, 46,
      0.368, 1.48, 0, 0.466, 0, 0, 0, 0, 0, 0.41346000406980465, 0, 0, 0,
    ],
  },
  {
    name: "ganyu",
    element: "cryo",
    level: 90,
    max_level: 90,
    cons: 6,
    weapon: {
      name: "amosbow",
      refine: 5,
      level: 90,
      max_level: 90,
    },
    talents: { attack: 10, skill: 10, burst: 10 },
    sets: { wandererstroupe: 4 },
    stats: [
      0, 0.19, 0, 5527, 0, 365, 0.513, 0.175, 99, 0.459, 1.337, 0, 0, 0, 0.466,
      0, 0, 0, 0, 0, 0, 0,
    ],
    snapshot: [
      0, 0.19, 0, 5527, 0, 365, 1.2591519980381722, 0.175, 179, 0.509,
      2.22100000333786, 0, 0, 0, 0.466, 0, 0, 0, 0, 0, 0, 0,
    ],
  },
  {
    name: "kazuha",
    element: "anemo",
    level: 90,
    max_level: 90,
    cons: 0,
    weapon: {
      name: "freedomsworn",
      refine: 1,
      level: 90,
      max_level: 90,
    },
    talents: { attack: 9, skill: 9, burst: 9 },
    sets: { viridescentvenerer: 4 },
    stats: [
      0, 0.306, 70, 5557, 0, 393, 0, 0.777, 490, 0.10900000000000001, 0.621, 0,
      0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
    ],
    snapshot: [
      0, 0.306, 70, 5557, 0, 393, 0.25, 0.777, 803.6607945205687,
      0.15900000000000003, 1.121, 0, 0, 0, 0, 0, 0.15, 0, 0, 0, 0, 0.1,
    ],
  },
  {
    name: "xiangling",
    element: "pyro",
    level: 90,
    max_level: 90,
    cons: 6,
    weapon: {
      name: "favoniuslance",
      refine: 5,
      level: 90,
      max_level: 90,
    },
    talents: { attack: 2, skill: 9, burst: 10 },
    sets: { emblemofseveredfate: 4 },
    stats: [
      0, 0.241, 107, 5258, 0, 344, 0.099, 0.628, 102, 0.56, 0.823, 0, 0.466, 0,
      0, 0, 0, 0, 0, 0, 0, 0,
    ],
    snapshot: [
      0, 0.241, 107, 5258, 0, 344, 0.349, 1.1342681795149587, 198,
      0.6100000000000001, 1.323, 0, 0.466, 0, 0, 0, 0, 0, 0, 0, 0, 0,
    ],
  },
];
