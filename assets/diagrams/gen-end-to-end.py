#!/usr/bin/env python3
"""Generate the large end-to-end IKU data-flow diagram for MATA SIDIK."""
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
.ic { width:13px; height:13px; vertical-align:-2px; }
.big { font-family:var(--font-display); font-weight:650; font-size:14px; line-height:1.25; }
.mono { font-family:var(--font-mono); }
.sub { font-family:var(--font-mono); font-size:8px; color:var(--color-muted); line-height:1.6; }
.mini { width:100%; border-collapse:collapse; font-family:var(--font-mono); font-size:7.6px; }
.mini th { text-align:left; font-weight:500; color:var(--color-muted); border-bottom:1px solid var(--color-rule-2); padding:2.5px 2px; }
.mini td { padding:2.5px 2px; border-bottom:1px solid var(--color-rule); vertical-align:top; }
.dash { width:100%; border-collapse:collapse; font-family:var(--font-mono); font-size:9.6px; }
.dash th { text-align:left; font-weight:500; color:var(--color-muted); border-bottom:1px solid var(--color-rule-2);
  padding:4px 6px; font-size:8.2px; letter-spacing:.04em; text-transform:uppercase; white-space:nowrap; }
.dash td { padding:4px 6px; border-bottom:1px solid var(--color-rule); }
.dash tbody tr:last-child td { border-bottom:0; }
.dash .num { text-align:right; }
.dash .c { text-align:center; }
.dash .ind { color:var(--color-ink); }
.dash .st { font-weight:700; font-size:10px; }
.st.ok { color:oklch(50% .15 150); }
.st.warn { color:oklch(68% .15 75); }
.st.bad { color:oklch(58% .19 25); }
.badge { display:inline-block; font-family:var(--font-mono); font-size:7px; line-height:1.35;
  padding:1px 5px; border-radius:20px; border:1px solid; white-space:nowrap; }
.b-ok { color:oklch(44% .12 150); border-color:oklch(70% .12 150); background:oklch(96% .04 150); }
.b-warn { color:oklch(50% .13 75); border-color:oklch(75% .13 75); background:oklch(96% .05 75); }
.b-final { color:var(--color-accent); border-color:var(--color-accent-soft); background:var(--color-accent-tint); }
.b-muted { color:var(--color-muted); border-color:var(--color-rule-2); background:var(--color-paper-2); }
.pick { display:inline-block; border:1px solid var(--color-accent); color:var(--color-accent);
  border-radius:4px; padding:1px 5px; font-family:var(--font-mono); font-size:7px; }
.calc { position:absolute; background:var(--color-paper); border:1px solid var(--color-rule-2);
  border-radius:6px; padding:7px 9px 7px 12px; }
.calc:before { content:""; position:absolute; left:0; top:7px; bottom:7px; width:3px;
  border-radius:2px; background:var(--color-g1); }
.calc.parent:before { background:var(--color-accent); }
.calc .label { font-family:var(--font-mono); font-size:7.4px; color:var(--color-muted); }
.calc .formula { font-family:var(--font-mono); font-size:8.4px; color:var(--color-ink-2); margin-top:2px; }
.calc .result { font-family:var(--font-display); font-size:20px; font-weight:700; color:var(--color-accent); line-height:1.1; }
.calc .result.small { font-size:16px; }
.tree { font-family:var(--font-mono); font-size:8.2px; line-height:1.9; }
.tree .parent { font-family:var(--font-display); font-size:10.5px; font-weight:650; }
.tree .row { display:flex; justify-content:space-between; gap:8px; align-items:center; }
.tree .indent { padding-left:14px; }
.strip { position:absolute; background:var(--color-paper-2); border:1px solid var(--color-rule-2);
  border-radius:8px; padding:8px 12px; display:flex; align-items:center; justify-content:space-between; gap:16px; }
