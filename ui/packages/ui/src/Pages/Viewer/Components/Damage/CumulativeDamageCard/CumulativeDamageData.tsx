import { TargetBucket, TargetBucketStats } from "@gcsim/types";
import { useMemo } from "react";

export interface Point {
  x: number;
  y: {
    min: number;
    max: number;
    q1: number;
    q2: number;
    q3: number;
  };
}

type OverTimeData = {
  data: Point[];
  duration: number;
  maxValue: number;
};

// takes two arrays of (potentially) varying length and returns element-wise sum in new array
function pairwiseAdd(a1: number[] | undefined, a2: number[] | undefined) {
  if (a1 == null && a1 == null) {
    return undefined;
  }
  if (a1 == null) {
    return a2;
  }
  if (a2 == null) {
    return a1;
  }
  const maxLength = Math.max(a1.length, a2.length);
  const newA: number[] = [];
  for (let i = 0; i < maxLength; i++) {
    newA[i] = (a1[i] ?? 0) + (a2[i] ?? 0);
  }
  return newA;
}

export function useData(
  graph: string,
  target: string,
  input?: TargetBucketStats
): OverTimeData {
  return useMemo(() => {
    if (input?.bucket_size == null || input.targets == null) {
      return { data: [], duration: 1, maxValue: 1 };
    }
    let targetData: TargetBucket | undefined;
    if (graph === "overall") {
      console.log(input.targets);
      targetData = Object.values(input.targets).reduce((acc, v) => {
        if (v.overall == null) {
          return acc;
        }
        return {
          min: pairwiseAdd(v.overall.min, acc.min),
          max: pairwiseAdd(v.overall.max, acc.max),
          q1: pairwiseAdd(v.overall.q1, acc.q1),
          q2: pairwiseAdd(v.overall.q2, acc.q2),
          q3: pairwiseAdd(v.overall.q3, acc.q3),
        };
      }, {} as TargetBucket);
    } else if (graph === "target") {
      targetData = input.targets[target].target;
    } else {
      return { data: [], duration: 1, maxValue: 1 };
    }
    if (targetData == null) {
      return { data: [], duration: 1, maxValue: 1 };
    }
    const targetMin = targetData.min;
    const targetMax = targetData.max;
    const targetQ1 = targetData.q1;
    const targetQ2 = targetData.q2;
    const targetQ3 = targetData.q3;
    if (
      targetMin == null ||
      targetMax == null ||
      targetQ1 == null ||
      targetQ2 == null ||
      targetQ3 == null
    ) {
      return { data: [], duration: 1, maxValue: 1 };
    }
    const maxBucketIndex = Math.max(
      targetMin.length,
      targetMax.length,
      targetQ1.length,
      targetQ2.length,
      targetQ3.length
    );
    const bucketSize = input.bucket_size;
    let max = 0;
    const data: Point[] = [];
    for (let i = 0; i < maxBucketIndex; i++) {
      const minVal = targetMin[i] ?? targetMin[targetMin.length - 1];
      const maxVal = targetMax[i] ?? targetMax[targetMax.length - 1];
      const q1Val = targetQ1[i] ?? targetQ1[targetQ1.length - 1];
      const q2Val = targetQ2[i] ?? targetQ2[targetQ2.length - 1];
      const q3Val = targetQ3[i] ?? targetQ3[targetQ3.length - 1];
      max = Math.max(max, minVal, maxVal, q1Val, q2Val, q3Val);
      data[i] = {
        x: (i * bucketSize) / 60,
        y: {
          min: minVal,
          max: maxVal,
          q1: q1Val,
          q2: q2Val,
          q3: q3Val,
        },
      };
    }

    if (data.length == 0) {
      return {
        data: [],
        duration: 1,
        maxValue: 1,
      };
    }

    const duration = ((data.length - 1) * bucketSize) / 60;

    return {
      data: data,
      duration: duration,
      maxValue: max,
    };
  }, [graph, input, target]);
}
