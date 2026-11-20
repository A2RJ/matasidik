#!/usr/bin/env python3
"""Generate BNBA v2 templates: Kejadian, Cakupan, Tahunan.

Reads the pilot SQLite master and writes three .xlsx files under
assets/templates/bnba/v2/. No application code is changed.

Requires: python3 -m pip install openpyxl
Run from any directory: python3 assets/templates/bnba/v2/gen_templates.py
"""

from __future__ import annotations

import sqlite3
from pathlib import Path

from openpyxl import Workbook
from openpyxl.comments import Comment
from openpyxl.formatting.rule import CellIsRule, FormulaRule
from openpyxl.styles import Alignment, Border, Font, PatternFill, Protection, Side
from openpyxl.utils import get_column_letter
from openpyxl.worksheet.datavalidation import DataValidation

ROOT = Path(__file__).resolve().parents[5]
DB = ROOT / "backend/data/matasidik.db"
OUT = ROOT / "assets/templates/bnba/arsip/v2"

# Palette
NAVY = "1F3864"
BLUE = "2F5597"
LIGHT_BLUE = "D9E2F3"
PALE_BLUE = "EEF3FA"
GREEN = "548235"
LIGHT_GREEN = "E2EFDA"
YELLOW = "FFF2CC"
ORANGE = "F4B183"
LIGHT_ORANGE = "FCE4D6"
GRAY = "F2F2F2"
MID_GRAY = "BFBFBF"
WHITE = "FFFFFF"
RED = "C00000"
LIGHT_RED = "F8CBAD"

THIN = Side(style="thin", color=MID_GRAY)
BORDER = Border(left=THIN, right=THIN, top=THIN, bottom=THIN)
FONT = "Calibri"
BODY = Font(name=FONT, size=10)
BODY_BOLD = Font(name=FONT, size=10, bold=True)
SMALL = Font(name=FONT, size=9)
SMALL_GRAY = Font(name=FONT, size=9, color="595959")
HEADER = Font(name=FONT, size=10, bold=True, color=WHITE)
TITLE = Font(name=FONT, size=18, bold=True, color=WHITE)
SUBTITLE = Font(name=FONT, size=10, color="404040")


def fill(color: str) -> PatternFill:
    return PatternFill("solid", fgColor=color)


def header_row(ws, row: int, headers: list[str], color: str = NAVY, height: int = 32):
    for i, h in enumerate(headers, 1):
        c = ws.cell(row=row, column=i, value=h)
        c.font = HEADER
        c.fill = fill(color)
        c.alignment = Alignment(horizontal="center", vertical="center", wrap_text=True)
        c.border = BORDER
    ws.row_dimensions[row].height = height


def body_cell(c, *, wrap=False, align="left", color=None, bold=False):
    c.font = BODY_BOLD if bold else BODY
    c.border = BORDER
    c.alignment = Alignment(horizontal=align, vertical="top", wrap_text=wrap)
    if color:
        c.fill = fill(color)


def section(ws, row: int, title: str, color: str = BLUE):
    ws.merge_cells(start_row=row, start_column=2, end_row=row, end_column=3)
    c = ws.cell(row=row, column=2, value=title)
    c.font = Font(name=FONT, size=11, bold=True, color=WHITE)
    c.fill = fill(color)
    c.alignment = Alignment(vertical="center")
    ws.row_dimensions[row].height = 22


def kv(ws, row: int, label: str, value: str, height: int = 20, value_color=None):
    a = ws.cell(row=row, column=2, value=label)
    b = ws.cell(row=row, column=3, value=value)
    a.font = BODY_BOLD
    a.fill = fill(PALE_BLUE)
    a.border = BORDER
    a.alignment = Alignment(vertical="top", wrap_text=True)
    b.font = BODY
    b.border = BORDER
    b.alignment = Alignment(vertical="top", wrap_text=True)
    if value_color:
        b.fill = fill(value_color)
    ws.row_dimensions[row].height = height


def add_list_dv(ws, rng: str, formula: str, title: str, prompt: str):
    dv = DataValidation(type="list", formula1=formula, allow_blank=True, showDropDown=False)
    dv.errorTitle = "Nilai tidak valid"
    dv.error = "Pilih nilai dari sheet Referensi atau ikuti format yang diminta."
    dv.promptTitle = title
    dv.prompt = prompt
    ws.add_data_validation(dv)
    dv.add(rng)


def add_decimal_dv(ws, rng: str, title: str, prompt: str, *, greater_than_zero=False):
    dv = DataValidation(
        type="decimal",
        operator="greaterThan" if greater_than_zero else "greaterThanOrEqual",
        formula1="0",
        allow_blank=True,
    )
    dv.errorTitle = "Angka tidak valid"
    dv.error = "Isi angka >= 0." if not greater_than_zero else "Isi angka > 0."
    dv.promptTitle = title
    dv.prompt = prompt
    ws.add_data_validation(dv)
    dv.add(rng)


def protect(ws):
    p = ws.protection
    p.sheet = True
    p.objects = False
    p.scenarios = False
    p.formatCells = False
    p.formatColumns = False
    p.formatRows = False
    p.insertColumns = False
    p.insertRows = False
    p.insertHyperlinks = False
    p.deleteColumns = False
    p.deleteRows = False
    p.sort = False
    p.autoFilter = False
    p.pivotTables = False


