import { CartesianGrid, Legend, Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis, } from "recharts";
import { SimResults } from "../DataType";
import { Trans, useTranslation } from "react-i18next";

const formatY = (value: string, index: number) => {
  const num = parseFloat(value);
  return Math.abs(num) > 999
    ? (Math.sign(num) * (Math.abs(num) / 1000)).toFixed(0) + "k"
    : (Math.sign(num) * Math.abs(num)).toFixed(0);
};

export default function DPSOverTime({ data }: { data: SimResults }) {
  useTranslation()

  let dmgOverTime: {
    time: number;
    mean: number;
    upper: number;
    lower: number;
  }[] = [];
  let max = 0;
  let maxDmg = 0;

  for (const key in data.damage_over_time) {
    const val = parseFloat(key);
    if (val === NaN) {
      continue;
    }
    if (max < val) {
      max = val;
    }
    let sd = 0;
    if (data.damage_over_time[key].sd) {
      sd = data.damage_over_time[key].sd!;
    }
    if (maxDmg < data.damage_over_time[key].mean + sd) {
      maxDmg = data.damage_over_time[key].mean + sd;
    }
    let lower = data.damage_over_time[key].mean - sd;
    if (lower < 0) {
      lower = 0;
    }
    dmgOverTime.push({
      time: val,
      mean: data.damage_over_time[key].mean,
      upper: data.damage_over_time[key].mean + sd,
      lower: lower,
    });
  }

  dmgOverTime.sort((a, b) => {
    return a.time - b.time;
  });
  return (
    <div>
      <span className="ml-2 mt-1 font-bold capitalize absolute top-0 left-0">
        <Trans>viewer.damage_dealt_over</Trans>
      </span>
      <ResponsiveContainer width="95%" height={288}>
        <LineChart data={dmgOverTime}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis
            dataKey="time"
            tick={{ fill: "white" }}
            type="number"
            domain={[0, max]}
            unit="s"
            tickCount={20}
          />
          <YAxis
            tick={{ fill: "white" }}
            domain={[0, maxDmg]}
            type="number"
            tickFormatter={formatY}
          />
          <Tooltip labelStyle={{ color: "black" }} />
          <Legend />
          <Line
            dot={false}
            type="linear"
            dataKey="mean"
            stroke="#8884d8"
            strokeWidth={1.5}
          />
          <Line
            dot={false}
            type="linear"
            dataKey="upper"
            stroke="#82ca9d"
            strokeDasharray="7 7"
          />
          <Line
            dot={false}
            type="linear"
            dataKey="lower"
            stroke="#F2B824"
            strokeDasharray="7 7"
          />
        </LineChart>
        {/* 
        <BarChart data={reactionsTriggered}>
          <Tooltip />
          <YAxis type="number" dataKey="value" tick={{ fill: "white" }} />
          <XAxis type="category" dataKey="name" tick={{ fill: "white" }} />
          <Bar dataKey="value" cx="50%" cy="50%" isAnimationActive={false}>
            {reactionsTriggered.map((entry, index) => (
              <Cell key={index} fill={COLORS[index % COLORS.length]} />
            ))}
          </Bar>
        </BarChart> */}
      </ResponsiveContainer>
    </div>
  );
}
