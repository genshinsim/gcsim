import { LogDetails } from "~src/Components/Viewer/parsev2";

export interface SimResults {
  schema_version?: Version
  max_iterations?: number

  initial_character?: string
  character_details?: CharacterDetail[]

  config_file?: string
  debug?: LogDetails[]
  debug_seed?: string

  statistics?: Statistics
}

export interface Version {
  major: number
  minor: number
}

export interface CharacterDetail {
  name: string
}

export interface Statistics {
  min_seed?: string
  max_seed?: string
  runtime?: number
  iterations?: number

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