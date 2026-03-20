import tagsData from "./tags.json";

/**
 * Information about a simulation tag.
 */
export interface TagInfo {
  display_name: string;
  blurb?: string;
  default_exclude?: boolean;
}

/**
 * Map of numeric tag IDs (as strings) to their metadata.
 */
export type TagMap = Record<string, TagInfo>;

/**
 * All known simulation tags, keyed by numeric tag ID.
 */
export const tags: TagMap = tagsData as TagMap;
