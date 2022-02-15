import { DebugItemView } from "./DebugItemView";
import { DebugRow } from "./parse";
import { useVirtual } from "react-virtual";
import AutoSizer from "react-virtualized-auto-sizer";
import React from "react";

const Row = ({ row }: { row: DebugRow }) => {
  const cols = row.slots.map((slot, ci) => {
    const events = slot.map((e, ei) => {
      return <DebugItemView item={e} key={ei} />;
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

  //map out each col
  return (
    <div className="flex flex-row" key={row.key}>
      <div
        className="text-right text-gray-100 border-b-2 border-gray-500"
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

export function Debugger({ data, team }: { data: DebugRow[]; team: string[] }) {
  const parentRef = React.useRef<HTMLDivElement>(null!);

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

  return (
    <div className="h-full m-2 p-2 rounded-md bg-gray-600 text-xs flex flex-col min-w-[60rem] min-h-[20rem]">
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
            <div className="flex flex-row debug-header">
              <div
                className="font-medium text-lg text-gray-100 border-b-2 border-gray-500 text-right"
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
                  <Row row={data[virtualRow.index]} />
                </div>
              ))}
            </div>
          </div>
        )}
      </AutoSizer>
    </div>
  );
}
