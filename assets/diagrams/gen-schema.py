#!/usr/bin/env python3
"""Generate the database schema stack-diagram.

The diagram is a table catalog: every table is a card with its columns,
types, and key notes. Important foreign-key paths are drawn as dashed arrows.
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
* { box-sizing: border-box; }
.schema-card { position:absolute; background:var(--color-paper); border:1px solid var(--color-rule-2);
  border-radius:6px; box-shadow:0 1px 2px oklch(24% 0.02 258 / .06); overflow:hidden; }
.schema-card .t-head { border-bottom:1px solid var(--color-rule-2); background:var(--color-paper-2);
  padding:6px 8px; }
.schema-card .t-name { font-family:var(--font-mono); font-size:10px; font-weight:700; color:var(--color-ink); }
.schema-card .t-desc { font-family:var(--font-mono); font-size:6.6px; color:var(--color-muted); margin-top:2px; line-height:1.35; }
.schema-card table { width:100%; border-collapse:collapse; font-family:var(--font-mono); font-size:6.7px; }
.schema-card th { text-align:left; font-weight:500; color:var(--color-muted); padding:3px 6px;
  border-bottom:1px solid var(--color-rule); font-size:6px; letter-spacing:.03em; text-transform:uppercase; }
.schema-card td { padding:2.4px 6px; border-bottom:1px solid var(--color-rule); vertical-align:top; line-height:1.25; }
.schema-card tbody tr:last-child td { border-bottom:0; }
.schema-card .col { color:var(--color-ink); white-space:nowrap; }
.schema-card .type { color:var(--color-muted); white-space:nowrap; }
.schema-card .note { color:var(--color-muted); }
.schema-card .key { color:var(--color-accent); font-weight:700; }
.col-head { position:absolute; font-family:var(--font-mono); font-size:8px; font-weight:700;
  letter-spacing:.08em; text-transform:uppercase; color:var(--color-accent); }
.legend-box { position:absolute; background:var(--color-paper); border:1px solid var(--color-rule-2);
  border-radius:6px; padding:8px 10px; font-family:var(--font-mono); font-size:7px; color:var(--color-muted); line-height:1.7; }
.legend-box b { color:var(--color-ink); }
.rel { fill:none; stroke:var(--color-accent); stroke-width:1.1; stroke-dasharray:3 3; opacity:.55; }
.rel-label { font-family:var(--font-mono); font-size:6.1px; fill:var(--color-muted); letter-spacing:.02em; }
"""

OUT = os.environ.get("MATA_DIAGRAM_DRAFTS", os.path.join(BASE, "drafts"))
os.makedirs(OUT, exist_ok=True)

CARD_W = 440
GAP = 14
COL_X = [30, 490, 950, 1410]
COL_HEAD = ["MASTER INTI", "TAKSONOMI & SUBMISSION", "SUBMISSION DETAIL & NILAI", "PROVENANCE, FAKTOR & AUDIT"]

