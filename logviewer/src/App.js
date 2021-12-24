import React from "react";
import {
  Dialog,
  H4,
  HTMLTable,
  Classes as CoreClasses,
  Checkbox,
  FormGroup,
  Tag,
  ButtonGroup,
  Button,
  Slider,
  Switch,
} from "@blueprintjs/core";
import { IconNames } from "@blueprintjs/icons";
import {
  Bar,
  BarChart,
  Legend,
  Line,
  LineChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";

function parseLogs() {
  let current = names.findIndex((e) => e === active);
  current++;
  const charCount = names.length;
  const str = logs.split(/\r?\n/);
  let data = [];
  let cols = [
    {
      Parts: [],
    },
    {
      Parts: [],
    },
    {
      Parts: [],
    },
    {
      Parts: [],
    },
    {
      Parts: [],
    },
  ];
  let lastFrame = -1;

  str.forEach((v) => {
    if (v === "") {
      return;
    }
    const d = JSON.parse(v);

    if (!("frame" in d)) {
      console.log("error, no frames: ", d);
      return;
    }

    if (!("event" in d)) {
      console.log("error, no event: ", d);
      return;
    }

    const event = d.event;

    let char = 0;

    if ("char" in d) {
      char = d.char + 1;
    }

    //if frame changed
    if (d.frame !== lastFrame) {
      //add it to the labels
      data.push({
        F: lastFrame,
        Cols: cols,
        Active: current,
      });
      lastFrame = d.frame;
      cols = [];
      for (var i = 0; i <= charCount; i++) {
        cols.push({ Parts: [] });
      }
    }

    let x = {
      F: d.frame,
      M: d.M,
      Raw: JSON.stringify(JSON.parse(v), null, 2),
      Event: event,
      Char: char,
      Intent: "none",
      Icon: "circle",
      Amount: 0,
      Right: "",
    };

    switch (x.Event) {
      case "damage":
        x.M +=
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
          x.M += " (" + extra.trim() + ")";
        }

        x.Icon = IconNames.FLAME;
        x.Intent = "primary";
        x.Amount = d.damage;
        x.Right = d.target;
        break;
      case "action":
        if (d.M.includes("executed") && d.action === "swap") {
          current = names.findIndex((e) => e === d.target);
          current++;
          x.M += " to " + d.target;
        }
        x.Icon = IconNames.PLAY;
        break;
      case "element":
        switch (d.M) {
          case "expired":
            x.M = d.old_ele + " expired";
            break;
          case "application":
            x.M =
              d.applied_ele +
              " applied" +
              (d.existing_ele === "" ? "" : " to " + d.existing_ele);
            break;
          case "refreshed":
            x.M = d.ele + " refreshed";
            break;
          default:
            x.M = d.M;
        }

        x.Icon = IconNames.FLASH;
        x.Intent = "warning";
        x.Right = d.target;
        break;
      case "energy":
        if (x.M.includes("particle")) {
          x.M =
            d.M +
            " from " +
            d.source +
            ", next: " +
            Math.round(d["post_recovery"]);
        }
        x.Icon = IconNames.IMPORT;
        x.Intent = "success";
        break;
      case "calc":
        x.Icon = IconNames.CALCULATOR;
        x.Right = d.target;
        break;
      case "character":
        x.Icon = IconNames.PERSON;
        break;
      case "snapshot":
        x.Icon = IconNames.CAMERA;
        break;
      case "snapshot_mods":
        x.Icon = IconNames.BUILD;
        break;
      case "heal":
        x.Icon = IconNames.HEART;
        break;
      case "hurt":
        x.Icon = IconNames.VIRUS;
        break;
      case "queue":
        x.Icon = IconNames.ADD_TO_ARTIFACT;
        break;
      case "shield":
        x.Icon = IconNames.SHIELD;
        break;
      case "hook":
        x.Icon = IconNames.PAPERCLIP;
        break;
      case "icd":
        x.Icon = IconNames.STOPWATCH;
        break;
      case "construct":
        x.Icon = IconNames.WRENCH;
        break;
      default:
        x.M = event + ": " + d.M;
    }

    if (!cols[char]) {
      console.log("?????");
      console.log(cols);
      console.log(char);
      console.log(d);
    }

    cols[char].Parts.push(x);
  });

  return data;
}

