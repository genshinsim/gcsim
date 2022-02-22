import { HTMLTable, Spinner } from "@blueprintjs/core";
import axios from "axios";
import React from "react";
import { Link } from "wouter";
import { DBItem } from "~src/types";

export function DB() {
  const [loading, setLoading] = React.useState<boolean>(true);
  const [data, setData] = React.useState<DBItem[]>([]);

  React.useEffect(() => {
    const url = "https://viewer.gcsim.workers.dev/gcsimdb";
    axios
      .get(url)
      .then((resp) => {
        console.log(resp.data);
        let data = resp.data;

        setData(data);
        setLoading(false);
      })
      .catch(function (error) {
        // handle error
        console.log(error);
        setLoading(false);
        setData([]);
      });
  }, []);

  if (loading) {
    return (
      <div className="m-2 text-center text-lg pt-2">
        <Spinner />
        Loading ...
      </div>
    );
  }

  if (data.length === 0) {
    return (
      <div className="m-2 text-center text-lg">
        Error loading database. No data found
      </div>
    );
  }

  const rows = data.map((e, i) => {
    const chars = e.team.map((char) => {
      return (
        <img
          src={"/images/avatar/" + char.name + ".png"}
          alt={char.name}
          className="wide:h-32 h-auto "
        />
      );
    });

    return (
      <tr key={i}>
        <td>
          <div className="grid grid-cols-4">{chars}</div>
        </td>
        <td>{e.author}</td>
        <td className=" text-center align-middle">{e.target_count}</td>
        <td>{e.dps.toFixed(0)}</td>
        <td>
          <Link href={"/viewer/share/" + e.viewer_key}>Viewer</Link>
        </td>
      </tr>
    );
  });

  return (
    <div className="m-2 text-center text-lg flex flex-col place-items-center">
      <HTMLTable>
        <thead>
          <tr>
            <th>Team</th>
            <th>Author</th>
            <th># of Targets</th>
            <th>DPS</th>
            <th>Viewer Link</th>
          </tr>
        </thead>
        <tbody>{rows}</tbody>
      </HTMLTable>
    </div>
  );
}
