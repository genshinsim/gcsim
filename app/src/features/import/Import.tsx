import {
  Dialog,
  Classes,
  Button,
  Callout,
  FileInput,
  TextArea,
  HTMLSelect,
  ControlGroup,
  InputGroup,
  H4,
} from "@blueprintjs/core";
import { RootState } from "app/store";
import React from "react";
import { useDispatch, useSelector } from "react-redux";
import { Artifact, DataSet, setData } from "./importSlice";

const keys = ["feather", "flower", "sand", "cup", "head"];

function Import({ isOpen, onClose }: { isOpen: boolean; onClose: () => void }) {
  const { data } = useSelector((state: RootState) => {
    return {
      data: state.import.data,
    };
  });
  const dispatch = useDispatch();

  const [file, setFile] = React.useState<File | null>(null);
  const [opt, setOpt] = React.useState<string>("Cocogoat");
  const [content, setContent] = React.useState<string>("");
  const [filter, setFilter] = React.useState<string>("");

  const handleParseData = () => {
    if (opt === "Cocogoat") {
      handleParseCocogoat();
    } else {
      handleParseGO();
    }
  };

  const handleParseGO = () => {
    let r = JSON.parse(content);

    if ("artifactDatabase" in r) {
      let result: DataSet = {
        feather: [],
        flower: [],
        sand: [],
        cup: [],
        head: [],
      };

      let db = r["artifactDatabase"];
      for (const [key] of Object.entries(db)) {
        let art = db[key];
        console.log(art);

        let mainKey = art.mainStatKey;
        if (mainKey.endsWith("_dmg_") && mainKey !== "physical_dmg_") {
          mainKey = "ele_dmg_";
        }

        //if value ends with _ then it's a percentage...
        let val = mainStatMap[art.numStars][mainKey][art.level];

        if (mainKey.endsWith("_")) {
          val = val / 100;
        }

        let slot = slotMapGO[art.slotKey];
        let ele: Artifact = {
          position: slot, // @ts-ignore
          setName: artMapGO[art.setKey],
          mainTag: {
            name: statMapGO[art.mainStatKey],
            value: val,
          },
          normalTags: [],
          comments:
            (art.location ? "," + art.location : "") +
            (art.id ? "," + art.id : ""),
        };
        for (let i = 0; i < art.substats.length; i++) {
          let subval = art.substats[i].value;
          if (art.substats[i].key.endsWith("_")) {
            subval = subval / 100;
          }

          ele.normalTags.push({
            name: statMapGO[art.substats[i].key],
            value: subval,
          });
        }

        result[slot].push(ele);
      }

      dispatch(setData(result));
    }
  };

  const handleParseCocogoat = () => {
    //data should be json string
    let r = JSON.parse(content);

    let result: DataSet = {
      feather: [],
      flower: [],
      sand: [],
      cup: [],
      head: [],
    };

    keys.forEach((k) => {
      if (k in r) {
        // console.log(k, r[k]);
        if (Array.isArray(r[k])) {
          for (let i = 0; i < r[k].length; i++) {
            let ele: Artifact = {
              position: k,
              setName: artMap[r[k][i].setName],
              mainTag: {
                name: statMap[r[k][i].mainTag.name],
                value: r[k][i].mainTag.value,
              },
              normalTags: [],
            };

            for (let j = 0; j < r[k][i].normalTags.length; j++) {
              ele.normalTags.push({
                name: statMap[r[k][i].normalTags[j].name],
                value: r[k][i].normalTags[j].value,
              });
            }

            result[k].push(ele);
          }
        }
      }
    });

    dispatch(setData(result));

    //loop through flower, feather, sand, cup, head

    // console.log(r);
  };

  const handleSelectFile = (event: React.ChangeEvent<HTMLInputElement>) => {
    // console.log(event);
    // console.log(event.target.files);
    var files = event.target.files;
    if (!files) return;
    if (files.length === 0) return;
    setFile(files[0]);

    //read the file
    const reader = new FileReader();
    reader.onload = (ev) => {
      // console.log(ev.target?.result);
      if (ev.target === null) return;
      var data = ev.target.result;
      if (data === null) {
        console.log("ERROR: no data from file");

        return;
      }
      if (data instanceof ArrayBuffer) {
        console.log("ERROR: data is type ArrayBuffer???");

        return;
      }

      setContent(data);
    };

    reader.readAsText(files[0]);
  };

  let result = "";
  //build result
  if (data !== null) {
    keys.forEach((k) => {
      /**
       * Expected output
       *
       *    #feather
       *    stats label=feather hp=??; #setname
       */
      data[k].forEach((art) => {
        // console.log(art);

        let line = "stats+=undefined ";
        line += "label=" + k + " ";
        line += art.mainTag.name + "=" + art.mainTag.value.toFixed(3);
        art.normalTags.forEach((sub) => {
          line += " " + sub.name + "=" + sub.value.toFixed(3);
        });
        line +=
          ";  #" + art.setName + (art.comments ? art.comments : "") + "\n";

        //search line
        if (line.includes(filter)) {
          result += line;
        }
      });
    });
  }

  return (
    <Dialog isOpen={isOpen} onClose={onClose} style={{ width: "80%" }}>
      <div>
        <div className={Classes.DIALOG_HEADER}>
          <H4>Import from GO/Cocogoat</H4>
        </div>
        <div className={Classes.DIALOG_BODY}>
          <Callout intent="primary">
            Paste your Genshin Optimizer/Cocogoat export below or use the file
            button to upload.
          </Callout>

          <TextArea
            fill
            rows={7}
            value={content}
            onChange={(e) => setContent(e.target.value)}
          />

          <FileInput
            fill
            text={file === null ? "Choose file.." : file.name}
            onInputChange={handleSelectFile}
          />
          <ControlGroup fill>
            <HTMLSelect
              onChange={(e) => {
                setOpt(e.currentTarget.value);
              }}
              value={opt}
            >
              <option label="Genshin Optimizer" value="GO" />
              <option label="Cocogoat" value="Cocogoat" />
            </HTMLSelect>
            <Button intent="primary" onClick={handleParseData}>
              Parse
            </Button>
          </ControlGroup>
          <br />
          {data ? (
            <div>
              <ControlGroup fill>
                <InputGroup
                  placeholder="Type to filter..."
                  value={filter}
                  onChange={(e) => setFilter(e.target.value)}
                />
              </ControlGroup>
              <pre>{result}</pre>
            </div>
          ) : null}
        </div>

        <div className={Classes.DIALOG_FOOTER}>
          <div className={Classes.DIALOG_FOOTER_ACTIONS}>
            <Button
              onClick={() => {
                navigator.clipboard.writeText(result);
              }}
              intent="primary"
            >
              Copy
            </Button>
            <Button onClick={onClose}>Close</Button>
          </div>
        </div>
      </div>
    </Dialog>
  );
}

