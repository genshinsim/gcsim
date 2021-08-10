import {
  Dialog,
  Classes,
  Button,
  TextArea,
  FormGroup,
  HTMLSelect,
  H4,
} from "@blueprintjs/core";
import React from "react";

function ArtifactBuilder({
  isOpen,
  onClose,
}: {
  isOpen: boolean;
  onClose: () => void;
}) {
  const [goblet, setGoblet] = React.useState<string>("atk%");
  const [sand, setSand] = React.useState<string>("atk%");
  const [circlet, setCirclet] = React.useState<string>("cr");

  let msg = "";

  msg +=
    "stats+=undefined label=main hp=4780 atk=311 " +
    goblet +
    "=" +
    mainStats[statMap[goblet]] +
    " " +
    sand +
    "=" +
    mainStats[statMap[sand]] +
    " " +
    circlet +
    "=" +
    mainStats[statMap[circlet]] +
    ";\n";

  msg +=
    "stats+=undefined label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 core=59 core%=.186;";

  return (
    <Dialog isOpen={isOpen} onClose={onClose} style={{ width: "50%" }}>
      <div>
        <div className={Classes.DIALOG_HEADER}>
          <H4>Generate Artifact Substats (Total)</H4>
        </div>
        <div className={Classes.DIALOG_BODY}>
          <div className="row">
            <div className="col-xs">
              <FormGroup label="Sand">
                <HTMLSelect
                  options={sandOpt}
                  value={sand}
                  fill
                  onChange={(e) => setSand(e.currentTarget.value)}
                />
              </FormGroup>
            </div>
            <div className="col-xs">
              <FormGroup label="Goblet">
                <HTMLSelect
                  options={gobletOpt}
                  value={goblet}
                  fill
                  onChange={(e) => setGoblet(e.currentTarget.value)}
                />
              </FormGroup>
            </div>
            <div className="col-xs">
              <FormGroup label="Circlet">
                <HTMLSelect
                  options={circletOpt}
                  value={circlet}
                  fill
                  onChange={(e) => setCirclet(e.currentTarget.value)}
                />
              </FormGroup>
            </div>
          </div>
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

export default ArtifactBuilder;

const sandOpt = [
  { value: "hp%" },
  { value: "atk%" },
  { value: "def%" },
  { value: "em" },
  { value: "er" },
];
const gobletOpt = [
  { value: "hp%" },

  { value: "atk%" },
  { value: "def%" },
  { value: "em" },
  { value: "phys%" },
  { value: "anemo%" },
  { value: "geo%" },
  { value: "electro%" },
  { value: "hydro%" },
  { value: "pyro%" },
  { value: "cryo%" },
  { value: "heal" },
];
const circletOpt = [
  { value: "hp%" },
  { value: "atk%" },
  { value: "def%" },
  { value: "em" },
  { value: "cr" },
  { value: "cd" },
];

const statMap: { [key: string]: string } = {
  hp: "hp",
  "hp%": "hp_",
  atk: "atk",
  "atk%": "atk_",
  def: "def",
  "def%": "def_",
  em: "eleMas",
  er: "enerRech_",
  cr: "critRate_",
  cd: "critDMG_",
  "phys%": "physical_dmg_",
  "anemo%": "ele_dmg_",
  "geo%": "ele_dmg_",
  "electro%": "ele_dmg_",
  "hydro%": "ele_dmg_",
  "pyro%": "ele_dmg_",
  "cryo%": "ele_dmg_",
  heal: "ele_dmg_",
};

const mainStats: {
  [key: string]: number;
} = {
  hp: 4780,
  atk: 311,
  hp_: 0.466,
  atk_: 0.466,
  def_: 0.583,
  physical_dmg_: 0.583,
  ele_dmg_: 0.466,
  eleMas: 187,
  enerRech_: 0.518,
  critRate_: 0.311,
  critDMG_: 0.622,
  heal_: 0.359,
};
