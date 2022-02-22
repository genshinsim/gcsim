import { StatKey } from "./goodTypes";

export const StatKeyToStatString: { [key in StatKey]: string } = {
  hp: "HP",
  hp_: "HPP",
  atk: "ATK",
  atk_: "ATKP",
  def: "DEF",
  def_: "DEFP",
  eleMas: "EM",
  enerRech_: "ER",
  heal_: "Heal",
  critRate_: "CR",
  critDMG_: "CD",
  physical_dmg_: "PhyP",
  anemo_dmg_: "AnemoP",
  geo_dmg_: "GeoP",
  electro_dmg_: "ElectroP",
  hydro_dmg_: "HydroP",
  pyro_dmg_: "PyroP",
  cryo_dmg_: "CryoP",
};