export default Import;

const statMap: { [key: string]: string } = {
  critical: "cr",
  criticalDamage: "cd",
  attackStatic: "atk",
  attackPercentage: "atk%",
  elementalMastery: "em",
  recharge: "er",
  lifeStatic: "hp",
  lifePercentage: "hp%",
  defendStatic: "def",
  defendPercentage: "def%",
  physicalBonus: "phys%",
  cureEffect: "heal",
  rockBonus: "geo%",
  windBonus: "anemo%",
  iceBonus: "cryo%",
  waterBonus: "hydro%",
  fireBonus: "pyro%",
  thunderBonus: "electro%",
};

const artMap: { [key: string]: string } = {
  archaicPetra: "archaic petra",
  blizzardStrayer: "blizzard strayer",
  bloodstainedChivalry: "bloodstained chivalry",
  crimsonWitch: "crimson witch of flames",
  gladiatorFinale: "gladiator's finale",
  heartOfDepth: "heart of depth",
  lavaWalker: "lavawalker",
  maidenBeloved: "",
  noblesseOblige: "noblesse oblige",
  retracingBolide: "retracing bolide",
  thunderSmoother: "thundersoother",
  thunderingFury: "thundering fury",
  viridescentVenerer: "viridescent venerer",
  wandererTroupe: "wanderer's troupe",
  berserker: "",
  braveHeart: "",
  defenderWill: "",
  exile: "",
  gambler: "",
  instructor: "",
  martialArtist: "",
  prayersForDestiny: "",
  prayersForIllumination: "",
  prayersForWisdom: "",
  prayersToSpringtime: "",
  resolutionOfSojourner: "",
  scholar: "",
  tinyMiracle: "",
  adventurer: "",
  luckyDog: "",
  travelingDoctor: "",
  tenacityOfTheMillelith: "",
  paleFlame: "pale flame",
};