def load_master():
    con = sqlite3.connect(DB)
    con.row_factory = sqlite3.Row
    q = con.execute
    categories = [dict(r) for r in q(
        "SELECT kode, nama, jenis FROM dashboard_category WHERE aktif=1 ORDER BY urutan, kode"
    )]
    labels = [dict(r) for r in q(
        """SELECT l.kode, l.nama, l.tipe, c.kode AS kategori_kode, c.nama AS kategori_nama,
                  p.kode AS parent_kode, l.level, l.allow_direct_input, l.is_rollup_only,
                  l.aggregation_method
           FROM label l
           LEFT JOIN dashboard_category c ON c.id = l.category_id
           LEFT JOIN label p ON p.id = l.parent_id
           WHERE l.aktif=1
           ORDER BY l.level, l.kode"""
    )]
    indicators = [dict(r) for r in q(
        """SELECT i.kode, i.nama, i.satuan, i.numerator_label, i.denominator_label,
                  l.kode AS iku_kode, l.nama AS iku_nama, c.kode AS kategori_kode,
                  c.nama AS kategori_nama, l.aggregation_method,
                  l.allow_direct_input, l.is_rollup_only
           FROM indikator i
           JOIN indikator_label il ON il.indikator_id = i.id AND il.is_primary = 1
           JOIN label l ON l.id = il.label_id
           LEFT JOIN dashboard_category c ON c.id = l.category_id
           WHERE i.aktif=1 AND l.aktif=1
           ORDER BY l.kode"""
    )]
    wilayah = [dict(r) for r in q(
        """SELECT w.kode, w.nama, w.level, p.kode AS parent_kode
           FROM wilayah w LEFT JOIN wilayah p ON p.id = w.parent_id
           ORDER BY w.kode"""
    )]
    faskes = [dict(r) for r in q(
        "SELECT id, kode, nama, jenis FROM faskes WHERE aktif=1 ORDER BY kode"
    )]
    periods = [dict(r) for r in q(
        "SELECT kode, tahun, bulan, status FROM periode ORDER BY kode"
    )]
    releases = [dict(r) for r in q(
        "SELECT tahun, status FROM master_release ORDER BY tahun DESC"
    )]
    con.close()
    return categories, labels, indicators, wilayah, faskes, periods, releases


def ref_table(ws, start_row: int, title: str, headers: list[str], rows: list[list], color: str = BLUE):
    ws.merge_cells(start_row=start_row, start_column=1, end_row=start_row, end_column=len(headers))
    h = ws.cell(row=start_row, column=1, value=title)
    h.font = Font(name=FONT, size=11, bold=True, color=WHITE)
    h.fill = fill(color)
    h.alignment = Alignment(vertical="center")
    ws.row_dimensions[start_row].height = 22
    header_row(ws, start_row + 1, headers, color, 28)
    for rr, rowvals in enumerate(rows, start_row + 2):
        for cc, val in enumerate(rowvals, 1):
            cell = ws.cell(row=rr, column=cc, value=val)
            body_cell(cell, wrap=True, align="left")
        ws.row_dimensions[rr].height = 22
    return start_row + 2, start_row + 1 + len(rows)


def note(ws, row: int, text: str, *, color=YELLOW, height=36, span=3, start_col=1):
    ws.merge_cells(start_row=row, start_column=start_col, end_row=row, end_column=start_col + span - 1)
    c = ws.cell(row=row, column=start_col, value=text)
    c.font = BODY
    c.fill = fill(color)
    c.border = BORDER
    c.alignment = Alignment(vertical="top", wrap_text=True)
    ws.row_dimensions[row].height = height


# Config ---------------------------------------------------------------------

