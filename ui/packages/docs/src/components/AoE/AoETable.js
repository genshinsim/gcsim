import React from "react";
import styled from "styled-components";
import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";
import character_data from "./character.dm.json";
import weapon_data from "./weapon.dm.json";
import artifact_data from "./artifact.dm.json";

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
      <tr key={e.name}>
        <TD>{e.name}</TD>
        <TD>{e.shape}</TD>
        <TD>{e.center}</TD>
        <TD>{e.offset_x ?? '-'}</TD>
        <TD>{e.offset_y ?? '-'}</TD>
        <TD>{e.radius ?? '-'}</TD>
        <TD>{e.fan_angle ?? '-'}</TD>
        <TD>{e.box_x ?? '-'}</TD>
        <TD>{e.box_y ?? '-'}</TD>
        <TD>{e.note ?? '-'}</TD>
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

const labels = {
  ["attack"]: "Normal",
  ["charge"]: "Charge",
  ["plunge"]: "Plunge",
  ["aim"]: "Aimed",
  ["skill"]: "Skill",
  ["burst"]: "Burst",
  ["asc"]: "Ascension",
  ["cons"]: "Constellation",
  ["-"]: "Other",
};

export default function AoETable({ item_key, data_src }) {
  let data = character_data;
  switch (data_src) {
    case "weapon":
      data = weapon_data;
      break;
    case "artifact":
      data = artifact_data;
      break;
  }
  if (!data) {
    return <div>No AoE data</div>;
  }
  if (!(item_key in data)) {
    return <div>No AoE data for {data_src}</div>;
  }
  let abil_data = {};
  data[item_key].forEach((e) => {
    abil_data[e.ability] = abil_data[e.ability] ?? [];
    abil_data[e.ability].push(e);
  });
  let tabs = [];
  for (let abil in abil_data) {
    tabs.push(
      <TabItem key={abil} value={abil} label={labels[abil]}>
        <AbilAoE data={abil_data[abil]} />
      </TabItem>
    );
  }
  if (tabs.length == 0) {
    return <div>No AoE data for {data_src}</div>;
  }
  return <Tabs>{tabs}</Tabs>;
}
