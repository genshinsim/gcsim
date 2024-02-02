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

export default function FieldsTable({ item_key, data_src }) {
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
    return <div>Does not have any fields</div>;
  }
  const rows = data[item_key].map((e) => {
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
      <tr key={item_key}>
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