# (column, type, note)
TABLES = {
    "users": (0, "Pengguna, username, password hash, dan posisi struktural.", [
        ("id", "TEXT", "PK"),
        ("nama", "TEXT", "nama tampilan"),
        ("username", "TEXT", "UNIQUE · login"),
        ("password_hash", "TEXT", "bcrypt"),
        ("role", "TEXT", "role dasar"),
        ("faskes_id", "TEXT", "FK → faskes"),
        ("bidang_id", "TEXT", "FK → bidang"),
        ("position_id", "TEXT", "FK → positions"),
        ("last_login_at", "TEXT", ""),
        ("aktif", "INTEGER", "soft delete"),
        ("created_at", "TEXT", "RFC3339"),
    ]),
    "positions": (0, "Jabatan struktural; role efektif user.", [
        ("id", "TEXT", "PK"),
        ("kode", "TEXT", "UNIQUE"),
        ("nama", "TEXT", "nama jabatan"),
        ("role", "TEXT", "role efektif"),
        ("level", "INTEGER", "level approval"),
        ("bidang_id", "TEXT", "FK → bidang"),
        ("faskes_id", "TEXT", "FK → faskes"),
        ("aktif", "INTEGER", ""),
        ("created_at", "TEXT", ""),
    ]),
    "auth_session": (0, "Refresh token PASETO; hash di DB.", [
        ("id", "TEXT", "PK"),
        ("user_id", "TEXT", "FK → users"),
        ("refresh_hash", "TEXT", "UNIQUE · sha256"),
        ("user_agent", "TEXT", ""),
        ("created_at", "TEXT", ""),
        ("expires_at", "TEXT", ""),
        ("revoked_at", "TEXT", ""),
    ]),
    "wilayah": (0, "Hierarki desa → kecamatan → kabupaten.", [
        ("id", "TEXT", "PK"),
        ("kode", "TEXT", "UNIQUE"),
        ("nama", "TEXT", ""),
        ("level", "TEXT", "desa|kecamatan|kabupaten"),
        ("parent_id", "TEXT", "FK → wilayah"),
    ]),
    "faskes": (0, "Puskesmas/faskes pengupload.", [
        ("id", "TEXT", "PK"),
        ("kode", "TEXT", "UNIQUE"),
        ("nama", "TEXT", ""),
        ("wilayah_id", "TEXT", "FK → wilayah"),
        ("jenis", "TEXT", "puskesmas"),
        ("aktif", "INTEGER", "soft delete"),
    ]),
    "bidang": (0, "Bidang/timkerja pemilik data.", [
        ("id", "TEXT", "PK"),
        ("kode", "TEXT", "UNIQUE"),
        ("nama", "TEXT", ""),
    ]),
    "periode": (0, "Bulan pelaporan; mengikat upload.", [
        ("id", "TEXT", "PK"),
        ("tahun", "INTEGER", ""),
        ("bulan", "INTEGER", "1–12"),
        ("kode", "TEXT", "UNIQUE · 2026-05"),
        ("status", "TEXT", "open|locked"),
        ("is_late_allowed", "INTEGER", "late policy"),
    ]),
    "master_release": (0, "Bundel master per tahun.", [
        ("id", "TEXT", "PK"),
        ("tahun", "INTEGER", ""),
        ("status", "TEXT", "draft|published|archived"),
        ("published_at", "TEXT", ""),
        ("published_by", "TEXT", "FK → users"),
    ]),
    "workflow_version": (0, "Versi workflow yang di-pin ke submission.", [
        ("id", "TEXT", "PK"),
        ("kode", "TEXT", "default_3_level"),
        ("version_no", "INTEGER", ""),
        ("steps_json", "TEXT", "definisi step"),
        ("status", "TEXT", "published"),
        ("published_at", "TEXT", ""),
    ]),
    "file_object": (0, "Metadata file asli di volume lokal.", [
        ("id", "TEXT", "PK"),
        ("original_name", "TEXT", ""),
        ("sha256", "TEXT", "hash konten"),
        ("path", "TEXT", "path lokal"),
        ("mime", "TEXT", ""),
        ("size", "INTEGER", "byte"),
        ("uploaded_by", "TEXT", "FK → users"),
        ("scan_status", "TEXT", "pending|clean|infected"),
        ("created_at", "TEXT", ""),
    ]),

    "label": (1, "Taksonomi + node IKU; hierarki lewat parent_id.", [
        ("id", "TEXT", "PK"),
        ("kode", "TEXT", "UNIQUE per release"),
        ("nama", "TEXT", ""),
        ("tipe", "TEXT", "IKU|IKK|SPM|RIBK|…"),
        ("parent_id", "TEXT", "FK → label"),
        ("master_release_id", "TEXT", "FK → master_release"),
        ("aggregation_method", "TEXT", "RATIO_OF_SUMS|…"),
        ("allow_direct_input", "INTEGER", "leaf input"),
        ("is_rollup_only", "INTEGER", "induk rollup"),
        ("pj_bidang_id", "TEXT", "FK → bidang"),
        ("iku_owner_id", "TEXT", "FK → users"),
        ("level", "INTEGER", "kedalaman"),
        ("path", "TEXT", "path hierarki"),
        ("aktif", "INTEGER", ""),
        ("berlaku_dari", "TEXT", ""),
        ("berlaku_sampai", "TEXT", ""),
    ]),
    "indikator": (1, "Definisi indikator yang diukur.", [
        ("id", "TEXT", "PK"),
        ("kode", "TEXT", "UNIQUE · IND-001"),
        ("nama", "TEXT", ""),
        ("satuan", "TEXT", "satuan baku"),
        ("definisi", "TEXT", "definisi operasional"),
        ("numerator_label", "TEXT", ""),
        ("denominator_label", "TEXT", ""),
        ("sifat_input", "TEXT", "puskesmas|kabupaten"),
        ("frekuensi", "TEXT", "bulanan|tahunan"),
        ("pj_bidang_id", "TEXT", "FK → bidang"),
        ("pj_timker", "TEXT", ""),
        ("aktif", "INTEGER", ""),
    ]),
    "indikator_label": (1, "Many-to-many indikator ↔ label.", [
        ("indikator_id", "TEXT", "PK/FK → indikator"),
        ("label_id", "TEXT", "PK/FK → label"),
        ("bobot", "REAL", "default 1"),
        ("is_primary", "INTEGER", ""),
    ]),
    "iku_komponen": (1, "Graf kalkulasi induk → anak/indikator.", [
        ("id", "TEXT", "PK"),
        ("parent_label_id", "TEXT", "FK → label"),
        ("tipe_komponen", "TEXT", "CHILD_IKU|INDIKATOR"),
        ("komponen_id", "TEXT", "label/indikator id"),
        ("bobot", "REAL", ""),
        ("sort", "INTEGER", ""),
    ]),
    "target": (1, "Target periode/scope; bulanan = tahunan/12.", [
        ("id", "TEXT", "PK"),
        ("label_id", "TEXT", "FK → label"),
        ("indikator_id", "TEXT", "FK → indikator"),
        ("periode_id", "TEXT", "FK → periode"),
        ("scope_type", "TEXT", "KABUPATEN|FASKES|…"),
        ("scope_id", "TEXT", ""),
        ("target_tahunan", "REAL", ""),
        ("target_bulanan", "REAL", "dihitung"),
    ]),
    "capaian_tahunan": (1, "Capaian historis tahunan; pembanding dashboard.", [
        ("id", "TEXT", "PK"),
        ("label_id", "TEXT", "FK → label"),
        ("indikator_id", "TEXT", "FK → indikator"),
        ("tahun", "INTEGER", "2025"),
        ("scope_type", "TEXT", "KABUPATEN"),
        ("scope_id", "TEXT", ""),
        ("nilai", "REAL", ""),
        ("catatan", "TEXT", ""),
        ("updated_by", "TEXT", "FK → users"),
        ("created_at", "TEXT", ""),
        ("updated_at", "TEXT", ""),
    ]),
    "dashboard_row": (1, "Kontrak baris output SPM.", [
        ("id", "TEXT", "PK"),
        ("kode", "TEXT", "UNIQUE · 01"),
        ("label_id", "TEXT", "FK → label"),
        ("indikator_id", "TEXT", "FK → indikator"),
        ("urutan", "INTEGER", ""),
        ("target", "REAL", "default 100"),
        ("baseline_json", "TEXT", "capaian TW1"),
        ("format_json", "TEXT", "format tampilan"),
    ]),
    "submission": (1, "Satu pengajuan BNBA; status workflow.", [
        ("id", "TEXT", "PK"),
        ("kode", "TEXT", "UNIQUE · S-2026-05-001"),
        ("tipe", "TEXT", "BNBA"),
        ("uploader_id", "TEXT", "FK → users"),
        ("faskes_id", "TEXT", "FK → faskes"),
        ("periode_id", "TEXT", "FK → periode"),
        ("file_id", "TEXT", "FK → file_object"),
        ("file_hash", "TEXT", "sha256"),
        ("status", "TEXT", "7 status formal"),
        ("step", "TEXT", "DRAFT_PKM|…"),
        ("master_release_id", "TEXT", "FK → master_release"),
        ("workflow_version_id", "TEXT", "FK → workflow_version"),
        ("idempotency_key", "TEXT", "UNIQUE partial"),
        ("created_at", "TEXT", ""),
        ("updated_at", "TEXT", ""),
    ]),

    "submission_revision": (2, "Revisi immutable; file lama tidak ditimpa.", [
        ("id", "TEXT", "PK"),
        ("submission_id", "TEXT", "FK → submission"),
        ("revision_no", "INTEGER", "UNIQUE per submission"),
        ("parent_revision_id", "TEXT", "FK → revision"),
        ("file_id", "TEXT", "FK → file_object"),
        ("created_by", "TEXT", "FK → users"),
        ("created_at", "TEXT", ""),
    ]),
    "submission_group": (2, "Pecahan approval per pemilik data.", [
        ("id", "TEXT", "PK"),
        ("submission_id", "TEXT", "FK → submission"),
        ("owner_bidang_id", "TEXT", "FK → bidang"),
        ("primary_label_id", "TEXT", "FK → label"),
        ("status", "TEXT", "active"),
        ("current_step", "TEXT", "step"),
    ]),
    "submission_row": (2, "Baris hasil parsing + klasifikasi.", [
        ("id", "TEXT", "PK"),
        ("submission_id", "TEXT", "FK → submission"),
        ("group_id", "TEXT", "FK → submission_group"),
        ("row_no", "INTEGER", "UNIQUE per submission"),
        ("kode_iku_raw", "TEXT", "dari file"),
        ("kode_indikator_raw", "TEXT", "dari file"),
        ("matched_indikator_id", "TEXT", "FK → indikator"),
        ("matched_label_id", "TEXT", "FK → label"),
        ("scope_type", "TEXT", "FASKES|DESA|…"),
        ("scope_code", "TEXT", ""),
        ("nilai", "REAL", ""),
        ("numerator", "REAL", ""),
        ("denominator", "REAL", ""),
        ("classification_status", "TEXT", "MATCHED|…"),
        ("error_message", "TEXT", "pesan per baris"),
        ("created_at", "TEXT", ""),
    ]),
    "submission_review": (2, "Keputusan tiap level approval.", [
        ("id", "TEXT", "PK"),
        ("submission_id", "TEXT", "FK → submission"),
        ("group_id", "TEXT", "FK → submission_group"),
        ("step", "TEXT", "VERIFY_DINAS|…"),
        ("reviewer_id", "TEXT", "FK → users"),
        ("action", "TEXT", "submit|approve|return"),
        ("notes", "TEXT", "wajib saat return"),
        ("details_json", "TEXT", "diff/komentar"),
        ("created_at", "TEXT", ""),
    ]),
    "job": (2, "Antrean background SQLite tanpa Redis.", [
        ("id", "TEXT", "PK"),
        ("type", "TEXT", "parse_submission"),
        ("payload_json", "TEXT", "submission_id/file_id"),
        ("dedupe_key", "TEXT", "opsional"),
        ("status", "TEXT", "queued|running|done|failed|dead"),
        ("attempts", "INTEGER", ""),
        ("max_attempts", "INTEGER", "default 5"),
        ("run_after", "TEXT", "backoff"),
        ("locked_by", "TEXT", "worker id"),
        ("locked_at", "TEXT", ""),
        ("heartbeat_at", "TEXT", ""),
        ("last_error", "TEXT", ""),
        ("progress_json", "TEXT", "progress UI"),
        ("created_at", "TEXT", ""),
        ("finished_at", "TEXT", ""),
    ]),
    "indicator_value": (2, "Nilai dasar satu indikator per baris.", [
        ("id", "TEXT", "PK"),
        ("submission_row_id", "TEXT", "FK/UNIQUE → row"),
        ("indikator_id", "TEXT", "FK → indikator"),
        ("periode_id", "TEXT", "FK → periode"),
        ("faskes_id", "TEXT", "FK → faskes"),
        ("scope_type", "TEXT", "FASKES|…"),
        ("scope_id", "TEXT", ""),
        ("numerator", "REAL", "jangan dibulatkan"),
        ("denominator", "REAL", "jangan dibulatkan"),
        ("value", "REAL", "num/den×100"),
        ("status", "TEXT", "COMPLETE|…"),
        ("data_version", "TEXT", "versi data"),
        ("created_at", "TEXT", ""),
    ]),

    "calculation_run": (3, "Satu eksekusi kalkulasi.", [
        ("id", "TEXT", "PK"),
        ("periode_id", "TEXT", "FK → periode"),
        ("scope_type", "TEXT", "KABUPATEN|FASKES"),
        ("scope_id", "TEXT", ""),
        ("trigger", "TEXT", "final|seed|target"),
        ("status", "TEXT", "running|done|failed"),
        ("started_at", "TEXT", ""),
        ("finished_at", "TEXT", ""),
        ("error", "TEXT", ""),
    ]),
    "iku_value": (3, "Nilai IKU materialized per label/periode/scope.", [
        ("id", "TEXT", "PK"),
        ("label_id", "TEXT", "FK → label"),
        ("periode_id", "TEXT", "FK → periode"),
        ("scope_type", "TEXT", "KABUPATEN|FASKES"),
        ("scope_id", "TEXT", ""),
        ("numerator", "REAL", "Σ numerator"),
        ("denominator", "REAL", "Σ denominator"),
        ("value", "REAL", "×100"),
        ("status", "TEXT", "COMPLETE|PARTIAL|NO_DATA|STALE"),
        ("formula_version", "TEXT", "v1"),
        ("calculation_run_id", "TEXT", "FK → calculation_run"),
        ("source_hash", "TEXT", "hash file sumber"),
        ("data_version", "TEXT", "v2026-05-xxxx"),
        ("updated_at", "TEXT", ""),
        ("UNIQUE", "—", "label+periode+scope"),
    ]),
    "iku_value_detail": (3, "Provenance komponen angka induk.", [
        ("id", "TEXT", "PK"),
        ("iku_value_id", "TEXT", "FK → iku_value"),
        ("component_type", "TEXT", "CHILD_IKU|INDIKATOR"),
        ("component_id", "TEXT", "label/indikator id"),
        ("component_label_id", "TEXT", "FK → label"),
        ("component_value", "REAL", ""),
        ("numerator", "REAL", ""),
        ("denominator", "REAL", ""),
        ("weight", "REAL", "default 1"),
        ("contribution", "REAL", "nilai × bobot"),
    ]),
    "faktor": (3, "Analisis kualitatif per IKU value.", [
        ("id", "TEXT", "PK"),
        ("iku_value_id", "TEXT", "FK/UNIQUE → iku_value"),
        ("periode_id", "TEXT", "FK → periode"),
        ("pendukung", "TEXT", "faktor positif"),
        ("penghambat", "TEXT", "faktor negatif"),
        ("created_by", "TEXT", "FK → users"),
        ("updated_at", "TEXT", ""),
    ]),
    "tindak_lanjut": (3, "TL ber-PJ dan berstatus.", [
        ("id", "TEXT", "PK"),
        ("iku_value_id", "TEXT", "FK → iku_value"),
        ("periode_id", "TEXT", "FK → periode"),
        ("masalah", "TEXT", ""),
        ("rekomendasi", "TEXT", ""),
        ("pj_user_id", "TEXT", "FK → users"),
        ("status", "TEXT", "open|in_progress|done|cancelled"),
        ("created_by", "TEXT", "FK → users"),
        ("created_at", "TEXT", ""),
        ("updated_at", "TEXT", ""),
    ]),
    "audit_log": (3, "Audit append-only; who/when/what/why.", [
        ("id", "TEXT", "PK"),
        ("actor_id", "TEXT", "FK → users"),
        ("action", "TEXT", "submission.finalized|…"),
        ("entity_type", "TEXT", "submission|iku_value|…"),
        ("entity_id", "TEXT", ""),
        ("before_json", "TEXT", "diff"),
        ("after_json", "TEXT", "diff"),
        ("reason", "TEXT", ""),
        ("request_id", "TEXT", "korelasi log"),
        ("impersonator_id", "TEXT", "FK → users · actingAs"),
        ("created_at", "TEXT", ""),
    ]),
}

