export interface DebugRow {
  f: number;
  key: number;
  slots: DebugItem[][];
  active: number;
}

export interface DebugItem {
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

export function parseLog(
  active: string,
  team: string[],
  log: string,
  selected: string[]
) {
  console.log("parsing log");
  //find initial active char
  let activeIndex = team.findIndex((e) => e === active);
  activeIndex++;

  let result: DebugRow[] = [];
  let slots: DebugItem[][] = [[], [], [], [], []];

  let lastFrame = -1;

  //   console.log(log);
  //split the logs by new line
  const lines = log.split(/\r?\n/);

  let rowKey = 0;
  //bool to check if there are elements added
  let added = false;

  lines.forEach((line) => {
    if (line === "") {
      return;
    }
    //parse json
    let d: any = {};
    try {
      d = JSON.parse(line);
    } catch (e) {
      console.log("error reading line: ", line, " skipping");
      console.log(e);
      return;
    }

    //if no frame then set frame to -1
    if (!("frame" in d)) {
      d.frame = -1;
    }

    //if no event then set event to sim
    if (!("event" in d)) {
      d.event = "sim";
    }

    //shift char index down by 1 (b/c index 0 is sim stuff)
    let index = 0;
    if ("char" in d) {
      index = d.char + 1;
    } else {
      d.char = 0;
    }

    //check if frame changed; if so append stuff
    if (d.frame !== lastFrame) {
      if (added) {
        result.push({
          key: rowKey,
          f: lastFrame,
          slots: slots,
          active: activeIndex,
        });
      }
      added = false;
      rowKey++;
      //reset
      lastFrame = d.frame;
      slots = [];
      for (var i = 0; i <= team.length; i++) {
        slots.push([]);
      }
    }

    //parse the data
    let e: DebugItem = {
      frame: d.frame,
      msg: d.M,
      raw: JSON.stringify(JSON.parse(line), null, 2),
      event: d.event,
      char: d.char,
      color: eventColor(d.event),
      icon: "circle",
      amount: 0,
      added: d.frame,
      ended: d.frame,
      target: "",
    };

    if (e.color === "") {
      e.color = "#6B7280";
    }

    //skip if event is not in selected
    if (selected.indexOf(e.event) == -1) {
      return;
    }

    //set icon/color etc... based one vent
    switch (e.event) {
      case "damage":
        e.msg +=
          " [" +
          Math.round(d.damage)
            .toString()
            .replace(/\B(?=(\d{3})+(?!\d))/g, ",") +
          "]";
        let extra = "";
        if (d.amp && d.amp !== "") {
          extra += d.amp;
        }
        if (d.crit) {
          extra += " crit";
        }
        if (extra !== "") {
          e.msg += " (" + extra.trim() + ")";
        }

        e.icon = "local_fire_department";
        e.amount = d.damage;
        e.target = d.target;
        break;
      case "queue":
        let msg = "";
        if (d.failed) {
          msg = `(${d.reason}): ${d.exec}`;
          if (msg.length > 40) {
            msg = msg.slice(0, 40) + "...";
          }
        }
        e.msg += msg;

        e.icon = "queue";
        break;
      case "action":
        if (d.M.includes("executed") && d.action === "swap") {
          activeIndex = team.findIndex((e) => e === d.target) + 1;
          e.msg += " to " + d.target;
        }

        if (d.M.includes("cooldown")) {
          // Add expiry frame to the end if exists
          switch (d.expiry) {
            case undefined:
              break;
            default:
              e.msg += strFrameWithSec(d.expiry);
              e.msg = d.type + " " + e.msg;
          }
        }
        e.icon = "play_arrow";
        break;
      case "element":
        switch (d.M) {
          case "expired":
            e.msg = d.old_ele + " expired";
            break;
          case "application":
            // console.log(d.existing);
            e.msg = d.applied_ele + " applied";
            if (d.existing) {
              e.msg += " to [";
              let before = d.existing.map((x: string) =>
                x.replace(/: (.+)/, " ($1)")
              );
              if (before.length > 0) {
                e.msg += before.join(" ");
              } else {
                e.msg += "no aura";
              }
              e.msg += "]";
            }
            if (d.after) {
              e.msg += " ➜ [";
              let after = d.after.map((x: string) =>
                x.replace(/: (.+)/, " ($1)")
              );
              if (after.length > 0) {
                e.msg += after.join(" ");
              } else {
                e.msg += "no aura";
              }
              e.msg += "]";
            }
            break;
          case "refreshed":
            e.msg = d.ele + " refreshed";
            break;
          default:
            e.msg = d.M;
        }

        e.icon = "bolt";
        e.target = d.target;
        break;
      case "energy":
        if (e.msg.includes("particle")) {
          e.msg =
            d.M +
            " from " +
            d.source +
            ", next: " +
            Math.round(d["post_recovery"]);
        }
        if (e.msg.includes("adding energy")) {
          e.msg += ` ${d["rec'd"]}`;
        }
        e.icon = "local_cafe";
        break;
      case "calc":
        e.icon = "calculate";
        e.target = d.target;

        break;
      case "character":
        e.icon = "person";
        break;
      case "snapshot":
        e.icon = "photo_camera";
        break;
      case "snapshot_mods":
        e.icon = "build";
        break;
      case "pre_damage_mods":
        e.icon = "dynamic_form";
        break;
      case "heal":
        e.icon = "healing";
        break;
      case "hurt":
        e.icon = "coronavirus";
        break;
      case "shield":
        e.icon = "shield";
        break;
      case "hook":
        e.icon = "attachment";
        break;
      case "icd":
        e.icon = "timer";
        break;
      case "construct":
        e.icon = "apartment";
        break;
      case "status":
        e.icon = "iso";

        // Add expiry frame to the end if exists
        switch (d.expiry) {
          case undefined:
            break;
          default:
            e.msg += strFrameWithSec(d.expiry);
            e.msg = d.key + " " + e.msg;
        }

        if (d.target != undefined) {
          e.target = d.target;
        }
        break;
      default:
        e.msg = e.event + ": " + e.msg;
    }

    //add it to slots
    // console.log(slots);
    // console.log(e.char);
    // console.log(d);
    slots[index].push(e);
    added = true;
  });

  // console.log(result);

  return result;
}

export function strFrameWithSec(frame: number): string {
  if (frame == -1) {
    return " [-1]";
  }
  let result =
    " [" + frame.toString() + " | " + (frame / 60).toFixed(2).toString() + "s]";
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
    case "action":
      return "#AB5F45";
    case "queue":
      return "";
    case "energy":
      return "#036345";
    case "character":
      return "";
    case "enemy":
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
