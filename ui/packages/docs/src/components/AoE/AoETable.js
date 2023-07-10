import React from "react";
import styled from "styled-components";
import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";
import data from "./data.json";

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

function AbilAoE({ data }) {
  const rows = data.map((e) => {
    return (
      <tr key={e.ability}>
        <TD>{e.ability}</TD>
        <TD>{e.shape}</TD>
        <TD>{e.center}</TD>
        <TD>{e.offsetX ?? '-'}</TD>
        <TD>{e.offsetY ?? '-'}</TD>
        <TD>{e.radius ?? '-'}</TD>
        <TD>{e.fanAngle ?? '-'}</TD>
        <TD>{e.boxX ?? '-'}</TD>
        <TD>{e.boxY ?? '-'}</TD>
        <TD>{e.notes ?? '-'}</TD>
      </tr>
    );
  });

  return (
    <div style={{ marginTop: "1rem", width: "100%" }}>
      <Table>
        <Thead>
          <tr>
            <TH>Ability</TH>
            <TH>Shape</TH>
            <TH>Center</TH>
            <TH>Offset X</TH>
            <TH>Offset Y</TH>
            <TH>Radius</TH>
            <TH>Fan Angle</TH>
            <TH>Box X</TH>
            <TH>Box Y</TH>
            <TH>Notes</TH>
          </tr>
        </Thead>
        <tbody>{rows}</tbody>
      </Table>
    </div>
  );
}

const abils = ["normal", "charge", "skill", "burst", "asc", "cons"];
const abilLabels = [
  "Normal",
  "Charge Attack",
  "Skill",
  "Burst",
  "Ascension",
  "Cons",
];

export default function AoETable({ character }) {
  if (!(character in data)) {
    return <div>No AoE data for character</div>;
  }
  const char = data[character];
  let tabs = [];
  let count = 0;
  abils.forEach((a, i) => {
    //skip if no data for this tab
    if (!(a in char)) {
      return;
    }
    count++;
    tabs.push(
      <TabItem value={a} label={abilLabels[i]} key={a}>
        <AbilAoE data={char[a]} />
      </TabItem>
    );
  });
  if (count == 0) {
    return <div>No AoE data for character</div>;
  }
  return <Tabs>{tabs}</Tabs>;
}
