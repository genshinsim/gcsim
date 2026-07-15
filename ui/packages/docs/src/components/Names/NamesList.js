import React from "react";
import character_data from "./character.dm.json";
import weapon_data from "./weapon.dm.json";
import artifact_data from "./artifact.dm.json";
import enemy_data from "./monster.dm.json";

export default function NamesList({ item_key, data_src }) {
  let data = character_data;
  switch (data_src) {
    case "weapon":
      data = weapon_data;
      break;
    case "artifact":
      data = artifact_data;
      break;
    case "monster":
      data = enemy_data;
      break;
  }
  if (data[item_key] === undefined) {
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
