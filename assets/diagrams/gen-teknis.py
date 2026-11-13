#!/usr/bin/env python3
"""Generate the technical flow/architecture stack diagrams for MATA SIDIK.

Outputs drafts consumed by build.sh:
  - alur-upload.html          pipa upload background (04 §3.2)
  - alur-workflow.html        state machine approval (04 §4.1)
  - alur-provenance.html      rollup IKU induk -> anak (04 §5.3)
  - alur-trace-spm.html       rantai trace SPM -> FILE_OBJECT (04 §5.5)
  - arsitektur-aplikasi.html  modular monolith + worker (04 §6.1, README B.8)
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
.node .name span, .node .role span { backdrop-filter: none; }
.panel { position:absolute; background:var(--color-paper); border:1px solid var(--color-rule-2);
  border-radius:8px; padding:9px 11px; box-shadow:0 1px 2px oklch(24% 0.02 258 / .04); }
.p-title { display:flex; align-items:center; gap:5px; font-family:var(--font-mono); font-size:8px;
  letter-spacing:.06em; text-transform:uppercase; color:var(--color-muted); margin-bottom:5px; }
.p-title img { width:12px; height:12px; }
.p-title .step { margin-left:auto; border:1px solid var(--color-rule-2); border-radius:20px;
  padding:0 5px; color:var(--color-muted); font-size:7px; letter-spacing:0; }
.big { font-family:var(--font-display); font-weight:650; font-size:14px; line-height:1.25; }
.mono { font-family:var(--font-mono); }
.sub { font-family:var(--font-mono); font-size:8px; color:var(--color-muted); line-height:1.6; }
.calc { position:absolute; background:var(--color-paper); border:1px solid var(--color-rule-2);
  border-radius:6px; padding:8px 10px 8px 13px; }
.calc:before { content:""; position:absolute; left:0; top:8px; bottom:8px; width:3px;
  border-radius:2px; background:var(--color-g1); }
.calc.parent:before { background:var(--color-accent); }
.calc .label { font-family:var(--font-mono); font-size:7.6px; color:var(--color-muted); }
.calc .result { font-family:var(--font-display); font-size:22px; font-weight:700;
  color:var(--color-accent); line-height:1.15; margin-top:1px; }
.calc .formula { font-family:var(--font-mono); font-size:8.4px; color:var(--color-ink-2); margin-top:1px; }
.badge { display:inline-block; font-family:var(--font-mono); font-size:7px; line-height:1.35;
  padding:1px 5px; border-radius:20px; border:1px solid; white-space:nowrap; }
.b-ok { color:oklch(44% .12 150); border-color:oklch(70% .12 150); background:oklch(96% .04 150); }
.b-warn { color:oklch(50% .13 75); border-color:oklch(75% .13 75); background:oklch(96% .05 75); }
.b-final { color:var(--color-accent); border-color:var(--color-accent-soft); background:var(--color-accent-tint); }
.b-muted { color:var(--color-muted); border-color:var(--color-rule-2); background:var(--color-paper-2); }
.dash { width:100%; border-collapse:collapse; font-family:var(--font-mono); font-size:9px; }
.dash th { text-align:left; font-weight:500; color:var(--color-muted); border-bottom:1px solid var(--color-rule-2);
  padding:5px 7px; font-size:7.7px; letter-spacing:.04em; text-transform:uppercase; white-space:nowrap; }
.dash td { padding:5px 7px; border-bottom:1px solid var(--color-rule); vertical-align:top; }
.dash tbody tr:last-child td { border-bottom:0; }
.dash .num { text-align:right; }
.dash .c { text-align:center; }
.dash .st { font-weight:700; font-size:10px; }
.st.ok { color:oklch(50% .15 150); }
.st.warn { color:oklch(68% .15 75); }
.st.bad { color:oklch(58% .19 25); }
.strip { position:absolute; background:var(--color-paper-2); border:1px solid var(--color-rule-2);
  border-radius:8px; padding:9px 13px; }
.strip .k { font-family:var(--font-mono); font-size:7.4px; letter-spacing:.08em; color:var(--color-accent);
  margin-bottom:3px; }
.strip .t { font-family:var(--font-mono); font-size:8px; color:var(--color-ink-2); line-height:1.7; }
.divider { border-top:1px solid var(--color-rule); margin:5px 0; }
"""
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


