import React from "react";
import styled from "styled-components";
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

export default function FieldsTable({ character }) {
  if (!(character in data)) {
    return <div>Character does not have any fields</div>;
  }
  const rows = data[character].map((e) => {
    if (e.fields.length < 0) {
      return;
    }
    const codes = e.fields.map((f, i) => {
      return (
        <div key={i}>
          <code>{f}</code>
        </div>
      );
    });
    return (
      <tr key={character}>
        <TD>{codes}</TD>
        <TD>{e.desc}</TD>
      </tr>
    );
  });
  return (
    <div style={{ marginTop: "1rem", width: "100%" }}>
      <Table>
        <Thead>
          <tr>
            <TH>Field</TH>
            <TH>Description</TH>
          </tr>
        </Thead>
        <tbody>{rows}</tbody>
      </Table>
      <i>If more than one field is available, then either field will work.</i>
    </div>
  );
}
