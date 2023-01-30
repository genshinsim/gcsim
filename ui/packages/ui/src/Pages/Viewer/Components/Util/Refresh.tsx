import { SimResults } from "@gcsim/types";
import { throttle } from "lodash-es";
import { useRef, useState } from "react";

const MAX_JITTER = 50;

export function useRefresh<T>(
      getter: (data: SimResults | null) => T, rate: number,
      data: SimResults | null): T {
  const refreshFunc = useRef(throttle(
      getter, rate + Math.random() * MAX_JITTER, { leading: false, trailing: true }));
  const [last, setLast] = useState<T | null>(null);

  if (data == null) {
    return getter(data);
  }

  if (last == null) {
    const next = getter(data);
    setLast(next);
    return next;
  }

  const next = refreshFunc.current(data);
  return next === undefined ? last : next;
}

// TODO: useRefreshWithTimer
// TODO: RefreshStatus component