import React from "react";
import styled from "styled-components";
import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";
import character_data from "./character_data.json";
import weapon_data from "./weapon_data.json";
import artifact_data from "./artifact_data.json";

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

const abils = ["normal", "charge", "plunge", "aim", "skill", "burst", "asc", "cons"];
const abilLabels = [
  "Normal",
  "Charge Attack",
  "Plunge",
  "Aimed Shot",
  "Skill",
  "Burst",
  "Ascension",
  "Cons"
];

export default function AoETable({ item_key, data_src }) {
  let data = character_data;
  let useAbils = true;
  switch (data_src) {
    case "weapon":
      data = weapon_data;
      useAbils = false;
      break;
    case "artifact":
      data = artifact_data;
      useAbils = false;
      break;
  }
  if (!data) {
    return <div>No AoE data</div>;
  }
  if (!(item_key in data)) {
    return <div>No AoE data for {data_src}</div>;
  }
  const item = data[item_key];
  let tabs = [];
  let count = 0;
  if (!useAbils) {
    if (item.length === 0) {
      return <div>No AoE data for {data_src}</div>;
    }
    return <AbilAoE data={item} />;
  } 
  abils.forEach((a, i) => {
    // skip if no data for this tab
    if (!(a in item)) {
      return;
    }
    count++;
    tabs.push(
      <TabItem value={a} label={abilLabels[i]} key={a}>
        <AbilAoE data={item[a]} />
      </TabItem>
    );
  });
  if (count == 0) {
    return <div>No AoE data for {data_src}</div>;
  }
  return <Tabs>{tabs}</Tabs>;
}
