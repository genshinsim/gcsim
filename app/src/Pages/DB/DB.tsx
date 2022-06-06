import {
  AnchorButton,
  Button,
  ButtonGroup,
  Callout,
  Checkbox,
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
import { Link, useLocation } from "wouter";
import {
  CharacterSelect,
  ICharacter,
  isTraveler,
} from "~src/Components/Character";
import { Viewport } from "~src/Components/Viewport";
import { IWeapon, WeaponSelect } from "~src/Components/Weapon";
import { useAppDispatch } from "~src/store";
import { DBCharInfo, DBItem } from "~src/types";
import { updateCfg } from "../Sim";
import { Trans, useTranslation } from "react-i18next";
import { Disclaimer } from "./Disclaimer";
import { useVirtual } from "react-virtual";
import AutoSizer from "react-virtualized-auto-sizer";

function CharTooltip({ char }: { char: DBCharInfo }) {
  let { t } = useTranslation();

  return (
    <div className="m-2 flex flex-col">
      <div className="ml-auto font-bold capitalize">{`${t(
        "game:character_names." + char.name
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
          {t("game:weapon_names." + char.weapon) + t("db.r") + char.refine}
        </div>
      </div>
      <div className="ml-auto">{`${t("db.er")}${char.er * 100 + 100}%`}</div>
    </div>
  );
}

function TeamCard({ row, setCfg }: { row: DBItem; setCfg: () => void }) {
  useTranslation();
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
    <div className="flex flex-row flex-wrap sm:flex-nowrap gap-y-1 w-full m-2 p-2 rounded-md bg-gray-700 place-items-center">
      <div className="flex flex-col sm:basis-1/4 xs:basis-full">
        <div className="grid grid-cols-4">{chars}</div>
        <div className="hidden basis-0 lg:block md:flex-1">
          <Trans>db.author</Trans>
          {row.author}
        </div>
      </div>
      <div className=" flex-1 overflow-hidden mb-auto pl-2 hidden lg:block">
        <div className="font-bold">
          <Trans>db.description</Trans>
        </div>
        {row.description.replace(/(.{150})..+/, "$1…")}
      </div>
      <div className="ml-auto flex flex-col mr-4 md:basis-60 basis-full">
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
          <a href={"https://github.com/genshinsim/gcsim/commits/" + row.hash}>
            {row.hash.substring(0, 8)}
          </a>
        </span>
      </div>
      <div>
        <ButtonGroup vertical>
          <Link href={"/viewer/share/" + row.viewer_key}>
            <AnchorButton small rightIcon="chart">
              <Trans>db.show_in_viewer</Trans>
            </AnchorButton>
          </Link>
          <Button small rightIcon="rocket-slant" onClick={setCfg}>
            <Trans>db.load_in_simulator</Trans>
          </Button>
          {/* <Button
            disabled
            small
            rightIcon="list-detail-view"
            onClick={() => {
              console.log("i do nothing yet");
            }}
          >
            <Trans>db.details</Trans>
          </Button> */}
        </ButtonGroup>
      </div>
    </div>
  );
}

const LOCALSTORAGE_KEY = "gcsim-viewer-cpy-cfg-settings";
const LOCALSTORAGE_DISC_KEY = "gcsim-db-disclaimer-show";

type DBViewProps = {
  db: DBItem[];
  setCfg: (cfg: string) => void;
};

export function DBView(props: DBViewProps) {
  const parentRef = React.useRef<HTMLDivElement>(null!);
  const rowVirtualizer = useVirtual({
    size: props.db.length,
    parentRef,
    keyExtractor: React.useCallback(
      (index: number) => {
        return index;
      },
      [props.db]
    ),
  });

  return (
    <div className="h-full w-full pl-2 pr-2">
      <AutoSizer defaultHeight={100}>
        {({ height, width }) => (
          <div
            ref={parentRef}
            style={{
              minHeight: "100px",
              height: height,
              width: width,
              overflow: "auto",
              position: "relative",
            }}
            id="resize-inner"
          >
            <div
              className="ListInner"
              style={{
                // Set the scrolling inner div of the parent to be the
                // height of all items combined. This makes the scroll bar work.
                height: `${rowVirtualizer.totalSize}px`,
                width: width - 50,
                position: "relative",
              }}
            >
              {
                // The meat and potatoes, an array of the virtual items
                // we currently want to render and their index in the original data.
              }
              {rowVirtualizer.virtualItems.map((virtualRow) => (
                <div
                  key={virtualRow.index}
                  // ref={virtualRow.measureRef}
                  ref={(el) => virtualRow.measureRef(el)}
                  style={{
                    position: "absolute",
                    top: 0,
                    left: 0,
                    width: "100%",
                    // Positions the virtual elements at the right place in container.
                    // minHeight: `${virtualRow.size - 10}px`,
                    transform: `translateY(${virtualRow.start}px)`,
                  }}
                  // id={"virtual-row-"+virtualRow.key}
                >
                  <TeamCard
                    row={props.db[virtualRow.index]}
                    key={virtualRow.index}
                    setCfg={() =>
                      props.setCfg(props.db[virtualRow.index].config)
                    }
                  />
                </div>
              ))}
            </div>
          </div>
        )}
      </AutoSizer>
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
  const [keepExistingTeam, setKeepExistingTeam] = React.useState<boolean>(
    () => {
      const saved = localStorage.getItem(LOCALSTORAGE_KEY);
      if (saved === "true") {
        return true;
      }
      return false;
    }
  );
  const [showDisclaimer, setShowDisclaimer] = React.useState<boolean>(() => {
    const saved = localStorage.getItem(LOCALSTORAGE_DISC_KEY);
    if (saved === "false") {
      return false;
    }
    return true;
  });

  const dispatch = useAppDispatch();
  const [_, setLocation] = useLocation();

  React.useEffect(() => {
    const url = "https://viewer.gcsim.workers.dev/gcsimdb";
    axios
      .get(url)
      .then((resp) => {
        console.log(resp.data);
        let data = resp.data;

        setData(data);
        parseFilterUrl();
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
    dispatch(updateCfg(cfg, keepExistingTeam));
    setLocation("/simulator");
    setCfg("");
  };

  const hideDisclaimer = () => {
    localStorage.setItem(LOCALSTORAGE_DISC_KEY, "false");
    setShowDisclaimer(false);
  };

  const updateFilterUrl = (type: string, data: Array<string>) => {
    const url = new URL(window.location.toString());
    url.searchParams.set(type, data.join(","));
    window.history.pushState({}, "", url);
  };

  const parseFilterUrl = () => {
    const url = new URL(window.location.toString());
    const chars = getSearchParamData(url, "chars");
    const weaps = getSearchParamData(url, "weaps");

    console.log(chars);

    if (chars) setCharFilter(chars);
    if (weaps) setWeapFilter(weaps);
  };

  const getSearchParamData = (url: URL, key: string) => {
    const data = url.searchParams.get(key);
    return data ? data.split(",") : null;
  };

  const addCharFilter = (char: ICharacter) => {
    setOpenAddChar(false);
    //add to array if not exist already if
    let key = char.key;
    if (isTraveler(key) && char.element != "none")
      key = "traveler" + char.element;
    if (charFilter.includes(key)) {
      return;
    }
    const next = [...charFilter];
    next.push(key);
    setCharFilter(next);

    updateFilterUrl("chars", next);
  };

  const removeCharFilter = (char: string) => {
    const idx = charFilter.indexOf(char);
    if (idx === -1) {
      return;
    }
    const next = [...charFilter];
    next.splice(idx, 1);
    setCharFilter(next);

    updateFilterUrl("chars", next);
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

    updateFilterUrl("weaps", next);
  };

  const removeWeapFilter = (weap: string) => {
    const idx = weapFilter.indexOf(weap);
    if (idx === -1) {
      return;
    }
    const next = [...weapFilter];
    next.splice(idx, 1);
    setWeapFilter(next);

    updateFilterUrl("weaps", next);
  };

  const handleToggleSelected = () => {
    localStorage.setItem(LOCALSTORAGE_KEY, keepExistingTeam ? "false" : "true");
    setKeepExistingTeam(!keepExistingTeam);
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
        {t("game:character_names." + e)}
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
        {t("game:weapon_names." + e)}
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
      ss += " " + t("game:character_names." + c.name);
      ss += " " + t("game:weapon_names." + c.weapon);
    });

    if (searchString !== "" && !ss.includes(searchString)) {
      return false;
    }

    return true;
  });

  return (
    <main className="flex flex-col h-full m-2 w-full xs:w-full sm:w-[640px] hd:w-full wide:w-[1160px] ml-auto mr-auto ">
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
        <div className="ml-auto flex flex-row gap-x-1">
          <Button
            intent="primary"
            onClick={() => {
              localStorage.setItem(LOCALSTORAGE_DISC_KEY, "true");
              setShowDisclaimer(true);
            }}
          >
            Show FAQs
          </Button>
          <InputGroup
            leftIcon="search"
            placeholder={t("db.type_to_search")}
            value={searchString}
            onChange={(e) => setSearchString(e.target.value)}
          ></InputGroup>
        </div>
      </div>
      <div className="border-b-2 mt-2 border-gray-300" />
      <div className="p-2 grow ">
        <DBView db={n} setCfg={setCfg} />
      </div>
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
          <Trans>viewer.load_this_configuration</Trans>
          <Callout intent="warning" className="mt-2">
            <Trans>viewer.this_will_overwrite</Trans>
          </Callout>
          <Checkbox
            label="Copy action list only (ignore character stats)"
            className="mt-2"
            checked={keepExistingTeam}
            onClick={handleToggleSelected}
          />
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
      <Disclaimer
        isOpen={showDisclaimer}
        onClose={() => setShowDisclaimer(false)}
        hideAlways={hideDisclaimer}
      />
    </main>
  );
}