def panel(x, y, w, h, title, body, icon=None, step=None):
    ic = f'<img src="%%ICON_{icon}%%" alt="">' if icon else ""
    st = f'<span class="step">{step}</span>' if step else ""
    return f'''    <div class="panel" style="left:{x}px;top:{y}px;width:{w}px;height:{h}px;">
      <div class="p-title">{ic}{html.escape(title)}{st}</div>
      {body}
    </div>'''


def card(x, y, w, h, label, result, formula, note, parent=False, badges=""):
    return f'''    <div class="calc {"parent" if parent else ""}" style="left:{x}px;top:{y}px;width:{w}px;height:{h}px;">
      <div class="label">{html.escape(label)}</div>
      <div class="result">{result}</div>
      <div class="formula">{formula}</div>
      <div class="sub">{note}</div>
      {badges}
    </div>'''


def badge(text, cls="b-muted"):
    return f'<span class="badge {cls}">{html.escape(text)}</span>'


def page(title, kicker, lede, w, h, groups, nodes, edges, legend, foot, extra=""):
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

{extra}

  </div>
  </div>

  <div class="legend">
{chr(10).join(f'    <span><span class="sw {c}"></span>{html.escape(t)}</span>' for c, t in legend)}
  </div>

  <div class="foot">
    <span>{html.escape(foot)}</span>
    <span>dibuat dengan <a href="https://github.com/caru-ini/stack-diagram" target="_blank" rel="noopener">stack-diagram</a></span>
  </div>

