import { IconPercent } from "../../../Components/Icons";
import { StatToIndexMap } from "../Components/util";

const regDec = new RegExp(/^(\d+)?(\.)?\d+$/);

export type subDisplayLine = {
  stat?: string;
  stat_?: string;
  label: string;
  val: number;
  val_: number;
  icon: React.ReactElement;
};

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
        <div className="w-4 fill-gray-100">{sub.icon}</div>
        <div className="ml-1">{sub.label}</div>
      </div>
      <div className="grid grid-cols-1 gap-y-1 pt-1 pb-1 hd:grid-cols-2 md:gap-1 md:gap-y-0 md:w-[400px]">
        {sub.stat ? (
          <div className="rounded-md flex flex-row">
            <input
              type="number"
              step="any"
              placeholder="enter amount"
              className="w-full p-2 rounded-l-md bg-gray-800 text-right focus:outline-none invalid:text-red-500"
              value={sub.val}
              onChange={(e) => {
                let val: number;
                //first we need to sanitize the value
                if (regDec.test(e.target.value)) {
                  // e.target.setCustomValidity("");
                  val = parseFloat(e.target.value);
                } else {
                  val = 0;
                }
                // e.target.setCustomValidity("invalid input");

                onChange(StatToIndexMap[sub.stat!], val);
              }}
            />
            <div className="p-1 pr-3 w-6 rounded-r-md bg-gray-800 items-center flex" />
          </div>
        ) : (
          <div />
        )}

        {sub.stat_ ? (
          <div className="rounded-md flex flex-row focus-within:ring focus-within:border-blue-300">
            <input
              type="number"
              step="any"
              placeholder="enter percentage"
              className="w-full p-2 rounded-l-md bg-gray-800 text-right focus:outline-none invalid:text-red-500"
              value={sub.val_}
              onChange={(e) => {
                let val: number;
                //first we need to sanitize the value
                if (regDec.test(e.target.value)) {
                  // e.target.setCustomValidity("");
                  val = parseFloat(e.target.value);
                } else {
                  val = 0;
                }
                // e.target.setCustomValidity("invalid input");

                onChange(StatToIndexMap[sub.stat_!], val / 100);
              }}
            />
            <div className="p-1 pr-3 w-6 rounded-r-md bg-gray-800 items-center flex">
              <IconPercent className="fill-gray-100" />
            </div>
          </div>
        ) : (
          <div />
        )}
      </div>
    </div>
  );
}
