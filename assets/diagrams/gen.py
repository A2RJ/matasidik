#!/usr/bin/env python3
"""Generate stack-diagram HTML drafts for MATA SIDIK."""
import html
import os

BASE = os.path.dirname(os.path.abspath(__file__))
TEMPLATE = os.environ.get(
    "STACK_DIAGRAM_TEMPLATE",
    "/tmp/opencode/stack-diagram/skills/stack-diagram/assets/template.html",
)
STYLE = open(TEMPLATE).read().split("<style>")[1].split("</style>")[0]
# Chromium full-page screenshots can duplicate nearby icons through backdrop-filter;
# the veil is only cosmetic, so disable the blur for these exports.
STYLE += "\n.node .name span, .node .role span { backdrop-filter: none; }\n"
OUT = os.environ.get("MATA_DIAGRAM_DRAFTS", os.path.join(BASE, "drafts"))
os.makedirs(OUT, exist_ok=True)


def node(cx, top, name, role, icon, cls="", width=132):
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
        s += f'\n      <text class="el{" a" if cls=="e" else ""}" x="{lx}" y="{ly}" text-anchor="{anchor}">{html.escape(label)}</text>'
    return s


def page(title, kicker, lede, w, h, groups, nodes, edges, legend, foot="icons: Lucide"):
    return f'''<!DOCTYPE html>
<html lang="id">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{html.escape(title)}</title>
<style>{STYLE}</style>
</head>
<body>
<div class="page">

  <header class="head">
    <div>
      <div class="kicker">{html.escape(kicker)}</div>
      <h1>{html.escape(title)}</h1>
    </div>
    <p class="lede">{html.escape(lede)}</p>
  </header>

  <div class="scroller">
  <div class="canvas" style="width:{w}px;height:{h}px;">

{chr(10).join(groups)}

    <svg class="edges" viewBox="0 0 {w} {h}" style="width:100%;height:100%" aria-hidden="true">
      <defs>
        <marker id="ar" viewBox="0 0 10 10" refX="8" refY="5" markerWidth="6.5" markerHeight="6.5" orient="auto-start-reverse">
          <path d="M0 0 L10 5 L0 10 z" fill="context-stroke"/>
        </marker>
      </defs>
{chr(10).join(edges)}
    </svg>

{chr(10).join(nodes)}

  </div>
  </div>

  <div class="legend">
{chr(10).join(f'    <span><span class="sw {c}"></span>{html.escape(t)}</span>' for c,t in legend)}
  </div>

  <div class="foot">
    <span>{html.escape(foot)}</span>
    <span>dibuat dengan <a href="https://github.com/caru-ini/stack-diagram" target="_blank" rel="noopener">stack-diagram</a></span>
  </div>

</div>
</body>
</html>
'''


# ---------------------------------------------------------------- 1. ALUR BISNIS
w, h = 1340, 600
groups = [
    group(20, 100, 200, 210, "Puskesmas", "ga", "pkm"),
    group(250, 100, 1060, 210, "Dinas Kesehatan", "gb", "dinas"),
]
nodes = [
    node(120, 160, "Input Data", "Puskesmas / PJ Program", "input", width=180),
    node(340, 160, "Verifikasi", "Petugas Dinas", "verif"),
    node(560, 160, "Validasi", "Kabid / Sekretaris", "valid"),
    node(780, 160, "Pengolahan", "Sistem + faktor Kabid", "engine", width=170),
    node(1000, 160, "Dashboard", "Semua peran", "dash"),
    node(1220, 160, "Keputusan", "Kepala Dinas + TL", "keputusan"),
]
edges = [
    edge("M156 190 H304", "e", "kirim + BNBA", 230, 182),
    edge("M376 190 H524", "e", "setuju", 450, 182),
    edge("M596 190 H744", "e", "final", 670, 182),
    edge("M816 190 H964", "e", "capaian", 890, 182),
    edge("M1036 190 H1184", "e", "arahan", 1110, 182),
    edge("M560 280 V370 H340 V280", "n dash", "kembalikan", 450, 362),
    edge("M340 280 V430 H120 V280", "n dash", "revisi", 230, 422),
    edge("M1220 280 V500 H110 V280", "n dash", "pemantauan TL + faktor", 665, 492),
]
legend = [("", "alur utama"), ("n dash", "pengembalian / siklus")]
open(f"{OUT}/alur-bisnis.html", "w").write(page(
    "Alur Bisnis MATA SIDIK",
    "MATA SIDIK / proses bisnis",
    "Enam tahap dari input Puskesmas sampai keputusan pimpinan, lengkap dengan jalur pengembalian data, faktor kualitatif, dan TL.",
    w, h, groups, nodes, edges, legend))

