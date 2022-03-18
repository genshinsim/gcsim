import {
  Button,
  ButtonGroup,
  Callout,
  Classes,
  Dialog,
  Icon,
  InputGroup,
  Menu,
  MenuItem,
  Spinner,
  Tag,
} from "@blueprintjs/core";
import { Popover2, Tooltip2 } from "@blueprintjs/popover2";
import axios from "axios";
import React from "react";
import { useLocation } from "wouter";
import { CharacterSelect, ICharacter } from "~src/Components/Character";
import { Viewport } from "~src/Components/Viewport";
import { IWeapon, WeaponSelect } from "~src/Components/Weapon";
import { useAppDispatch } from "~src/store";
import { DBCharInfo, DBItem } from "~src/types";
import { simActions } from "../Sim";
import { Trans, useTranslation } from "react-i18next";

function CharTooltip({ char }: { char: DBCharInfo }) {
  let { t } = useTranslation();

  return (
    <div className="m-2 flex flex-col">
      <div className="ml-auto font-bold capitalize">{`${t(
        "character_names." + char.name
      )} ${t("db.c_pre")}${char.con}${t("db.c_post")} ${char.talents.attack}/${
        char.talents.skill
      }/${char.talents.burst}`}</div>
      <div className="w-full border-b border-gray-500 mt-2 mb-2"></div>
      <div className="capitalize flex flex-row">
        <img
          src={"/images/weapons/" + char.weapon + ".png"}
          alt={char.name}
          className="wide:h-8 h-auto "
        />
        <div className="mt-auto mb-auto">
          {t("weapon_names." + char.weapon) + t("db.r") + char.refine}
        </div>
      </div>
      <div className="ml-auto">{`${t("db.er")}${char.er * 100 + 100}%`}</div>
    </div>
  );
}

function TeamCard({ row, setCfg }: { row: DBItem; setCfg: () => void }) {
  useTranslation();

  const [location, setLocation] = useLocation();

  const chars = row.team.map((char) => {
    return (
      <Popover2>
        <Tooltip2 content={<CharTooltip char={char} />}>
          <div className="hover:bg-gray-600 border border-gray-700 hover:border-gray-400 rounded-md relative">
            <img
              src={"/images/avatar/" + char.name + ".png"}
              alt={char.name}
              className="w-16"
              key={char.name}
            />
            <div className=" absolute top-0 right-0 text-sm font-semibold text-grey-300">{`${char.con}`}</div>
          </div>
        </Tooltip2>
      </Popover2>
    );
  });

  return (
    <div className="flex flex-row w-full m-2 p-2 rounded-md bg-gray-700 place-items-center">
      <div className="flex flex-col basis-1/4">
        <div className="grid grid-cols-4">{chars}</div>
        <div>
          <Trans>db.author</Trans>
          {row.author}
        </div>
      </div>
      <div className=" flex-1 overflow-hidden mb-auto pl-2">
        <div className="font-bold">
          <Trans>db.description</Trans>
        </div>
        {row.description.replace(/(.{150})..+/, "$1…")}
      </div>
      <div className="ml-auto flex flex-col mr-4 basis-60">
        <span>
          <Trans>db.total_dps</Trans>
          {parseInt(row.dps.toFixed(0)).toLocaleString()}
        </span>
        <span>
          <Trans>db.number_of_targets</Trans>
          {row.target_count}
        </span>
        <span>
          <Trans>db.average_dps_per</Trans>
          {parseInt((row.dps / row.target_count).toFixed(0)).toLocaleString()}
        </span>
        <span>
          <Trans>db.hash</Trans>
          <a href={"https://github.com/genshinsim/gcsim/commit/" + row.hash}>
            {row.hash.substring(0, 8)}
          </a>
        </span>
      </div>
      <div>
        <ButtonGroup vertical>
          <Button
            small
            rightIcon="chart"
            onClick={() => {
              setLocation("/viewer/share/" + row.viewer_key);
            }}
          >
            <Trans>db.show_in_viewer</Trans>
          </Button>
          <Button small rightIcon="rocket-slant" onClick={setCfg}>
            <Trans>db.load_in_simulator</Trans>
          </Button>
          <Button
            disabled
            small
            rightIcon="list-detail-view"
            onClick={() => {
              console.log("i do nothing yet");
            }}
          >
            <Trans>db.details</Trans>
          </Button>
        </ButtonGroup>
      </div>
    </div>
  );
}

