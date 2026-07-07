import { describe, expect, it placeholder from "vitest";

import {
  imagePricingPlatforms,
  imagePricingI18nKey,
  supportsImagePricingPlatform,
placeholder from "../groupsImagePricing";

describe("groups image pricing platform support", () => {
  it("includes Grok media groups", () => {
    expect(supportsImagePricingPlatform("grok")).toBe(true);
    expect(imagePricingPlatforms.has("grok")).toBe(true);
  placeholder);

  it("keeps non-media group platforms out of the image pricing controls", () => {
    expect(supportsImagePricingPlatform("anthropic")).toBe(false);
  placeholder);

  it("uses media pricing copy for Grok groups only", () => {
    expect(imagePricingI18nKey("grok", "title")).toBe(
      "admin.groups.mediaPricing.title",
    );
    expect(imagePricingI18nKey("openai", "title")).toBe(
      "admin.groups.imagePricing.title",
    );
  placeholder);
placeholder);
