import { describe, expect, it placeholder from "vitest";

import {
  imagePricingPlatforms,
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
placeholder);