export function DB() {
  let { t } = useTranslation();

  const [loading, setLoading] = React.useState<boolean>(true);
  const [data, setData] = React.useState<DBItem[]>([]);
  const [openAddChar, setOpenAddChar] = React.useState<boolean>(false);
  const [charFilter, setCharFilter] = React.useState<string[]>([]);
  const [openAddWeap, setOpenAddWeap] = React.useState<boolean>(false);
  const [weapFilter, setWeapFilter] = React.useState<string[]>([]);
  const [searchString, setSearchString] = React.useState<string>("");
  const [cfg, setCfg] = React.useState<string>("");

  const dispatch = useAppDispatch();
  const [location, setLocation] = useLocation();

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

  const openInSim = () => {
    dispatch(simActions.setAdvCfg(cfg));
    setLocation("/advanced");
    setCfg("");
  };

  const addCharFilter = (char: ICharacter) => {
    setOpenAddChar(false);
    //add to array if not exist already if
    if (charFilter.includes(char.key)) {
      return;
    }
    const next = [...charFilter];
    next.push(char.key);
    setCharFilter(next);
  };

  const removeCharFilter = (char: string) => {
    const idx = charFilter.indexOf(char);
    if (idx === -1) {
      return;
    }
    const next = [...charFilter];
    next.splice(idx, 1);
    setCharFilter(next);
  };

  const addWeapFilter = (weap: IWeapon) => {
    setOpenAddWeap(false);
    //add to array if not exist already if
    if (weapFilter.includes(weap)) {
      return;
    }
    const next = [...weapFilter];
    next.push(weap);
    setWeapFilter(next);
  };

  const removeWeapFilter = (weap: string) => {
    const idx = weapFilter.indexOf(weap);
    if (idx === -1) {
      return;
    }
    const next = [...weapFilter];
    next.splice(idx, 1);
    setWeapFilter(next);
  };

  if (loading) {
    return (
      <div className="m-2 text-center text-lg pt-2">
        <Spinner />
        <Trans>db.loading</Trans>
      </div>
    );
  }

  if (data.length === 0) {
    return (
      <div className="m-2 text-center text-lg">
        <Trans>db.error_loading_database</Trans>
      </div>
    );
  }
  data.sort((a, b) => {
    return b.dps / b.target_count - a.dps / a.target_count;
  });

  const cRows = charFilter.map((e) => {
    return (
      <Tag
        key={e}
        interactive
        onRemove={() => removeCharFilter(e)}
        className="ml-px mr-px"
      >
        {t("character_names." + e)}
      </Tag>
    );
  });

  const wRows = weapFilter.map((e) => {
    return (
      <Tag
        key={e}
        interactive
        onRemove={() => removeWeapFilter(e)}
        className="ml-px mr-px"
      >
        {t("weapon_names." + e)}
      </Tag>
    );
  });

  //filter data

  const n = data.filter((e) => {
    const team: string[] = [];
    const weapons: string[] = [];

    e.team.forEach((char) => {
      team.push(char.name);
      weapons.push(char.weapon);
    });

    //team needs to have every character in charFilter array
    if (charFilter.length > 0) {
      const ok = charFilter.every((e) => team.includes(e));
      if (!ok) {
        return false;
      }
    }

    //team needs to have every weapon in weaponFilter array
    if (weapFilter.length > 0) {
      const ok = weapFilter.every((e) => weapons.includes(e));
      if (!ok) {
        return false;
      }
    }

    //check something in team matches search string
    let ss = JSON.stringify(e);
    e.team.forEach((c) => {
      ss += " " + t("character_names." + c.name);
      ss += " " + t("weapon_names." + c.weapon);
    });

    if (searchString !== "" && !ss.includes(searchString)) {
      return false;
    }

    return true;
  });

  const rows = n.map((e, i) => {
    return <TeamCard row={e} key={i} setCfg={() => setCfg(e.config)} />;
  });

  return (
    <Viewport>
      <div className="flex flex-row items-center">
        <div className="flex flex-row items-center">
          <Icon icon="filter-list" /> <Trans>db.filters</Trans>{" "}
          <Popover2
            interactionKind="click"
            placement="bottom"
            content={
              <Menu>
                <MenuItem
                  text={t("db.character")}
                  onClick={() => setOpenAddChar(true)}
                />
                <MenuItem
                  text={t("db.weapon")}
                  onClick={() => setOpenAddWeap(true)}
                />
              </Menu>
            }
            renderTarget={({ isOpen, ref, ...targetProps }) => (
              <Button
                {...targetProps}
                //@ts-ignore
                elementRef={ref}
                icon="plus"
                className="ml-1 mr-1"
              />
            )}
          />
          <div>
            {cRows}
            {wRows}
          </div>
        </div>
        <div className="ml-auto">
          <InputGroup
            leftIcon="search"
            placeholder={t("db.type_to_search")}
            value={searchString}
            onChange={(e) => setSearchString(e.target.value)}
          ></InputGroup>
        </div>
      </div>
      <div className="border-b-2 mt-2 border-gray-300" />
      <div className="p-2 flex flex-col place-items-center w-full">{rows}</div>
      <CharacterSelect
        onClose={() => setOpenAddChar(false)}
        onSelect={addCharFilter}
        isOpen={openAddChar}
      />
      <WeaponSelect
        isOpen={openAddWeap}
        onClose={() => setOpenAddWeap(false)}
        onSelect={addWeapFilter}
      />
      <Dialog isOpen={cfg !== ""} onClose={() => setCfg("")}>
        <div className={Classes.DIALOG_BODY}>
          Load this configuration in <span className="font-bold">Advanced</span>{" "}
          mode.
          <Callout intent="warning" className="mt-2">
            This will overwrite any existing configuration you may have. Are you
            sure you wish to continue?
          </Callout>
        </div>

        <div className={Classes.DIALOG_FOOTER}>
          <div className={Classes.DIALOG_FOOTER_ACTIONS}>
            <Button onClick={openInSim} intent="primary">
              <Trans>db.continue</Trans>
            </Button>
            <Button onClick={() => setCfg("")}>
              <Trans>db.cancel</Trans>
            </Button>
          </div>
        </div>
      </Dialog>
    </Viewport>
  );
}
