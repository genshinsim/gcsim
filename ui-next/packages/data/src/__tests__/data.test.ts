import { describe, expect, it } from "vitest";
import type { LatestCharsMap, TagInfo, TagMap } from "../index";
import { latestChars, tags } from "../index";

describe("@gcsim/data exports", () => {
  it("exports tags as a non-empty object", () => {
    expect(tags).toBeDefined();
    expect(typeof tags).toBe("object");
    expect(Object.keys(tags).length).toBeGreaterThan(0);
  });

  it("exports latestChars as a non-empty object", () => {
    expect(latestChars).toBeDefined();
    expect(typeof latestChars).toBe("object");
    expect(Object.keys(latestChars).length).toBeGreaterThan(0);
  });
});

describe("tags data shape", () => {
  it("each tag has a display_name string", () => {
    for (const [_id, tag] of Object.entries(tags)) {
      expect(tag.display_name).toBeDefined();
      expect(typeof tag.display_name).toBe("string");
    }
  });

  it("blurb is a string when present", () => {
    for (const [_id, tag] of Object.entries(tags)) {
      if (tag.blurb !== undefined) {
        expect(typeof tag.blurb).toBe("string");
      }
    }
  });

  it("contains the gcsim tag at id 1", () => {
    expect(tags["1"]).toBeDefined();
    expect(tags["1"].display_name).toBe("gcsim");
    expect(tags["1"].blurb).toBeDefined();
  });
});

describe("latestChars data shape", () => {
  it("each version maps to an array of strings", () => {
    for (const [_version, chars] of Object.entries(latestChars)) {
      expect(Array.isArray(chars)).toBe(true);
      for (const char of chars) {
        expect(typeof char).toBe("string");
      }
    }
  });

  it("version keys follow v-prefix pattern", () => {
    for (const version of Object.keys(latestChars)) {
      expect(version).toMatch(/^v\d+\.\d+$/);
    }
  });

  it("contains v2.38 with expected characters", () => {
    expect(latestChars["v2.38"]).toBeDefined();
    expect(latestChars["v2.38"]).toContain("varesa");
    expect(latestChars["v2.38"]).toContain("luminepyro");
  });
});

describe("type exports compile correctly", () => {
  it("TagInfo type is usable", () => {
    const info: TagInfo = { display_name: "Test" };
    expect(info.display_name).toBe("Test");
  });

  it("TagMap type is usable", () => {
    const map: TagMap = { "99": { display_name: "Custom" } };
    expect(map["99"].display_name).toBe("Custom");
  });

  it("LatestCharsMap type is usable", () => {
    const map: LatestCharsMap = { "v1.0": ["amber"] };
    expect(map["v1.0"]).toContain("amber");
  });
});
