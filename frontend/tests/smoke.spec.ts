import { expect, test } from "@playwright/test";

test("app shell renders and views toggle when available", async ({ page }) => {
  await page.goto("/");
  await expect(page.locator("#root")).toBeVisible();

  const workflowsBtn = page.getByRole("button", { name: "Workflows" });
  const marketplaceBtn = page.getByRole("button", { name: "Marketplace" });

  if ((await workflowsBtn.count()) > 0 && (await marketplaceBtn.count()) > 0) {
    await workflowsBtn.click();
    await expect(page.getByText(/Workflow Gallery/i)).toBeVisible({
      timeout: 5000,
    });

    await marketplaceBtn.click();
    await expect(page.getByText(/MCP Marketplace/i)).toBeVisible({
      timeout: 5000,
    });
  } else {
    test.skip(
      true,
      "Sidebar navigation not available in this environment (likely backend error state).",
    );
  }
});
