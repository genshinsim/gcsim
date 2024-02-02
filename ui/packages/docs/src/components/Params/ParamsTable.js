import React from "react";
import styled from "styled-components";
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

export default function ParamsTable({ item_key, data_src }) {
  let data = character_data;
  switch (data_src) {
    case "weapon":
      data = weapon_data;
      break;
    case "artifact":
      data = artifact_data;
      break;
  }
  if (!(item_key in data) || data[item_key].length === 0) {
    return <div>Does not have any ability params</div>;
  }
  const rows = data[item_key].map((e) => {
    return (
      <tr key={item_key}>
        <TD>
          <code>{e.ability}</code>
        </TD>
        <TD>
          <code>{e.param}</code>
        </TD>
        <TD>{e.desc}</TD>
      </tr>
    );
  });
  return (
    <div style={{ marginTop: "1rem", width: "100%" }}>
      <Table>
        <Thead>
          <tr>
            <TH>Ability</TH>
            <TH>Param</TH>
            <TH>Description</TH>
          </tr>
        </Thead>
        <tbody>{rows}</tbody>
      </Table>
    </div>
  );
}
