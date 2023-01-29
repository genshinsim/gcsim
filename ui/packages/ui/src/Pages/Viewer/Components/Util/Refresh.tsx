import { SimResults } from "@gcsim/types";
import { throttle } from "lodash-es";
import { useRef } from "react";


export function useRefresh<T>(
      getter: (data: SimResults | null) => T, rate: number,
      data: SimResults | null) {
  const refreshFunc = useRef(throttle(getter, rate, { leading: true, trailing: true }));
  return refreshFunc.current(data);
}

// TODO: useRefreshWithTimer
// TODO: RefreshStatus component