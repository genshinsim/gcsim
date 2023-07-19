import { BucketStats, FloatStat } from "@gcsim/types";
import { useMemo } from "react";

export interface Point {
  x: number;
  y: FloatStat;
}

type OverTimeData = {
  data: Point[];
  duration: number;
  maxValue: number;
}

export function useData(input?: BucketStats): OverTimeData {
  return useMemo(() => {
    if (input?.bucket_size == null || input.buckets == null) {
      return { data: [], duration: 1, maxValue: 1 };
    }

    const bucketSize = input.bucket_size;
    let max = 0;
    const data: Point[] = input.buckets.map((v, i) => {
      max = Math.max(max, v.max ?? 0);
      return {
        x: (i * bucketSize) / 60,
        y: v,
      };
    });

    if (data.length == 0) {
      return {
        data: [],
        duration: 1,
        maxValue: 1,
      };
    }

    const duration = Math.floor(((data.length-1) * bucketSize) / 60);
    if (duration < data[data.length-1].x) {
      data.pop();
    }

    return {
      data: data,
      duration: duration,
      maxValue: max,
    };
  }, [input]);
}