RELATIONS = [
    ("positions", "users", "position_id"),
    ("users", "auth_session", "refresh"),
    ("faskes", "submission", "faskes_id"),
    ("periode", "submission", "periode_id"),
    ("file_object", "submission", "file_id"),
    ("master_release", "label", "release"),
    ("label", "label", "parent"),
    ("indikator", "indikator_label", "label"),
    ("label", "indikator_label", "indikator"),
    ("label", "iku_komponen", "parent"),
    ("label", "capaian_tahunan", "capaian"),
    ("indikator", "capaian_tahunan", "capaian"),
    ("users", "capaian_tahunan", "update"),
    ("label", "dashboard_row", "iku"),
    ("submission", "submission_row", "baris"),
    ("submission", "submission_review", "review"),
    ("submission_row", "indicator_value", "nilai"),
    ("indikator", "indicator_value", "indikator"),
    ("calculation_run", "iku_value", "run"),
    ("iku_value", "iku_value_detail", "detail"),
    ("iku_value", "faktor", "faktor"),
    ("iku_value", "tindak_lanjut", "TL"),
    ("users", "audit_log", "aktor"),
]


def card_height(cols):
    return 64 + len(cols) * 14.4


def table_card(name, desc, cols, x, y):
    rows = ""
    for col, typ, note in cols:
        note_class = "note key" if note in ("PK", "UNIQUE") or note.startswith("PK") or note.startswith("FK") else "note"
        rows += (
            f'<tr><td class="col">{html.escape(col)}</td>'
            f'<td class="type">{html.escape(typ)}</td>'
            f'<td class="{note_class}">{html.escape(note)}</td></tr>'
        )
    return f'''<div class="schema-card" style="left:{x}px;top:{y}px;width:{CARD_W}px;height:{card_height(cols)}px;">
  <div class="t-head"><div class="t-name">{html.escape(name)}</div><div class="t-desc">{html.escape(desc)}</div></div>
  <table><thead><tr><th>Kolom</th><th>Tipe</th><th>Keterangan</th></tr></thead><tbody>{rows}</tbody></table>
</div>'''


