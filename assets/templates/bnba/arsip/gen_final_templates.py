#!/usr/bin/env python3
"""Generate final BNBA templates for all units and annual data.

Requires: python3 -m pip install openpyxl
Run from any directory: python3 assets/templates/bnba/final/gen_final_templates.py

Outputs:
- assets/templates/bnba/final/BNBA-Bulanan-Semua-Satuan.xlsx
- assets/templates/bnba/final/BNBA-Tahunan.xlsx
Copies both to frontend/static/templates/.
"""

from __future__ import annotations

import shutil
import sqlite3
from pathlib import Path

from openpyxl import Workbook
from openpyxl.comments import Comment
from openpyxl.formatting.rule import CellIsRule, FormulaRule
from openpyxl.styles import Alignment, Border, Font, PatternFill, Protection, Side
from openpyxl.utils import get_column_letter
from openpyxl.worksheet.datavalidation import DataValidation

ROOT = Path(__file__).resolve().parents[4]
DB = ROOT / "backend/data/matasidik.db"
OUT = ROOT / "assets/templates/bnba/final"
STATIC = ROOT / "frontend/static/templates"

NAVY = "1F3864"
BLUE = "2F5597"
LIGHT_BLUE = "D9E2F3"
PALE_BLUE = "EEF3FA"
GREEN = "548235"
LIGHT_GREEN = "E2EFDA"
YELLOW = "FFF2CC"
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


def header_row(ws, row: int, headers: list[str], color: str = NAVY, height: int = 34):
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


def kv(ws, row: int, label: str, value: str, height: int = 22):
    a = ws.cell(row=row, column=2, value=label)
    b = ws.cell(row=row, column=3, value=value)
    a.font = BODY_BOLD
    a.fill = fill(PALE_BLUE)
    a.border = BORDER
    a.alignment = Alignment(vertical="top", wrap_text=True)
    b.font = BODY
    b.border = BORDER
    b.alignment = Alignment(vertical="top", wrap_text=True)
    ws.row_dimensions[row].height = height


def add_list_dv(ws, rng: str, formula: str, title: str, prompt: str):
    dv = DataValidation(type="list", formula1=formula, allow_blank=True, showDropDown=False)
    dv.errorTitle = "Nilai tidak valid"
    dv.error = "Pilih nilai dari sheet Referensi atau ikuti format yang diminta."
    dv.promptTitle = title
    dv.prompt = prompt
    ws.add_data_validation(dv)
    dv.add(rng)