PATTERNS = [
    {
        "key": "kejadian",
        "file": "BNBA-1-Kejadian-Bulanan.xlsx",
        "title": "BNBA KEJADIAN BULANAN",
        "subtitle": "Untuk kejadian/count bulanan: kematian ibu, kematian balita, dan kejadian lain. Isi jumlah kasus, bukan persentase.",
        "headers": ["periode", "kode_faskes", "kode_desa", "kode_indikator", "jumlah", "catatan"],
        "input_cols": [1, 2, 3, 4, 5, 6],
        "numeric_cols": [5],
        "numeric_label": "jumlah",
        "unit": "Kasus",
        "tipe_input": "KEJADIAN",
        "frekuensi": "BULANAN",
        "is_monthly": True,
        "is_annual": False,
        "has_indicators": False,
        "examples": [
            ["2026-05", "PKM-TALIWANG", "5201.01.01", "CONTOH-KEJADIAN-IBU", 1,
             "ILUSTRASI — kode contoh, bukan master"],
            ["2026-05", "PKM-TALIWANG", "5201.01.02", "CONTOH-KEJADIAN-BALITA", 0,
             "ILUSTRASI — kode contoh, bukan master"],
        ],
    },
    {
        "key": "cakupan",
        "file": "BNBA-2-Cakupan-Bulanan.xlsx",
        "title": "BNBA CAKUPAN BULANAN",
        "subtitle": "Untuk cakupan/rasio bulanan: SPM dan IKU cakupan. Isi numerator dan denominator; hasil % dihitung sistem.",
        "headers": ["periode", "kode_faskes", "kode_desa", "kode_indikator", "numerator", "denominator", "catatan"],
        "input_cols": [1, 2, 3, 4, 5, 6, 7],
        "numeric_cols": [5, 6],
        "numeric_label": "numerator/denominator",
        "unit": "%",
        "tipe_input": "CAKUPAN",
        "frekuensi": "BULANAN",
        "is_monthly": True,
        "is_annual": False,
        "has_indicators": True,
        "examples": [
            ["2026-05", "PKM-TALIWANG", "5201.01.01", "IKU-2026-001.01", 92, 120,
             "ILUSTRASI — angka fiktif"],
            ["2026-05", "PKM-TALIWANG", "5201.01.01", "IKU-2026-001.02", 64, 80,
             "ILUSTRASI — angka fiktif"],
            ["2026-05", "PKM-TALIWANG", "5201.01.02", "SPM-2026-04", 60, 80,
             "ILUSTRASI — angka fiktif"],
        ],
    },
    {
        "key": "tahunan",
        "file": "BNBA-3-Tahunan.xlsx",
        "title": "DATA TAHUNAN",
        "subtitle": "Untuk indikator tahunan: UHH, AKIP, dokumen/laporan tahunan. Satu baris = satu indikator + satu scope + satu tahun.",
        "headers": ["tahun", "scope_type", "kode_scope", "kode_indikator", "nilai", "catatan"],
        "input_cols": [1, 2, 3, 4, 5, 6],
        "numeric_cols": [5],
        "numeric_label": "nilai",
        "unit": "sesuai indikator",
        "tipe_input": "TAHUNAN",
        "frekuensi": "TAHUNAN",
        "is_monthly": False,
        "is_annual": True,
        "has_indicators": False,
        "examples": [
            ["2026", "KABUPATEN", "5201", "CONTOH-UHH", 73.5,
             "ILUSTRASI — kode contoh, bukan master"],
            ["2026", "KABUPATEN", "5201", "CONTOH-AKIP", 85,
             "ILUSTRASI — kode contoh, bukan master"],
            ["2026", "KABUPATEN", "5201", "CONTOH-DOKUMEN", 12,
             "ILUSTRASI — kode contoh, bukan master"],
        ],
    },
]

def build_workbook(cfg, master):
    categories, labels, master_indicators, wilayah, faskes, periods, releases = master
    wb = Workbook()
    ws_data = wb.active
    ws_data.title = "Data"
    ws_petunjuk = wb.create_sheet("Petunjuk")
    ws_ref = wb.create_sheet("Referensi")
    ws_contoh = wb.create_sheet("Contoh")

    # Indicator rows for this pattern.
    indicator_rows = []
    if cfg["key"] == "cakupan":
        for i in master_indicators:
            indicator_rows.append([
                i["kode"], i["nama"], i["iku_kode"], i["iku_nama"], i["kategori_kode"] or "",
                i["satuan"], "%", "CAKUPAN", "BULANAN", "HIGHER_IS_BETTER", "RATIO_OF_SUMS", "DESA",
            ])

    ranges = write_referensi(ws_ref, cfg, master, indicator_rows)
    write_data(ws_data, cfg, ranges, indicator_rows)
    write_petunjuk(ws_petunjuk, cfg, master, indicator_rows)
    write_contoh(ws_contoh, cfg)

    wb.properties.title = f"{cfg['title']} — MATA SIDIK"
    wb.properties.subject = cfg["subtitle"]
    wb.properties.creator = "MATA SIDIK"
    wb.properties.description = (
        "Template v2 dari master pilot. Pisahkan pola input; jangan mencampur kejadian, cakupan, dan tahunan. "
        "Angka di sheet Contoh fiktif."
    )
    wb.properties.keywords = "BNBA, MATA SIDIK, template, kejadian, cakupan, tahunan"
    wb.active = 0
    out = OUT / cfg["file"]
    wb.save(out)
    print(f"Wrote {out} ({out.stat().st_size} bytes)")
    return out