.strip .l { display:flex; align-items:center; gap:9px; }
.strip .k { font-family:var(--font-mono); font-size:7.5px; letter-spacing:.08em; color:var(--color-accent); }
.strip .t { font-family:var(--font-display); font-size:10px; font-weight:600; }
.strip .r { font-family:var(--font-mono); font-size:7.5px; color:var(--color-muted); text-align:right; line-height:1.7; }
.chain { display:flex; flex-wrap:wrap; align-items:center; gap:4px; font-family:var(--font-mono); font-size:7.6px; }
.chain .box { border:1px solid var(--color-rule-2); background:var(--color-paper); border-radius:4px; padding:2px 6px; }
.chain .ar { color:var(--color-accent); }
.trace { font-family:var(--font-mono); font-size:7.4px; }
.trace .line { display:flex; align-items:center; gap:4px; margin-bottom:4px; }
.trace .tag { width:58px; flex:none; color:var(--color-muted); font-size:7.2px; }
.trace .box { border:1px solid var(--color-rule-2); background:var(--color-paper); border-radius:4px; padding:2px 6px; }
.trace .box.hi { border-color:var(--color-accent); color:var(--color-accent); background:var(--color-accent-tint); }
.trace .ar { color:var(--color-accent); }
.trace .extra { color:var(--color-muted); }
.src { font-family:var(--font-mono); font-size:6.6px; color:var(--color-muted);
  border:1px solid var(--color-rule-2); border-radius:3px; padding:0 3px; margin-left:4px; }
.spin { display:inline-block; width:12px; height:12px; border:2px solid var(--color-rule-2);
  border-top-color:var(--color-accent); border-radius:50%; vertical-align:-2px; }
.divider { border-top:1px solid var(--color-rule); margin:5px 0; }
.g-modes { position:absolute; top:9px; right:10px; display:flex; gap:3px; }
.mode { font-family:var(--font-mono); font-size:6.2px; letter-spacing:.07em; line-height:1.3;
  padding:1px 4px; border-radius:20px; border:1px solid; white-space:nowrap; }
.mode.user { color:var(--color-ink-2); border-color:var(--color-rule-2); background:var(--color-paper); }
.mode.system { color:var(--color-accent); border-color:var(--color-accent-soft); background:var(--color-accent-tint); }
.flow { width:100%; border-collapse:collapse; font-family:var(--font-mono); font-size:7.6px; }
.flow th { text-align:left; font-weight:500; color:var(--color-muted); border-bottom:1px solid var(--color-rule-2);
  padding:3px 3px; font-size:7px; letter-spacing:.04em; text-transform:uppercase; }
.flow td { padding:5px 3px; border-bottom:1px solid var(--color-rule); vertical-align:top; line-height:1.45; }
.flow tbody tr:last-child td { border-bottom:0; }
.flow .stage { font-weight:700; color:var(--color-ink); }
.flow .actor { color:var(--color-muted); font-size:7px; }
.flow .next { color:var(--color-accent); white-space:nowrap; }
.flow .time { color:var(--color-muted); font-size:6.8px; }
.steps { display:flex; flex-direction:column; gap:4px; }
.step { display:flex; gap:6px; font-size:8.1px; line-height:1.45; }
.step .n { flex:none; width:14px; height:14px; border-radius:50%; background:var(--color-accent-tint);
  color:var(--color-accent); font-family:var(--font-mono); font-size:7px;
  display:flex; align-items:center; justify-content:center; }