positions = {}
cards = []
col_y = [110, 110, 110, 110]
for name, (col, desc, cols) in TABLES.items():
    x = COL_X[col]
    y = col_y[col]
    positions[name] = (x, y, CARD_W, card_height(cols), col)
    cards.append(table_card(name, desc, cols, x, y))
    col_y[col] = y + card_height(cols) + GAP

H = max(col_y) + 60
W = COL_X[-1] + CARD_W + 30

headers = "".join(
    f'<div class="col-head" style="left:{COL_X[i]}px;top:78px;">{html.escape(COL_HEAD[i])}</div>'
    for i in range(4)
)


def relation_path(a, b, label):
    ax, ay, aw, ah, acol = positions[a]
    bx, by, bw, bh, bcol = positions[b]
    if acol == bcol:
        # Route through the right gutter of this column.
        gutter = ax + aw + 10
        x1, y1 = ax + aw, ay + ah / 2
        x2, y2 = bx + bw, by + bh / 2
        mid_y = (y1 + y2) / 2
        d = f"M{x1} {y1} H{gutter} V{y2} H{x2}"
        lx, ly = gutter + 4, mid_y
    else:
        if acol < bcol:
            x1, y1 = ax + aw, ay + ah / 2
            x2, y2 = bx, by + bh / 2
            gutter = ax + aw + 14
        else:
            x1, y1 = ax, ay + ah / 2
            x2, y2 = bx + bw, by + bh / 2
            gutter = ax - 14
        d = f"M{x1} {y1} H{gutter} V{y2} H{x2}"
        lx, ly = gutter + 4, (y1 + y2) / 2
    return f'''      <path class="rel" d="{d}" marker-end="url(#rel-ar)"/>
      <text class="rel-label" x="{lx}" y="{ly - 3}">{html.escape(label)}</text>'''