def write_referensi(ws, cfg, master, indicator_rows):
    categories, labels, master_indicators, wilayah, faskes, periods, releases = master
    ws.sheet_view.showGridLines = False
    ws.sheet_properties.tabColor = GREEN
    widths = {"A": 22, "B": 34, "C": 18, "D": 24, "E": 16, "F": 16, "G": 16, "H": 22, "I": 14, "J": 18, "K": 18, "L": 14}
    for col, w in widths.items():
        ws.column_dimensions[col].width = w

    ws.merge_cells("A1:L1")
    t = ws["A1"]
    t.value = "REFERENSI MASTER — JANGAN DIUBAH"
    t.font = TITLE
    t.fill = fill(GREEN)
    t.alignment = Alignment(vertical="center")
    ws.row_dimensions[1].height = 32

    ws.merge_cells("A2:L2")
    c = ws["A2"]
    c.value = f"Sumber: master pilot `backend/data/matasidik.db` per 2026-11-19. Template: {cfg['title']}. Daftar ini untuk validasi, bukan data produksi."
    c.font = SUBTITLE
    c.fill = fill(LIGHT_GREEN)
    c.alignment = Alignment(vertical="center", wrap_text=True)
    ws.row_dimensions[2].height = 30

    ranges = {}
    row = 4

    periode_rows = [[p["kode"], f'Bulan {p["bulan"]:02d} tahun {p["tahun"]}'] for p in periods]
    if cfg["is_annual"]:
        # Annual periods are not in the pilot table yet; use one year row for the template.
        periode_rows = [["2026", "Tahun 2026 (contoh; periode tahunan belum ada di pilot)"]]
    start, end = ref_table(ws, row, "1. Periode", ["kode", "keterangan"], periode_rows)
    ranges["periode"] = (start, end)

    row = end + 3
    start, end = ref_table(
        ws, row, "2. Faskes", ["kode", "nama", "id_internal"],
        [[f["kode"], f["nama"], f["id"]] for f in faskes],
    )
    ranges["faskes"] = (start, end)

    row = end + 3
    start, end = ref_table(
        ws, row, "3. Wilayah", ["kode", "nama", "level", "kode_parent"],
        [[w["kode"], w["nama"], w["level"], w["parent_kode"] or "—"] for w in wilayah],
    )
    ranges["wilayah"] = (start, end)

    row = end + 3
    desa = [w for w in wilayah if w["level"] == "desa"]
    start, end = ref_table(
        ws, row, "4. Desa (dipakai dropdown kode_desa)", ["kode", "nama", "kode_kecamatan"],
        [[w["kode"], w["nama"], w["parent_kode"] or "—"] for w in desa],
    )
    ranges["desa"] = (start, end)

    row = end + 3
    scope_rows = []
    for w in wilayah:
        scope_type = "KABUPATEN" if w["level"] == "kabupaten" else w["level"].upper()
        scope_rows.append([w["kode"], w["nama"], scope_type, "Gunakan kode wilayah ini"])
    start, end = ref_table(
        ws, row, "5. Scope (dipakai dropdown kode_scope)", ["kode", "nama", "level", "catatan"],
        scope_rows,
    )
    ranges["scope"] = (start, end)

    row = end + 3
    start, end = ref_table(
        ws, row, "6. Satuan", ["kode", "nama", "contoh"],
        [
            ["%", "Persentase", "76,67%"],
            ["Kasus", "Jumlah kasus/kejadian", "2 kasus"],
            ["tahun", "Tahun", "73,5 tahun"],
            ["Poin", "Poin/indeks", "85 poin"],
            ["Rasio", "Rasio non-persen", "1 : 1.200"],
            ["Dokumen", "Jumlah dokumen", "12 dokumen"],
            ["Laporan", "Jumlah laporan", "4 laporan"],
            ["Berita Acara", "Jumlah berita acara", "2 berita acara"],
            ["orang", "Jumlah orang (numerator)", "92 orang"],
        ],
    )
    ranges["satuan"] = (start, end)

    row = end + 3
    start, end = ref_table(
        ws, row, "7. Tipe Input", ["kode", "nama", "pola"],
        [
            ["KEJADIAN", "Kejadian/count bulanan", "Desa + bulan + jumlah"],
            ["CAKUPAN", "Cakupan/rasio bulanan", "Desa + bulan + numerator/denominator"],
            ["TAHUNAN", "Tahunan", "Tahun + scope + nilai"],
        ],
    )
    ranges["tipe_input"] = (start, end)

    row = end + 3
    start, end = ref_table(
        ws, row, "8. Arah Indikator", ["kode", "nama"],
        [
            ["HIGHER_IS_BETTER", "Makin besar makin baik"],
            ["LOWER_IS_BETTER", "Makin kecil makin baik"],
        ],
    )
    ranges["arah"] = (start, end)

    row = end + 3
    start, end = ref_table(
        ws, row, "9. Agregasi", ["kode", "nama", "dipakai untuk"],
        [
            ["SUM", "Jumlahkan", "Kejadian/count"],
            ["RATIO_OF_SUMS", "Jumlahkan num/den lalu rasio", "Cakupan/persentase"],
            ["LAST_VALUE", "Ambil nilai terakhir", "UHH/AKIP tahunan"],
        ],
    )
    ranges["agregasi"] = (start, end)

    row = end + 3
    if indicator_rows:
        start, end = ref_table(
            ws, row, "10. Indikator untuk template ini",
            ["kode", "nama", "kode_iku", "nama_iku", "kategori", "satuan_master", "satuan_hasil",
             "tipe_input", "frekuensi", "arah", "agregasi", "scope"],
            indicator_rows,
        )
        ranges["indikator"] = (start, end)
        row = end + 3
    else:
        ws.merge_cells(start_row=row, start_column=1, end_row=row, end_column=12)
        c = ws.cell(row=row, column=1, value="10. Indikator untuk template ini — BELUM ADA DI MASTER PILOT")
        c.font = Font(name=FONT, size=11, bold=True, color=WHITE)
        c.fill = fill("C00000")
        c.alignment = Alignment(vertical="center")
        ws.row_dimensions[row].height = 22
        note(
            ws, row + 1,
            "Master pilot belum memiliki indikator untuk pola ini. Tambahkan indikator lewat Master/Pengajuan IKU "
            "sebelum memakai template. Kode di tabel 'Contoh kode desain' hanya ilustrasi dan TIDAK VALID untuk upload.",
            color=YELLOW, span=12, height=44,
        )
        ranges["indikator"] = None
        row += 3

    row = max(row, 4)
    if not indicator_rows:
        example_rows = []
        if cfg["key"] == "kejadian":
            example_rows = [
                ["CONTOH-KEJADIAN-IBU", "Jumlah Kematian Ibu", "CONTOH-IKU-KEMATIAN", "Contoh IKU Kematian",
                 "IKU", "Kasus", "Kasus", "KEJADIAN", "BULANAN", "LOWER_IS_BETTER", "SUM", "DESA"],
                ["CONTOH-KEJADIAN-BALITA", "Jumlah Kematian Balita", "CONTOH-IKU-KEMATIAN", "Contoh IKU Kematian",
                 "IKU", "Kasus", "Kasus", "KEJADIAN", "BULANAN", "LOWER_IS_BETTER", "SUM", "DESA"],
            ]
        elif cfg["key"] == "tahunan":
            example_rows = [
                ["CONTOH-UHH", "Usia Harapan Hidup", "CONTOH-IKU-UHH", "Contoh IKU UHH",
                 "IKU", "tahun", "tahun", "TAHUNAN", "TAHUNAN", "HIGHER_IS_BETTER", "LAST_VALUE", "KABUPATEN"],
                ["CONTOH-AKIP", "Nilai AKIP", "CONTOH-IKU-AKIP", "Contoh IKU AKIP",
                 "IKU", "Poin", "Poin", "TAHUNAN", "TAHUNAN", "HIGHER_IS_BETTER", "LAST_VALUE", "KABUPATEN"],
                ["CONTOH-DOKUMEN", "Jumlah Dokumen", "CONTOH-IKU-DOKUMEN", "Contoh IKU Dokumen",
                 "IKU", "Dokumen", "Dokumen", "TAHUNAN", "TAHUNAN", "HIGHER_IS_BETTER", "SUM", "KABUPATEN"],
            ]
        start, end = ref_table(
            ws, row, "11. Contoh kode desain — BUKAN MASTER, JANGAN DIUNGGAH",
            ["kode", "nama", "kode_iku", "nama_iku", "kategori", "satuan_master", "satuan_hasil",
             "tipe_input", "frekuensi", "arah", "agregasi", "scope"],
            example_rows, color="7F7F7F",
        )
        row = end + 3

    if cfg["key"] == "cakupan":
        ref_note = (
            "Catatan: IND-001 s.d. IND-003 di master pilot memiliki satuan 'orang' (numerator), tetapi hasil IKU/SPM-nya "
            "adalah persentase. Template Cakupan menampilkan satuan hasil '%'. Normalisasi satuan master perlu diputuskan Dinas. "
            "Jangan mengetik kode baru di file ini; tambahkan lewat Master."
        )
    elif cfg["key"] == "kejadian":
        ref_note = (
            "Catatan: template ini hanya untuk kejadian/count. Indikator kematian ibu/balita belum ada di master pilot; "
            "tambahkan lewat Master sebelum memakai template. Kode CONTOH-* bukan master dan tidak boleh diunggah. "
            "Jangan mengetik kode baru di file ini."
        )
    else:
        ref_note = (
            "Catatan: template ini untuk data tahunan. UHH adalah snapshot tahunan dan tidak boleh dijumlahkan seperti count. "
            "Indikator tahunan belum ada di master pilot; tambahkan lewat Master. Kode CONTOH-* bukan master dan tidak boleh diunggah."
        )
    ws.merge_cells(start_row=row, start_column=1, end_row=row, end_column=12)
    c = ws.cell(row=row, column=1, value=ref_note)
    c.font = SMALL
    c.fill = fill(YELLOW)
    c.border = BORDER
    c.alignment = Alignment(vertical="top", wrap_text=True)
    ws.row_dimensions[row].height = 56
    return ranges


