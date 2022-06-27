import { DebugRow, DebugItem, eventColor, strFrameWithSec } from "./parse";

type LogDetails = {
  char_index: number;
  ended: number;
  event: string;
  frame: number;
  msg: string;
  logs: { [key in string]: any };
  ordering?: { [key: string]: number };
};

type endedStatus = {
  [key: number]: DebugItem[][];
};

export function parseLogV2(
  active: string,
  team: string[],
  log: string,
  selected: string[]
) {
  let activeIndex = team.findIndex((e) => e === active);
  activeIndex++; // +1 since we set the first field to be sim slot

  let result: DebugRow[] = [];
  let slots: DebugItem[][] = [[], [], [], [], []];
  let ended: endedStatus = {};
  let lastFrame = -1;
  let finalFrame = -1;

  //we just need to parse it here
  /**
        char_index: -1
        ended: 0
        event: "hook"
        frame: 0
        logs: (3) [{…}, {…}, {…}]
        msg: "hook added"
   */

  let lines: LogDetails[] = [];
  try {
    lines = JSON.parse(log);
  } catch (e) {
    console.warn("error parsing debug log (v2)");
    console.warn(e);
    return [];
  }

  // let rowKey = 0;
  //bool to check if there are elements added
  let added = false;

  //map out all ended
  lines.forEach((line) => {
    if (line.ended > line.frame) {
      if (line.event !== "status") {
        return;
      }
      if (!line.msg.includes("added")) {
        return;
      }
      if (!(line.ended in ended)) {
        let slots = [];
        for (var i = 0; i <= team.length; i++) {
          slots.push([]);
        }
        ended[line.ended] = slots;
      }
      const key = line.logs["key"];
      if (key === null || key === "") {
        return;
      }
      const index = line.char_index + 1;
      let e: DebugItem = {
        frame: line.frame,
        msg: key + " expired" + strFrameWithSec(line.frame),
        raw: JSON.stringify(line, null, 2),
        event: line.event,
        char: index,
        color: eventColor(line.event),
        icon: "iso",
        amount: 0,
        target: "",
        added: line.frame,
        ended: line.ended,
      };
      ended[e.ended][index].push(e);
    }
  });

  lines.forEach((line) => {
    const index = line.char_index + 1;
    if (line.frame > finalFrame) {
      finalFrame = line.frame;
    }

    if (line.frame !== lastFrame) {
      if (added) {
        result.push({
          key: lastFrame,
          f: lastFrame,
          slots: slots,
          active: activeIndex,
        });
      }
      added = false;
      // rowKey++;
      //reset
      lastFrame = line.frame;
      slots = [];
      for (var i = 0; i <= team.length; i++) {
        slots.push([]);
      }
    }

    //make a copy of line and sort by ordering if ordering exist (then purge ordering)

    let logLines: { key: string; val: any }[] = [];
    //convert logs into array
    for (const key in line.logs) {
      logLines.push({ key: key, val: line.logs[key] });
    }
    //sort
    if (line.ordering) {
      logLines.sort((a, b) => {
        let ao = line.ordering![a.key] || 0;
        let bo = line.ordering![b.key] || 0;
        return ao - bo;
      });
      //get rid of ordering
      delete line.ordering;
    }
    //store it back
    line.logs = {};
    logLines.forEach((e) => {
      line.logs[e.key] = e.val;
    });

    let e: DebugItem = {
      frame: line.frame,
      msg: line.msg,
      raw: JSON.stringify(line, null, 2),
      event: line.event,
      char: index,
      color: eventColor(line.event),
      icon: "circle",
      amount: 0,
      target: "",
      added: line.frame,
      ended: line.ended,
    };

    if (e.color === "") {
      e.color = "#6B7280";
    }

    const d = line.logs;
    //hightlight active char
    if (
      e.event === "action" &&
      line.msg.includes("executed") &&
      d.action === "swap"
    ) {
      activeIndex = e.char;
    }
    //skip if event is not in selected
    if (selected.indexOf(e.event) == -1) {
      return;
    }
    //set icon/color etc... based on event
    switch (e.event) {
      case "damage":
        //grab dmg amount
        const dmg = Math.round(line.logs["damage"])
          .toString()
          .replace(/\B(?=(\d{3})+(?!\d))/g, ",");
        e.msg += ` [${dmg}]`;
        let extra = "";
        const amp = line.logs["amp"] ? line.logs["amp"] : "";
        if (amp && amp !== "") {
          extra += amp;
        }
        const crit = line.logs["crit"] ? line.logs["crit"] : "";
        if (crit) {
          extra += " crit";
        }
        if (extra !== "") {
          e.msg += " (" + extra.trim() + ")";
        }

        e.icon = "local_fire_department";
        e.amount = line.logs[d.damage];
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
        if (line.msg.includes("executed") && d.action === "swap") {
          e.msg += " to " + d.target;
        }

        if (line.msg.includes("cooldown")) {
          // Add expiry frame to the end if exists
          switch (d.expiry) {
            case undefined:
              break;
            default:
              e.msg += strFrameWithSec(d.expiry);
              e.msg = d.type + " " + e.msg;
          }
        }
        //trim "executed "
        e.msg = e.msg.replace("executed ", "");
        e.icon = "play_arrow";
        break;
      case "hitlag":
        e.icon = "severe_cold";
        break;
      case "enemy":
        e.icon = "mood_bad";
        break;
      case "user":
        e.icon = "comment";
        break;
      case "element":
        switch (line.msg) {
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
              }
              e.msg += "]";
            } else {
              e.msg += " [no aura]";
            }
            if (d.after) {
              e.msg += " ➜ [";
              let after = d.after.map((x: string) =>
                x.replace(/: (.+)/, " ($1)")
              );
              if (after.length > 0) {
                e.msg += after.join(" ");
              }
              e.msg += "]";
            } else {
              e.msg += " ➜ [no aura]";
            }
            break;
          case "refreshed":
            e.msg = d.ele + " refreshed";
            break;
          default:
            e.msg = line.msg;
        }

        e.icon = "bolt";
        e.target = d.target;
        break;
      case "energy":
        if (e.msg.includes("particle")) {
          e.msg =
            line.msg +
            " from " +
            d.source +
            ", next: " +
            Math.floor(d["post_recovery"]);
        }
        if (e.msg.includes("adding energy")) {
          let amt = d["rec'd"];
          if (typeof amt === "number") {
            amt = amt.toFixed(2);
          }
          e.msg =
            "adding " +
            amt +
            " energy from " +
            d.source +
            ", next: " +
            Math.floor(d["post_recovery"]);
        }
        if (d["post_recovery"] == d["max_energy"] && d["max_energy"]) {
          e.msg += " (max)";
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

        // this hacky but i don't care
        if (e.ended === e.frame && line.msg.includes("refreshed")) {
          let idx = lines.findIndex((a) => {
            return (
              a.event === "status" &&
              line.char_index === a.char_index &&
              !a.logs.overwrite &&
              a.logs.key === line.logs.key &&
              line.frame >= a.frame &&
              line.frame < a.ended
            );
          });
          if (idx !== -1) {
            e.added = lines[idx].frame;
            e.ended = lines[idx].ended;
          }
        }

        if (d.target != undefined) {
          e.target = d.target;
        }
        break;
      default:
        e.msg = e.event + ": " + e.msg;
    }

    slots[index].push(e);
    added = true;
  });
  if (added) {
    result.push({
      key: lastFrame,
      f: lastFrame,
      slots: slots,
      active: activeIndex,
    });
  }

  //append in ended status
  console.log(ended);

  if (selected.indexOf("status") != -1) {
    for (let f = -1; f <= finalFrame; f++) {
      if (!(f in ended)) {
        continue;
      }

      let idx = result.findIndex((e) => e.f === f);
      if (idx !== -1) {
        for (let j = 0; j < result[idx].slots.length; j++) {
          result[idx].slots[j].push(...ended[f][j]);
        }
        continue;
      }

      // TODO: set active correctly instead of using last one. only matters if action log is disabled
      let active = 0;
      let insertAt = result.findIndex((r, i) => {
        active = result[i].active;
        return r.f > f;
      });

      const x = {
        key: f,
        f: f,
        slots: ended[f],
        active: active,
      };
      if (insertAt !== -1) {
        result.splice(insertAt, 0, x);
      } else {
        result.push(x);
      }
    }
  }

  //   console.log(result);

  return result;
}
