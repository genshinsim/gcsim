import { FloatStat } from "@gcsim/types";
import { range, unzip } from "lodash-es";
import { useMemo } from "react";

export interface CumulativePoint {
  x: number;
  y: FloatStat[];
}

type ChartData = {
  data: CumulativePoint[];
  keys: number[];
  duration: number;
}

export function useData(input?: FloatStat[][], bucketSize?: number, names?: string[]): ChartData {
  return useMemo(() => {
    if (!input || !bucketSize || !names) {
      return { data: [], keys: [], duration: 1 };
    }

    const data: CumulativePoint[] = unzip(input).map((v, i) => {
      return {
        x: (i * bucketSize) / 60,
        y: v
      };
    });

    if (data.length == 0) {
      return {
        data: [],
        keys: [],
        duration: 1,
      };
    }
    
    const duration = Math.floor(((data.length-1) * bucketSize) / 60);
    if (duration < data[data.length-1].x) {
      data.pop();
    }

    return {
      data: data,
      keys: range(names.length),
      duration: duration,
    };
  }, [input, bucketSize, names]);
}