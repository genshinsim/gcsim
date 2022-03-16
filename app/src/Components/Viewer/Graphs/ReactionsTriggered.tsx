import { Bar, BarChart, Cell, ResponsiveContainer, Tooltip, XAxis, YAxis, } from "recharts";
import { COLORS } from "./Graphs";
import { SimResults } from "../DataType";
import { Trans } from "~node_modules/react-i18next";

export default function ReactionsTriggered({ data }: { data: SimResults }) {
  let reactionsTriggered: { name: string; value: number }[] = [];

  for (const key in data.reactions_triggered) {
    reactionsTriggered.push({
      name: key,
      value: data.reactions_triggered[key].mean,
    });
  }
  return (
    <div>
      <span className="ml-2 mt-1 font-bold capitalize absolute top-0 left-0">
        <Trans>viewer.reactions_triggered</Trans>
      </span>
      <ResponsiveContainer width="95%" height={288}>
        <BarChart data={reactionsTriggered}>
          <Tooltip />
          <YAxis type="number" dataKey="value" tick={{ fill: "white" }} />
          <XAxis type="category" dataKey="name" tick={{ fill: "white" }} />
          <Bar dataKey="value" cx="50%" cy="50%" isAnimationActive={false}>
            {reactionsTriggered.map((entry, index) => (
              <Cell key={index} fill={COLORS[index % COLORS.length]} />
            ))}
          </Bar>
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}
