import { SampleItemView } from "./SampleItemView";
import { SampleItem, SampleRow } from "./parse";
import { useVirtual } from "react-virtual";
import AutoSizer from "react-virtualized-auto-sizer";
import React, { useState } from "react";
import { Button, ButtonGroup, Card, FormGroup, InputGroup, Intent } from "@blueprintjs/core";
import { Sample } from "@gcsim/types";
import { saveAs } from "file-saver";
import { AdvancedPreset, AllSampleOptions, DebugPreset, DefaultSampleOptions, SimplePreset, VerbosePreset } from "./SampleOptions";
import { Options } from "./Options";
import { bytesToBase64 } from "@gcsim/utils";
import Pako from "pako";

type buffSetting = {
  start: number;
  end: number;
  show: boolean;
};

const Row = ({
  row,
  highlight,
  showBuffDuration,
}: {
  row: SampleRow;
  highlight: buffSetting;
  showBuffDuration: (e: SampleItem) => void;
}) => {
  const cols = row.slots.map((slot, ci) => {
    const events = slot.map((e, ei) => {
      return (
        <SampleItemView item={e} key={ei} showBuffDuration={showBuffDuration} />
      );
    });

    return (
      <div
        key={ci}
        className={
          row.active == ci
            ? "border-l-2 border-gray-500 bg-gray-400	"
            : "border-l-2 border-gray-500"
        }
      >
        {events}
      </div>
    );
  });

  const hl =
    highlight.show && row.f >= highlight.start && row.f <= highlight.end;

  //map out each col
  return (
    <div className="flex flex-row" key={row.key}>
      <div
        className={
          hl
            ? "text-right text-gray-100 border-b-2 border-gray-500 bg-blue-500"
            : "text-right text-gray-100 border-b-2 border-gray-500"
        }
        style={{ minWidth: "100px" }}
      >
        <div>{`${row.f} | ${(row.f / 60).toFixed(2)}s`}</div>
      </div>
      <div className="grid grid-cols-5 flex-grow border-b-2 border-gray-500">
        {cols}
      </div>
      <div style={{ width: "20px", minWidth: "20px" }} />
    </div>
  );
};

type SampleOptionsProps = {
  settings: string[];
  setSettings: (val: string[]) => void;
}

const SampleOptions = ({ settings, setSettings }: SampleOptionsProps) => {
  const [isOpen, setOpen] = useState(false);

  const toggle = (t: string) => {
    const i = settings.indexOf(t);
    const next = [...settings];
    if (i === -1) {
      next.push(t);
    } else {
      next.splice(i, 1);
    }
    setSettings(next);
  };

  const presets = (opt: "simple" | "advanced" | "verbose" | "debug") => {
    switch (opt) {
      case "simple":
        setSettings(SimplePreset);
        return;
      case "advanced":
        setSettings(AdvancedPreset);
        return;
      case "verbose":
        setSettings(VerbosePreset);
        return;
      case "debug":
        setSettings(DebugPreset);
        return;
    }
  };

  return (
    <>
      <Button
          onClick={() => setOpen(true)}
          icon="cog"
          text={"Settings"} />
      <Options
          isOpen={isOpen}
          handleClose={() => setOpen(false)}
          handleClear={() => setSettings([])}
          handleResetDefault={() => setSettings(DefaultSampleOptions)}
          handleToggle={toggle}
          handleSetPresets={presets}
          selected={settings}
          options={AllSampleOptions} />
    </>
  );
};

let lastSearchIndex = 0;

type SamplerProps = {
  sample: Sample;
  data: SampleRow[];
  team: string[];
  searchable: { [key: number]: string[] };
  settings: string[];
  setSettings: (val: string[]) => void;
}

function SamplerUI({ sample, data, team, searchable, settings, setSettings }: SamplerProps) {
  const parentRef = React.useRef<HTMLDivElement>(null);
  const searchRef = React.useRef<HTMLInputElement>(null);
  const [hl, sethl] = React.useState<buffSetting>({
    start: 0,
    end: 0,
    show: false,
  });

  const handleShowBuffDuration = (e: SampleItem) => {
    // const show = hl.show;
    const next = {
      show: true,
      start: e.added,
      end: e.ended,
    };
    sethl(next);
  };

  const rowVirtualizer = useVirtual({
    size: data.length,
    parentRef,
    keyExtractor: React.useCallback(
      (index: number) => {
        return data[index].f;
      },
      [data]
    ),
  });

  const char = team.map((c) => {
    return (
      <div
        key={c}
        className="capitalize text-lg font-medium text-gray-100 border-l-2 border-b-2 border-gray-500"
      >
        {c}
      </div>
    );
  });

  const searchAndScroll = (val: string) => {
    const total = Object.keys(searchable).length;
    for (let index = lastSearchIndex; index < total; index++) {
      for (const msg of searchable[index]) {
        const lowerMsg = msg.toLowerCase()
        if (lowerMsg.indexOf(val.toLowerCase()) > -1) {
          console.log(index, lastSearchIndex);
          lastSearchIndex = index + 1;
          rowVirtualizer.scrollToIndex(index, { align: "start" });
          return;
        }
      }
    }
  };

  return (
    <div className="flex flex-col overflow-x-auto h-[80vh]">
      <div className="flex justify-between">
        <FormGroup label="Search" inline>
          <InputGroup
            type="text"
            inputRef={searchRef}
            rightElement={
              <FormGroup>
                <Button
                  icon="arrow-down"
                  intent="warning"
                  onClick={() => {
                    if (searchRef.current != null) {
                      searchAndScroll(searchRef.current.value);
                    }
                  }}
                />
                <Button
                  icon="reset"
                  intent="warning"
                  onClick={() => {
                    if (searchRef.current != null) {
                      searchRef.current.value = "";
                    }
                    lastSearchIndex = 0;
                    rowVirtualizer.scrollToIndex(0);
                  }}
                />
              </FormGroup>
            }
          />
        </FormGroup>
        <ButtonGroup className="mb-[15px]">
          <SampleOptions settings={settings} setSettings={setSettings} />
          <Button
              icon="bring-data"
              text="Download"
              intent={Intent.SUCCESS}
              onClick={() => {
                const out = bytesToBase64(Pako.deflate(JSON.stringify(sample)));
                const blob = new Blob([out], { type: "application/base64" });
                saveAs(blob, "sample");
              }} />
        </ButtonGroup>
      </div>
      <Card className="flex-auto !bg-gray-600 !text-xs min-w-[60rem] ">
        <AutoSizer disableWidth={true}>
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
              <div className="flex flex-row sample-header">
                <div
                  className={
                    "font-medium text-lg text-gray-100 border-b-2 border-gray-500 text-right "
                  }
                  style={{ minWidth: "100px" }}
                >
                  F | Sec
                </div>
                <div className="grid grid-cols-5 flex-grow">
                  <div className="font-medium text-lg text-gray-100 border-l-2 border-b-2 border-gray-500">
                    Sim
                  </div>
                  {char}
                </div>
                <div style={{ width: "20px", minWidth: "20px" }} />
              </div>
              <div
                className="ListInner"
                style={{
                  // Set the scrolling inner div of the parent to be the
                  // height of all items combined. This makes the scroll bar work.
                  height: `${rowVirtualizer.totalSize}px`,
                  width: "100%",
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
                    <Row
                      row={data[virtualRow.index]}
                      highlight={hl}
                      showBuffDuration={handleShowBuffDuration}
                    />
                  </div>
                ))}
              </div>
            </div>
          )}
        </AutoSizer>
      </Card>
    </div>
  );
}

export const Sampler = React.memo(SamplerUI);