.flow .badge { font-size:6.8px; padding:0 4px; }
.legend .mode { font-size:8px; padding:1px 6px; }
"""

OUT = os.environ.get("MATA_DIAGRAM_DRAFTS", os.path.join(BASE, "drafts"))
os.makedirs(OUT, exist_ok=True)

W, H = 1860, 1810
STAGE_TOP, STAGE_H = 110, 820
STAGE_W = 420
XS = [30, 490, 950, 1410]


MODE_LABEL = {"user": "USER", "system": "SISTEM"}


def group(x, y, w, h, label, cls="gb", icon=None, modes=()):
    ic = f'<img src="%%ICON_{icon}%%" alt="">' if icon else ""
    mode_html = ""
    if modes:
        mode_html = '<div class="g-modes">' + "".join(
            f'<span class="mode {m}">{MODE_LABEL[m]}</span>' for m in modes
        ) + "</div>"
    return f'''    <div class="group {cls}" style="left:{x}px;top:{y}px;width:{w}px;height:{h}px;">
      <div class="g-label">{ic}{html.escape(label)}</div>{mode_html}
    </div>'''


def node(cx, top, name, role, icon, cls=""):
    return f'''    <div class="node" style="left:{cx - 66}px;top:{top}px">
      <div class="tile {cls}"><img src="%%ICON_{icon}%%" alt=""></div>
      <div class="name"><span>{html.escape(name)}</span></div><div class="role"><span>{html.escape(role)}</span></div>
    </div>'''


def panel(x, y, w, h, title, body, icon=None):
    ic = f'<img src="%%ICON_{icon}%%" alt="">' if icon else ""
    return f'''    <div class="panel" style="left:{x}px;top:{y}px;width:{w}px;height:{h}px;">
      <div class="p-title">{ic}{html.escape(title)}</div>
      {body}
    </div>'''


def badge(text, cls="b-muted"):
    return f'<span class="badge {cls}">{html.escape(text)}</span>'


def status(kind):
    sym = "✓" if kind == "ok" else "●"
    return f'<span class="st {kind}">{sym}</span>'


def mini(headers, rows):
    th = "".join(f"<th>{html.escape(h)}</th>" for h in headers)
    body = "".join("<tr>" + "".join(f"<td>{c}</td>" for c in row) + "</tr>" for row in rows)
    return f'<table class="mini"><thead><tr>{th}</tr></thead><tbody>{body}</tbody></table>'


def flow_table(rows):
    head = "<tr><th>Tahap / aktor</th><th>Aksi</th><th>Status</th><th>Berikutnya</th></tr>"
    body = ""
    for stage, actor, action, st, time_, nxt in rows:
        body += (
            "<tr>"
            f'<td><div class="stage">{html.escape(stage)}</div>'
            f'<div class="actor">{html.escape(actor)}</div></td>'
            f"<td>{html.escape(action)}</td>"
            f'<td>{st}<div class="time">{html.escape(time_)}</div></td>'
            f'<td class="next">{html.escape(nxt)}</td>'
            "</tr>"
        )
    return (
        '<table class="flow"><colgroup><col style="width:100px"><col style="width:80px">'
        '<col style="width:80px"><col></colgroup>'
        f"<thead>{head}</thead><tbody>{body}</tbody></table>"
    )


def calc(x, y, w, h, label, formula, result, note="", parent=False, small=False):
    return f'''    <div class="calc {"parent" if parent else ""}" style="left:{x}px;top:{y}px;width:{w}px;height:{h}px;">
      <div class="label">{html.escape(label)}</div>
      <div class="result {"small" if small else ""}">{result}</div>
      <div class="formula">{formula}</div>
      <div class="sub">{note}</div>
    </div>'''


def edge(x1, x2, y, label):
    mid = (x1 + x2) / 2
    return f'''      <path class="e" d="M{x1} {y} H{x2}" marker-end="url(#ar)"/>
      <text class="el a" x="{mid}" y="{y - 8}" text-anchor="middle">{html.escape(label)}</text>'''


def edge_down(x, y1, y2, label):
    mid = (y1 + y2) / 2
    return f'''      <path class="e" d="M{x} {y1} V{y2}" marker-end="url(#ar)"/>
      <text class="el a" x="{x + 9}" y="{mid + 3}" text-anchor="start">{html.escape(label)}</text>'''


groups = [
    group(XS[0], STAGE_TOP, STAGE_W, STAGE_H, "1 · Upload Puskesmas", "ga", "upload", ("user",)),
    group(XS[1], STAGE_TOP, STAGE_W, STAGE_H, "2 · Validasi & Preview", "gb", "klasifikasi", ("system", "user")),
    group(XS[2], STAGE_TOP, STAGE_W, STAGE_H, "3 · Workflow 3 Level", "gb", "workflow", ("user", "system")),
    group(XS[3], STAGE_TOP, STAGE_W, STAGE_H, "4 · Kalkulasi IKU", "gb", "calc", ("system",)),
    group(30, 970, 1800, 580, "5 · Dashboard IKU & SPM", "gb", "dash", ("system", "user")),
]

nodes = [
    node(XS[0] + STAGE_W / 2, 142, "Upload BNBA", "Siti · PKM Taliwang", "upload"),
    node(XS[1] + STAGE_W / 2, 142, "Validasi + Klasifikasi", "sistem → user", "klasifikasi"),
    node(XS[2] + STAGE_W / 2, 142, "Approval Berjenjang", "user + lock sistem", "workflow"),
    node(XS[3] + STAGE_W / 2, 142, "Hitung Nilai IKU", "otomatis · FINAL", "calc"),
    node(930, 985, "Dashboard SPM", "sistem → user", "dash"),
]

edges = [
    edge(XS[0] + STAGE_W, XS[1], 560, "upload"),
    edge(XS[1] + STAGE_W, XS[2], 560, "submit"),
    edge(XS[2] + STAGE_W, XS[3], 560, "final"),
    edge_down(XS[3] + STAGE_W / 2, 930, 970, "publish"),
]

x0 = XS[0] + 15
x1 = XS[1] + 15
x2 = XS[2] + 15
x3 = XS[3] + 15

panels = []
panels.append(panel(x0, 270, 390, 112, "File upload", f'''
  <div class="big">BNBA_2026-05.xlsx</div>
  <div class="sub">Mei 2026 · 3 baris · 18 KB<br>uploader: Siti (PKM Taliwang)<br>hash: 9f2c…a71</div>
''', "upload"))
file_table = mini(
    ["row", "kode_iku", "indikator", "num", "den", "periode"],
    [
        ["R1", "IKU-2026-001.01", "IND-001", "450", "500", "2026-05"],
        ["R2", "IKU-2026-001.02", "IND-002", "380", "400", "2026-05"],
        ["R3", "IKU-2026-001", "IND-003", "20", "20", "2026-05"],
    ],
)
column_help = mini(
    ["Kolom", "Arti & fungsi"],
    [
        ["row", "nomor baris; rujukan error/koreksi per baris"],
        ["kode_iku", "IKU tujuan; harus leaf, menentukan rollup"],
        ["indikator", "identitas indikator yang diukur"],
        ["num", "pembilang: realisasi; wajib untuk rasio"],
        ["den", "penyebut: sasaran; wajib untuk rasio"],
        ["periode", "Mei 2026; mengikat ke master release"],
    ],
)
panels.append(panel(x0, 400, 390, 300, "Isi file BNBA", file_table + f'''
  <div class="divider"></div>
  <div class="p-title" style="margin-bottom:3px">Arti kolom</div>
  {column_help}
  <div class="divider"></div>
  <div class="sub">R1–R2: IKU anak (leaf), valid.<br>
  R3: IKU induk (rollup-only) — tidak boleh input langsung.<br>
  <b>Kunci:</b> kode_iku + kode_indikator + periode.<br>
  Kode tidak dikenal → karantina, bukan auto-create.</div>
''', "laporan"))
panels.append(panel(x0, 720, 390, 80, "Hasil awal", f'''
  <div>{badge("2 VALID", "b-ok")} {badge("1 PERLU KOREKSI", "b-warn")}</div>
  <div class="sub" style="margin-top:5px">File tetap tersimpan; belum masuk workflow.</div>
''', "check"))
panels.append(panel(x0, 820, 390, 90, "File asli", f'''
  <div class="sub">Hash + file asli disimpan. Upload ulang<br>dibuat sebagai revisi baru, bukan timpa.</div>
''', "disk"))

panels.append(panel(x1, 270, 390, 142, "State 1 · Halaman hasil (proses)", f'''
  <div style="display:flex;align-items:center;gap:8px">
    <span class="spin"></span><span class="big" style="font-size:12px">Sedang dikerjakan…</span>
  </div>
  <div class="sub" style="margin-top:7px">Validasi file ✓ &nbsp; Klasifikasi ✓<br>
  Simpan sementara ✓ &nbsp; Siapkan preview <span style="color:var(--color-accent)">berjalan</span></div>
  <div class="sub">Halaman hasil refresh otomatis (polling 2 dtk).</div>
''', "loading"))
panels.append(panel(x1, 427, 390, 280, "State 2 · Preview — IKU bisa diubah", mini(
    ["row", "IKU leaf", "nilai", "status"],
    [
        ["R1", ".01 Ibu Hamil", "450/500", badge("MATCHED", "b-ok")],
        ["R2", ".02 Ibu Bersalin", "380/400", badge("MATCHED", "b-ok")],
        ["R3", ".03 Bayi Baru Lahir", "20/20", badge("MATCHED", "b-ok")],
    ],
) + f'''
  <div class="divider"></div>
  <div class="sub">R3 semula terdeteksi IKU induk; dipilih ulang ke <b>.03</b>.<br>
  Nilai 20/20 masuk ke .03. Dropdown hanya berisi<br>
  IKU leaf yang valid untuk indikator tersebut.</div>
  <div style="margin-top:5px"><span class="pick">Ubah IKU ▾</span> <span class="sub">hanya saat preview</span></div>
  <div class="divider"></div>
  <div class="sub"><b>Status:</b> MATCHED masuk preview; MISMATCH,<br>
  AMBIGUOUS, DUPLIKAT → karantina. IKU_ROLLUP_ONLY<br>
  tidak boleh input langsung.</div>
''', "tree"))
panels.append(panel(x1, 722, 390, 80, "Aksi", f'''
  <div>{badge("SEMUA BARIS VALID", "b-ok")}</div>
  <div class="sub" style="margin-top:4px">Tombol <b>Kirim</b> aktif setelah tidak ada blocker.</div>
''', "check"))
panels.append(panel(x1, 817, 390, 93, "Batas preview", f'''
  <div class="sub">Ubah IKU hanya di preview. Setelah submit,<br>perubahan mengikuti revisi + approval.</div>
''', "warning"))

flow_rows = [
    ("DRAFT_PKM", "Siti · Puskesmas/PJ", "simpan / kirim", badge("TERKIRIM"), "10:02", "VERIFY_DINAS"),
    ("VERIFY_DINAS", "Rina · Dinas", "approve / return", badge("DISETUJUI", "b-ok"), "10:15", "VALIDATE_KABID"),
    ("VALIDATE_KABID", "dr. Ahmad · Kabid", "approve / return", badge("FINAL", "b-final"), "10:30", "FINAL · sistem"),
    ("FINAL", "Sistem", "lock + hitung", badge("TERKUNCI", "b-final"), "10:31", "Kalkulasi IKU"),
]
panels.append(panel(x2, 270, 390, 240, "Alur status & tahap berikutnya", flow_table(flow_rows) + f'''
  <div class="divider"></div>
  <div class="sub">Urutan = alur normal. Return → kembali ke pengirim<br>
  dengan catatan wajib. Contoh: R3 dikoreksi ke .03<br>
  sebelum dikirim, lalu disetujui hingga FINAL.</div>
''', "workflow"))
panels.append(panel(x2, 525, 390, 110, "Hasil workflow", f'''
  <div>{badge("FINAL · TERKUNCI", "b-final")}</div>
  <div class="sub" style="margin-top:6px">Setiap aksi menyimpan aktor, waktu, catatan,<br>
  diff baris, dan audit log. Koreksi setelah final<br>
  hanya lewat <b>special change</b> + dual approval.</div>
''', "audit"))
panels.append(panel(x2, 650, 390, 100, "Faktor & narasi (Kabid)", f'''
  <div>{badge("INPUT KABID", "b-muted")} {badge("PASCA-FINAL", "b-final")}</div>
  <div class="sub" style="margin-top:5px">Faktor pendukung/penghambat = analisis kualitatif; bukan TL.<br>
  Diisi Kabid/Sekretaris setelah data final; bukan input angka.<br>
  Tidak mengubah kalkulasi.</div>
''', "faktor"))
panels.append(panel(x2, 765, 390, 140, "Jika ditolak", f'''
  <div>{badge("RETURN + CATATAN", "b-warn")}</div>
  <div class="sub" style="margin-top:5px">Status kembali ke level pengirim; resubmission<br>
  menjadi revisi baru dengan diff lama vs baru.<br>
  Catatan revisi wajib per baris/field.</div>
  <div class="divider"></div>
  <div>{badge("SPECIAL CHANGE", "b-muted")} <span class="sub">alasan + dual approval</span></div>
''', "kembali"))

panels.append(panel(x3, 270, 390, 120, "Langkah otomatis", f'''
  <div class="steps">
    <div class="step"><span class="n">1</span><div>Ambil <b>numerator & denominator</b> dari 3 baris final.</div></div>
    <div class="step"><span class="n">2</span><div>Hitung tiap anak: 450/500, 380/400, 20/20 → 90,00%, 95,00%, 100,00%.</div></div>
    <div class="step"><span class="n">3</span><div>Jumlahkan untuk induk: <b>850/920 × 100 = 92,39%</b>; bukan rata-rata persen.</div></div>
  </div>
  <div style="margin-top:6px">{badge("TANPA INPUT MANUAL", "b-final")}</div>
''', "calc"))
panels.append(panel(x3, 405, 390, 150, "Rincian komponen", mini(
    ["Komponen", "Num / Den", "Capaian"],
    [
        [".01 Ibu Hamil", "450 / 500", "90,00%"],
        [".02 Ibu Bersalin", "380 / 400", "95,00%"],
        [".03 Bayi Baru Lahir", "20 / 20", "100,00%"],
        ["<b>Jumlah</b>", "<b>850 / 920</b>", "<b>92,39%</b>"],
    ],
) + f'''
  <div class="divider"></div>
  <div class="sub">Nilai anak dihitung lebih dulu; induk memakai<br>
  penjumlahan numerator/denominator, bukan rata-rata.</div>
''', "tree"))
panels.append(calc(x3, 570, 390, 110, "IKU-2026-001 · INDUK", "(450+380+20) / (500+400+20) × 100", "92,39%", "850 / 920 · RATIO_OF_SUMS, bukan rata-rata persentase", parent=True, small=False))
panels.append(panel(x3, 695, 390, 110, "Target & status", f'''
  <div class="sub">Target 100% · gap <b>-7,61 pp</b></div>
  <div style="margin-top:6px">{badge("PERLU PERHATIAN", "b-warn")} {badge("DATA FINAL", "b-final")}</div>
  <div class="sub" style="margin-top:5px">Status dihitung dari target periode; capaian di bawah target ditandai perlu perhatian.</div>
''', "target"))
panels.append(panel(x3, 820, 390, 90, "Yang disimpan", f'''
  <div>{badge("COMPLETE", "b-ok")} <span class="sub">3/3 komponen final</span></div>
  <div class="sub" style="margin-top:5px">Simpan 850/920 + calculation_run_id + source_hash.<br>
  Komponen kurang → PARTIAL, bukan nol.</div>
''', "disk"))

panels.append(panel(45, 1335, 590, 195, "Pohon IKU", f'''
  <div style="display:flex;gap:14px">
    <div class="tree" style="flex:1.1">
      <div class="row"><span class="parent">IKU-2026-001 Kesehatan Ibu</span></div>
      <div class="row"><span style="font-size:14px;font-weight:700;color:var(--color-accent)">92,39%</span>{badge("PERHATIAN", "b-warn")}</div>
      <div class="row"><span class="sub">target 100% · gap -7,61 pp</span></div>
      <div class="divider"></div>
      <div class="row"><span class="indent">├─ .01 Ibu Hamil</span><span>90,00% {badge("PERHATIAN", "b-warn")}</span></div>
      <div class="row"><span class="indent">├─ .02 Ibu Bersalin</span><span>95,00% {badge("OK", "b-ok")}</span></div>
      <div class="row"><span class="indent">└─ .03 Bayi Baru Lahir</span><span>100,00% {badge("OK", "b-ok")}</span></div>
    </div>
    <div style="flex:1;border-left:1px solid var(--color-rule);padding-left:12px">
      <div class="sub">Klik angka → buka <b>IKU_VALUE / DETAIL</b>:<br>
      450/500 + 380/400 + 20/20 = 850/920.</div>
      <div class="divider"></div>
      <div class="sub">Mapping output:<br>
      .01 → SPM 01<br>.02 → SPM 02<br>.03 → SPM 03</div>
    </div>
  </div>
  <div class="sub" style="margin-top:6px">Nilai final dibaca dari IKU_VALUE; drill-down menuju file.</div>
''', "tree"))
panels.append(panel(655, 1335, 380, 195, "Metadata dashboard", f'''
  <div class="sub">data version: <b>v2026-05-001</b><br>
  dihitung: 10:31:02 · sumber: 3 baris final<br>
  status: <b>FINAL</b> · tampil ke pimpinan</div>
  <div style="margin-top:6px">{badge("DRILL-DOWN AKTIF", "b-final")}</div>
  <div class="divider"></div>
  <div class="sub">Setiap baris membawa iku_value_id,<br>
  indicator_value_id, source_hash, data_version,<br>
  calculation_run_id.</div>
  <div class="sub" style="margin-top:4px">Dashboard membaca nilai materialized,<br>
  bukan menghitung ulang saat request.</div>
''', "dash"))
panels.append(panel(1055, 1335, 380, 195, "Jejak angka", f'''
  <div class="sub">92,39% ← 850/920 ← 3 IKU anak ←<br>
  file hash 9f2c…a71 ← uploader Siti.</div>
  <div class="divider"></div>
  <div class="sub"><b>R1</b> .01 450/500 → SPM 01<br>
  <b>R2</b> .02 380/400 → SPM 02<br>
  <b>R3</b> .03 20/20 → SPM 03</div>
  <div class="divider"></div>
  <div class="sub">Klik baris SPM → drill-down sampai file.</div>
''', "audit"))
panels.append(panel(1455, 1335, 360, 195, "Pimpinan", f'''
  <div>{badge("ARAHAN + TL", "b-final")}</div>
  <div class="sub" style="margin-top:5px">Lihat tren + faktor + gap; beri arahan dan<br>
  tetapkan tindak lanjut ber-PJ dan status.</div>
  <div class="p-title" style="margin-top:7px;margin-bottom:3px">Alur TL</div>
  <div class="chain">
    <span class="box">Masalah</span><span class="ar">→</span>
    <span class="box">Rekomendasi</span><span class="ar">→</span>
    <span class="box">PJ</span><span class="ar">→</span>
    <span class="box">Status</span>
  </div>
  <div class="sub" style="margin-top:4px">Ditetapkan pimpinan; ditindaklanjuti<br>
  Kabid/Sekretaris; dipantau.</div>
  <div class="divider"></div>
  <div class="sub">Klik baris SPM → telusuri sampai file.<br>
  Publikasi agregasi mengikuti IKU owner.</div>
''', "keputusan"))

# Format dashboard SPM: 12 baris yang dapat ditelusuri sampai file.
dash_rows = [
    ("01", "Pelayanan Kesehatan Ibu Hamil", "100%", "100%", "ok", "27,3%", "ok", "90,00%", "warn", "R1"),
    ("02", "Pelayanan Kesehatan Ibu Bersalin", "100%", "96%", "warn", "26,2%", "ok", "95,00%", "ok", "R2"),
    ("03", "Pelayanan Kesehatan Bayi Baru Lahir", "100%", "100%", "ok", "28,1%", "ok", "100,00%", "ok", "R3"),
]
dash_body = ""
for no, nama, target, c2025, s2025, ctw1, stw1, cmei, smei, src in dash_rows:
    dash_body += (
        f'<tr><td class="num">{no}</td><td class="ind">{html.escape(nama)}'
        f'<span class="src">{src}</span></td>'
        f'<td class="num">{target}</td><td class="num">{c2025}</td><td class="c">{status(s2025)}</td>'
        f'<td class="num">{ctw1}</td><td class="c">{status(stw1)}</td>'
        f'<td class="num">{cmei}</td><td class="c">{status(smei)}</td></tr>'
    )

dashboard = f'''    <div class="panel" style="left:45px;top:1090px;width:1770px;height:225px;border-left:3px solid var(--color-accent);">
      <div class="p-title"><img src="%%ICON_dash%%" alt="">Format dashboard SPM — 12 baris yang dapat ditelusuri
        <span class="badge b-final" style="margin-left:auto">DATA FINAL · TRACEABLE</span></div>
      <div style="display:flex;gap:18px;align-items:flex-start">
        <table class="dash" style="width:1340px;flex:none">
          <colgroup>
            <col style="width:38px"><col><col style="width:70px">
            <col style="width:105px"><col style="width:52px">
            <col style="width:120px"><col style="width:52px">
            <col style="width:110px"><col style="width:52px">
          </colgroup>
          <thead><tr>
            <th class="num">No</th><th>Indikator</th><th class="num">Target</th>
            <th class="num">Capaian 2025</th><th class="c">Status</th>
            <th class="num">Capaian TW1 2026</th><th class="c">Status</th>
            <th class="num">Capaian MEI 2026</th><th class="c">Status</th>
          </tr></thead>
          <tbody>{dash_body}</tbody>
        </table>
        <div style="flex:1;border-left:1px solid var(--color-rule);padding-left:14px">
          <div class="p-title" style="margin-bottom:4px">Cara baca</div>
          <div class="sub" style="line-height:1.9">
            <span class="st ok">✓</span> sesuai target<br>
            <span class="st warn">●</span> perlu perhatian<br>
            <span class="st bad">●</span> kritis
          </div>
          <div class="divider"></div>
          <div class="sub"><b>R1/R2/R3</b> = baris file BNBA.<br>
          Capaian MEI 2026 = hasil kalkulasi IKU leaf.<br>
          Capaian 2025/TW1 = baseline dashboard SPM.</div>
        </div>
      </div>
      <div class="sub" style="margin-top:6px"><b>R1</b> .01 450/500 → SPM 01 · <b>R2</b> .02 380/400 → SPM 02 · <b>R3</b> .03 20/20 → SPM 03.
      Tiga baris terisi; sembilan baris lain mengikuti kontrak yang sama. Klik baris → drill-down ke IKU_VALUE / DETAIL → SUBMISSION_ROW → FILE_OBJECT.</div>
    </div>'''


top_strip = f'''    <div class="strip" style="left:30px;top:30px;width:1800px;height:56px;">
      <div class="l"><img src="%%ICON_table%%" class="ic" alt="">
        <div><div class="k">ALUR DATA IKU · MEI 2026</div>
        <div class="t">PKM Taliwang · 3 baris upload · 3 IKU anak · 1 induk · 3 baris SPM</div></div>
      </div>
      <div class="r">alur proses: upload → validasi → workflow → kalkulasi → dashboard SPM</div>
    </div>'''

lineage = f'''    <div class="panel" style="left:30px;top:1570px;width:1800px;height:200px;">
      <div class="p-title"><img src="%%ICON_audit%%" alt="">Jejak 3 baris SPM — dari file upload ke dashboard</div>
      <div class="trace" style="margin-top:3px">
        <div class="line"><span class="tag">SPM 01</span>
          <span class="box">FILE R1</span><span class="ar">→</span>
          <span class="box">SUB_ROW R1 · MATCHED</span><span class="ar">→</span>
          <span class="box">IND-001 · 450/500</span><span class="ar">→</span>
          <span class="box">IKU .01 · 90,00%</span><span class="ar">→</span>
          <span class="box hi">SPM 01 · 90,00%</span>
        </div>
        <div class="line"><span class="tag">SPM 02</span>
          <span class="box">FILE R2</span><span class="ar">→</span>
          <span class="box">SUB_ROW R2 · MATCHED</span><span class="ar">→</span>
          <span class="box">IND-002 · 380/400</span><span class="ar">→</span>
          <span class="box">IKU .02 · 95,00%</span><span class="ar">→</span>
          <span class="box hi">SPM 02 · 95,00%</span>
        </div>
        <div class="line"><span class="tag">SPM 03</span>
          <span class="box">FILE R3</span><span class="ar">→</span>
          <span class="box">SUB_ROW R3 · MATCHED <span class="badge b-warn">EX-INDUK</span></span><span class="ar">→</span>
          <span class="box">IND-003 · 20/20</span><span class="ar">→</span>
          <span class="box">IKU .03 · 100,00%</span><span class="ar">→</span>
          <span class="box hi">SPM 03 · 100,00%</span>
        </div>
      </div>
      <div class="divider"></div>
      <div class="sub"><b>Rollup induk:</b> IKU-2026-001 · 92,39% · 850/920 ← .01 450/500 + .02 380/400 + .03 20/20.
      &nbsp; <b>Sumber:</b> FILE_OBJECT 9f2c…a71 → SUBMISSION S-2026-05-001 → REVISION 1.</div>
      <div class="sub"><b>Tabel:</b> SUBMISSION_ROW · INDICATOR_VALUE · IKU_VALUE · IKU_VALUE_DETAIL · AUDIT_LOG &nbsp;·&nbsp;
      <b>Aktor:</b> uploader Siti · verifikator Rina · Kabid dr. Ahmad.</div>
    </div>'''

html_doc = f'''<!DOCTYPE html>
<html lang="id">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Alur Data IKU End-to-End — MATA SIDIK</title>
<style>{STYLE}</style>
</head>
<body>
<div class="page">

  <header class="head">
    <div>
      <div class="kicker">MATA SIDIK / alur data IKU</div>
      <h1>Dari Upload BNBA sampai Dashboard IKU</h1>
    </div>
    <p class="lede">Satu alur utuh: upload 3 baris, klasifikasi ke IKU leaf, koreksi di preview, approval 3 level, kalkulasi induk–anak, dan 3 baris output SPM yang bisa ditelusuri sampai file. Label <b>USER</b>/<b>SISTEM</b> menandai siapa yang mengerjakan tiap tahap.</p>
  </header>

  <div class="scroller">
  <div class="canvas" style="width:{W}px;height:{H}px;">

{top_strip}

{chr(10).join(groups)}

    <svg class="edges" viewBox="0 0 {W} {H}" style="width:100%;height:100%" aria-hidden="true">
      <defs>
        <marker id="ar" viewBox="0 0 10 10" refX="8" refY="5" markerWidth="6.5" markerHeight="6.5" orient="auto-start-reverse">
          <path d="M0 0 L10 5 L0 10 z" fill="context-stroke"/>
        </marker>
      </defs>
{chr(10).join(edges)}
    </svg>

{chr(10).join(nodes)}

{chr(10).join(panels)}

{dashboard}

{lineage}

  </div>
  </div>

  <div class="legend">
    <span><span class="sw"></span>alur data</span>
    <span><span class="sw n dash"></span>koreksi / pengembalian (di luar alur utama)</span>
    <span><span class="mode user">USER</span> aksi pengguna</span>
    <span><span class="mode system">SISTEM</span> proses otomatis</span>
  </div>

  <div class="foot">
    <span>angka: Mei 2026 · PKM Taliwang · 3 baris SPM yang ditelusuri ke file</span>
    <span>icons: Lucide · dibuat dengan <a href="https://github.com/caru-ini/stack-diagram" target="_blank" rel="noopener">stack-diagram</a></span>
  </div>

</div>
</body>
</html>
'''

open(os.path.join(OUT, "alur-data-iku.html"), "w").write(html_doc)
print("end-to-end draft written")