</div>
</body>
</html>
'''


# ------------------------------------------------------------------ 1. UPLOAD
w, h = 1280, 620
groups = []
nodes = [
    node(150, 160, "Unggah BNBA", "Svelte · multipart", "upload"),
    node(470, 160, "Gin API", "autentikasi · cakupan", "server"),
    node(790, 160, "PostgreSQL 16", "SUBMISSION · FILE_OBJECT · JOB", "engine"),
    node(1110, 160, "Worker Go", "urai · validasi · klasifikasi", "klasifikasi"),
    node(470, 450, "Pratinjau & pemantauan", "2 dtk · ubah IKU", "check"),
    node(790, 450, "Volume lokal", "/data/files · sha256", "disk"),
]
edges = [
    edge("M184 190 H436", "e", "POST · 202 job_id", 305, 182),
    edge("M504 190 H756", "e", "transaksi: 3 tabel", 630, 182),
    edge("M824 190 H1076", "e", "JOB · SKIP LOCKED", 950, 182),
    edge("M470 264 V446", "e", "halaman hasil", 462, 350, anchor="end"),
    edge("M500 446 V264", "n dash", "pemantauan GET /jobs/{id}", 508, 350, anchor="start"),
    edge("M1096 264 V320 H790 V264", "n", "BARIS + NILAI + progres", 943, 312),
    edge("M504 220 H680 V480 H756", "n", "simpan berkas sha256", 672, 472, anchor="end"),
    edge("M1124 264 V480 H824", "n", "baca berkas", 970, 472),
]
legend = [("", "alur utama"), ("n dash", "pemantauan / respons")]
open(f"{OUT}/alur-upload.html", "w").write(page(
    "Pipa Unggah Latar Belakang",
    "MATA SIDIK / penerimaan data",
    "Dari multipart ke baris terklasifikasi: API menyimpan berkas + job, worker memproses di latar belakang, peramban memantau status.",
    w, h, groups, nodes, edges, legend,
    "contoh: BNBA Mei 2026 · 3 baris · tanpa Redis · tabel job PostgreSQL · tanpa penyimpanan objek"))


# ------------------------------------------------------------------ 2. WORKFLOW
w, h = 1360, 620
groups = [
    group(70, 100, 1020, 260, "Alur utama pengajuan", "gb", "workflow"),
    group(960, 370, 200, 180, "Pasca-final", "ga", "audit"),
]
nodes = [
    node(160, 160, "DRAFT_PKM", "Puskesmas / PJ", "draf"),
    node(460, 160, "VERIFY_DINAS", "Verifikator Dinas", "verif"),
    node(760, 160, "VALIDATE_KABID", "Kabid / Sekretaris", "valid"),
    node(1060, 160, "FINAL", "kunci + hitung", "final"),
    node(1060, 430, "SPECIAL_CHANGE", "Kabid + Admin", "ubah"),
]
edges = [
    edge("M194 190 H426", "e", "kirim", 310, 182),
    edge("M494 190 H726", "e", "setujui", 610, 182),
    edge("M794 190 H1026", "e", "setujui + catatan", 910, 182),
    edge("M1060 264 V426", "n dash", "ubah paksa + alasan", 1052, 345, anchor="end"),
    edge("M1090 426 V264", "n dash", "versi baru", 1098, 345, anchor="start"),
    edge("M460 264 V330 H160 V264", "n dash", "revisi + catatan", 310, 322),
    edge("M760 264 V400 H460 V264", "n dash", "revisi", 610, 392),
]
legend = [("", "alur utama"), ("n dash", "pengembalian / pasca-final")]
open(f"{OUT}/alur-workflow.html", "w").write(page(
    "Alur Persetujuan Berjenjang",
    "MATA SIDIK / alur kerja",
    "Bawaan tiga tingkat dengan jalur pengembalian eksplisit; final mengunci data, perubahan pasca-final lewat kewenangan khusus.",
    w, h, groups, nodes, edges, legend,
    "versi alur kerja dipatok per pengajuan · tidak ada persetujuan otomatis saat SLA lewat"))


# --------------------------------------------------------------- 3. PROVENANCE
w, h = 1400, 560
extra = [
    card(530, 110, 340, 130, "IKU-2026-001 · INDUK", "92,39%", "850 / 920 · target 100%",
         "LENGKAP · gap -7,61 pp", parent=True,
         badges=f'<div style="margin-top:5px">{badge("PERHATIAN", "b-warn")} {badge("3 IKU ANAK", "b-final")}</div>'),
    card(90, 350, 330, 130, "IKU-2026-001.01 · Ibu Hamil", "90,00%", "450 / 500 · bobot 1",
         "SPM 01 · perlu perhatian",
         badges=f'<div style="margin-top:5px">{badge("PERHATIAN", "b-warn")}</div>'),
    card(535, 350, 330, 130, "IKU-2026-001.02 · Ibu Bersalin", "95,00%", "380 / 400 · bobot 1",
         "SPM 02 · sesuai target",
         badges=f'<div style="margin-top:5px">{badge("SESUAI", "b-ok")}</div>'),
    card(980, 350, 330, 130, "IKU-2026-001.03 · Bayi Baru Lahir", "100,00%", "20 / 20 · bobot 1",
         "SPM 03 · setelah reklasifikasi",
         badges=f'<div style="margin-top:5px">{badge("SESUAI", "b-ok")} {badge("DARI INDUK", "b-final")}</div>'),
    f'''    <div class="strip" style="left:30px;top:500px;width:1340px;height:45px;">
      <div class="k">RATIO_OF_SUMS · BUKAN RATA-RATA PERSENTASE</div>
      <div class="t">(450 + 380 + 20) / (500 + 400 + 20) = 850 / 920 = 92,39% &nbsp;·&nbsp;
      tiap anak disimpan di <b>IKU_VALUE_DETAIL</b> (component_type IKU anak, pembilang, penyebut, bobot) &nbsp;·&nbsp;
      .01 → SPM 01 · .02 → SPM 02 · .03 → SPM 03</div>
    </div>''',
]
edges = [
    edge("M600 240 V300 H255 V344", "n", "IKU anak · bobot 1", 420, 292),
    edge("M700 240 V344", "n", "bobot 1", 708, 300, anchor="start"),
    edge("M800 240 V300 H1145 V344", "n", "IKU anak · bobot 1", 975, 292),
]
legend = [("n", "komponen rollup IKU")]
open(f"{OUT}/alur-provenance.html", "w").write(page(
    "Jejak Angka Nilai IKU Induk",
    "MATA SIDIK / kalkulasi",
    "Nilai induk dihitung dari komponen eksplisit; setiap anak menyumbang numerator dan denominator, bukan rata-rata persentase.",
    w, h, [], [], edges, legend,
    "contoh: Mei 2026 · IKU-2026-001 · tiga IKU anak · target 100%", extra="\n".join(extra)))


# ------------------------------------------------------------------ 4. TRACE
w, h = 1860, 660
chain = [
    ("DASHBOARD_ROW", "SPM 01", "Ibu Hamil · target 100%", "iku_value_id", "table"),
    ("IKU_VALUE", "90,00%", "450/500 · LENGKAP", "iku_value_id + source_hash", "calc"),
    ("INDICATOR_VALUE", "IND-001", "450/500 · FINAL", "indicator_value_id", "check"),
    ("SUBMISSION_ROW", "R1 · COCOK", "row_no 1 · cakupan PKM", "submission_row_id", "upload"),
    ("SUBMISSION_REVISION", "REV 1", "tidak berubah · file_id", "revision_id", "audit"),
    ("SUBMISSION", "S-2026-05-001", "Mei 2026 · PKM Taliwang", "submission_id", "laporan"),
    ("FILE_OBJECT", "9f2c…a71", "BNBA_2026-05.xlsx", "file_id + uploaded_by", "disk"),
]
extra = []
for i, (tname, tval, tmeta, tid, icon) in enumerate(chain):
    x = 30 + i * 260
    extra.append(panel(x, 110, 220, 130, tname, f'''
      <div class="big">{html.escape(tval)}</div>
      <div class="sub">{html.escape(tmeta)}</div>
      <div class="divider"></div>
      <div class="sub">{html.escape(tid)}</div>
