import React from "react";
import styled from "styled-components";
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

const actions = [
  "attack",
  "charge",
  "aim",
  "skill",
  "burst",
  "low_plunge",
  "high_plunge",
  "dash",
  "jump",
  "walk",
  "swap",
];

export default function ActionsTable({ item_key }) {
  if (!(item_key in character_data) || character_data[item_key].length === 0) {
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
          {actions.map((a) => {
            const action = character_data[item_key].find(
              (item) => item.ability === a
            );
            return (
              <tr key={item_key}>
                <TD><code>{a}</code></TD>
                <TD align="center">
                  {action
                    ? action.legal === undefined
                      ? "⚠"
                      : action.legal
                      ? "✔"
                      : "❌"
                    : "❌"}
                </TD>
                <TD>{action?.notes ?? "-"}</TD>
              </tr>
            );
          })}
        </tbody>
      </Table>
    </div>
  );
}