const statMapGO: { [key: string]: string } = {
  hp: "hp",
  hp_: "hp%",
  atk: "atk",
  atk_: "atk%",
  def: "def",
  def_: "def%",
  eleMas: "em",
  enerRech_: "er",
  critRate_: "cr",
  critDMG_: "cd",
  physical_dmg_: "phys%",
  anemo_dmg_: "anemo%",
  geo_dmg_: "geo%",
  electro_dmg_: "electro%",
  hydro_dmg_: "hydro%",
  pyro_dmg_: "pyro%",
  cryo_dmg_: "cryo%",
  heal_: "heal",
};

const artMapGO: { [key: string]: string } = {
  Adventurer: "",
  ArchaicPetra: "archaic petra",
  Berserker: "",
  BlizzardStrayer: "blizzard strayer",
  BloodstainedChivalry: "bloodstained chivalry",
  BraveHeart: "",
  CrimsonWitchOfFlames: "crimson witch of flames",
  DefendersWill: "",
  Gambler: "",
  GladiatorsFinale: "gladiator's finale",
  HeartOfDepth: "heart of depth",
  Instructor: "",
  Lavawalker: "lavawalker",
  LuckyDog: "",
  MaidenBeloved: "",
  MartialArtist: "",
  NoblesseOblige: "noblesse oblige",
  PrayersForDestiny: "",
  PrayersForIllumination: "",
  PrayersForWisdom: "",
  PrayersToSpringtime: "",
  ResolutionOfSojourner: "",
  RetracingBolide: "retracing bolide",
  Scholar: "",
  TheExile: "",
  ThunderingFury: "thundering fury",
  Thundersoother: "thundersoother",
  TinyMiracle: "",
  TravelingDoctor: "",
  ViridescentVenerer: "viridescent venerer",
  WanderersTroupe: "wanderer's troupe",
  TenacityOfTheMillelith: "",
  PaleFlame: "pale flame",
};

const slotMapGO: { [key: string]: string } = {
  flower: "flower",
  plume: "feather",
  sands: "sand",
  goblet: "cup",
  circlet: "head",
};

