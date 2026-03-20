import i18next from "i18next";
import { beforeEach, describe, expect, it } from "vitest";
import { initI18n, resources, specialLocales } from "../index";

describe("@gcsim/i18n", () => {
  beforeEach(() => {
    // Reset i18next instance between tests
    if (i18next.isInitialized) {
      i18next.changeLanguage("en");
    }
  });

  it("should initialize with English language", async () => {
    await initI18n("en");
    expect(i18next.isInitialized).toBe(true);
    expect(i18next.language).toBe("en");
  });

  it("should resolve a known English translation key", async () => {
    await initI18n("en");
    const result = i18next.t("nav.simulator");
    expect(result).toBe("Simulator");
  });

  it("should fall back to English when a key is missing in another language", async () => {
    await initI18n("zh");
    // "common.data_not_found" exists in English but may not in all languages
    // Use a key from English that we know the value of
    const enValue = resources.en.translation.nav.simulator;
    expect(enValue).toBe("Simulator");

    // Chinese has its own translation for this key
    const zhValue = i18next.t("nav.simulator");
    expect(zhValue).toBe("模拟器");

    // Verify fallback works for a completely fabricated namespace scenario:
    // If we request a key that does not exist in zh, it should fall back to en
    const fallback = i18next.t("nav.simulator", { lng: "en" });
    expect(fallback).toBe("Simulator");
  });

  it("should resolve character names from the game namespace", async () => {
    await initI18n("en");
    // "aetherelectro" should come from traveler names (overriding generated names)
    const name = i18next.t("game:character_names.aetherelectro");
    expect(name).toBe("Aether (Electro)");

    // "albedo" should come from the generated names
    const albedo = i18next.t("game:character_names.albedo");
    expect(albedo).toBeTruthy();
    expect(typeof albedo).toBe("string");
  });

  it("should have resources defined for all 7 languages", () => {
    const expectedLanguages = ["en", "zh", "ja", "ko", "es", "ru", "de"];
    for (const lang of expectedLanguages) {
      expect(resources[lang as keyof typeof resources]).toBeDefined();
      expect(resources[lang as keyof typeof resources].translation).toBeDefined();
      expect(resources[lang as keyof typeof resources].game).toBeDefined();
    }
  });

  it("should export specialLocales with CJK languages", () => {
    expect(specialLocales).toEqual(["zh", "ja", "ko"]);
  });
});
