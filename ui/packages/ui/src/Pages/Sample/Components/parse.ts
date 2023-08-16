export interface SampleRow {
  f: number;
  key: number;
  slots: SampleItem[][];
  active: number;
}

export interface SampleItem {
  frame: number;
  event: string;
  char: number;
  msg: string;
  raw: string;
  color: string;
  icon: string;
  amount: number;
  added: number;
  ended: number;
  target: "";
}

export function strFrameWithSec(frame: number): string {
  if (frame == -1) {
    return " [-1]";
  }
  const result = " [" + frame.toString() + " | " + (frame / 60).toFixed(2).toString() + "s]";
  return result;
}

export function eventColor(eve: string): string {
  switch (eve) {
    case "procs":
      return "";
    case "damage":
      return "#2563EB";
    case "hurt":
      return "";
    case "heal":
      return "";
    case "calc":
      return "#9D174D";
    case "reaction":
      return "";
    case "element":
      return "#3F60A6";
    case "snapshot":
      return "#6366F1";
    case "snapshot_mods":
      return "#818CF8";
    case "pre_damage_mods":
      return "#818CF8";
    case "status":
      return "#902D89";
    case "cooldown":
      return "#0D9488"; // tailwind teal-600
    case "action":
      return "#AB5F45";
    case "hitlag":
      return "#A27B5C";
    case "user":
      return "#5F7161";
    case "enemy":
      return "#632626";
    case "queue":
      return "";
    case "energy":
      return "#036345";
    case "character":
      return "";
    case "hook":
      return "";
    case "sim":
      return "";
    case "task":
      return "";
    case "artifact":
      return "";
    case "weapon":
      return "";
    case "shield":
      return "";
    case "construct":
      return "";
    case "icd":
      return "";
    default:
      return "gray-500";
  }
}