const mainStatMap: {
  [key: number]: {
    [key: string]: number[];
  };
} = {
  1: {
    hp: [129, 178, 227, 275, 324],
    atk: [8, 12, 15, 18, 21],
    hp_: [3.1, 4.3, 5.5, 6.7, 7.9],
    atk_: [3.1, 4.3, 5.5, 6.7, 7.9],
    def_: [3.9, 5.4, 6.9, 8.4, 9.9],
    physical_dmg_: [3.9, 5.4, 6.9, 8.4, 9.9],
    ele_dmg_: [3.1, 4.3, 5.5, 6.7, 7.9],
    eleMas: [13, 17, 22, 27, 32],
    enerRech_: [3.5, 4.8, 6.1, 7.5, 8.8],
    critRate_: [2.1, 2.9, 3.7, 4.5, 5.3],
    critDMG_: [4.2, 5.8, 7.4, 9.0, 10.5],
    heal_: [2.4, 3.3, 4.3, 5.2, 6.1],
  },
  2: {
    hp: [258, 331, 404, 478, 551, 624, 697, 770, 843],
    atk: [17, 22, 26, 31, 36, 41, 45, 50, 55],
    hp_: [4.2, 5.4, 6.6, 7.8, 9, 10.1, 11.3, 12.5, 13.7],
    atk_: [4.2, 5.4, 6.6, 7.8, 9, 10.1, 11.3, 12.5, 13.7],
    def_: [5.2, 6.7, 8.2, 9.7, 11.2, 12.7, 14.2, 15.6, 17.1],
    physical_dmg_: [5.2, 6.7, 8.2, 9.7, 11.2, 12.7, 14.2, 15.6, 17.1],
    ele_dmg_: [4.2, 5.4, 6.6, 7.8, 9, 10.1, 11.3, 12.5, 13.7],
    eleMas: [17, 22, 26, 31, 36, 41, 45, 50, 55],
    enerRech_: [4.7, 6, 7.3, 8.6, 9.9, 11.3, 12.6, 13.9, 15.2],
    critRate_: [2.8, 3.6, 4.4, 5.2, 6, 6.8, 7.6, 8.3, 9.1],
    critDMG_: [5.6, 7.2, 8.8, 10.4, 11.9, 13.5, 15.1, 16.7, 18.3],
    heal_: [3.2, 4.1, 5.1, 6, 6.9, 7.8, 8.7, 9.6, 10.5],
  },
  3: {
    hp: [
      430, 552, 674, 796, 918, 1040, 1162, 1283, 1405, 1527, 1649, 1771, 1893,
    ],
    atk: [28, 36, 44, 52, 60, 68, 76, 84, 91, 99, 107, 115, 123],
    hp_: [
      5.2, 6.7, 8.2, 9.7, 11.2, 12.7, 14.2, 15.6, 17.1, 18.6, 20.1, 21.6, 23.1,
    ],
    atk_: [
      5.2, 6.7, 8.2, 9.7, 11.2, 12.7, 14.2, 15.6, 17.1, 18.6, 20.1, 21.6, 23.1,
    ],
    def_: [
      6.6, 8.4, 10.3, 12.1, 14.0, 15.8, 17.7, 19.6, 21.4, 23.3, 25.1, 27.0,
      28.8,
    ],
    physical_dmg_: [
      6.6, 8.4, 10.3, 12.1, 14.0, 15.8, 17.7, 19.6, 21.4, 23.3, 25.1, 27.0,
      28.8,
    ],
    ele_dmg_: [
      5.2, 6.7, 8.2, 9.7, 11.2, 12.7, 14.2, 15.6, 17.1, 18.6, 20.1, 21.6, 23.1,
    ],
    eleMas: [21, 27, 33, 39, 45, 51, 57, 63, 69, 75, 80, 86, 92],
    enerRech_: [
      5.8, 7.5, 9.1, 10.8, 12.4, 14.1, 15.7, 17.4, 19.0, 20.7, 22.3, 24.0, 25.6,
    ],
    critRate_: [
      3.5, 4.5, 5.5, 6.5, 7.5, 8.4, 9.4, 10.4, 11.4, 12.4, 13.4, 14.4, 15.4,
    ],
    critDMG_: [
      7.0, 9.0, 11.0, 12.9, 14.9, 16.9, 18.9, 20.9, 22.8, 24.8, 26.8, 28.8,
      30.8,
    ],
    heal_: [
      4.0, 5.2, 6.3, 7.5, 8.6, 9.8, 10.9, 12.0, 13.2, 14.3, 15.5, 16.6, 17.8,
    ],
  },
  4: {
    hp: [
      645, 828, 1011, 1194, 1377, 1559, 1742, 1925, 2108, 2291, 2474, 2657,
      2839, 3022, 3205, 3388, 3571,
    ],
    atk: [
      42, 54, 66, 78, 90, 102, 113, 125, 137, 149, 161, 173, 185, 197, 209, 221,
      232,
    ],
    hp_: [
      6.3, 8.1, 9.9, 11.6, 13.4, 15.2, 17.0, 18.8, 20.6, 22.3, 24.1, 25.9, 27.7,
      29.5, 31.3, 33.0, 34.8,
    ],
    atk_: [
      6.3, 8.1, 9.9, 11.6, 13.4, 15.2, 17.0, 18.8, 20.6, 22.3, 24.1, 25.9, 27.7,
      29.5, 31.3, 33.0, 34.8,
    ],
    def_: [
      7.9, 10.1, 12.3, 14.6, 16.8, 19.0, 21.2, 23.5, 25.7, 27.9, 30.2, 32.4,
      34.6, 36.8, 39.1, 41.3, 43.5,
    ],
    physical_dmg_: [
      7.9, 10.1, 12.3, 14.6, 16.8, 19.0, 21.2, 23.5, 25.7, 27.9, 30.2, 32.4,
      34.6, 36.8, 39.1, 41.3, 43.5,
    ],
    ele_dmg_: [
      6.3, 8.1, 9.9, 11.6, 13.4, 15.2, 17.0, 18.8, 20.6, 22.3, 24.1, 25.9, 27.7,
      29.5, 31.3, 33.0, 34.8,
    ],
    eleMas: [
      25, 32, 39, 47, 54, 61, 68, 75, 82, 89, 97, 104, 111, 118, 125, 132, 139,
    ],
    enerRech_: [
      7.0, 9.0, 11.0, 12.9, 14.9, 16.9, 18.9, 20.9, 22.8, 24.8, 26.8, 28.8,
      30.8, 32.8, 34.7, 36.7, 38.7,
    ],
    critRate_: [
      4.2, 5.4, 6.6, 7.8, 9.0, 10.1, 11.3, 12.5, 13.7, 14.9, 16.1, 17.3, 18.5,
      19.7, 20.8, 22.0, 23.2,
    ],
    critDMG_: [
      8.4, 10.8, 13.1, 15.5, 17.9, 20.3, 22.7, 25.0, 27.4, 29.8, 32.2, 34.5,
      36.9, 39.3, 41.7, 44.1, 46.4,
    ],
    heal_: [
      4.8, 6.2, 7.6, 9.0, 10.3, 11.7, 13.1, 14.4, 15.8, 17.2, 18.6, 19.9, 21.3,
      22.7, 24.0, 25.4, 26.8,
    ],
  },
  5: {
    hp: [
      717, 920, 1123, 1326, 1530, 1733, 1936, 2139, 2342, 2545, 2749, 2952,
      3155, 3358, 3561, 3764, 3967, 4171, 4374, 4577, 4780,
    ],
    atk: [
      47, 60, 73, 86, 100, 113, 126, 139, 152, 166, 179, 192, 205, 219, 232,
      245, 258, 272, 285, 298, 311,
    ],
    hp_: [
      7.0, 9.0, 11.0, 12.9, 14.9, 16.9, 18.9, 20.9, 22.8, 24.8, 26.8, 28.8,
      30.8, 32.8, 34.7, 36.7, 38.7, 40.7, 42.7, 44.6, 46.6,
    ],
    atk_: [
      7.0, 9.0, 11.0, 12.9, 14.9, 16.9, 18.9, 20.9, 22.8, 24.8, 26.8, 28.8,
      30.8, 32.8, 34.7, 36.7, 38.7, 40.7, 42.7, 44.6, 46.6,
    ],
    def_: [
      8.7, 11.2, 13.7, 16.2, 18.6, 21.1, 23.6, 26.1, 28.6, 31, 33.5, 36, 38.5,
      40.9, 43.4, 45.9, 48.4, 50.8, 53.3, 55.8, 58.3,
    ],
    physical_dmg_: [
      8.7, 11.2, 13.7, 16.2, 18.6, 21.1, 23.6, 26.1, 28.6, 31, 33.5, 36, 38.5,
      40.9, 43.4, 45.9, 48.4, 50.8, 53.3, 55.8, 58.3,
    ],
    ele_dmg_: [
      7.0, 9.0, 11.0, 12.9, 14.9, 16.9, 18.9, 20.9, 22.8, 24.8, 26.8, 28.8,
      30.8, 32.8, 34.7, 36.7, 38.7, 40.7, 42.7, 44.6, 46.6,
    ],
    eleMas: [
      28, 36, 44, 52, 60, 68, 76, 84, 91, 99, 107, 115, 123, 131, 139, 147, 155,
      163, 171, 179, 187,
    ],
    enerRech_: [
      7.8, 10.0, 12.2, 14.4, 16.6, 18.8, 21.0, 23.2, 25.4, 27.6, 29.8, 32.0,
      34.2, 36.4, 38.6, 40.8, 43.0, 45.2, 47.4, 49.6, 51.8,
    ],
    critRate_: [
      4.7, 6.0, 7.3, 8.6, 9.9, 11.3, 12.6, 13.9, 15.2, 16.6, 17.9, 19.2, 20.5,
      21.8, 23.2, 24.5, 25.8, 27.1, 28.4, 29.8, 31.1,
    ],
    critDMG_: [
      9.3, 12.0, 14.6, 17.3, 19.9, 22.5, 25.2, 27.8, 30.5, 33.1, 35.7, 38.4,
      41.0, 43.7, 46.3, 49.0, 51.6, 54.2, 56.9, 59.5, 62.2,
    ],
    heal_: [
      5.4, 6.9, 8.4, 10.0, 11.5, 13.0, 14.5, 16.1, 17.6, 19.1, 20.6, 22.1, 23.7,
      25.2, 26.7, 28.2, 29.8, 31.3, 32.8, 34.3, 35.9,
    ],
  },
};
