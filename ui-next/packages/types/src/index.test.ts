import { describe, expect, it } from "vitest";
import type { Sim } from "./index.js";
import {
  CharacterSchema,
  DescriptiveStatsSchema,
  db,
  EnemySchema,
  OverviewStatsSchema,
  WeaponSchema as ProtoWeaponSchema,
  preview,
  SampleSchema,
  SimulationResultSchema,
  SimulationStatisticsSchema,
  share,
  VersionSchema,
} from "./index.js";

describe("@gcsim/types", () => {
  describe("generated proto types — model", () => {
    it("exports SimulationResult schema", () => {
      expect(SimulationResultSchema).toBeDefined();
      expect(SimulationResultSchema.typeName).toBe("model.SimulationResult");
    });

    it("exports SimulationStatistics schema", () => {
      expect(SimulationStatisticsSchema).toBeDefined();
      expect(SimulationStatisticsSchema.typeName).toBe("model.SimulationStatistics");
    });

    it("exports Character schema", () => {
      expect(CharacterSchema).toBeDefined();
      expect(CharacterSchema.typeName).toBe("model.Character");
    });

    it("exports Enemy schema", () => {
      expect(EnemySchema).toBeDefined();
    });

    it("exports Weapon schema", () => {
      expect(ProtoWeaponSchema).toBeDefined();
    });

    it("exports Version schema", () => {
      expect(VersionSchema).toBeDefined();
    });

    it("exports OverviewStats schema", () => {
      expect(OverviewStatsSchema).toBeDefined();
    });

    it("exports DescriptiveStats schema", () => {
      expect(DescriptiveStatsSchema).toBeDefined();
    });

    it("exports Sample schema", () => {
      expect(SampleSchema).toBeDefined();
    });
  });

  describe("generated proto types — backend", () => {
    it("exports share namespace", () => {
      expect(share.ShareEntrySchema).toBeDefined();
    });

    it("exports db namespace", () => {
      expect(db.EntrySchema).toBeDefined();
    });

    it("exports preview namespace", () => {
      expect(preview.GetRequestSchema).toBeDefined();
    });
  });

  describe("custom sim interfaces (compile-time checks)", () => {
    it("SimResults type is usable", () => {
      const result: Sim.SimResults = {
        sim_version: "1.0.0",
        config_file: "test config",
        statistics: {
          iterations: 1000,
          dps: { mean: 50000, sd: 1000, min: 40000, max: 60000 },
        },
      };
      expect(result.sim_version).toBe("1.0.0");
      expect(result.statistics?.dps?.mean).toBe(50000);
    });

    it("Character type is usable", () => {
      const char: Sim.Character = {
        name: "hutao",
        level: 90,
        element: "pyro",
        max_level: 90,
        cons: 1,
        weapon: { name: "staffofhoma", refine: 1, level: 90, max_level: 90 },
        talents: { attack: 10, skill: 10, burst: 10 },
        stats: [0, 0, 0],
        snapshot: [0, 0, 0],
        sets: { crimsonwitchofflames: 4 },
      };
      expect(char.name).toBe("hutao");
    });

    it("StatusType type is usable", () => {
      const status: Sim.StatusType = "done";
      expect(status).toBe("done");
    });

    it("ParsedResult type is usable", () => {
      const parsed: Sim.ParsedResult = {
        characters: [],
        errors: [],
        player_initial_pos: { x: 0, y: 0, r: 0 },
      };
      expect(parsed.characters).toHaveLength(0);
    });
  });
});
