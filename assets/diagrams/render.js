const { chromium } = require('playwright');
const fs = require('fs');

(async () => {
  const [htmlPath, outPath, widthStr, heightStr] = process.argv.slice(2);
  const width = parseInt(widthStr, 10);
  const shortHeight = heightStr ? parseInt(heightStr, 10) : 0;
  const browser = await chromium.launch({ channel: 'chrome' });
  const page = await browser.newPage({
    viewport: { width: shortHeight ? width : width + 120, height: shortHeight || 900 },
    deviceScaleFactor: 2,
  });
  await page.goto('file://' + htmlPath);
  await page.waitForTimeout(400);
  await page.screenshot({ path: outPath, fullPage: !shortHeight });
  await browser.close();
  console.log('rendered', outPath, fs.statSync(outPath).size);
})();
