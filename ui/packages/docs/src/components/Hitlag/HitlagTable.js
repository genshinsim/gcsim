import React from "react";
import styled from "styled-components";
import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";
import character_data from "./character.dm.json";

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

function AbilHitlag({ data }) {
  const rows = data.map((e) => {
    return (
      <tr key={e.name}>
        <TD>{e.name}</TD>
        <TD>{e.time}</TD>
        <TD>{e.scale}</TD>
        <TD>{e.defense_halt ? "true" : "false"}</TD>
        <TD>{e.deployable ? "true" : "false"}</TD>
      </tr>
    );
  });

  return (
    <div style={{ marginTop: "1rem", width: "100%" }}>
      <Table>
        <Thead>
          <tr>
            <TH>Ability</TH>
            <TH>Halt Time</TH>
            <TH>Scale</TH>
            <TH>Defense Halt</TH>
            <TH>Deployable</TH>
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

export default function HitlagTable({ item_key }) {
  let data = character_data;
  if (!(item_key in data)) {
    return <div>No hitlag data for character</div>;
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
        <AbilHitlag data={abil_data[abil]} />
      </TabItem>
    );
  }
  if (tabs.length == 0) {
    return <div>No hitlag data for character</div>;
  }
  return <Tabs>{tabs}</Tabs>;
}
