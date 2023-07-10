import React from "react";
import data from "./data.json";

export default function NamesList({ character }) {
  if (!(character in data)) {
    return (
      <div>
        Character does not have any known names
      </div>
    );
  }
  if (data[character].length === 0) {
    return (
      <div>
        <ul>
          <li key='0'>{character}</li>
        </ul>
      </div>
    );
  }
  const rows = [character, ...data[character]].map((e, i) => {
    return <li key={i}>{e}</li>;
  });
  return (
    <div>
      <ul>
        {rows}
      </ul>
    </div>
  );
}
