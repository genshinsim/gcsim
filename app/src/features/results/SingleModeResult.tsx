import { Callout, Card, Elevation, H5 } from "@blueprintjs/core";
import React from "react";
import {
  Bar,
  BarChart,
  Cell,
  Legend,
  Pie,
  PieChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { SingleModeSummary } from "./resultsSlice";

const COLORS = [
  "#2965CC",
  "#29A634",
  "#D99E0B",
  "#D13913",
  "#8F398F",
  "#00B3A4",
  "#DB2C6F",
  "#9BBF30",
  "#96622D",
  "#7157D9",
];

// const CHAR_COLORS = ["#1D7324", "#AFC15A", "#93B9A2", "#1F4B99"];

// const CHAR_COLORS = [
//   "#003f5c",
//   "#2f4b7c",
//   "#665191",
//   "#a05195",
//   "#d45087",
//   "#f95d6a",
//   "#ff7c43",
//   "#ffa600",
// ];

const CHAR_COLORS = ["#4472C4", "#ED7D31", "#A5A5A5", "#70AD47"];

// const AURA_VAL = [
//   "Pyro",
//   "Hydro",
//   "Cryo",
//   "Electro",
//   "Geo",
//   "Anemo",
//   "Dendro",
//   "Physical",
//   "Frozen",
//   "EC",
//   "NoElement",
// ];

const RADIAN = Math.PI / 180;

const renderCustomizedLabel = ({
  cx,
  cy,
  midAngle,
  innerRadius,
  outerRadius,
  percent,
  index,
}: {
  cx: any;
  cy: any;
  midAngle: any;
  innerRadius: any;
  outerRadius: any;
  percent: any;
  index: any;
}) => {
  const radius = innerRadius + (outerRadius - innerRadius) * 0.5;
  const x = cx + radius * Math.cos(-midAngle * RADIAN);
  const y = cy + radius * Math.sin(-midAngle * RADIAN);

  return (
    <text
      x={x}
      y={y}
      fill="white"
      textAnchor={x > cx ? "start" : "end"}
      dominantBaseline="central"
    >
      {`${(percent * 100).toFixed(0)}%`}
    </text>
  );
};

function SingleModeResult({ data }: { data: SingleModeSummary }) {
  const [charSelected, setCharSelected] = React.useState<string>("");

  let dmg: { name: string; value: number }[] = [];
  let dmgDetail: { name: string; value: number }[] = [];
  let useCount: { name: string; value: number }[] = [];
  let useCountDetails: { name: string; value: number }[] = [];
  let fieldTime: { name: string; value: number }[] = [];
  let reactionCount: { name: string; value: number }[] = [];

  let index = -1;

  //dmg
  data.char_names.forEach((char, i) => {
    let total = 0;
    if (char === charSelected) {
      index = i;
    }
    //add up dmg per char?
    for (const [key, val] of Object.entries(data.damage_by_char[i])) {
      let v = Math.round((val * 60) / data.sim_duration);
      if (char === charSelected) {
        dmgDetail.push({
          name: key,
          value: v,
        });
      }
      total += v;
    }
    dmg.push({
      name: char,
      value: total,
    });
    //check abil usage
    total = 0;
    for (const [key, val] of Object.entries(data.abil_usage_count_by_char[i])) {
      if (char === charSelected) {
        useCountDetails.push({
          name: key,
          value: val,
        });
      }
      total += val;
    }
    useCount.push({
      name: char,
      value: total,
    });
    //check field time
    fieldTime.push({
      name: char,
      value: data.char_active_time[i] / 60,
    });
  });

  for (const [key, val] of Object.entries(data.reactions_triggered)) {
    reactionCount.push({
      name: key,
      value: val,
    });
  }

  return (
    <div>
      <Card style={{ margin: "5px" }} elevation={Elevation.TWO}>
        <H5>Damage Breakdown (in damage per second)</H5>
        <div className="row">
          <div className="col-xs-4">
            <ResponsiveContainer width="95%" height={400}>
              <PieChart>
                <Tooltip
                  formatter={(value: number, name: string) => {
                    return [
                      "" +
                        value.toFixed(2) +
                        " (" +
                        ((100 * value) / data.dps).toFixed(2) +
                        "%)",
                      name,
                    ];
                  }}
                />
                <Legend verticalAlign="top" height={36} />
                <Pie
                  data={dmg}
                  dataKey="value"
                  nameKey="name"
                  cx="50%"
                  cy="50%"
                  outerRadius={100}
                  labelLine={false}
                  label={renderCustomizedLabel}
                  onClick={(e: any) => setCharSelected(e.name)}
                >
                  {dmg.map((entry, index) => (
                    <Cell fill={CHAR_COLORS[index % CHAR_COLORS.length]} />
                  ))}
                </Pie>
              </PieChart>
            </ResponsiveContainer>
          </div>

          {charSelected === "" ? null : (
            <div className="col-xs-8">
              <ResponsiveContainer width="95%" height={400}>
                <BarChart data={dmgDetail} layout="vertical">
                  <Tooltip />
                  <XAxis type="number" dataKey="value" hide />
                  <YAxis type="category" dataKey="name" width={150} />
                  <Bar dataKey="value" cx="50%" cy="50%">
                    {dmgDetail.map((entry, index) => (
                      <Cell fill={COLORS[index % COLORS.length]} />
                    ))}
                  </Bar>
                </BarChart>
              </ResponsiveContainer>
            </div>
          )}
        </div>

        <Callout intent="primary">
          In total, the team did {Math.round(data.dps)} damage per second, over
          the course of {data.sim_duration / 60} seconds.
          <br />
          {index !== -1 ? (
            <span>
              {charSelected +
                " did " +
                Math.round(dmg[index].value) +
                " damage per second, representing " +
                ((100 * dmg[index].value) / data.dps).toFixed(2) +
                "% of the total team dps."}
            </span>
          ) : null}
        </Callout>
      </Card>
      <Card style={{ margin: "5px" }} elevation={Elevation.TWO}>
        <H5>Ability Usage Count (count)</H5>
        <div className="row">
          <div className="col-xs-4">
            <ResponsiveContainer width="95%" height={400}>
              <PieChart>
                <Tooltip />
                <Legend verticalAlign="top" height={36} />
                <Pie
                  data={useCount}
                  dataKey="value"
                  nameKey="name"
                  cx="50%"
                  cy="50%"
                  outerRadius={100}
                  label
                  onClick={(e: any) => setCharSelected(e.name)}
                >
                  {useCount.map((entry, index) => (
                    <Cell fill={CHAR_COLORS[index % CHAR_COLORS.length]} />
                  ))}
                </Pie>
              </PieChart>
            </ResponsiveContainer>
          </div>

          {charSelected === "" ? null : (
            <div className="col-xs-8">
              <ResponsiveContainer width="95%" height={400}>
                <BarChart data={useCountDetails} layout="vertical">
                  <Tooltip />
                  <XAxis type="number" dataKey="value" hide />
                  <YAxis type="category" dataKey="name" width={150} />
                  <Bar dataKey="value" cx="50%" cy="50%">
                    {useCountDetails.map((entry, index) => (
                      <Cell fill={COLORS[index % COLORS.length]} />
                    ))}
                  </Bar>
                </BarChart>
              </ResponsiveContainer>
            </div>
          )}
        </div>
      </Card>
      <Card style={{ margin: "5px" }} elevation={Elevation.TWO}>
        <H5>Character Field Time (in seconds)</H5>
        <div className="row">
          <div className="col-xs-4">
            <ResponsiveContainer width="95%" height={400}>
              <PieChart>
                <Tooltip />
                <Legend verticalAlign="top" height={36} />
                <Pie
                  data={fieldTime}
                  dataKey="value"
                  nameKey="name"
                  cx="50%"
                  cy="50%"
                  outerRadius={100}
                  label
                  onClick={(e: any) => setCharSelected(e.name)}
                >
                  {fieldTime.map((entry, index) => (
                    <Cell fill={CHAR_COLORS[index % CHAR_COLORS.length]} />
                  ))}
                </Pie>
              </PieChart>
            </ResponsiveContainer>
          </div>
        </div>
      </Card>
      <Card style={{ margin: "5px" }} elevation={Elevation.TWO}>
        <H5>Reactions (count)</H5>
        <div className="row">
          <div className="col-xs-4">
            <ResponsiveContainer width="95%" height={400}>
              <BarChart data={reactionCount} layout="vertical">
                <Tooltip />
                <XAxis type="number" dataKey="value" hide />
                <YAxis type="category" dataKey="name" width={150} />
                <Bar dataKey="value" cx="50%" cy="50%">
                  {reactionCount.map((entry, index) => (
                    <Cell fill={COLORS[index % COLORS.length]} />
                  ))}
                </Bar>
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>
      </Card>
    </div>
  );
}

export default SingleModeResult;
