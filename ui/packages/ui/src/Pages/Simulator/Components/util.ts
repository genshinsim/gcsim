export function maxLvlToAsc(lvl: number): number {
  switch (lvl) {
    case 90:
      return 6;
    case 80:
      return 5;
    case 70:
      return 4;
    case 60:
      return 3;
    case 50:
      return 2;
    case 40:
      return 1;
    default:
      return 0;
  }
}

export function ascToMaxLvl(asc: number): number {
  switch (asc) {
    case 6:
      return 90;
    case 5:
      return 80;
    case 4:
      return 70;
    case 3:
      return 60;
    case 2:
      return 50;
    case 1:
      return 40;
    default:
      return 20;
  }
}

export function ascLvlMin(asc: number): number {
  switch (asc) {
    case 1:
      return 20;
    case 2:
      return 40;
    case 3:
      return 50;
    case 4:
      return 60;
    case 5:
      return 70;
    case 6:
      return 80;
  }
  return 1;
}

export function ascLvlMax(asc: number): number {
  switch (asc) {
    case 0:
      return 20;
    case 1:
      return 40;
    case 2:
      return 50;
    case 3:
      return 60;
    case 4:
      return 70;
    case 5:
      return 80;
    case 6:
      return 90;
  }
  return 0;
}

export const StatToIndexMap: { [key in string]: number } = {
  DEFP: 1,
  DEF: 2,
  HP: 3,
  HPP: 4,
  ATK: 5,
  ATKP: 6,
  ER: 7,
  EM: 8,
  CR: 9,
  CD: 10,
  Heal: 11,
  PyroP: 12,
  HydroP: 13,
  CryoP: 14,
  ElectroP: 15,
  AnemoP: 16,
  GeoP: 17,
  DendroP: 18,
  PhyP: 19,
};

export const GOODStatToIndexMap: { [key in string]: number } = {
  def_: 1,
  def: 2,
  hp: 3,
  hp_: 4,
  atk: 5,
  atk_: 6,
  enerRech_: 7,
  eleMas: 8,
  critRate_: 9,
  critDMG_: 10,
  heal: 11,
  pyro_dmg_: 12,
  hydro_dmg_: 13,
  cryo_dmg_: 14,
  electro_dmg_: 15,
  anemo_dmg_: 16,
  geo_dmg_: 17,
  dendro_dmg_: 18,
  physical_dmg_: 19,
};

export const DMElementToKey: { [key in string]: string } = {
  Electric: "electro",
  Fire: "pyro",
  Ice: "cryo",
  Water: "hydro",
  Grass: "dendro",
  Wind: "anemo",
  Rock: "geo",
};
