import { StatToIndexMap } from "~src/Components/Character";
import { subDisplayLine } from ".";

import { IconPercent } from "./Icons";
const regDec = new RegExp(/^(\d+)?(\.)?\d+$/);

export function StatRow({
  sub,
  onChange,
}: {
  sub: subDisplayLine;
  onChange: (index: number, value: number) => void;
}) {
  return (
    <div className="flex flex-row place-items-center ">
      <div className="flex-grow flex flex-row items-center mr-auto pt-1 pb-1">
        <div className="w-5 fill-gray-100">{sub.icon}</div>
        <div className="ml-1">{sub.label}</div>
      </div>
      <div className="grid grid-cols-1 gap-y-1 pt-1 pb-1 sm:grid-cols-2 sm:gap-1 sm:gap-y-0  ">
        {sub.stat_ ? (
          <div className="rounded-md flex flex-row focus-within:ring focus-within:border-blue-300">
            <input
              type="number"
              step="any"
              placeholder="enter percentage"
              className="p-2 rounded-l-md bg-gray-800 text-right focus:outline-none invalid:text-red-500"
              value={sub.val_}
              onChange={(e) => {
                const val = e.target.value;
                //first we need to sanitize the value
                if (regDec.test(val)) {
                  e.target.setCustomValidity("");
                  onChange(StatToIndexMap[sub.stat_!], parseFloat(val));
                  return;
                }
                e.target.setCustomValidity("invalid input");
              }}
            />
            <div className="p-1 pr-3 w-6 rounded-r-md bg-gray-800 items-center flex">
              <IconPercent className="fill-gray-100" />
            </div>
          </div>
        ) : (
          <div />
        )}

        {sub.stat ? (
          <div className="rounded-md flex flex-row  ">
            <input
              type="number"
              step="any"
              placeholder="enter amount"
              className="p-2 rounded-l-md bg-gray-800 text-right focus:outline-none invalid:text-red-500"
              value={sub.val}
              onChange={(e) => {
                const val = e.target.value;
                //first we need to sanitize the value
                if (regDec.test(val)) {
                  e.target.setCustomValidity("");
                  onChange(StatToIndexMap[sub.stat!], parseFloat(val) / 100);
                  return;
                }
                e.target.setCustomValidity("invalid input");
              }}
            />
            <div className="p-1 pr-3 w-6 rounded-r-md bg-gray-800 items-center flex" />
          </div>
        ) : (
          <div />
        )}
      </div>
    </div>
  );
}
