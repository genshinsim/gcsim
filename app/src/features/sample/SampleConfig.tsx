import {
  Dialog,
  Classes,
  H4,
  InputGroup,
  HTMLTable,
  Button,
  Tag,
  Intent,
  Alert,
} from "@blueprintjs/core";
import React from "react";
import { setConfig } from "features/sim/simSlice";
import { useDispatch } from "react-redux";
import { xlxqbtfish } from "./data/national";
import { ningandco } from "./data/geo";
import { ganyuAimOnly } from "./data/ganyu";
import { eulaBasic, eulaBennett } from "./data/eula";
import { hutao4TF } from "./data/hutao";
import { bennett4tf, bennett4tfzajef } from "./data/bennett";
import { dilucvape } from "./data/diluc";
import { xiaomango } from "./data/xiao";

export interface PremadeConfig {
  name: string;
  description: string;
  characters: string[];
  tags: string[];
  data: string;
}

function SampleConfig({
  isOpen,
  onClose,
}: {
  isOpen: boolean;
  onClose: () => void;
}) {
  const [openAlert, setOpenAlert] = React.useState<boolean>(false);
  const [cfg, setCfg] = React.useState<string>("");
  const [filter, setFilter] = React.useState<string>("");

  const dispatch = useDispatch();

  let rows = sample.map((ele, i) => {
    if (filter !== "" && !JSON.stringify(ele).includes(filter)) {
      return null;
    }

    let chars = "";
    ele.characters.forEach((c) => (chars += c + " "));

    const tags = ele.tags.map((t) => {
      return (
        <Tag round key={t}>
          {t}
        </Tag>
      );
    });

    return (
      <tr key={i}>
        <td>
          <Button
            icon="download"
            onClick={() => {
              setCfg(ele.data);
              setOpenAlert(true);
            }}
          >
            Load
          </Button>
        </td>
        <td>{ele.name}</td>
        <td>{ele.description}</td>
        <td>{chars}</td>
        <td>{tags}</td>
      </tr>
    );
  });

  return (
    <div>
      <Dialog isOpen={isOpen} onClose={onClose} style={{ width: "60%" }}>
        <div>
          <div className={Classes.DIALOG_HEADER}>
            <H4>Load Sample Config</H4>
          </div>
          <div className={Classes.DIALOG_BODY}>
            <InputGroup
              placeholder="Type here to filter"
              value={filter}
              onChange={(e) => setFilter(e.target.value)}
            />
            <HTMLTable width="100%">
              <thead>
                <tr>
                  <th></th>
                  <th style={{ width: "15%" }}>Name</th>
                  <th style={{ width: "25%" }}>Desc</th>
                  <th style={{ width: "25%" }}>Characters</th>
                  <th style={{ width: "25%" }}>Tags</th>
                </tr>
              </thead>
              <tbody>{rows}</tbody>
            </HTMLTable>
          </div>
        </div>
      </Dialog>
      <Alert
        cancelButtonText="Cancel"
        confirmButtonText="Overwrite"
        icon="trash"
        intent={Intent.DANGER}
        isOpen={openAlert}
        onCancel={() => setOpenAlert(false)}
        onConfirm={() => {
          dispatch(setConfig(cfg));
          setCfg("");
          setOpenAlert(false);
          onClose();
        }}
      >
        <p>
          Are you sure you want to load this config? This will overwrite any
          existing config you have. This operation cannot be reversed.
        </p>
      </Alert>
    </div>
  );
}

export default SampleConfig;

const sample: PremadeConfig[] = [
  xlxqbtfish,
  ningandco,
  ganyuAimOnly,
  eulaBasic,
  eulaBennett,
  hutao4TF,
  bennett4tf,
  bennett4tfzajef,
  dilucvape,
  xiaomango,
];
