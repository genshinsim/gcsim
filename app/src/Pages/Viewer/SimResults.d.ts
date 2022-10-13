export interface SimResults {
  statistics?: Statistics
}

export interface Statistics {
  duration?: FloatStat
  dps?: FloatStat
  rps?: FloatStat
  eps?: FloatStat
  hps?: FloatStat
  sps?: FloatStat
}

export interface FloatStat {
  min?: number
  max?: number
  mean?: number
  sd?: number
  q1?: number
  q2?: number
  q3?: number
}