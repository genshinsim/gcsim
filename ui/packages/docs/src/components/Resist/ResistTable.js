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

export default function ResistTable({ item_key, data_src }) {
  let data = enemy_data;
  if (!(item_key in data) || data[item_key].length === 0) {
    return <div>No Resist data</div>;
  }
  const rows = Object.entries(data[item_key]).map((e) => {
    const [element, resist] = e;
    return (
      <tr key={item_key}>
        <TD><code>{element}</code></TD>
        <TD>{(resist * 100).toFixed(0)}%</TD>
      </tr>
    );
  });
  return (
    <div style={{ marginTop: "1rem", width: "100%" }}>
      <Table>
        <Thead>
          <tr>
            <TH>Element</TH>
            <TH>Resist (%)</TH>
          </tr>
        </Thead>
        <tbody>{rows}</tbody>
      </Table>
    </div>
  );
}
