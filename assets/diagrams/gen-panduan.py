#!/usr/bin/env python3
"""Generate the user-guide overview flow (canvas-only stack-diagram).

Output: drafts/alur-panduan.html

The page at /user-guide uses one PNG for every role module. The green start
frame is an HTML overlay in the app, positioned with percentages that match
this canvas (1680x440). Keep CANVAS_W/CANVAS_H in sync with the app overlay.
"""
import html
import os

BASE = os.path.dirname(os.path.abspath(__file__))
TEMPLATE = os.environ.get(
    "STACK_DIAGRAM_TEMPLATE",
    "/tmp/opencode/stack-diagram/skills/stack-diagram/assets/template.html",
)
STYLE = open(TEMPLATE).read().split("<style>")[1].split("</style>")[0]
STYLE += """
body { padding: 0; background: var(--color-paper); }
.canvas { border: 0; border-radius: 0; }
.note { position:absolute; border:1px solid var(--color-rule-2); background:var(--color-paper);
  border-radius:8px; padding:10px 12px; font-family:var(--font-mono); font-size:9.5px;
  line-height:1.65; color:var(--color-neutral); }
.note b { display:block; font-family:var(--font-display); font-size:12px; color:var(--color-ink);
  margin-bottom:2px; letter-spacing:0; }
.topline { position:absolute; display:flex; align-items:center; gap:10px; border:1px solid var(--color-rule-2);
  background:var(--color-paper-2); border-radius:8px; padding:8px 13px; font-family:var(--font-mono);
  font-size:9.5px; line-height:1.5; color:var(--color-ink-2); }
.topline b { color:var(--color-ink); }
.topline img { width:18px; height:18px; }
"""
OUT = os.environ.get("MATA_DIAGRAM_DRAFTS", os.path.join(BASE, "drafts"))
os.makedirs(OUT, exist_ok=True)

CANVAS_W, CANVAS_H = 1680, 440


def node(cx, top, name, role, icon, cls="", width=200):
    left = cx - width // 2
    style = f"left:{left}px;top:{top}px"
    if width != 132:
        style += f";width:{width}px"
    return f'''    <div class="node" style="{style}">
      <div class="tile {cls}"><img src="%%ICON_{icon}%%" alt=""></div>
      <div class="name"><span>{html.escape(name)}</span></div><div class="role"><span>{html.escape(role)}</span></div>
    </div>'''


def group(left, top, w, h, label, cls="gb", icon=None):
    ic = f'<img src="%%ICON_{icon}%%" alt="">' if icon else ""
    return f'''    <div class="group {cls}" style="left:{left}px;top:{top}px;width:{w}px;height:{h}px;">
      <div class="g-label">{ic}{html.escape(label)}</div>
    </div>'''


def edge(path, cls="e", label=None, lx=0, ly=0, anchor="middle"):
    s = f'      <path class="{cls}" d="{path}" marker-end="url(#ar)"/>'
    if label:
        s += f'\n      <text class="el{" a" if cls == "e" else ""}" x="{lx}" y="{ly}" text-anchor="{anchor}">{html.escape(label)}</text>'
    return s


def note(x, y, w, h, title, body):
    return f'''    <div class="note" style="left:{x}px;top:{y}px;width:{w}px;height:{h}px;">
      <b>{html.escape(title)}</b>{html.escape(body)}
    </div>'''


groups = [
    group(20, 90, 320, 330, "Puskesmas / PJ Program", "ga", "pkm"),
    group(350, 90, 300, 330, "Verifikator Dinas", "gb", "verif"),
    group(660, 90, 300, 330, "Kabid / Sekretaris", "gb", "kabid"),
    group(970, 90, 300, 330, "Sistem MATA SIDIK", "gb", "engine"),
    group(1280, 90, 380, 330, "Pimpinan", "gb", "pimpinan"),
]

nodes = [
    # Jalur utama per role
    node(180, 160, "Input BNBA", "unggah + pratinjau + koreksi IKU", "upload"),
    node(500, 160, "Verifikasi", "periksa baris + berkas + catatan", "verif"),
    node(810, 160, "Validasi final", "final + catatan + kunci data", "valid"),
    node(1120, 160, "Kalkulasi & dashboard", "RATIO_OF_SUMS + penelusuran", "engine", width=210),
    node(1470, 160, "Dashboard & arahan", "baca capaian + tetapkan TL", "dash", width=230),
]

notes = [
    note(40, 280, 280, 100, "Alur peran",
         "Masuk → pilih faskes → unggah → tunggu proses → pratinjau → koreksi IKU anak → kirim."),
    note(370, 280, 260, 100, "Alur peran",
         "Buka pengajuan Dikirim → periksa baris/berkas → setujui ke Kabid atau revisi + catatan."),
    note(680, 280, 260, 100, "Alur peran",
         "Buka pengajuan Diverifikasi → validasi → final + catatan → isi faktor."),
    note(990, 280, 260, 100, "Alur sistem",
         "Klasifikasi deterministik → alur kerja → kalkulasi → SPM + jejak angka."),
    note(1300, 280, 340, 100, "Alur peran",
         "Baca dashboard → telusuri angka → beri arahan → tetapkan & pantau TL."),
]

edges = [
    edge("M280 190 H400", "e", "kirim", 340, 182),
    edge("M600 190 H710", "e", "setujui", 655, 182),
    edge("M910 190 H1015", "e", "final", 963, 182),
    edge("M1225 190 H1355", "e", "sajikan", 1290, 182),
    # Return paths inside the bottom of the main lane groups
    edge("M500 250 V405 H180 V250", "n dash", "revisi + catatan", 340, 397),
    edge("M810 250 V415 H500 V250", "n dash", "perlu diperiksa", 655, 407),
]

topline = f'''    <div class="topline" style="left:20px;top:20px;width:1640px;height:50px;">
      <img src="%%ICON_user%%" alt="">
      <span><b>Masuk & sesi (semua peran)</b> → pilih periode → hak akses mengikuti peran efektif jabatan. Peran Puskesmas/PJ dibatasi faskes sendiri.</span>
    </div>'''

legend = [
    ("", "alur utama"),
    ("n dash", "revisi / dukungan lintas peran"),
]

html_out = f'''<!DOCTYPE html>
<html lang="id">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Alur Lengkap MATA SIDIK — Panduan Pengguna</title>
<style>{STYLE}</style>
</head>
<body>
  <div class="canvas" style="width:{CANVAS_W}px;height:{CANVAS_H}px;">

{chr(10).join(groups)}

    <svg class="edges" viewBox="0 0 {CANVAS_W} {CANVAS_H}" style="width:100%;height:100%" aria-hidden="true">
      <defs>
        <marker id="ar" viewBox="0 0 10 10" refX="8" refY="5" markerWidth="6.5" markerHeight="6.5" orient="auto-start-reverse">
          <path d="M0 0 L10 5 L0 10 z" fill="context-stroke"/>
        </marker>
      </defs>
{chr(10).join(edges)}
    </svg>

{topline}

{chr(10).join(nodes)}

{chr(10).join(notes)}

  </div>
</body>
</html>
'''
open(f"{OUT}/alur-panduan.html", "w").write(html_out)
print("alur-panduan draft written")
