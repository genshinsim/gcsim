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
import { AvgModeSummary } from "./resultsSlice";

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
const CHAR_COLORS = ["#4472C4", "#ED7D31", "#A5A5A5", "#70AD47"];

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

function AverageModeResult({ data }: { data: AvgModeSummary }) {
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
      let v = Math.round(val.mean * 100) / 100;
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
      let v = Math.round(val.mean * 100) / 100;
      if (char === charSelected) {
        useCountDetails.push({
          name: key,
          value: v,
        });
      }
      total += v;
    }
    useCount.push({
      name: char,
      value: total,
    });
    //check field time
    fieldTime.push({
      name: char,
      value: Math.round((100 * data.char_active_time[i].mean) / 60) / 100,
    });
  });

  for (const [key, val] of Object.entries(data.reactions_triggered)) {
    reactionCount.push({
      name: key,
      value: Math.round(val.mean * 100) / 100,
    });
  }

  return (
    <div>
      <Card style={{ margin: "5px" }} elevation={Elevation.TWO}>
        <H5>Damage Breakdown (in damage per second, on average)</H5>
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
                      ((100 * value) / data.dps.mean).toFixed(2) +
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
          In total, the team did {data.dps.mean.toFixed(2)} damage per second
          (on average, over {data.iter} iterations [min:{" "}
          {data.dps.min.toFixed(2)}, max: {data.dps.max.toFixed(2)}, std dev:{" "}
          {data.dps.sd?.toFixed(2)}]), over the course of{" "}
          {data.sim_duration.mean.toFixed(2)} seconds.
          <br />
          {index !== -1 ? (
            <span>
              {charSelected +
                " did " +
                Math.round(dmg[index].value) +
                " damage per second, representing " +
                ((100 * dmg[index].value) / data.dps.mean).toFixed(2) +
                "% of the total team dps."}
            </span>
          ) : null}
        </Callout>
      </Card>
      <Card style={{ margin: "5px" }} elevation={Elevation.TWO}>
        <H5>Ability Usage Count (count, on average)</H5>
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
        <H5>Character Field Time (in seconds, on average)</H5>
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
        <H5>Reactions (count, on average)</H5>
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
      <Card style={{ margin: "5px" }} elevation={Elevation.TWO}>
        <H5>Text Summary</H5>
        <pre>{data.text}</pre>
      </Card>
    </div>
  );
}

export default AverageModeResult;