function App() {
  const [bin, setBin] = React.useState(60);
  const [cumul, setCumul] = React.useState(false);
  const [bar, setBar] = React.useState(false);
  const [f, setF] = React.useState(0);
  const [row, setRow] = React.useState(0);
  const [col, setCol] = React.useState(0);
  const [show, setShow] = React.useState(false);
  const [logEvents, setLogEvents] = React.useState([
    "damage",
    "hurts",
    "action",
    "energy",
    "element",
  ]);

  var toggleLogEvent = (val) => {
    var next = [];
    var found = false;
    for (var i = 0; i < logEvents.length; i++) {
      if (logEvents[i] === val) {
        found = true;
        continue;
      }
      next.push(logEvents[i]);
    }

    if (!found) {
      next.push(val);
    }

    setLogEvents(next);
  };

  const data = parseLogs();

  if (data.length === 0) {
    return <div>No data to show</div>;
  }

  let hist = [{ x: bin / 60, value: 0 }];
  let histCounter = 0;
  let last = 0;

  const rows = data.map((e, i) => {
    let hasData = false;
    const columns = e.Cols.map((c, j) => {
      const labels = c.Parts.map((l, k) => {
        if (l.Event === "damage") {
          //add it to hist
          if (l.F < histCounter * bin) {
            hist[histCounter].value += l.Amount;
          } else {
            last = hist[histCounter] ? hist[histCounter].value : 0;
            histCounter++;
            hist.push({
              x: (histCounter * bin) / 60,
              value: cumul ? last + l.Amount : l.Amount,
            });
          }
        }
        if (logEvents.includes(l.Event)) {
          hasData = true;
          return (
            <div key={i + "-" + j + "-" + k}>
              <Tag
                onClick={() => {
                  setF(i);
                  setRow(k);
                  setCol(j);
                  setShow(true);
                }}
                style={{
                  marginTop: "2px",
                  marginBottom: "2px",
                }}
                interactive
                icon={l.Icon}
                multiline
                intent={l.Intent}
                fill
                rightIcon={<span>{l.Right}</span>}
              >
                {l.M}
              </Tag>
            </div>
          );
        }

        return null;
      });
      return (
        <td
          key={i + "-" + j}
          style={{
            backgroundColor: e.Active === j ? "#D1F26D" : "transparent",
          }}
        >
          {labels}
        </td>
      );
    });
    if (hasData) {
      return (
        <tr key={i}>
          <td>{e.F}</td>
          {columns}
        </tr>
      );
    }
    return null;
  });

  const charCount = names.length + 1;
  const widthPer = Math.round(90 / charCount).toString() + "%";

  const n = names.map((e, i) => {
    return (
      <th key={i} style={{ width: widthPer, backgroundColor: "#CED9E0" }}>
        {e}
      </th>
    );
  });

  return (
    <div className="App">
      <div style={{ marginLeft: "20px", marginRight: "20px" }}>
        <div className="row">
          <div className="col-xs-10">
            <div className="box">
              <H4>Damage Graph</H4>
              <ResponsiveContainer width="95%" height={400}>
                {bar ? (
                  <BarChart data={hist}>
                    <XAxis dataKey="x" />
                    <YAxis />
                    <Tooltip />
                    <Legend />
                    <Bar type="monotone" dataKey="value" fill="#82ca9d" />
                  </BarChart>
                ) : (
                  <LineChart data={hist}>
                    <XAxis dataKey="x" />
                    <YAxis />
                    <Tooltip />
                    <Legend />
                    <Line type="monotone" dataKey="value" stroke="#82ca9d" />
                  </LineChart>
                )}
              </ResponsiveContainer>
              <H4>Log</H4>
              <HTMLTable condensed bordered style={{ width: "100%" }}>
                <thead>
                  <tr>
                    <th style={{ backgroundColor: "#CED9E0" }}>F</th>
                    <th style={{ width: widthPer, backgroundColor: "#CED9E0" }}>
                      Sim
                    </th>
                    {n}
                  </tr>
                </thead>
                <tbody>{rows}</tbody>
              </HTMLTable>
            </div>
          </div>
          <div className="col-xs-2">
            <H4>Graph Options</H4>
            <div className="stick">
              <FormGroup helperText="How many frames to bin damage amount">
                <Slider
                  min={60}
                  max={600}
                  stepSize={10}
                  labelStepSize={600}
                  onChange={(val) => setBin(val)}
                  value={bin}
                  vertical={false}
                />
              </FormGroup>
              <Switch
                checked={cumul}
                onChange={(e) => setCumul(e.currentTarget.checked)}
              >
                Show cumulative
              </Switch>
              <Switch
                checked={bar}
                innerLabel="line"
                innerLabelChecked="bar"
                onChange={(e) => setBar(e.currentTarget.checked)}
              >
                Toggle bar/line graph
              </Switch>
              <H4>Log Options</H4>
              <FormGroup helperText="which logs should be shown">
                <Checkbox
                  checked={logEvents.includes("procs")}
                  label="procs"
                  onChange={(e) => toggleLogEvent("procs")}
                />
                <Checkbox
                  checked={logEvents.includes("damage")}
                  label="damage"
                  onChange={(e) => toggleLogEvent("damage")}
                />
                <Checkbox
                  checked={logEvents.includes("pre_damage_mods")}
                  label="pre_damage_mods"
                  onChange={(e) => toggleLogEvent("pre_damage_mods")}
                />
                <Checkbox
                  checked={logEvents.includes("hurt")}
                  label="hurt"
                  onChange={(e) => toggleLogEvent("hurt")}
                />
                <Checkbox
                  checked={logEvents.includes("heal")}
                  label="heal"
                  onChange={(e) => toggleLogEvent("heal")}
                />
                <Checkbox
                  checked={logEvents.includes("calc")}
                  label="calc"
                  onChange={(e) => toggleLogEvent("calc")}
                />
                <Checkbox
                  checked={logEvents.includes("reaction")}
                  label="reaction"
                  onChange={(e) => toggleLogEvent("reaction")}
                />
                <Checkbox
                  checked={logEvents.includes("element")}
                  label="element"
                  onChange={(e) => toggleLogEvent("element")}
                />
                <Checkbox
                  checked={logEvents.includes("snapshot")}
                  label="snapshot"
                  onChange={(e) => toggleLogEvent("snapshot")}
                />
                <Checkbox
                  checked={logEvents.includes("snapshot_mods")}
                  label="mods (snapshot)"
                  onChange={(e) => toggleLogEvent("snapshot_mods")}
                />
                <Checkbox
                  checked={logEvents.includes("status")}
                  label="status"
                  onChange={(e) => toggleLogEvent("status")}
                />
                <Checkbox
                  checked={logEvents.includes("action")}
                  label="action"
                  onChange={(e) => toggleLogEvent("action")}
                />
                <Checkbox
                  checked={logEvents.includes("queue")}
                  label="queue"
                  onChange={(e) => toggleLogEvent("queue")}
                />
                <Checkbox
                  checked={logEvents.includes("energy")}
                  label="energy"
                  onChange={(e) => toggleLogEvent("energy")}
                />
                <Checkbox
                  checked={logEvents.includes("character")}
                  label="character"
                  onChange={(e) => toggleLogEvent("character")}
                />
                <Checkbox
                  checked={logEvents.includes("enemy")}
                  label="enemy"
                  onChange={(e) => toggleLogEvent("enemy")}
                />
                <Checkbox
                  checked={logEvents.includes("hook")}
                  label="hook"
                  onChange={(e) => toggleLogEvent("hook")}
                />
                <Checkbox
                  checked={logEvents.includes("sim")}
                  label="sim"
                  onChange={(e) => toggleLogEvent("sim")}
                />
                <Checkbox
                  checked={logEvents.includes("task")}
                  label="task"
                  onChange={(e) => toggleLogEvent("task")}
                />
                <Checkbox
                  checked={logEvents.includes("artifact")}
                  label="artifact"
                  onChange={(e) => toggleLogEvent("artifact")}
                />
                <Checkbox
                  checked={logEvents.includes("weapon")}
                  label="weapon"
                  onChange={(e) => toggleLogEvent("weapon")}
                />
                <Checkbox
                  checked={logEvents.includes("shield")}
                  label="shield"
                  onChange={(e) => toggleLogEvent("shield")}
                />
                <Checkbox
                  checked={logEvents.includes("construct")}
                  label="construct"
                  onChange={(e) => toggleLogEvent("construct")}
                />
                <Checkbox
                  checked={logEvents.includes("icd")}
                  label="icd"
                  onChange={(e) => toggleLogEvent("icd")}
                />
              </FormGroup>
              <ButtonGroup vertical fill>
                <Button intent="danger" onClick={() => setLogEvents([])}>
                  Clear Options
                </Button>
                <Button
                  onClick={() =>
                    setLogEvents([
                      "damage",
                      "hurts",
                      "action",
                      "energy",
                      "element",
                    ])
                  }
                >
                  Show Defaults
                </Button>
              </ButtonGroup>
            </div>
          </div>
        </div>
        <Dialog
          isOpen={show}
          canEscapeKeyClose={true}
          canOutsideClickClose={true}
          onClose={() => setShow(false)}
        >
          <div className={CoreClasses.DIALOG_BODY}>
            <pre>
              {data[f].Cols[col].Parts[row]
                ? data[f].Cols[col].Parts[row].Raw
                : ""}
            </pre>
          </div>
        </Dialog>
      </div>
    </div>
  );
}

export default App;

const active = "{{.Active}}";
const nameString = `{{.Team}}`;
const names = nameString.split(",");
const logs = `{{.Log}}`;