def write_data(ws, cfg, ranges, indicator_rows):
    headers = cfg["headers"]
    header_row(ws, 1, headers, NAVY, 34)
    ws.sheet_view.showGridLines = False
    ws.sheet_properties.tabColor = NAVY

    widths = {
        "periode": 12, "tahun": 10, "kode_faskes": 18, "kode_desa": 14,
        "scope_type": 14, "kode_scope": 14, "kode_indikator": 22,
        "jumlah": 12, "numerator": 14, "denominator": 14, "nilai": 12, "catatan": 34,
    }
    cl = {h: get_column_letter(i) for i, h in enumerate(headers, 1)}
    for h in headers:
        ws.column_dimensions[cl[h]].width = widths.get(h, 14)

    comments = {
        "periode": "Wajib YYYY-MM. Pilih periode bulanan dari Referensi.",
        "tahun": "Wajib YYYY. Jangan isi bulan/triwulan.",
        "kode_faskes": "Kode faskes pengunggah. Satu file = satu faskes.",
        "kode_desa": "Kode desa resmi. Nama desa dapat dilihat di Referensi.",
        "scope_type": "DESA, KECAMATAN, atau KABUPATEN. Faskes bukan scope; faskes adalah pengunggah.",
        "kode_scope": "Kode wilayah resmi sesuai scope_type.",
        "kode_indikator": "Pilih indikator dari Referensi. Kode IKU/nama/satuan ditentukan sistem.",
        "jumlah": "Jumlah kasus/kejadian. Satuan: Kasus. Bukan persentase.",
        "numerator": "Jumlah orang/layanan yang dilayani. Jangan dibulatkan.",
        "denominator": "Jumlah sasaran. Wajib > 0 untuk cakupan.",
        "nilai": "Nilai tahunan sesuai satuan indikator: tahun, Poin, Dokumen, Laporan, atau Berita Acara.",
        "catatan": "Catatan baris, boleh kosong.",
    }
    for i, h in enumerate(headers, 1):
        if h in comments:
            ws.cell(row=1, column=i).comment = Comment(comments[h], "MATA SIDIK")

    numeric_cols = set(cfg["numeric_cols"])
    for r in range(2, 502):
        for i, h in enumerate(headers, 1):
            cell = ws.cell(row=r, column=i)
            body_cell(cell, align="right" if i in numeric_cols else "left",
                      color=WHITE if h == "catatan" else YELLOW)
            if h in ("periode", "tahun", "kode_faskes", "kode_desa", "scope_type", "kode_scope", "kode_indikator"):
                cell.number_format = "@"
            elif i in numeric_cols:
                cell.number_format = "0.###"
            cell.protection = Protection(locked=False)
        ws.row_dimensions[r].height = 18

    add_list_dv(ws, f'A2:A501', f"=Referensi!$A${ranges['periode'][0]}:$A${ranges['periode'][1]}",
                "Periode", "Pilih periode dari Referensi.")
    if cfg["is_monthly"]:
        add_list_dv(ws, "B2:B501", f"=Referensi!$A${ranges['faskes'][0]}:$A${ranges['faskes'][1]}",
                    "Faskes", "Pilih kode faskes pengunggah.")
        add_list_dv(ws, "C2:C501", f"=Referensi!$A${ranges['desa'][0]}:$A${ranges['desa'][1]}",
                    "Desa", "Pilih kode desa dari Referensi.")
        if cfg["has_indicators"] and ranges.get("indikator"):
            add_list_dv(ws, "D2:D501", f"=Referensi!$A${ranges['indikator'][0]}:$A${ranges['indikator'][1]}",
                        "Indikator", "Pilih indikator yang sesuai pola template.")
        add_decimal_dv(ws, "E2:E501", cfg["numeric_label"], "Isi angka >= 0.")
        if cfg["key"] == "cakupan":
            add_decimal_dv(ws, "F2:F501", "Denominator", "Wajib > 0 untuk cakupan.", greater_than_zero=True)
    else:
        add_list_dv(ws, "B2:B501", '"DESA,KECAMATAN,KABUPATEN"', "Scope Type", "Pilih level wilayah; faskes bukan scope.")
        add_list_dv(ws, "C2:C501", f"=Referensi!$A${ranges['scope'][0]}:$A${ranges['scope'][1]}",
                    "Kode Scope", "Pilih kode wilayah sesuai scope_type.")
        if cfg["has_indicators"] and ranges.get("indikator"):
            add_list_dv(ws, "D2:D501", f"=Referensi!$A${ranges['indikator'][0]}:$A${ranges['indikator'][1]}",
                        "Indikator", "Pilih indikator tahunan dari Referensi.")
        add_decimal_dv(ws, "E2:E501", "Nilai", "Isi angka >= 0 sesuai satuan indikator.")

    red_fill = PatternFill(start_color=LIGHT_RED, end_color=LIGHT_RED, fill_type="solid")
    orange_fill = PatternFill(start_color=LIGHT_ORANGE, end_color=LIGHT_ORANGE, fill_type="solid")
    if cfg["key"] == "cakupan":
        ws.conditional_formatting.add("F2:F501", CellIsRule(operator="lessThanOrEqual", formula=["0"], fill=red_fill))
        for col in ("E", "F"):
            ws.conditional_formatting.add(f"{col}2:{col}501", FormulaRule(formula=[f'AND($D2<>"",{col}2="")'], fill=orange_fill))
    else:
        ws.conditional_formatting.add("E2:E501", CellIsRule(operator="lessThan", formula=["0"], fill=red_fill))
        ws.conditional_formatting.add("E2:E501", FormulaRule(formula=['AND($D2<>"",E2="")'], fill=orange_fill))

    ws.freeze_panes = "A2"
    ws.auto_filter.ref = f"A1:{get_column_letter(len(headers))}501"
    ws.print_title_rows = "1:1"
    ws.page_setup.orientation = "landscape"
    ws.page_setup.fitToWidth = 1
    ws.sheet_properties.pageSetUpPr.fitToPage = True
    protect(ws)

