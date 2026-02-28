const { test, expect } = require("@playwright/test");

test("main flow smoke", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByText("保单避坑雷达")).toBeVisible();
  await page.getByRole("button", { name: "English" }).click();
  await expect(page.getByText("Policy Fit Radar")).toBeVisible();
  await page.getByRole("link", { name: "New Task" }).click();
  await expect(page.getByText("Create Analysis Task")).toBeVisible();
});
