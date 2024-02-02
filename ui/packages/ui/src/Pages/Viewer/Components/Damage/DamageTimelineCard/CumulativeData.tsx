import { CharacterBucketStats, FloatStat } from "@gcsim/types";
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

export function useData(input?: CharacterBucketStats, names?: string[]): ChartData {
  return useMemo(() => {
    if (!input?.characters || input.bucket_size == null || !names) {
      return { data: [], keys: [], duration: 1 };
    }

    const bucket_size = input.bucket_size;
    const points = input.characters.map(x => x.buckets);

    const data: CumulativePoint[] = unzip(points).map((v, i) => {
      return {
        x: (i * bucket_size) / 60,
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
    
    const duration = ((data.length - 1) * bucket_size) / 60;

    return {
      data: data,
      keys: range(names.length),
      duration: duration,
    };
  }, [input, names]);
}