def write_petunjuk(ws, cfg, master, indicator_rows):
    _, _, _, _, _, _, releases = master
    ws.sheet_view.showGridLines = False
    ws.sheet_properties.tabColor = BLUE
    ws.column_dimensions["A"].width = 3
    ws.column_dimensions["B"].width = 30
    ws.column_dimensions["C"].width = 105

    ws.merge_cells("B1:C1")
    t = ws["B1"]
    t.value = cfg["title"]
    t.font = TITLE
    t.fill = fill(NAVY)
    t.alignment = Alignment(vertical="center")
    ws.row_dimensions[1].height = 36

    ws.merge_cells("B2:C2")
    c = ws["B2"]
    c.value = cfg["subtitle"]
    c.font = SUBTITLE
    c.fill = fill(LIGHT_BLUE)
    c.alignment = Alignment(vertical="center", wrap_text=True)
    ws.row_dimensions[2].height = 40

    ws.merge_cells("B3:C3")
    c = ws["B3"]
    c.value = "STATUS: TEMPLATE V2 — review. Belum mengubah kode aplikasi; upload final menunggu penyesuaian kode di docs/05 §10.9–§10.11."
    c.font = Font(name=FONT, size=10, bold=True, color=RED)
    c.fill = fill(YELLOW)
    c.alignment = Alignment(vertical="center", wrap_text=True)
    ws.row_dimensions[3].height = 28

    r = 5
    section(ws, r, "1. Kapan memakai template ini")
    r += 1
    kv(ws, r, "Pola", cfg["tipe_input"]); r += 1
    kv(ws, r, "Frekuensi", cfg["frekuensi"]); r += 1
    kv(ws, r, "Satuan hasil", cfg["unit"]); r += 1
    kv(ws, r, "Sheet unggah", "Data"); r += 1
    kv(ws, r, "Jangan pakai untuk", "Pola lain: gunakan file BNBA-1/2/3 yang sesuai."); r += 1
    r += 1

    section(ws, r, "2. Cara isi")
    r += 1
    steps = [
        ("Langkah 1", "Buka sheet Data. Jangan mengubah header atau urutan kolom."),
        ("Langkah 2", "Isi satu baris = satu " + ("desa + bulan + indikator" if cfg["is_monthly"] else "tahun + scope + indikator") + "."),
        ("Langkah 3", "Pilih kode indikator dari dropdown. Nama dan satuan dapat dilihat di sheet Referensi; tidak perlu diisi di Data." if cfg["has_indicators"] else "Master pola ini belum ada. Tambahkan indikator lewat Master; setelah itu kode indikator dipilih dari dropdown."),
        ("Langkah 4", "Isi angka pada kolom input: " + cfg["numeric_label"] + "."),
        ("Langkah 5", "Jangan menambahkan total/subtotal, merged cell, atau formula."),
        ("Langkah 6", "Simpan .xlsx, unggah, lalu periksa preview dan baris bermasalah."),
    ]
    for label, text in steps:
        kv(ws, r, label, text, height=32); r += 1
    r += 1

    section(ws, r, "3. Aturan wajib")
    r += 1
    rules = [
        ("Satu file", "Satu faskes + satu bulan." if cfg["is_monthly"] else "Satu tahun; tidak ada bulan."),
        ("Satu baris", "Satu indikator + satu scope."),
        ("Kode resmi", "Gunakan kode dari sheet Referensi; jangan mengetik kode baru."),
        ("Blank ≠ 0", "Blank berarti belum ada data; 0 berarti nol yang dilaporkan."),
        ("Tanpa PII", "Jangan masukkan nama, NIK, atau identitas individu."),
        ("Tanpa formula", "Jangan mengubah kolom otomatis atau menambah formula."),
    ]
    if cfg["key"] == "kejadian":
        rules.append(("Count", "Isi jumlah kasus; jangan isi numerator/denominator atau persentase."))
        rules.append(("Arah", "Kematian: makin kecil makin baik; jangan memakai rumus capaian persen."))
    elif cfg["key"] == "cakupan":
        rules.append(("Rasio", "Isi numerator dan denominator; jangan isi persentase jadi."))
        rules.append(("Denominator", "Wajib > 0; jika 0/kosong, baris tidak valid."))
    else:
        rules.append(("Tahunan", "Jangan membuat bulan/triwulan; target tahunan jangan dibagi 12."))
        rules.append(("Snapshot", "UHH snapshot: jangan dijumlahkan lintas tahun."))
    for label, text in rules:
        kv(ws, r, label, text, height=32); r += 1
    r += 1

    section(ws, r, "4. Satuan")
    r += 1
    unit_text = {
        "kejadian": "Semua baris pada template ini memakai satuan Kasus. Jangan menulis 'orang' atau '%'.",
        "cakupan": "Hasil cakupan selalu %. Numerator adalah jumlah orang/layanan; denominator adalah sasaran. "
                   "Master IND-001 s.d. IND-003 masih bertuliskan 'orang' untuk numerator; itu perlu dinormalisasi.",
        "tahunan": "Satuan mengikuti indikator: tahun (UHH), Poin (AKIP), Dokumen, Laporan, Berita Acara. "
                   "Jangan memakai % untuk semua indikator.",
    }[cfg["key"]]
    kv(ws, r, "Prinsip", unit_text, height=50); r += 1
    r += 1

    section(ws, r, "5. Kesiapan kode pilot")
    r += 1
    readiness = [
        ("Parser", "Template minimal; parser pilot masih membutuhkan kode_iku/scope. Upload menunggu penyesuaian kode."),
        ("Master", "Indikator pola ini belum ada di master pilot; tambahkan dulu." if not cfg["has_indicators"] else "Indikator cakupan tersedia di master pilot; satuan hasil %."),
        ("Kalkulasi", "Mesin pilot masih hanya KABUPATEN/FASKES dan menghitung semua input langsung sebagai rasio."),
        ("Periode", "Periode tahunan belum ada di tabel periode pilot." if cfg["is_annual"] else "Periode bulanan tersedia; per triwulan/tahun diturunkan."),
        ("Dashboard", "Dashboard per desa dan format non-% belum aktif."),
    ]
    for label, text in readiness:
        kv(ws, r, label, text, height=40); r += 1
    r += 1

    section(ws, r, "6. File lain")
    r += 1
    kv(ws, r, "BNBA-1", "Kejadian/count bulanan: kematian ibu/balita, jumlah kasus.", height=28); r += 1
    kv(ws, r, "BNBA-2", "Cakupan/rasio bulanan: SPM/IKU cakupan, numerator/denominator.", height=28); r += 1
    kv(ws, r, "BNBA-3", "Tahunan: UHH, AKIP, dokumen/laporan tahunan, satu nilai per tahun.", height=28); r += 1
    kv(ws, r, "Referensi", "Kode master, satuan, tipe input, arah, agregasi.", height=24); r += 1
    r += 1

    section(ws, r, "7. Catatan")
    r += 1
    kv(ws, r, "Contoh", "Angka dan kode CONTOH-* di sheet Contoh hanya ilustrasi; bukan data produksi.", height=32); r += 1
    kv(ws, r, "Tanpa kode", "Template ini tidak mengubah kode aplikasi. Penyesuaian kode ada di docs/05.", height=30); r += 1
    kv(ws, r, "Proteksi", "Sheet dilindungi untuk mengurangi salah edit; hanya sel input yang bisa diubah.", height=30); r += 1
    kv(ws, r, "Sumber", "Master pilot per 2026-11-19; file Excel asli Proses Bisnis tidak diubah.", height=30); r += 1