# ---------------------------------------------------------------- 2. STATUS DATA
w, h = 1240, 460
groups = [
    group(770, 90, 420, 210, "Pasca-final + audit", "gb", "audit"),
]
nodes = [
    node(140, 150, "Draf", "baru diinput", "draf"),
    node(380, 150, "Dikirim", "menunggu verifikasi", "kirim"),
    node(620, 150, "Diverifikasi", "lolos pemeriksaan Dinas", "verif"),
    node(860, 150, "Final", "terkunci", "final"),
    node(1100, 150, "Diubah Khusus", "kewenangan + audit", "ubah"),
]
edges = [
    edge("M176 180 H314", "e", "kirim", 245, 172),
    edge("M416 180 H554", "e", "verifikasi", 485, 172),
    edge("M656 180 H794", "e", "sahkan", 725, 172),
    edge("M896 168 H1064", "e", "kewenangan khusus", 980, 160),
    edge("M1064 192 H896", "n dash", "versi baru", 980, 204),
    edge("M380 256 V340 H140 V256", "n dash", "kembalikan", 260, 332),
    edge("M620 256 V400 H392 V256", "n dash", "perlu diperiksa", 506, 392),
]
legend = [("", "alur status"), ("n dash", "pengembalian / versi")]
open(f"{OUT}/status-data.html", "w").write(page(
    "Siklus Status Data",
    "MATA SIDIK / alur kerja",
    "Tujuh status dengan jalur pengembalian yang eksplisit; data final hanya berubah lewat kewenangan khusus dan jejak audit.",
    w, h, groups, nodes, edges, legend))

# ---------------------------------------------------------------- 3. ALUR PERAN
w, h = 1360, 800
groups = [
    group(20, 80, 240, 660, "Puskesmas / PJ Program", "gb", "pkm"),
    group(280, 80, 240, 660, "Verifikator Dinas", "gb", "verif"),
    group(540, 80, 240, 660, "Kabid / Sekretaris", "gb", "kabid"),
    group(800, 80, 240, 660, "Sistem MATA SIDIK", "gb", "engine"),
    group(1060, 80, 240, 660, "Pimpinan", "gb", "pimpinan"),
]
nodes = [
    # row 1 — alur utama
    node(140, 170, "Input + BNBA", "Puskesmas", "pkm"),
    node(400, 170, "Verifikasi", "petugas dinas", "verif"),
    node(660, 170, "Validasi final", "kabid / sekretaris", "valid"),
    node(920, 170, "Olah data", "sistem + agregasi", "engine"),
    node(1180, 170, "Dashboard", "capaian + tren", "dash"),
    # row 2 — pengembalian & arahan
    node(140, 390, "Perbaiki data", "revisi", "perbaiki"),
    node(400, 390, "Kembalikan", "+ catatan", "kembali"),
    node(660, 390, "Kembalikan", "ke verifikator", "kembali"),
    node(1180, 390, "Arahan + TL", "pimpinan → PJ", "keputusan"),
    # row 3 — penyelesaian
    node(140, 610, "Kirim ulang", "data revisi", "kirim"),
    node(400, 610, "Setujui", "verifikasi ulang", "valid"),
    node(660, 610, "Isi faktor", "pendukung / penghambat", "faktor"),
    node(920, 610, "Perbarui dashboard", "faktor + capaian", "dash"),
    node(1180, 610, "Pantau TL", "status penyelesaian", "tl"),
]
edges = [
    # row 1
    edge("M176 200 H364", "e", "kirim", 270, 192),
    edge("M436 200 H624", "e", "setuju", 530, 192),
    edge("M696 200 H884", "e", "final", 790, 192),
    edge("M956 200 H1144", "e", "sajikan", 1050, 192),
    # cabang pengembalian
    edge("M400 276 V384", "n dash", "jika belum", 410, 330),
    edge("M660 276 V384", "n dash", "jika perlu diperiksa", 670, 330),
    edge("M660 420 H436", "n dash", "perlu diperiksa", 548, 412),
    edge("M400 420 H176", "n dash", "revisi", 288, 412),
    # row 3
    edge("M176 640 H364", "e", "kirim ulang", 270, 632),
    edge("M436 640 H624", "e", "setuju", 530, 632),
    edge("M696 640 H884", "e", "faktor", 790, 632),
    # TL ditetapkan pimpinan, lalu dipantau
    edge("M1180 496 V604", "n dash", "TL", 1190, 550),
    # vertikal & elbow
    edge("M1180 276 V384", "n dash", "arahan", 1190, 330),
    edge("M140 496 V604", "n dash", "kirim ulang", 150, 550),
    edge("M920 276 V530 H650 V604", "n dash", "minta faktor", 785, 522),
]
legend = [("", "alur utama"), ("n dash", "pengembalian / arahan / TL")]
open(f"{OUT}/alur-peran.html", "w").write(page(
    "Alur Peran MATA SIDIK",
    "MATA SIDIK / serah terima",
    "Swimlane lima peran: siapa mengirim, memeriksa, memvalidasi, mengolah, mengisi faktor, serta menetapkan dan memantau TL.",
    w, h, groups, nodes, edges, legend))