rel_svg = "\n".join(relation_path(a, b, label) for a, b, label in RELATIONS)

legend = f'''<div class="legend-box" style="left:30px;top:{H - 46}px;width:1860px;">
  <b>Relasi:</b> garis putus-putus = FK utama. &nbsp; <b>PK</b> primary key · <b>FK</b> foreign key · <b>UNIQUE</b> kunci unik.
  &nbsp; Tabel <b>job</b> memakai payload JSON untuk menunjuk submission; tabel <b>dashboard_row</b> memetakan output SPM ke label IKU.
  &nbsp; Enum status mengikuti docs/rules.md §6.1; nilai status IKU mengikuti §7.2.
</div>'''

html_doc = f'''<!DOCTYPE html>
<html lang="id">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Skema Database MATA SIDIK</title>
<style>{STYLE}</style>
</head>
<body>
<div class="page">
  <header class="head">
    <div>
      <div class="kicker">MATA SIDIK / database</div>
      <h1>Skema Database</h1>
    </div>
    <p class="lede">Katalog tabel, kolom, kunci, dan relasi utama. Implementasi pilot memakai SQLite; tipe dan nama kolom disiapkan agar mudah dimigrasikan ke PostgreSQL.</p>
  </header>
  <div class="scroller">
  <div class="canvas" style="width:{W}px;height:{H}px;">
    {headers}
    <svg class="edges" viewBox="0 0 {W} {H}" style="width:100%;height:100%" aria-hidden="true">
      <defs>
        <marker id="rel-ar" viewBox="0 0 10 10" refX="8" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
          <path d="M0 0 L10 5 L0 10 z" fill="context-stroke"/>
        </marker>
      </defs>
      {rel_svg}
    </svg>
    {chr(10).join(cards)}
    {legend}
  </div>
  </div>
  <div class="legend">
    <span><span class="sw"></span>relasi kunci</span>
    <span>SQLite pilot · repository pattern · siap migrasi PostgreSQL</span>
  </div>
  <div class="foot">
    <span>schema.md · 30 tabel · 1 file SQLite + volume file lokal</span>
    <span>icons: Lucide · dibuat dengan <a href="https://github.com/caru-ini/stack-diagram" target="_blank" rel="noopener">stack-diagram</a></span>
  </div>
</div>
</body>
</html>
'''
open(os.path.join(OUT, "schema-db.html"), "w").write(html_doc)
print("schema draft written")
