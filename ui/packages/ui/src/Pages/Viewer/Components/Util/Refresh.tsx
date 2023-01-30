import { SimResults } from "@gcsim/types";
import { throttle } from "lodash-es";
import { useRef, useState } from "react";

const MAX_JITTER = 50;

// TODO: optional runnning pass to immediately flush if not running?
export function useRefresh<T>(
    getter: (data: SimResults | null) => T,
    rate: number,
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

export function useRefreshWithTimer<T>(
    getter: (data: SimResults | null) => T,
    rate: number,
    data: SimResults | null,
    running = true): [T, number] {
  const [last, setLast] = useState<[T, number] | null>(null);
  const refreshRate = useRef(rate + Math.random() * MAX_JITTER);
  const refreshFunc = useRef(throttle((data: SimResults | null) => {
    return [getter(data), Date.now() + refreshRate.current];
  }, refreshRate.current, { leading: false, trailing: true }));

  if (data == null || !running) {
    return [getter(data), 0];
  }

  if (last == null) {
    const next: [T, number] = [getter(data), Date.now() + refreshRate.current];
    setLast(next);
    return next;
  }

  const next = refreshFunc.current(data);
  return next === undefined ? last : next as [T, number];
}