import React from "react";
import character_data from "./character_data.json";
import weapon_data from "./weapon_data.json";
import artifact_data from "./artifact_data.json";

export default function NamesList({ item_key, data_src }) {
  let data = character_data;
  switch (data_src) {
    case "weapon":
      data = weapon_data;
      break;
    case "artifact":
      data = artifact_data;
      break;
  }
  if (!(item_key in data)) {
    return <div>Does not have any known names</div>;
  }
  if (data[item_key].length === 0) {
    return (
      <div>
        <ul>
          <li key='0'>{item_key}</li>
        </ul>
      </div>
    );
  }
  const rows = [item_key, ...data[item_key]].map((e, i) => {
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
