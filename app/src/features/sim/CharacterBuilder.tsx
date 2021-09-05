import {
  Dialog,
  Classes,
  Button,
  ControlGroup,
  InputGroup,
  TextArea,
  H4,
} from "@blueprintjs/core";
import React from "react";
import genshindb from "genshin-db";

const validLvl = /^(?:\d{1,2})?\+?$/;

const ascStatMap: { [key: string]: string } = {
  albedo: "geo%",
  amber: "atk%",
  barbara: "hp%",
  beidou: "electro%",
  bennett: "er",
  chongyun: "atk%",
  diluc: "cr",
  diona: "cryo%",
  eula: "cd",
  fischl: "atk%",
  ganyu: "cd",
  hutao: "cd",
  jean: "heal",
  kaeya: "er",
  keqing: "cd",
  klee: "pyro%",
  lisa: "em",
  mona: "er",
  ningguang: "geo%",
  noelle: "def%",
  qiqi: "heal",
  razor: "phys%",
  rosaria: "atk%",
  sucrose: "anemo%",
  tartaglia: "hydro%",
  traveler: "atk%",
  venti: "er",
  xiangling: "em",
  xiao: "cr",
  xingqiu: "atk%",
  xinyan: "atk%",
  yanfei: "pyro%",
  zhongli: "geo%",
};

const subMap: { [key: string]: string } = {
  ATK: "atk%",
  "ATK%": "atk%",
  "HP%": "hp%",
  HP: "hp%",
  "DEF%": "def%",
  "Energy Recharge%": "er",
  "Energy Recharge": "er",
  "CRIT Rate%": "cr",
  "CRIT DMG%": "cd",
  "CRIT Rate": "cr",
  "CRIT DMG": "cd",
  "Elemental Mastery": "em",
  "Healing%": "heal",
  Healing: "heal",
  "Physical DMG%": "phys%",
  "Electro DMG%": "electro%",
  "Geo DMG%": "geo%",
  "Pyro DMG%": "pyro%",
  "Cryo DMG%": "cryo%",
  "Hydro DMG%": "hydro%",
  "Anemo DMG%": "anemo%",
  "Physical DMG": "phys%",
  "Electro DMG": "electro%",
  "Geo DMG": "geo%",
  "Pyro DMG": "pyro%",
  "Cryo DMG": "cryo%",
  "Hydro DMG": "hydro%",
  "Anemo DMG": "anemo%",
};

function CharacterBuilder({
  isOpen,
  onClose,
}: {
  isOpen: boolean;
  onClose: () => void;
}) {
  const [char, setChar] = React.useState<string>("");
  const [charLvl, setCharLvl] = React.useState<string>("");
  const [weapon, setWeapon] = React.useState<string>("");
  const [weaponLvl, setWeaponLvl] = React.useState<string>("");
  const [msg, setMsg] = React.useState<string>("");

  const handleFind = () => {
    let c = genshindb.characters(char);
    console.log(c);

    let result = "";
    let cname = "";

    if (c) {
      if (!validLvl.test(charLvl)) {
        setMsg("invalid char lvl");
        return;
      }
      cname = c.name.toLowerCase();

      let stat: genshindb.StatResult;

      if (/^\d{1,2}\+$/.test(charLvl)) {
        const lvl = parseInt(charLvl.slice(0, -1));
        stat = c.stats(lvl, "+");
      } else {
        const lvl = parseInt(charLvl);
        stat = c.stats(lvl);
      }

      result +=
        "char+=" +
        cname +
        " ele=" +
        c.element.toLowerCase() +
        " lvl=" +
        stat.level +
        " hp=" +
        stat.hp?.toFixed(3) +
        " atk=" +
        stat.attack?.toFixed(3) +
        " def=" +
        stat.defense?.toFixed(3) +
        " " +
        ascStatMap[cname] +
        "=" +
        stat.specialized?.toFixed(3) +
        " cr=.05 cd=0.5 cons=0 talent=1,1,1;\n";

      console.log(JSON.stringify(stat, null, 2));
    } else {
      result = "character not found\n";
    }

    let w = genshindb.weapons(weapon);
    if (w) {
      if (!validLvl.test(weaponLvl)) {
        result += "invalid weapon lvl";
        setMsg(result);
        return;
      }

      console.log(w);

      let stat: genshindb.StatResult;

      if (/^\d{1,2}\+$/.test(weaponLvl)) {
        const lvl = parseInt(weaponLvl.slice(0, -1));
        stat = w.stats(lvl, "+");
      } else {
        const lvl = parseInt(weaponLvl);
        stat = w.stats(lvl);
      }

      result +=
        "weapon+=" +
        cname +
        ' label="' +
        w.name.toLowerCase() +
        '"' +
        " atk=" +
        stat.attack?.toFixed(3) +
        " " +
        subMap[w.substat] +
        "=" +
        stat.specialized?.toFixed(3) +
        " refine=1;\n";

      console.log(JSON.stringify(stat, null, 2));
    } else {
      result += "weapon not found";
    }

    setMsg(result);
  };

  return (
    <Dialog isOpen={isOpen} onClose={onClose} style={{ width: "50%" }}>
      <div>
        <div className={Classes.DIALOG_HEADER}>
          <H4>Generate Character Config</H4>
        </div>
        <div className={Classes.DIALOG_BODY}>
          <ControlGroup fill vertical={false}>
            <InputGroup
              placeholder="Character"
              value={char}
              onChange={(e) => setChar(e.target.value)}
            />
            <InputGroup
              placeholder="Character level (80+ for asc)"
              value={charLvl}
              intent={validLvl.test(charLvl) ? "none" : "danger"}
              onChange={(e) => setCharLvl(e.target.value)}
            />
            <InputGroup
              placeholder="Weapon"
              value={weapon}
              onChange={(e) => setWeapon(e.target.value)}
            />
            <InputGroup
              placeholder="Weapon level (80+ for asc)"
              value={weaponLvl}
              intent={validLvl.test(charLvl) ? "none" : "danger"}
              onChange={(e) => setWeaponLvl(e.target.value)}
            />
            <Button icon="arrow-right" onClick={handleFind} />
          </ControlGroup>
          <br />
          {msg !== "" ? <TextArea value={msg} fill rows={10}></TextArea> : null}
        </div>
        <div className={Classes.DIALOG_FOOTER}>
          <div className={Classes.DIALOG_FOOTER_ACTIONS}>
            <Button
              onClick={() => {
                navigator.clipboard.writeText(msg);
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

export default CharacterBuilder;
