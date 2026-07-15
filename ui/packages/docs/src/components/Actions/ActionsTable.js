import React from "react";
import styled from "styled-components";
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

export default function ActionsTable({ item_key }) {
  let data = character_data;
  if (!(item_key in data) || data[item_key].length === 0) {
    return <div>Does not have any known legal actions</div>;
  }
  return (
    <div style={{ marginTop: "1rem", width: "100%" }}>
      <Table>
        <Thead>
          <tr>
            <TH>Ability</TH>
            <TH>Legal</TH>
            <TH>Notes</TH>
          </tr>
        </Thead>
        <tbody>
          {data[item_key].map((e) => {
            return (
              <tr key={item_key}>
                <TD><code>{e.ability}</code></TD>
                <TD align="center">{e.invalid ? "❌" : (e.note ? "⚠" : "✔")}</TD>
                <TD>{e.note ?? "-"}</TD>
              </tr>
            );
          })}
        </tbody>
      </Table>
    </div>
  );
}