''', icon, step=str(i + 1)))

rows = [
    ("01", "Pelayanan Kesehatan Ibu Hamil", ".01", "R1", "IND-001", "450/500", "90,00%", "warn", "BERKAS R1 → SPM 01"),
    ("02", "Pelayanan Kesehatan Ibu Bersalin", ".02", "R2", "IND-002", "380/400", "95,00%", "ok", "BERKAS R2 → SPM 02"),
    ("03", "Pelayanan Kesehatan Bayi Baru Lahir", ".03", "R3", "IND-003", "20/20", "100,00%", "ok", "BERKAS R3 → SPM 03"),
]
trs = ""
for no, nama, leaf, baris, ind, nd, val, st, trace in rows:
    st_html = f'<span class="st {st}">{"✓" if st == "ok" else "●"}</span>'
    trs += (f'<tr><td class="num">{no}</td><td>{html.escape(nama)}</td><td class="num">{leaf}</td>'
            f'<td class="num">{baris}</td><td class="num">{ind}</td><td class="num">{nd}</td>'
            f'<td class="num">{val}</td><td class="c">{st_html}</td><td>{html.escape(trace)}</td></tr>')

extra.append(panel(30, 280, 1800, 210, "3 baris SPM — satu berkas, tiga IKU anak", f'''
      <table class="dash">
        <colgroup><col style="width:44px"><col><col style="width:70px"><col style="width:60px">
        <col style="width:90px"><col style="width:90px"><col style="width:90px"><col style="width:52px"><col style="width:190px"></colgroup>
        <thead><tr>
          <th class="num">SPM</th><th>Indikator</th><th class="num">IKU anak</th><th class="num">Baris</th>
          <th class="num">Indikator</th><th class="num">Num/den</th><th class="num">Capaian</th>
          <th class="c">Status</th><th>Penelusuran</th>
        </tr></thead>
        <tbody>{trs}</tbody>
      </table>
      <div class="divider"></div>
      <div class="sub">Sumber sama untuk ketiganya: <b>FILE_OBJECT 9f2c…a71</b> → <b>SUBMISSION S-2026-05-001</b> →
      <b>REVISION 1</b> · pengunggah Siti (PKM Taliwang) · periode Mei 2026. R3 semula <b>IKU_ROLLUP_ONLY</b>,
      dikoreksi ke <b>.03</b> di pratinjau dan tercatat di <b>AUDIT_LOG</b>.</div>
