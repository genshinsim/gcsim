type Props = {
  label: string;
  onChange?: (val: number) => void;
  value: number;
  min?: number;
  max?: number;
  stepSize?: number;
  integerOnly?: boolean;
};

const regDec = new RegExp(/^(\d+)?(\.)?\d+$/);
const regInt = new RegExp(/^(\d+)$/);

export function NumberInput({
  label,
  onChange = (val: number) => {},
  value,
  min = 0,
  max = 100,
  stepSize = 1,
  integerOnly = true,
}: Props) {
  let check = regDec;
  let parse = parseFloat;
  if (integerOnly) {
    check = regInt;
    parse = parseInt;
  }

  return (
    <div className="rounded-md flex flex-row place-items-center">
      <span className="font-bold flex-grow">{label}</span>
      <div className="flex flex-row">
        <input
          type="number"
          step="any"
          placeholder="enter amount"
          className="p-2 rounded-l-md bg-gray-800 text-right focus:outline-none invalid:text-red-500"
          value={value}
          onChange={(e) => {
            const val = e.target.value;
            if (!check.test(val)) {
              // e.target.setCustomValidity("invalid input");
              onChange(min);
              return;
            }
            let v = parse(val);
            if (v > max) {
              v = max;
            }
            if (v < min) {
              v = min;
            }

            e.target.setCustomValidity("");
            onChange(v);
          }}
        />
        <div className="rounded-r-md flex flex-col">
          <button
            className="bg-gray-800 w-12 rounded-tr-md focus:outline-none hover:bg-gray-900"
            disabled={value === max}
            onClick={() => {
              const v = value + stepSize;
              if (v < min || v > max) {
                return;
              }
              onChange(v);
            }}
          >
            <span aria-hidden="true" className="bp3-icon bp3-icon-chevron-up">
              <svg
                data-icon="chevron-up"
                width="16"
                height="16"
                viewBox="0 0 16 16"
              >
                <path
                  d="M12.71 9.29l-4-4C8.53 5.11 8.28 5 8 5s-.53.11-.71.29l-4 4a1.003 1.003 0 001.42 1.42L8 7.41l3.29 3.29c.18.19.43.3.71.3a1.003 1.003 0 00.71-1.71z"
                  fillRule="evenodd"
                ></path>
              </svg>
            </span>
          </button>
          <button
            className="bg-gray-800 w-12 rounded-br-md focus:outline-none hover:bg-gray-900"
            disabled={value === min}
            onClick={() => {
              const v = value - stepSize;
              if (v < min || v > max) {
                return;
              }
              onChange(v);
            }}
          >
            <span aria-hidden="true" className="bp3-icon bp3-icon-chevron-down">
              <svg
                data-icon="chevron-down"
                width="16"
                height="16"
                viewBox="0 0 16 16"
              >
                <path
                  d="M12 5c-.28 0-.53.11-.71.29L8 8.59l-3.29-3.3a1.003 1.003 0 00-1.42 1.42l4 4c.18.18.43.29.71.29s.53-.11.71-.29l4-4A1.003 1.003 0 0012 5z"
                  fillRule="evenodd"
                ></path>
              </svg>
            </span>
          </button>
        </div>
      </div>
    </div>
  );
}