def add_number_dv(ws, rng: str, title: str, prompt: str, *, positive=False):
    dv = DataValidation(
        type="decimal", operator="greaterThan" if positive else "greaterThanOrEqual",
        formula1="0", allow_blank=True,
    )
    dv.errorTitle = "Angka tidak valid"
    dv.error = "Isi angka > 0." if positive else "Isi angka >= 0."
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
           WHERE l.aktif=1 ORDER BY l.level, l.kode"""
    )]
    indicators = [dict(r) for r in q(
        """SELECT i.kode, i.nama, i.satuan,
                  l.kode AS iku_kode, l.nama AS iku_nama, c.kode AS kategori_kode,
                  c.nama AS kategori_nama, l.aggregation_method,
                  pb.nama AS pj_bidang_nama, pt.nama AS pj_timker_nama
           FROM indikator i
           JOIN indikator_label il ON il.indikator_id=i.id AND il.is_primary=1
           JOIN label l ON l.id=il.label_id
           LEFT JOIN dashboard_category c ON c.id=l.category_id
           LEFT JOIN bidang pb ON pb.id=i.pj_bidang_id
           LEFT JOIN bidang pt ON pt.id=i.pj_timker_id
           WHERE i.aktif=1 AND l.aktif=1 ORDER BY l.kode"""
    )]
    indicator_labels = [dict(r) for r in q(
        """SELECT i.kode AS indikator_kode, i.nama AS indikator_nama, i.satuan,
                  l.kode AS iku_kode, l.nama AS iku_nama,
                  c.kode AS kategori_kode, il.is_primary,
                  l.allow_direct_input, l.is_rollup_only
           FROM indikator_label il
           JOIN indikator i ON i.id=il.indikator_id
           JOIN label l ON l.id=il.label_id
           LEFT JOIN dashboard_category c ON c.id=l.category_id
           WHERE i.aktif=1 AND l.aktif=1
           ORDER BY i.kode, il.is_primary DESC, l.kode"""
    )]
    wilayah = [dict(r) for r in q(
        """SELECT w.kode, w.nama, w.level, p.kode AS parent_kode
           FROM wilayah w LEFT JOIN wilayah p ON p.id=w.parent_id ORDER BY w.kode"""
    )]
    faskes = [dict(r) for r in q(
        "SELECT id, kode, nama FROM faskes WHERE aktif=1 ORDER BY kode"
    )]
    periods = [dict(r) for r in q(
        "SELECT kode, tahun, bulan FROM periode ORDER BY kode"
    )]
    releases = [dict(r) for r in q(
        "SELECT tahun FROM master_release ORDER BY tahun DESC"
    )]
    con.close()
    return categories, labels, indicators, indicator_labels, wilayah, faskes, periods, releases


def ref_table(ws, start_row: int, title: str, headers: list[str], rows: list[list], color: str = BLUE):
    ws.merge_cells(start_row=start_row, start_column=1, end_row=start_row, end_column=len(headers))
    c = ws.cell(row=start_row, column=1, value=title)
    c.font = Font(name=FONT, size=11, bold=True, color=WHITE)
    c.fill = fill(color)
    c.alignment = Alignment(vertical="center")
    ws.row_dimensions[start_row].height = 22
    header_row(ws, start_row + 1, headers, color, 28)
    for rr, rowvals in enumerate(rows, start_row + 2):
        for cc, val in enumerate(rowvals, 1):
            cell = ws.cell(row=rr, column=cc, value=val)
            body_cell(cell, wrap=True)
        ws.row_dimensions[rr].height = 22
    return start_row + 2, start_row + 1 + len(rows)


def note(ws, row: int, text: str, *, color=YELLOW, height=40, span=12, start_col=1):
    ws.merge_cells(start_row=row, start_column=start_col, end_row=row, end_column=start_col + span - 1)
    c = ws.cell(row=row, column=start_col, value=text)
    c.font = BODY
    c.fill = fill(color)
    c.border = BORDER
    c.alignment = Alignment(vertical="top", wrap_text=True)
    ws.row_dimensions[row].height = height


# Satuan mapping from sheet 2. The app must use tipe_nilai + Agregasi, not
# a single percentage formula for everything.
UNIT_ROWS = [
    ["(tahun)", "SCALAR", "nilai", "LAST_VALUE", "UHH; input 1x/tahun, kabupaten"],
    ["Kasus", "COUNT", "jumlah", "SUM", "Kematian ibu/balita; kejadian"],
    ["%", "RATIO", "numerator + denominator", "RATIO_OF_SUMS", "Prevalensi/cakupan"],
    ["Persentase", "RATIO", "numerator + denominator", "RATIO_OF_SUMS", "Cakupan"],
    ["(%)", "RATIO", "numerator + denominator", "RATIO_OF_SUMS", "Cakupan"],
    ["Poin", "SCALAR", "nilai", "LAST_VALUE", "AKIP/indeks"],
    ["Rasio", "RATIO_NON_PERSEN", "numerator + denominator", "RATIO_OF_SUMS", "Rasio tenaga kesehatan"],
    ["Dokumen", "COUNT", "jumlah", "SUM", "Jumlah dokumen"],
    ["Laporan", "COUNT", "jumlah", "SUM", "Jumlah laporan"],
    ["Berita Acara", "COUNT", "jumlah", "SUM", "Jumlah berita acara"],
    ["orang", "COUNT/RATIO", "jumlah atau numerator/denominator", "SUM atau RATIO_OF_SUMS", "Tergantung indikator"],
]

POLA_ROWS = [
    ["BULANAN_RASIO", "Rasio/persentase bulanan", "Sheet Rasio", "numerator + denominator"],
    ["BULANAN_JUMLAH", "Jumlah/kejadian bulanan", "Sheet Jumlah", "jumlah"],
    ["TAHUNAN_NILAI", "Tahunan/scalar", "Sheet Data (template tahunan)", "nilai"],
]


def indicator_rows(indicators):
    rows = []
    for i in indicators:
        satuan = i["satuan"]
        if satuan.strip() in ("%", "Persentase", "(%)"):
            tipe = "RATIO"
            pola = "BULANAN_RASIO"
        elif satuan.strip().lower() == "orang":
            tipe = "RATIO"
            pola = "BULANAN_RASIO"
        else:
            tipe = "RATIO"
            pola = "BULANAN_RASIO"
        rows.append([
            i["kode"], i["nama"], i["iku_kode"], i["iku_nama"], i["kategori_kode"] or "",
            satuan, "%" if tipe == "RATIO" else satuan, tipe, pola, i["aggregation_method"],
        ])
    return rows


def mapping_rows(indicator_labels):
    out = []
    for x in indicator_labels:
        out.append([
            x["indikator_kode"], x["indikator_nama"], x["iku_kode"], x["iku_nama"],
            x["kategori_kode"] or "", "Ya" if x["is_primary"] else "Tidak",
            "Ya" if (x["allow_direct_input"] and not x["is_rollup_only"]) else "Tidak",
            x["satuan"],
        ])
    return out


def write_referensi(ws, *, annual: bool, master):
    categories, labels, indicators, indicator_labels, wilayah, faskes, periods, _ = master
    ws.sheet_view.showGridLines = False
    ws.sheet_properties.tabColor = GREEN
    widths = {"A": 22, "B": 34, "C": 20, "D": 18, "E": 20, "F": 18, "G": 16, "H": 24, "I": 16, "J": 16}
    for col, w in widths.items():
        ws.column_dimensions[col].width = w

    ws.merge_cells("A1:J1")
    c = ws["A1"]
    c.value = "REFERENSI MASTER & SATUAN — JANGAN DIUBAH"
    c.font = TITLE
    c.fill = fill(GREEN)
    c.alignment = Alignment(vertical="center")
    ws.row_dimensions[1].height = 32

    ws.merge_cells("A2:J2")
    c = ws["A2"]
    c.value = "Sumber: master pilot `backend/data/matasidik.db` per 2026-11-19; pola satuan dari sheet 2 `Proses Bisnis Mata Sidik.xlsx`. Daftar dinamis; regenerate setelah master berubah."
    c.font = SUBTITLE
    c.fill = fill(LIGHT_GREEN)
    c.alignment = Alignment(vertical="center", wrap_text=True)
    ws.row_dimensions[2].height = 34

    ranges = {}
    row = 4

    start, end = ref_table(ws, row, "1. Satuan dan pola nilai", ["satuan", "tipe_nilai", "field_wajib", "agregasi", "contoh"], UNIT_ROWS)
    ranges["satuan"] = (start, end)

    row = end + 3
    start, end = ref_table(ws, row, "2. Pola input", ["kode", "nama", "sheet", "field"], POLA_ROWS)
    ranges["pola"] = (start, end)

    row = end + 3
    start, end = ref_table(ws, row, "3. Kategori dashboard", ["kode", "nama", "jenis"],
                           [[c["kode"], c["nama"], c["jenis"]] for c in categories])
    ranges["kategori"] = (start, end)

    row = end + 3
    if annual:
        years = sorted({str(p["tahun"]) for p in periods})
        periode_rows = [[y, f"Tahun {y}"] for y in years]
    else:
        periode_rows = [[p["kode"], f'Bulan {p["bulan"]:02d} tahun {p["tahun"]}'] for p in periods]
    start, end = ref_table(ws, row, "4. Periode", ["kode", "keterangan"], periode_rows)
    ranges["periode"] = (start, end)

    row = end + 3
    start, end = ref_table(ws, row, "5. Faskes pengunggah", ["kode", "nama", "id_internal"], [[f["kode"], f["nama"], f["id"]] for f in faskes])
    ranges["faskes"] = (start, end)

    row = end + 3
    scope_rows = []
    for w in wilayah:
        scope_rows.append([w["level"].upper() if w["level"] != "kabupaten" else "KABUPATEN", w["kode"], w["nama"], w["level"]])
    start, end = ref_table(ws, row, "6. Scope wilayah (DESA/KECAMATAN/KABUPATEN)", ["scope_type", "kode", "nama", "level"], scope_rows)
    ranges["scope"] = (start, end)

    row = end + 3
    ind_rows = indicator_rows(indicators)
    start, end = ref_table(
        ws, row, "7. Indikator master pilot (selalu ratio saat ini)",
        ["kode", "nama", "kode_iku", "nama_iku", "kategori", "satuan_master", "satuan_hasil", "tipe_nilai", "pola", "agregasi"],
        ind_rows,
    )
    ranges["indikator"] = (start, end)

    row = end + 3
    start, end = ref_table(
        ws, row, "8. Mapping semua label IKU ↔ indikator (many-to-many)",
        ["kode_indikator", "nama_indikator", "kode_iku", "nama_iku", "kategori", "is_primary", "is_input", "satuan"],
        mapping_rows(indicator_labels),
    )
    ranges["mapping"] = (start, end)

    row = end + 3
    start, end = ref_table(
        ws, row, "9. PJ Bidang/Timker (bukan label)",
        ["kode_indikator", "nama_indikator", "pj_bidang", "pj_timker"],
        [[i["kode"], i["nama"], i["pj_bidang_nama"] or "—", i["pj_timker_nama"] or "—"] for i in indicators],
    )
    ranges["pj"] = (start, end)

    row = end + 3
    if annual:
        note(ws, row, (
            "Belum ada indikator tahunan di master pilot. Tambahkan UHH/AKIP/dokumen/laporan tahunan lewat Master "
            "sebelum memakai sheet Data. Contoh di sheet Contoh bukan master. Satuan tahun: (tahun), Poin, Dokumen, Laporan, Berita Acara, Kasus."
        ), height=54)
    else:
        note(ws, row, (
            "Sheet Rasio memakai numerator/denominator. Sheet Jumlah memakai jumlah. "
            "Satuan Kasus/Dokumen/Laporan/Berita Acara bukan persentase. "
            "IND-001 s.d. IND-003 memakai 'orang' sebagai numerator, bukan satuan hasil; hasilnya tetap %."
        ), height=54)
    return ranges


def write_data_sheet(ws, *, mode: str, ranges, has_indicator: bool, headers: list[str]):
    header_row(ws, 1, headers, NAVY, 36)
    ws.sheet_view.showGridLines = False
    ws.sheet_properties.tabColor = NAVY
    widths = {
        "periode": 12, "tahun": 10, "kode_faskes": 18, "kode_kategori": 14,
        "kode_iku": 22, "kode_indikator": 20, "scope_type": 14, "scope_code": 16,
        "numerator": 12, "denominator": 14, "jumlah": 12, "nilai": 12, "catatan": 34,
    }
    cl = {h: get_column_letter(i) for i, h in enumerate(headers, 1)}
    for h in headers:
        ws.column_dimensions[cl[h]].width = widths.get(h, 14)

    comments = {
        "periode": "Wajib YYYY-MM; pilih periode bulanan.",
        "tahun": "Wajib YYYY; data 1x per tahun. Jangan isi bulan.",
        "kode_faskes": "Kode faskes pengunggah/pemilik data.",
        "kode_kategori": "Opsional; harus cocok dengan kategori label utama.",
        "kode_iku": "Satu/beberapa kode IKU leaf, pisahkan ;. Label pertama = utama/konteks.",
        "kode_indikator": "Kode indikator dari master; menentukan satuan dan pola.",
        "scope_type": "DESA/KECAMATAN/KABUPATEN. Faskes bukan scope.",
        "scope_code": "Kode wilayah resmi sesuai scope_type.",
        "numerator": "Untuk rasio: jumlah pembilang. Jangan diisi untuk count.",
        "denominator": "Untuk rasio: wajib > 0. Jangan diisi untuk count.",
        "jumlah": "Untuk count: angka jumlah. Jangan diisi untuk rasio.",
        "nilai": "Untuk tahunan: nilai sesuai satuan (tahun/Poin/Dokumen/Laporan/Berita Acara).",
        "catatan": "Catatan baris, boleh kosong.",
    }
    for i, h in enumerate(headers, 1):
        if h in comments:
            ws.cell(row=1, column=i).comment = Comment(comments[h], "MATA SIDIK")

    numeric = {"numerator", "denominator", "jumlah", "nilai"}
    for r in range(2, 502):
        for i, h in enumerate(headers, 1):
            cell = ws.cell(row=r, column=i)
            body_cell(cell, align="right" if h in numeric else "left", color=WHITE if h == "catatan" else YELLOW)
            if h in ("periode", "tahun", "kode_faskes", "kode_kategori", "kode_iku", "kode_indikator", "scope_type", "scope_code"):
                cell.number_format = "@"
            elif h in numeric:
                cell.number_format = "0.###"
            cell.protection = Protection(locked=False)
        ws.row_dimensions[r].height = 18

    add_list_dv(ws, "A2:A501", f"=Referensi!$A${ranges['periode'][0]}:$A${ranges['periode'][1]}", "Periode", "Pilih periode dari Referensi.")
    add_list_dv(ws, "B2:B501", f"=Referensi!$A${ranges['faskes'][0]}:$A${ranges['faskes'][1]}", "Faskes", "Pilih kode faskes.")
    add_list_dv(ws, "C2:C501", f"=Referensi!$A${ranges['kategori'][0]}:$A${ranges['kategori'][1]}", "Kategori", "Opsional.")
    # kode_iku (D) sengaja tanpa dropdown karena boleh multi ;.
    if has_indicator:
        add_list_dv(ws, "E2:E501", f"=Referensi!$A${ranges['indikator'][0]}:$A${ranges['indikator'][1]}", "Indikator", "Pilih indikator.")
    add_list_dv(ws, "F2:F501", '"DESA,KECAMATAN,KABUPATEN"', "Scope Type", "Pilih level wilayah.")
    add_list_dv(ws, "G2:G501", f"=Referensi!$A${ranges['scope'][0]}:$A${ranges['scope'][1]}", "Scope", "Pilih kode wilayah.")
    if mode == "rasio":
        add_number_dv(ws, "H2:H501", "Numerator", "Isi angka >= 0.")
        add_number_dv(ws, "I2:I501", "Denominator", "Wajib > 0 untuk rasio.", positive=True)
    elif mode == "jumlah":
        add_number_dv(ws, "H2:H501", "Jumlah", "Isi angka >= 0; satuan mengikuti indikator.")
    else:
        add_number_dv(ws, "H2:H501", "Nilai", "Isi angka >= 0; satuan mengikuti indikator.")

    red = PatternFill(start_color=LIGHT_RED, end_color=LIGHT_RED, fill_type="solid")
    orange = PatternFill(start_color=LIGHT_ORANGE, end_color=LIGHT_ORANGE, fill_type="solid")
    if mode == "rasio":
        ws.conditional_formatting.add("I2:I501", CellIsRule(operator="lessThanOrEqual", formula=["0"], fill=red))
        ws.conditional_formatting.add("H2:H501", FormulaRule(formula=['AND($E2<>"",H2="")'], fill=orange))
        ws.conditional_formatting.add("I2:I501", FormulaRule(formula=['AND($E2<>"",I2="")'], fill=orange))
    else:
        ws.conditional_formatting.add("H2:H501", CellIsRule(operator="lessThan", formula=["0"], fill=red))
        ws.conditional_formatting.add("H2:H501", FormulaRule(formula=['AND($E2<>"",H2="")'], fill=orange))

    ws.freeze_panes = "A2"
    ws.auto_filter.ref = f"A1:{get_column_letter(len(headers))}501"
    ws.print_title_rows = "1:1"
    ws.page_setup.orientation = "landscape"
    ws.page_setup.fitToWidth = 1
    ws.sheet_properties.pageSetUpPr.fitToPage = True
    protect(ws)


def write_petunjuk(ws, *, annual: bool):
    ws.sheet_view.showGridLines = False
    ws.sheet_properties.tabColor = BLUE
    ws.column_dimensions["A"].width = 3
    ws.column_dimensions["B"].width = 30
    ws.column_dimensions["C"].width = 110

    ws.merge_cells("B1:C1")
    c = ws["B1"]
    c.value = "DATA TAHUNAN — 1× PER TAHUN" if annual else "BNBA BULANAN — SEMUA SATUAN"
    c.font = TITLE
    c.fill = fill(NAVY)
    c.alignment = Alignment(vertical="center")
    ws.row_dimensions[1].height = 36

    ws.merge_cells("B2:C2")
    c = ws["B2"]
    c.value = (
        "Khusus data tahunan: UHH, AKIP, dokumen/laporan tahunan. Satu baris = satu tahun + scope + indikator."
        if annual else
        "Dua pola input: Rasio (numerator/denominator) dan Jumlah/kejadian (jumlah). Satuan mengikuti master."
    )
    c.font = SUBTITLE
    c.fill = fill(LIGHT_BLUE)
    c.alignment = Alignment(vertical="center", wrap_text=True)
    ws.row_dimensions[2].height = 36

    ws.merge_cells("B3:C3")
    c = ws["B3"]
    c.value = "STATUS: FORMAT FINAL — dukungan parser/calculation untuk COUNT/SCALAR masih menunggu fase kode berikutnya. Template rasio yang diterima saat ini: template-bnba-dinamis.xlsx."
    c.font = Font(name=FONT, size=10, bold=True, color=RED)
    c.fill = fill(YELLOW)
    c.alignment = Alignment(vertical="center", wrap_text=True)
    ws.row_dimensions[3].height = 34

    r = 5
    if annual:
        section(ws, r, "1. Kapan memakai template ini"); r += 1
        for label, text in [
            ("Frekuensi", "Hanya untuk indikator yang diinput 1× per tahun; jangan dipakai untuk data bulanan."),
            ("Contoh", "UHH (tahun), AKIP (Poin), dokumen/laporan/berita acara tahunan."),
            ("Tanpa bulan", "Tidak ada bulan/triwulan. Jangan membagi target tahunan dengan 12."),
            ("Snapshot", "UHH snapshot; jangan dijumlahkan lintas tahun."),
        ]:
            kv(ws, r, label, text, 34); r += 1
        r += 1
        section(ws, r, "2. Cara isi sheet Data"); r += 1
        for label, text in [
            ("Langkah 1", "Buka sheet Data. Jangan mengubah header."),
            ("Langkah 2", "Isi tahun, kode_faskes, kode_iku, kode_indikator, scope_type, scope_code, nilai."),
            ("Langkah 3", "Gunakan kode dari Referensi; kode_iku boleh beberapa dengan ;."),
            ("Langkah 4", "Isi satu nilai per indikator/scope/tahun. Jangan membuat baris bulanan."),
            ("Langkah 5", "Jangan menambahkan total, formula, atau merged cell."),
        ]:
            kv(ws, r, label, text, 32); r += 1
        r += 1
        section(ws, r, "3. Satuan tahunan"); r += 1
        for label, text in [
            ("(tahun)", "UHH: isi angka tahun, mis. 73,5."),
            ("Poin", "AKIP/indeks: isi poin, mis. 85."),
            ("Dokumen/Laporan/Berita Acara", "Isi jumlah, bukan persentase."),
            ("Kasus", "Bila indikator tahunan berupa jumlah kasus, isi jumlah."),
        ]:
            kv(ws, r, label, text, 30); r += 1
        r += 1
    else:
        section(ws, r, "1. Dua pola input"); r += 1
        for label, text in [
            ("Sheet Rasio", "Untuk %, Persentase, (%), dan Rasio non-persen. Isi numerator + denominator."),
            ("Sheet Jumlah", "Untuk Kasus, Dokumen, Laporan, Berita Acara, dan orang (count). Isi jumlah."),
            ("Pilih satu", "Jangan mengisi baris yang sama di kedua sheet."),
            ("satuan", "Satuan mengikuti master indikator; lihat Referensi."),
        ]:
            kv(ws, r, label, text, 36); r += 1
        r += 1
        section(ws, r, "2. Cara isi sheet Rasio"); r += 1
        for label, text in [
            ("Langkah 1", "Buka sheet Rasio."),
            ("Langkah 2", "Isi periode, faskes, kode_iku, kode_indikator, scope_type, scope_code."),
            ("Langkah 3", "Isi numerator dan denominator; denominator wajib > 0."),
            ("Langkah 4", "Jangan mengisi jumlah; jangan mengisi persentase jadi."),
        ]:
            kv(ws, r, label, text, 30); r += 1
        r += 1
        section(ws, r, "3. Cara isi sheet Jumlah"); r += 1
        for label, text in [
            ("Langkah 1", "Buka sheet Jumlah."),
            ("Langkah 2", "Isi periode, faskes, kode_iku, kode_indikator, scope_type, scope_code."),
            ("Langkah 3", "Isi jumlah sesuai satuan indikator (Kasus/Dokumen/Laporan/Berita Acara)."),
            ("Langkah 4", "Jangan mengisi numerator/denominator; jangan mengubah jumlah menjadi persentase."),
        ]:
            kv(ws, r, label, text, 32); r += 1
        r += 1

    section(ws, r, "4. Aturan wajib"); r += 1
    for label, text in [
        ("Scope", "Hanya DESA/KECAMATAN/KABUPATEN. Faskes adalah pengunggah (kode_faskes), bukan scope."),
        ("Satu grain", "Simpan level terendah; jangan menulis desa + kecamatan + kabupaten sekaligus."),
        ("Blank ≠ 0", "Blank = belum ada data; 0 = nol yang dilaporkan."),
        ("Tanpa PII", "Jangan masukkan nama, NIK, atau identitas individu."),
        ("Tanpa formula", "Jangan mengubah header atau menambah formula."),
    ]:
        kv(ws, r, label, text, 32); r += 1
    r += 1

    section(ws, r, "5. Kesiapan kode"); r += 1
    for label, text in [
        ("Saat ini", "Parser/classifier/calculation pilot masih rasio. Count/scalar belum dihitung."),
        ("Rasio", "Gunakan template-bnba-dinamis.xlsx untuk upload yang diterima sekarang."),
        ("Count/scalar", "Template ini menyiapkan format; perlu tipe_nilai + SUM/LAST_VALUE di kode."),
    ]:
        kv(ws, r, label, text, 36); r += 1
    r += 1
    section(ws, r, "6. Catatan"); r += 1
    for label, text in [
        ("Contoh", "Angka dan kode CONTOH-* hanya ilustrasi; bukan data produksi."),
        ("Master", "Indikator di luar master pilot harus ditambahkan lewat Master; jangan mengetik kode baru."),
        ("Regenerate", "Setelah master berubah, regenerate template agar Referensi ikut berubah."),
    ]:
        kv(ws, r, label, text, 32); r += 1


def write_contoh(ws, *, annual: bool, ranges):
    ws.sheet_view.showGridLines = False
    ws.sheet_properties.tabColor = GREEN
    if annual:
        headers = ["tahun", "kode_faskes", "kode_kategori", "kode_iku", "kode_indikator", "scope_type", "scope_code", "nilai", "catatan"]
        rows = [
            ["2026", "PKM-TALIWANG", "IKU", "CONTOH-IKU-UHH", "CONTOH-UHH", "KABUPATEN", "5201", 73.5, "ILUSTRASI — kode contoh, bukan master"],
            ["2026", "PKM-TALIWANG", "IKU", "CONTOH-IKU-AKIP", "CONTOH-AKIP", "KABUPATEN", "5201", 85, "ILUSTRASI — kode contoh, bukan master"],
            ["2026", "PKM-TALIWANG", "IKU", "CONTOH-IKU-DOK", "CONTOH-DOKUMEN", "KABUPATEN", "5201", 12, "ILUSTRASI — kode contoh, bukan master"],
        ]
        ws.merge_cells("A1:I1")
        t = ws["A1"]; t.value = "CONTOH DATA TAHUNAN — ANGKA FIKTIF"; t.font = TITLE; t.fill = fill(GREEN); t.alignment = Alignment(vertical="center")
        ws.merge_cells("A2:I2")
        c = ws["A2"]; c.value = "Kode CONTOH-* bukan master dan tidak boleh diunggah. Ganti dengan kode dari Referensi setelah indikator ditambahkan."; c.font = SUBTITLE; c.fill = fill(LIGHT_GREEN); c.alignment = Alignment(vertical="center", wrap_text=True); ws.row_dimensions[2].height = 30
        header_row(ws, 4, headers, GREEN, 34)
        for i, h in enumerate(headers, 1):
            ws.column_dimensions[get_column_letter(i)].width = {"tahun":10,"kode_faskes":18,"kode_kategori":14,"kode_iku":20,"kode_indikator":20,"scope_type":14,"scope_code":16,"nilai":12,"catatan":34}.get(h, 14)
        for rr, vals in enumerate(rows, 5):
            for cc, val in enumerate(vals, 1):
                cell = ws.cell(row=rr, column=cc, value=val)
                body_cell(cell, wrap=True, align="right" if cc == 8 else "left", color=LIGHT_GREEN)
            ws.row_dimensions[rr].height = 24
        note(ws, 9, "Jangan mengunggah sheet Contoh. Data tahunan tidak memakai bulan/triwulan; target tahunan tidak dibagi 12.", span=9, height=32)
        return

    # Monthly: two example tables.
    rasio_headers = ["periode", "kode_faskes", "kode_kategori", "kode_iku", "kode_indikator", "scope_type", "scope_code", "numerator", "denominator", "catatan"]
    jumlah_headers = ["periode", "kode_faskes", "kode_kategori", "kode_iku", "kode_indikator", "scope_type", "scope_code", "jumlah", "catatan"]
    ws.merge_cells("A1:J1")
    t = ws["A1"]; t.value = "CONTOH BULANAN — ANGKA FIKTIF"; t.font = TITLE; t.fill = fill(GREEN); t.alignment = Alignment(vertical="center")
    ws.merge_cells("A2:J2")
    c = ws["A2"]; c.value = "Baris rasio memakai numerator/denominator; baris jumlah memakai jumlah. Kode CONTOH-* bukan master."; c.font = SUBTITLE; c.fill = fill(LIGHT_GREEN); c.alignment = Alignment(vertical="center", wrap_text=True); ws.row_dimensions[2].height = 30
    ws.merge_cells("A4:J4"); c = ws["A4"]; c.value = "Contoh Rasio"; c.font = Font(name=FONT, size=11, bold=True, color=WHITE); c.fill = fill(BLUE); c.alignment = Alignment(vertical="center")
    header_row(ws, 5, rasio_headers, BLUE, 30)
    rasio_rows = [
        ["2026-05", "PKM-TALIWANG", "IKU", "IKU-2026-001.01", "IND-001", "DESA", "5201.01.01", 92, 120, "ILUSTRASI — angka fiktif"],
        ["2026-05", "PKM-TALIWANG", "SPM", "SPM-2026-04", "IND-004", "KECAMATAN", "5201.01", 60, 80, "ILUSTRASI — angka fiktif"],
    ]
    for rr, vals in enumerate(rasio_rows, 6):
        for cc, val in enumerate(vals, 1):
            cell = ws.cell(row=rr, column=cc, value=val)
            body_cell(cell, wrap=True, align="right" if cc in (8,9) else "left", color=LIGHT_GREEN)
        ws.row_dimensions[rr].height = 24
    start = 9
    ws.merge_cells(start_row=start, start_column=1, end_row=start, end_column=9)
    c = ws.cell(row=start, column=1, value="Contoh Jumlah (count) — kode CONTOH-* bukan master")
    c.font = Font(name=FONT, size=11, bold=True, color=WHITE); c.fill = fill("7F7F7F"); c.alignment = Alignment(vertical="center")
    header_row(ws, start + 1, jumlah_headers, "7F7F7F", 30)
    jumlah_rows = [
        ["2026-05", "PKM-TALIWANG", "IKU", "CONTOH-IKU-KEMATIAN", "CONTOH-KASUS-IBU", "DESA", "5201.01.01", 2, "ILUSTRASI — kode contoh, bukan master"],
    ]
    for rr, vals in enumerate(jumlah_rows, start + 2):
        for cc, val in enumerate(vals, 1):
            cell = ws.cell(row=rr, column=cc, value=val)
            body_cell(cell, wrap=True, align="right" if cc == 8 else "left", color=LIGHT_GREEN)
        ws.row_dimensions[rr].height = 24
    note(ws, start + 4, "Jangan mengunggah sheet Contoh. Untuk upload saat ini, gunakan sheet Rasio/template-bnba-dinamis.xlsx; count menunggu kode.", span=10, height=32)


def build_monthly(master):
    wb = Workbook()
    ws_petunjuk = wb.active; ws_petunjuk.title = "Petunjuk"
    ws_rasio = wb.create_sheet("Rasio")
    ws_jumlah = wb.create_sheet("Jumlah")
    ws_ref = wb.create_sheet("Referensi")
    ws_contoh = wb.create_sheet("Contoh")
    ranges = write_referensi(ws_ref, annual=False, master=master)
    write_data_sheet(ws_rasio, mode="rasio", ranges=ranges, has_indicator=True, headers=[
        "periode", "kode_faskes", "kode_kategori", "kode_iku", "kode_indikator", "scope_type", "scope_code", "numerator", "denominator", "catatan"
    ])
    write_data_sheet(ws_jumlah, mode="jumlah", ranges=ranges, has_indicator=False, headers=[
        "periode", "kode_faskes", "kode_kategori", "kode_iku", "kode_indikator", "scope_type", "scope_code", "jumlah", "catatan"
    ])
    write_petunjuk(ws_petunjuk, annual=False)
    write_contoh(ws_contoh, annual=False, ranges=ranges)
    wb.active = 0
    out = OUT / "BNBA-Bulanan-Semua-Satuan.xlsx"
    wb.save(out)
    print("wrote", out)


def build_annual(master):
    wb = Workbook()
    ws_petunjuk = wb.active; ws_petunjuk.title = "Petunjuk"
    ws_data = wb.create_sheet("Data")
    ws_ref = wb.create_sheet("Referensi")
    ws_contoh = wb.create_sheet("Contoh")
    ranges = write_referensi(ws_ref, annual=True, master=master)
    write_data_sheet(ws_data, mode="tahunan", ranges=ranges, has_indicator=False, headers=[
        "tahun", "kode_faskes", "kode_kategori", "kode_iku", "kode_indikator", "scope_type", "scope_code", "nilai", "catatan"
    ])
    write_petunjuk(ws_petunjuk, annual=True)
    write_contoh(ws_contoh, annual=True, ranges=ranges)
    wb.active = 0
    out = OUT / "BNBA-Tahunan.xlsx"
    wb.save(out)
    print("wrote", out)


def main():
    master = load_master()
    OUT.mkdir(parents=True, exist_ok=True)
    STATIC.mkdir(parents=True, exist_ok=True)
    build_monthly(master)
    build_annual(master)
    for name in ["BNBA-Bulanan-Semua-Satuan.xlsx", "BNBA-Tahunan.xlsx"]:
        shutil.copy2(OUT / name, STATIC / name)


if __name__ == "__main__":
    main()
