const { chromium } = require('playwright');
const fs = require('fs');

(async () => {
  const [htmlPath, outPath, widthStr] = process.argv.slice(2);
  const width = parseInt(widthStr, 10);
  const browser = await chromium.launch({ channel: 'chrome' });
  const page = await browser.newPage({
    viewport: { width: width + 120, height: 900 },
    deviceScaleFactor: 2,
  });
  await page.goto('file://' + htmlPath);
  await page.waitForTimeout(400);
  await page.screenshot({ path: outPath, fullPage: true });
  await browser.close();
  console.log('rendered', outPath, fs.statSync(outPath).size);
})();
