import React from "react";
import styled from "styled-components";
import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";
import character_data from "./character_data.json";

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
      <tr key={e.ability}>
        <TD>{e.ability}</TD>
        <TD>{e.hitHaltTime}</TD>
        <TD>{e.hitHaltTimeScale}</TD>
        <TD>{e.canBeDefenseHalt ? "true" : "false"}</TD>
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

const abils = ["normal", "charge", "aim", "skill", "burst", "asc", "cons"];
const abilLabels = [
  "Normal",
  "Charge Attack",
  "Aimed Shot",
  "Skill",
  "Burst",
  "Ascension",
  "Cons",
];

export default function HitlagTable({ item_key }) {
  if (!(item_key in character_data)) {
    return <div>No hitlag data for character</div>;
  }
  const char = character_data[item_key];
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
        <AbilHitlag data={char[a]} />
      </TabItem>
    );
  });
  if (count == 0) {
    return <div>No hitlag data for character</div>;
  }
  return <Tabs>{tabs}</Tabs>;
}
