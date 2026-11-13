const { chromium } = require('playwright');
const fs = require('fs');

(async () => {
  const [htmlPath, outPath, widthStr, heightStr] = process.argv.slice(2);
  const width = parseInt(widthStr, 10);
  const height = parseInt(heightStr, 10);
  const browser = await chromium.launch({ channel: 'chrome' });
  const page = await browser.newPage({
    viewport: { width, height },
    deviceScaleFactor: 2,
  });
  await page.goto('file://' + htmlPath);
  await page.waitForTimeout(250);
  await page.evaluate(() => document.fonts.ready);
  await page.locator('.canvas').screenshot({ path: outPath });
  await browser.close();
  console.log('rendered', outPath, fs.statSync(outPath).size);
})();