def write_contoh(ws, cfg):
    ws.sheet_view.showGridLines = False
    ws.sheet_properties.tabColor = GREEN
    headers = cfg["headers"]
    ws.merge_cells(start_row=1, start_column=1, end_row=1, end_column=len(headers))
    t = ws.cell(row=1, column=1, value=f"CONTOH — {cfg['title']} — ANGKA FIKTIF")
    t.font = TITLE
    t.fill = fill(GREEN)
    t.alignment = Alignment(vertical="center")
    ws.row_dimensions[1].height = 34

    ws.merge_cells(start_row=2, start_column=1, end_row=2, end_column=len(headers))
    c = ws.cell(row=2, column=1, value=(
        "Contoh ini hanya menunjukkan bentuk baris yang benar. Jangan mengunggah sheet Contoh. "
        "Kode CONTOH-* bukan master; ganti dengan kode dari Referensi setelah indikator ditambahkan."
    ))
    c.font = SUBTITLE
    c.fill = fill(LIGHT_GREEN)
    c.alignment = Alignment(vertical="center", wrap_text=True)
    ws.row_dimensions[2].height = 34

    header_row(ws, 4, headers, GREEN, 32)
    widths = {
        "periode": 12, "tahun": 10, "kode_faskes": 18, "kode_desa": 14, "nama_desa": 22,
        "scope_type": 14, "kode_scope": 14, "nama_scope": 24,
        "kode_indikator": 20, "nama_indikator": 32, "satuan": 10,
        "jumlah": 10, "numerator": 12, "denominator": 14, "nilai": 10,
        "catatan": 30, "kode_iku": 22, "scope_code": 14,
    }
    for i, h in enumerate(headers, 1):
        ws.column_dimensions[get_column_letter(i)].width = widths.get(h, 14)
    for rr, rowvals in enumerate(cfg["examples"], 5):
        for cc, val in enumerate(rowvals, 1):
            cell = ws.cell(row=rr, column=cc, value=val)
            body_cell(cell, wrap=True, align="center" if cc in cfg["numeric_cols"] else "left", color=LIGHT_GREEN)
            if cc in cfg["numeric_cols"]:
                cell.number_format = "0.###"
        ws.row_dimensions[rr].height = 24

    note(
        ws, 5 + len(cfg["examples"]) + 1,
        "Untuk upload, salin pola ini ke sheet Data dan ganti kode CONTOH-* dengan kode dari Referensi. "
        "Jangan mengunggah sheet Contoh; parser tidak membacanya.",
        color=YELLOW, span=len(headers), height=34,
    )


def main():
    master = load_master()
    OUT.mkdir(parents=True, exist_ok=True)
    for cfg in PATTERNS:
        build_workbook(cfg, master)


if __name__ == "__main__":
    main()
