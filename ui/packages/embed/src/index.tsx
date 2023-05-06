import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App";

import "@blueprintjs/core/lib/css/blueprint.css";
import "./index.css";

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);

if (import.meta.env.DEV) {
  data = JSON.stringify({
    schema_version: { major: 4, minor: 0 },
    sim_version: "134360fcbd3cece77bfd9280b1c133544546b31c",
    build_date: "2022-11-16T14:35:06Z",
    modified: false,
    initial_character: "rosaria",
    key_type: "dev",
    character_details: [
      {
        name: "zhongli",
        element: "geo",
        level: 90,
        max_level: 90,
        cons: 0,
        weapon: { name: "favoniuslance", refine: 3, level: 90, max_level: 90 },
        talents: { attack: 9, skill: 9, burst: 9 },
        sets: { archaicpetra: 2, noblesseoblige:2 },
        stats: [
          0, 0.124, 39.36, 5287.88, 0.0992, 344.08, 0.6644, 0.1102, 39.64,
          0.642, 0.7944, 0, 0, 0, 0, 0, 0, 0.466, 0, 0, 0, 0,
        ],
        snapshot: [
          0, 0, 868.6609026633635, 21440.72679049085, 0, 1702.098977990197, 0,
          0.4164681795149587, 39.64, 0.6920000000000001, 1.2944, 0, 0, 0, 0, 0,
          0, 0.9039999876022339, 0, 0, 0, 0,
        ],
      },
      {
        name: "rosaria",
        element: "cryo",
        level: 90,
        max_level: 90,
        cons: 6,
        weapon: { name: "deathmatch", refine: 1, level: 90, max_level: 90 },
        talents: { attack: 9, skill: 9, burst: 9 },
        sets: { noblesseoblige: 2 },
        stats: [
          0, 0.124, 39.36, 5287.88, 0.0992, 344.08, 0.5652, 0.2204, 39.64,
          0.3972, 1.284, 0, 0, 0, 0.466, 0, 0, 0, 0, 0, 0, 0,
        ],
        snapshot: [
          0, 0, 837.2003625042938, 18795.568527226642, 0, 1764.2095791084657, 0,
          0.2204, 39.64, 0.81471998079896, 1.784, 0, 0, 0, 0.466, 0, 0, 0, 0, 0,
          0, 0,
        ],
      },
      {
        name: "albedo",
        element: "geo",
        level: 90,
        max_level: 90,
        cons: 0,
        weapon: {
          name: "cinnabarspindle",
          refine: 5,
          level: 90,
          max_level: 90,
        },
        talents: { attack: 9, skill: 9, burst: 9 },
        sets: { huskofopulentdreams: 1 },
        stats: [
          0, 0.831, 39.36, 5287.88, 0.0992, 344.08, 0.0992, 0.1102, 39.64,
          0.642, 0.7944, 0, 0, 0, 0, 0, 0, 0.466, 0, 0, 0, 0,
        ],
        snapshot: [
          0, 0, 2510.7315418868984, 19825.44170599099, 0, 1119.5661270387823, 0,
          0.1102, 39.64, 0.6920000000000001, 1.2944, 0, 0, 0, 0, 0, 0,
          0.7539999876022339, 0, 0, 0, 0,
        ],
      },
      {
        name: "ganyu",
        element: "cryo",
        level: 90,
        max_level: 90,
        cons: 0,
        weapon: {
          name: "prototypecrescent",
          refine: 5,
          level: 90,
          max_level: 90,
        },
        talents: { attack: 9, skill: 9, burst: 9 },
        sets: { blizzardstrayer: 4 },
        stats: [
          0, 0.124, 39.36, 5287.88, 0.0992, 344.08, 0.5652, 0.2204, 39.64,
          0.3972, 1.284, 0, 0, 0, 0.466, 0, 0, 0, 0, 0, 0, 0,
        ],
        snapshot: [
          0, 0, 747.7212110718447, 16056.444526993902, 0, 2014.9707397558066, 0,
          0.2204, 39.64, 0.4472, 2.16800000333786, 0, 0, 0, 0.616, 0, 0, 0, 0,
          0, 0, 0,
        ],
      },
    ],
    target_details: [
      {
        level: 100,
        hp: 100000000,
        resist: {
          anemo: 0.1,
          cryo: 0.1,
          dendro: 0.1,
          electro: 0.1,
          geo: 0.1,
          hydro: 0.1,
          physical: 0.1,
          pyro: 0.1,
        },
        particle_drop_threshold: 0,
        particle_drop_count: 0,
        particle_element: "electro",
      },
    ],
    simulator_settings: {
      damage_mode: true,
      enable_hitlag: true,
      def_halt: true,
      iterations: 1000,
      delays: {
        skill: 0,
        burst: 0,
        attack: 0,
        charge: 0,
        aim: 0,
        dash: 0,
        jump: 0,
        swap: 12,
      },
    },
    energy_settings: {
      active: true,
      once: false,
      start: 480,
      end: 720,
      amount: 1,
      last_energy_drop: 0,
    },
    config_file:
      'zhongli char lvl=90/90 cons=0 talent=9,9,9;\nzhongli add weapon="favoniuslance" refine=3 lvl=90/90;\nzhongli add set="archaicpetra" count=4;\nzhongli add stats hp=4780 atk=311 atk%=0.466 geo%=0.466 cr=0.311;\nzhongli add stats def=39.36 def%=0.124 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.1984 er=0.1102 em=39.64 cr=0.331 cd=0.7944 ;\t\t\n\t\t\t\t\t\t\t\t\t\t\t\nrosaria char lvl=90/90 cons=6 talent=9,9,9; \nrosaria add weapon="deathmatch" refine=1 lvl=90/90;\nrosaria add set="no" count=5;\nrosaria add stats hp=4780 atk=311 atk%=0.466 cryo%=0.466 cd=0.622 ; #main\nrosaria add stats def=39.36 def%=0.124 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.0992 er=0.2204 em=39.64 cr=0.3972 cd=0.662 ;\t\n\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\nalbedo char lvl=90/90 cons=0 talent=9,9,9;\nalbedo add weapon="cinnabarspindle" refine=5 lvl=90/90;\nalbedo add set="huskofopulentdreams" count=4;\nalbedo add stats hp=4780 atk=311 def%=0.583 geo%=0.466 cr=0.311;\nalbedo add stats def%=0.248 def=39.36 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.0992 er=0.1102 em=39.64 cr=0.331 cd=0.7944;\t\t\n\nganyu char lvl=90/90 cons=0 talent=9,9,9;\nganyu add weapon="prototypecrescent" refine=5 lvl=90/90;\nganyu add set="blizzardstrayer" count=4;\nganyu add stats hp=4780 atk=311 atk%=0.466 cryo%=0.466 cd=0.622;\nganyu add stats def=39.36 def%=0.124 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.0992 er=0.2204 em=39.64 cr=0.3972 cd=0.662 ;\t\t\t\n\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\n\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\noptions swap_delay=12 debug=true iteration=1000 workers=50 mode=sl;\ntarget lvl=100 resist=0.1 hp=100000000 pos=0,0;\nenergy every interval=480,720 amount=1;\n\nactive rosaria;\nfor let i = 0; i < 5; i = i + 1 {\n  rosaria skill, burst;\n  zhongli skill[hold=1], dash;\n  albedo skill, attack;\n  ganyu aim[weakspot=1], skill, burst[radius=2], aim[weakspot=1];\n  zhongli attack, burst;\n  rosaria skill;\n  ganyu aim[weakspot=1], skill, aim, aim[weakspot=1];\n}\n\n\n\n',
    sample_seed: "9321530717894557214",
    statistics: {
      min_seed: "17511453863270297604",
      max_seed: "16590787742054998978",
      p25_seed: "10342512838490931276",
      p50_seed: "6256492013046413457",
      p75_seed: "6915323472330700161",
      runtime: -1668644170330000100,
      iterations: 1000,
      duration: {
        min: 109.15,
        max: 117.36666666666666,
        mean: 109.6288166666667,
        sd: 0.8699182411913554,
        q1: 109.15,
        q2: 109.15,
        q3: 110.35,
        histogram: [
          689, 7, 8, 3, 232, 3, 7, 0, 19, 0, 3, 1, 23, 0, 2, 0, 1, 1, 0, 0, 0,
          0, 0, 0, 0, 0, 0, 1,
        ],
      },
      dps: {
        min: 41758.01777098939,
        max: 48304.68597091312,
        mean: 45699.118762647755,
        sd: 883.7075217467534,
        q1: 45175.17127470411,
        q2: 45770.57592006592,
        q3: 46283.13842074638,
        histogram: [
          1, 0, 2, 3, 7, 5, 10, 19, 30, 57, 66, 109, 121, 152, 140, 114, 78, 43,
          26, 11, 5, 1,
        ],
      },
      rps: {
        min: 0.5970149253731344,
        max: 0.6962895098488319,
        mean: 0.6626574918041097,
        sd: 0.02012836623446547,
        q1: 0.6504809894640403,
        q2: 0.668804397617957,
        q3: 0.6779661016949152,
        histogram: [
          5, 12, 9, 7, 16, 84, 33, 43, 155, 116, 141, 57, 107, 209, 6,
        ],
      },
      eps: {
        min: 9.745338725209981,
        max: 12.039390083667184,
        mean: 11.09109144548229,
        sd: 0.3460660559942403,
        q1: 10.871434367546094,
        q2: 11.105713450797372,
        q3: 11.341981141144373,
        histogram: [
          1, 0, 2, 10, 14, 25, 22, 57, 84, 112, 135, 145, 120, 113, 76, 45, 31,
          5, 3,
        ],
      },
      hps: {
        min: 0,
        max: 0,
        mean: 0,
        sd: 0,
        q1: 0,
        q2: 0,
        q3: 0,
        histogram: [1000],
      },
      sps: {
        min: 29460.188987698704,
        max: 32858.854869579955,
        mean: 31707.304444102556,
        sd: 695.35085208854,
        q1: 31306.35217439808,
        q2: 31903.941873400454,
        q3: 32254.69524259954,
        histogram: [
          5, 14, 6, 12, 16, 92, 39, 34, 148, 117, 53, 147, 102, 209, 6,
        ],
      },
      total_damage: {
        min: 4707510.132134999,
        max: 5272456.473725167,
        mean: 5009510.368807329,
        sd: 81751.46648573084,
        q1: 0,
        q2: 0,
        q3: 0,
      },
      warnings: {
        target_overlap: false,
        insufficient_energy: false,
        insufficient_stamina: false,
        swap_cd: false,
        skill_cd: false,
      },
      failed_actions: [
        {
          insufficient_energy: {
            min: 0,
            max: 3.066666666666667,
            mean: 0.02779999999999998,
            sd: 0.2052907798813778,
          },
          insufficient_stamina: { min: 0, max: 0, mean: 0, sd: 0 },
          swap_cd: { min: 0, max: 0, mean: 0, sd: 0 },
          skill_cd: { min: 0, max: 0, mean: 0, sd: 0 },
        },
        {
          insufficient_energy: {
            min: 0,
            max: 3.95,
            mean: 0.016616666666666637,
            sd: 0.16584129210634133,
          },
          insufficient_stamina: { min: 0, max: 0, mean: 0, sd: 0 },
          swap_cd: { min: 0, max: 0, mean: 0, sd: 0 },
          skill_cd: { min: 0, max: 0, mean: 0, sd: 0 },
        },
        {
          insufficient_energy: { min: 0, max: 0, mean: 0, sd: 0 },
          insufficient_stamina: { min: 0, max: 0, mean: 0, sd: 0 },
          swap_cd: {
            min: 0.8333333333333333,
            max: 0.8333333333333333,
            mean: 0.8333333333333333,
            sd: 0,
          },
          skill_cd: { min: 0, max: 0, mean: 0, sd: 0 },
        },
        {
          insufficient_energy: {
            min: 0,
            max: 8.216666666666667,
            mean: 0.4343999999999999,
            sd: 0.8157732619446366,
          },
          insufficient_stamina: { min: 0, max: 0, mean: 0, sd: 0 },
          swap_cd: { min: 0, max: 0, mean: 0, sd: 0 },
          skill_cd: { min: 0, max: 0, mean: 0, sd: 0 },
        },
      ],
    },
  });
}

interface HasErr{ 
  err: string
}

//TODO(kyle): parsed should be typed model.ISimulationResult | HasErr
export const parsed: any | HasErr = JSON.parse(data);
