import React from "react";
import styled from "styled-components";
import enemy_data from "./enemy_data.json";

const Table = styled.table`
  border-collapse: collapse;
  width: 100%;
`;

const Thead = styled.thead`
  background-color: #333;
`;

const TD = styled.td`
  padding: 0.5rem;
`;

const TH = styled.th`
  padding: 0.5rem;
`;

const elements = [
  "clear",
  "pyro",
  "hydro",
  "dendro",
  "electro",
  "anemo",
  "cryo",
  "geo",
];

function ParticleInfo(item_key, { drop_id, hp_percent }) {
  const count = (drop_id / 10) % 10;
  const element = elements[drop_id % 10];
  return (
    <tr key={item_key}>
      <TD>{(hp_percent * 100).toFixed()}%</TD>
      <TD>{count}</TD>
      <TD><code>{element}</code></TD>
    </tr>
  );
}

export default function ParticleTable({ item_key, data_src }) {
  let data = enemy_data;
  if (!(item_key in data) || data[item_key].length === 0) {
    return <div>No Particle data</div>;
  }
  const rows = data[item_key].map((e) => ParticleInfo(item_key, e));
  return (
    <div style={{ marginTop: "1rem", width: "100%" }}>
      <Table>
        <Thead>
          <tr>
            <TH>HP% thereshold</TH>
            <TH>Count</TH>
            <TH>Element</TH>
          </tr>
        </Thead>
        <tbody>{rows}</tbody>
      </Table>
    </div>
  );
}