# ---------------------------------------------------------------- 4. MODEL DATA
w, h = 1320, 800
groups = [
    group(40, 80, 1240, 330, "Master & referensi", "ga", "table"),
    group(40, 460, 1240, 300, "Transaksi, audit, faktor & TL", "gb", "laporan"),
]
nodes = [
    # master
    node(130, 170, "WILAYAH", "desa / kecamatan", "wilayah"),
    node(340, 170, "FASKES", "puskesmas / faskes", "pkm"),
    node(550, 170, "USER", "pengguna", "user"),
    node(760, 170, "ROLE", "hak akses", "role"),
    node(970, 170, "PERIODE", "bulan / tahun", "periode"),
    node(1180, 170, "INDIKATOR", "versi master", "table"),
    # transaksi
    node(130, 550, "TINDAK_LANJUT", "masalah → rekomendasi → PJ → status", "tl", width=170),
    node(340, 550, "FAKTOR", "kualitatif · pendukung/penghambat", "faktor", width=170),
    node(550, 550, "LAPORAN", "nilai per indikator", "laporan"),
    node(760, 550, "AUDIT_LOG", "siapa / kapan / apa", "audit"),
    node(1180, 550, "TARGET", "tahunan / bulanan", "target"),
]
edges = [
    # master relations
    edge("M166 200 H304", "n", "1:N", 235, 192),
    edge("M724 200 H586", "n", "1:N", 655, 192),
    # indikator -> target
    edge("M1180 276 V544", "e", "1:N", 1188, 437, anchor="start"),
    # indikator -> laporan
    edge("M1180 276 V500 H560 V544", "e", "diukur", 870, 492),
    # faskes -> laporan
    edge("M340 276 V520 H575 V544", "n", "1:N", 460, 512),
    # user -> laporan
    edge("M540 276 V544", "n",  "masukan", 548, 350, anchor="start"),
    # periode -> laporan
    edge("M970 276 V480 H525 V544", "n", "periode", 785, 472),
    # laporan relations
    edge("M514 575 H376", "n", "faktor", 445, 567),
    edge("M586 575 H724", "n", "audit", 655, 567),
    edge("M550 656 V720 H130 V656", "n", "TL", 340, 712),
]
legend = [("", "relasi utama"), ("n", "relasi pendukung 1:N")]
open(f"{OUT}/model-data.html", "w").write(page(
    "Model Data MATA SIDIK",
    "MATA SIDIK / model data",
    "Entitas master dipisahkan dari transaksi; setiap laporan menyimpan target, faktor kualitatif, audit, dan TL ber-PJ.",
    w, h, groups, nodes, edges, legend))

print("drafts written")