''', "tree"))

extra.append(f'''    <div class="strip" style="left:30px;top:520px;width:1800px;height:90px;">
      <div class="k">KONTRAK PENELUSURAN</div>
      <div class="t"><b>iku_value_id</b> · <b>indicator_value_id</b> · <b>source_hash</b> · <b>data_version</b> · <b>calculation_run_id</b>
      &nbsp;·&nbsp; rollup induk: 92,39% = 850/920 ← .01 450/500 + .02 380/400 + .03 20/20</div>
      <div class="t">API: <span class="mono">GET /api/v1/dashboard/spm?periode=2026-05</span> →
      <span class="mono">GET /api/v1/dashboard/spm/{{kode}}/trace?periode=2026-05</span> →
      <span class="mono">GET /api/v1/files/{{id}}</span></div>
    </div>''')

edges = []
# arrows between the seven chain panels (40px gaps)
for i in range(6):
    x1 = 30 + i * 260 + 220
    x2 = x1 + 36
    edges.append(edge(f"M{x1} 185 H{x2}"))

legend = [("", "rantai penelusuran antar tabel")]
open(f"{OUT}/alur-trace-spm.html", "w").write(page(
    "Penelusuran Output SPM sampai Berkas",
    "MATA SIDIK / lineage",
    "Setiap sel SPM adalah proyeksi nilai terhitung: telusuri dari baris dashboard sampai berkas unggahan aslinya.",
    w, h, [], [], edges, legend,
    "contoh: Mei 2026 · 3 baris SPM · satu berkas BNBA · penelusuran lewat ID kontrak", extra="\n".join(extra)))


# ------------------------------------------------------------- 5. ARCHITECTURE
w, h = 1440, 760
groups = [
    group(20, 100, 240, 260, "Peramban", "ga", "user"),
    group(300, 70, 1100, 640, "Docker Compose · satu host", "gb", "container"),
    group(590, 110, 180, 560, "Image Go", "ga", "code"),
    group(960, 110, 180, 220, "Basis data", "ga", "engine"),
    group(960, 430, 180, 220, "Penyimpanan berkas", "ga", "disk"),
]
nodes = [
    node(140, 170, "Svelte 5", "build statis · satu origin", "code"),
    node(400, 170, "Proxy balik", "TLS · static + /api/v1", "shield"),
    node(680, 170, "Gin API", "cmd/api · REST", "server"),
    node(1050, 170, "PostgreSQL 16", "data + JOB + outbox", "engine"),
    node(680, 490, "Worker Go", "cmd/worker · siklus job", "workflow"),
    node(1050, 490, "Volume lokal", "/data/files · sha256", "disk"),
]
edges = [
    edge("M184 200 H356", "e", "HTTPS", 270, 192),
    edge("M444 200 H616", "e", "/api/v1 + static", 530, 192),
    edge("M724 200 H976", "e", "SQL · pgx", 850, 192),
    edge("M724 220 H860 V500 H1016", "n", "alirkan berkas", 868, 450, anchor="start"),
    edge("M1050 264 V360 H680 V486", "n", "JOB · SKIP LOCKED", 800, 352),
    edge("M724 540 H900 V220 H1016", "n", "BARIS + NILAI + audit", 906, 380, anchor="start"),
    edge("M724 540 H1016", "n", "baca/tulis berkas", 870, 532),
]
legend = [("", "alur utama"), ("n", "koneksi internal")]
open(f"{OUT}/arsitektur-aplikasi.html", "w").write(page(
    "Arsitektur Aplikasi",
    "MATA SIDIK / penerapan",
    "Monolit modular + worker dalam satu image: dua perintah, PostgreSQL sebagai basis data + antrean job, berkas di volume lokal.",
    w, h, groups, nodes, edges, legend,
    "tanpa Redis · tanpa penyimpanan objek · Docker Compose: api, worker, postgres, reverse proxy"))

print("technical drafts written")
