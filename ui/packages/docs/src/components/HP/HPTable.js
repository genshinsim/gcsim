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

export default function HPTable({ item_key, data_src }) {
  let data = enemy_data;
  if (!(item_key in data) || data[item_key].length === 0) {
    return <div>No HP data</div>;
  }
  const rows = data[item_key].map((e) => {
    return (
      <tr key={item_key}>
        <TD>{e.level}</TD>
        <TD>{e.hp.toFixed()}</TD>
      </tr>
    );
  });
  return (
    <div style={{ marginTop: "1rem", width: "100%" }}>
      <Table>
        <Thead>
          <tr>
            <TH>Level</TH>
            <TH>HP</TH>
          </tr>
        </Thead>
        <tbody>{rows}</tbody>
      </Table>
    </div>
  );
}
