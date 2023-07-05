import React from "react";
import data from "./data.json";

export default function FieldsTable({ character }) {
  if (!(character in data)) {
    return <div>Character does not have any known issues</div>;
  }
  if (data[character].length === 0) {
    return <div>Character does not have any known issues</div>;
  }
  const rows = data[character].map((e, i) => {
    return <li key={i}>{e}</li>;
  });
  return (
    <div>
      <ul>{rows}</ul>
    </div>
  );
}
