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
import { RootState } from "app/store";
import React from "react";
import { useSelector } from "react-redux";
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
import { selectLogs } from "./debugSlice";

function Debug() {
  const { data, names } = useSelector((state: RootState) => {
    return {
      data: selectLogs(state),
      names: state.debug.names,
    };
  });

  const [bin, setBin] = React.useState<number>(60);
  const [cumul, setCumul] = React.useState<boolean>(false);
  const [bar, setBar] = React.useState<boolean>(false);
  const [f, setF] = React.useState<number>(0);
  const [row, setRow] = React.useState<number>(0);
  const [col, setCol] = React.useState<number>(0);
  const [show, setShow] = React.useState<boolean>(false);
  const [logEvents, setLogEvents] = React.useState<Array<string>>([
    "damage",
    "hurts",
    "action",
    "energy",
    "element",
  ]);

  var toggleLogEvent = (val: string) => {
    var next: string[] = [];
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

  if (data.length === 0) {
    return (
      <div>
        No results. Please run sim first. Note that average mode will not
        generate debug logs
      </div>
    );
  }

  let hist: Array<{ x: number; value: number }> = [{ x: bin / 60, value: 0 }];
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
  //console.log(hist);

  return (
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
                checked={logEvents.includes("hurt")}
                label="hurt"
                onChange={(e) => toggleLogEvent("hurt")}
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
  );
}

export default Debug;
