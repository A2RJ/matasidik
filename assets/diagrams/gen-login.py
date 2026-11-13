#!/usr/bin/env python3
"""Generate the simple four-step approval strip used on the login page."""
import os

BASE = os.path.dirname(os.path.abspath(__file__))
TEMPLATE = os.environ.get(
    "STACK_DIAGRAM_TEMPLATE",
    "/tmp/opencode/stack-diagram/skills/stack-diagram/assets/template.html",
)
STYLE = open(TEMPLATE).read().split("<style>")[1].split("</style>")[0]
STYLE += """
* { box-sizing:border-box; }
html, body { height:220px; margin:0; padding:0; overflow:hidden; background:var(--color-paper); }
.simple-canvas { position:relative; width:1200px; height:220px; background:var(--color-paper); }
.step-card { position:absolute; top:36px; width:260px; height:148px; border:2px solid var(--color-ink);
  background:var(--color-paper); padding:14px 16px; box-shadow:4px 4px 0 0 var(--color-ink); }
.step-card.final { background:#FFDB33; }
.step-no { font-family:var(--font-mono); font-size:10px; font-weight:700; color:var(--color-muted); }
.step-title { margin-top:5px; font-family:var(--font-display); font-size:20px; font-weight:700; line-height:1.1; }
.step-text { margin-top:12px; font-family:var(--font-mono); font-size:9.5px; line-height:1.55; color:var(--color-ink-2); }
"""

OUT = os.environ.get("MATA_DIAGRAM_DRAFTS", os.path.join(BASE, "drafts"))
os.makedirs(OUT, exist_ok=True)

steps = [
    ("1", "Unggah", "Puskesmas mengirim BNBA.", False),
    ("2", "Verifikasi", "Dinas memeriksa & menyetujui/mengembalikan.", False),
    ("3", "Validasi", "Kabid menyetujui final.", False),
    ("4", "Final", "Sistem mengunci + menghitung.", True),
]
xs = [30, 330, 630, 930]
cards = ""
for (no, title, text, final), x in zip(steps, xs):
    cards += f'''<div class="step-card{" final" if final else ""}" style="left:{x}px;">
      <div class="step-no">{no}.</div>
      <div class="step-title">{title}</div>
      <div class="step-text">{text}</div>
    </div>'''

arrows = ""
for x in [290, 590, 890]:
    arrows += f'<path d="M{x} 110 H{x + 34}" stroke="#111" stroke-width="2" marker-end="url(#ar)"/>'

html_doc = f'''<!DOCTYPE html>
<html lang="id">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Alur Persetujuan MATA SIDIK</title>
<style>{STYLE}</style>
</head>
<body>
  <div class="simple-canvas">
    <svg class="edges" viewBox="0 0 1200 220" style="width:100%;height:100%" aria-hidden="true">
      <defs>
        <marker id="ar" viewBox="0 0 10 10" refX="8" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
          <path d="M0 0 L10 5 L0 10 z" fill="#111"/>
        </marker>
      </defs>
      {arrows}
    </svg>
    {cards}
  </div>
</body>
</html>
'''
open(os.path.join(OUT, "alur-login.html"), "w").write(html_doc)
print("login draft written")
