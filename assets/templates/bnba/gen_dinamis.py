#!/usr/bin/env python3
"""Generate the generic, dynamic BNBA upload template for MATA SIDIK.

Requires: python3 -m pip install openpyxl
Run from any directory: python3 assets/templates/bnba/gen_dinamis.py

This template follows the current pilot contract:
periode, kode_faskes, kode_kategori, kode_iku, kode_indikator,
scope_type, scope_code, nilai, numerator, denominator, catatan.

It is intentionally indicator-agnostic: codes come from the master.
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

ROOT = Path(__file__).resolve().parents[3]
DB = ROOT / "backend/data/matasidik.db"
OUT = ROOT / "assets/templates/bnba/template-bnba-dinamis.xlsx"
VERSION = "v3-dinamis-2026-11-19"

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
    dv.error = "Isi angka > 0." if greater_than_zero else "Isi angka >= 0."
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
        """SELECT i.kode, i.nama, i.satuan,
                  l.kode AS iku_kode, l.nama AS iku_nama, c.kode AS kategori_kode,
                  c.nama AS kategori_nama, l.aggregation_method,
                  pb.nama AS pj_bidang_nama, pt.nama AS pj_timker_nama
           FROM indikator i
           JOIN indikator_label il ON il.indikator_id = i.id AND il.is_primary = 1
           JOIN label l ON l.id = il.label_id
           LEFT JOIN dashboard_category c ON c.id = l.category_id
           LEFT JOIN bidang pb ON pb.id = i.pj_bidang_id
           LEFT JOIN bidang pt ON pt.id = i.pj_timker_id
           WHERE i.aktif=1 AND l.aktif=1
           ORDER BY l.kode"""
    )]
    indicator_labels = [dict(r) for r in q(
        """SELECT i.kode AS indikator_kode, i.nama AS indikator_nama, i.satuan,
                  l.kode AS iku_kode, l.nama AS iku_nama,
                  c.kode AS kategori_kode, c.nama AS kategori_nama,
                  il.is_primary, il.bobot,
                  l.allow_direct_input, l.is_rollup_only
           FROM indikator_label il
           JOIN indikator i ON i.id = il.indikator_id
           JOIN label l ON l.id = il.label_id
           LEFT JOIN dashboard_category c ON c.id = l.category_id
           WHERE i.aktif=1 AND l.aktif=1
           ORDER BY i.kode, il.is_primary DESC, l.kode"""
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
    return categories, labels, indicators, indicator_labels, wilayah, faskes, periods, releases


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


def main():
    categories, labels, indicators, indicator_labels, wilayah, faskes, periods, releases = load_master()
    leaf_labels = [l for l in labels if l["allow_direct_input"] and not l["is_rollup_only"]]
    leaf_labels = [l for l in leaf_labels if any(x["iku_kode"] == l["kode"] for x in indicator_labels)]
    master_year = releases[0]["tahun"] if releases else 2026

    wb = Workbook()
    ws_data = wb.active
    ws_data.title = "Data"
    ws_petunjuk = wb.create_sheet("Petunjuk")
    ws_ref = wb.create_sheet("Referensi")
    ws_contoh = wb.create_sheet("Contoh")

    headers = [
        "periode", "kode_faskes", "kode_kategori", "kode_iku", "kode_indikator",
        "scope_type", "scope_code", "nilai", "numerator", "denominator", "catatan",
    ]
    header_row(ws_data, 1, headers, NAVY, 36)
    ws_data.sheet_view.showGridLines = False
    ws_data.sheet_properties.tabColor = NAVY

    widths = {
        "periode": 12, "kode_faskes": 18, "kode_kategori": 14, "kode_iku": 22,
        "kode_indikator": 16, "scope_type": 14, "scope_code": 16, "nilai": 10,
        "numerator": 12, "denominator": 14, "catatan": 34,
    }
    cl = {h: get_column_letter(i) for i, h in enumerate(headers, 1)}
    for h in headers:
        ws_data.column_dimensions[cl[h]].width = widths[h]

    comments = {
        "periode": "Wajib YYYY-MM. Harus sama dengan periode yang dipilih saat upload.",
        "kode_faskes": "Kode faskes pengunggah. Satu file = satu faskes.",
        "kode_kategori": "Opsional; jika diisi, harus cocok dengan kategori label.",
        "kode_iku": "Satu atau beberapa kode IKU, pisahkan dengan titik koma (;). Label pertama = utama/konteks. Semua harus leaf dan terhubung ke kode_indikator.",
        "kode_indikator": "Kode indikator dinamis dari master; harus terhubung ke kode_iku.",
        "scope_type": "DESA, KECAMATAN, atau KABUPATEN. Faskes bukan scope; faskes adalah pengunggah/pemilik data.",
        "scope_code": "Kode wilayah resmi sesuai scope_type (desa/kecamatan/kabupaten).",
        "nilai": "Opsional; untuk rasio biarkan kosong dan isi numerator/denominator.",
        "numerator": "Angka mentah; jangan dibulatkan.",
        "denominator": "Wajib > 0 untuk indikator rasio.",
        "catatan": "Catatan baris, boleh kosong.",
    }
    for i, h in enumerate(headers, 1):
        if h in comments:
            ws_data.cell(row=1, column=i).comment = Comment(comments[h], "MATA SIDIK")

    numeric_cols = {"nilai", "numerator", "denominator"}
    for r in range(2, 502):
        for i, h in enumerate(headers, 1):
            cell = ws_data.cell(row=r, column=i)
            body_cell(cell, align="right" if h in numeric_cols else "left",
                      color=WHITE if h == "catatan" else YELLOW)
            if h in ("periode", "kode_faskes", "kode_kategori", "kode_iku", "kode_indikator", "scope_type", "scope_code"):
                cell.number_format = "@"
            elif h in numeric_cols:
                cell.number_format = "0.###"
            cell.protection = Protection(locked=False)
        ws_data.row_dimensions[r].height = 18

    # ------------------------- Referensi -------------------------
    ws_ref.sheet_view.showGridLines = False
    ws_ref.sheet_properties.tabColor = GREEN
    widths_ref = {"A": 22, "B": 34, "C": 18, "D": 24, "E": 16, "F": 16, "G": 16, "H": 22, "I": 14, "J": 18, "K": 18, "L": 14}
    for col, w in widths_ref.items():
        ws_ref.column_dimensions[col].width = w

    ws_ref.merge_cells("A1:L1")
    t = ws_ref["A1"]
    t.value = "REFERENSI MASTER — JANGAN DIUBAH"
    t.font = TITLE
    t.fill = fill(GREEN)
    t.alignment = Alignment(vertical="center")
    ws_ref.row_dimensions[1].height = 32

    ws_ref.merge_cells("A2:L2")
    c = ws_ref["A2"]
    c.value = f"Sumber: master pilot `backend/data/matasidik.db` per 2026-11-19. Tahun master: {master_year}. Daftar ini dinamis; regenerate setelah master berubah."
    c.font = SUBTITLE
    c.fill = fill(LIGHT_GREEN)
    c.alignment = Alignment(vertical="center", wrap_text=True)
    ws_ref.row_dimensions[2].height = 30

    ranges = {}
    row = 4

    start, end = ref_table(
        ws_ref, row, "1. Periode",
        ["kode", "keterangan"],
        [[p["kode"], f'Bulan {p["bulan"]:02d} tahun {p["tahun"]}'] for p in periods],
    )
    ranges["periode"] = (start, end)

    row = end + 3
    start, end = ref_table(
        ws_ref, row, "2. Faskes",
        ["kode", "nama", "id_internal"],
        [[f["kode"], f["nama"], f["id"]] for f in faskes],
    )
    ranges["faskes"] = (start, end)

    row = end + 3
    start, end = ref_table(
        ws_ref, row, "3. Wilayah",
        ["kode", "nama", "level", "kode_parent"],
        [[w["kode"], w["nama"], w["level"], w["parent_kode"] or "—"] for w in wilayah],
    )
    ranges["wilayah"] = (start, end)

    row = end + 3
    scope_rows = []
    for w in wilayah:
        scope_rows.append([w["level"].upper() if w["level"] != "kabupaten" else "KABUPATEN", w["kode"], w["nama"], w["level"], "Gunakan kode ini untuk filter wilayah"])
    start, end = ref_table(
        ws_ref, row, "4. Scope — level dan kode untuk filter",
        ["scope_type", "kode", "nama", "level", "catatan"],
        scope_rows,
    )
    ranges["scope"] = (start, end)

    row = end + 3
    start, end = ref_table(
        ws_ref, row, "5. Kategori",
        ["kode", "nama", "jenis"],
        [[c["kode"], c["nama"], c["jenis"]] for c in categories],
    )
    ranges["kategori"] = (start, end)

    row = end + 3
    start, end = ref_table(
        ws_ref, row, "6. IKU Leaf (dinamis; hanya leaf yang boleh diinput)",
        ["kode", "nama", "kategori", "level", "aggregation_method"],
        [[l["kode"], l["nama"], l["kategori_kode"] or "", l["level"], l["aggregation_method"]] for l in leaf_labels],
    )
    ranges["iku"] = (start, end)

    row = end + 3
    start, end = ref_table(
        ws_ref, row, "7. Indikator (dinamis)",
        ["kode", "nama", "kode_iku_utama", "kategori", "satuan", "agregasi"],
        [[i["kode"], i["nama"], i["iku_kode"], i["kategori_kode"] or "", i["satuan"], i["aggregation_method"]] for i in indicators],
    )
    ranges["indikator"] = (start, end)

    row = end + 3
    start, end = ref_table(
        ws_ref, row, "8. Mapping Semua Label IKU ↔ Indikator (many-to-many)",
        ["kode_indikator", "nama_indikator", "kode_iku", "nama_iku", "kategori", "is_primary", "is_input", "satuan"],
        [[x["indikator_kode"], x["indikator_nama"], x["iku_kode"], x["iku_nama"],
          x["kategori_kode"] or "", "Ya" if x["is_primary"] else "Tidak",
          "Ya" if (x["allow_direct_input"] and not x["is_rollup_only"]) else "Tidak", x["satuan"]]
         for x in indicator_labels],
    )
    ranges["mapping"] = (start, end)

    row = end + 3
    start, end = ref_table(
        ws_ref, row, "9. Tipe Nilai (status implementasi)",
        ["kode", "nama", "status"],
        [
            ["RATIO", "Rasio/persentase: numerator/denominator", "Didukung parser/classifier pilot"],
            ["COUNT", "Jumlah kasus/dokumen/laporan", "Belum didukung mesin pilot"],
            ["SCALAR", "Nilai tunggal: tahun/poin/indeks", "Belum didukung mesin pilot"],
        ],
    )
    ranges["tipe_nilai"] = (start, end)

    row = end + 3
    start, end = ref_table(
        ws_ref, row, "10. PJ Bidang/Timker (bukan label IKU)",
        ["kode_indikator", "nama_indikator", "pj_bidang", "pj_timker"],
        [[i["kode"], i["nama"], i["pj_bidang_nama"] or "—", i["pj_timker_nama"] or "—"] for i in indicators],
    )
    ranges["pj"] = (start, end)

    row = end + 3
    note(
        ws_ref, row,
        "Catatan: template ini sengaja tidak mengunci daftar indikator. Kode di sheet ini berasal dari master dan akan berubah saat master berubah. "
        "Tabel mapping menampilkan semua label per indikator; pada master pilot saat ini masih 1:1, sehingga baru terlihat satu baris per indikator. "
        "Setelah label ditambahkan, satu indikator akan muncul beberapa baris di mapping, tetapi tetap satu baris di Data. "
        "Jangan mengetik kode baru di file ini.",
        color=YELLOW, span=12, height=44,
    )

    # ------------------------- Data validations -------------------------
    add_list_dv(ws_data, "A2:A501", f"=Referensi!$A${ranges['periode'][0]}:$A${ranges['periode'][1]}",
                "Periode", "Pilih periode YYYY-MM.")
    add_list_dv(ws_data, "B2:B501", f"=Referensi!$A${ranges['faskes'][0]}:$A${ranges['faskes'][1]}",
                "Faskes", "Pilih kode faskes pengunggah.")
    add_list_dv(ws_data, "C2:C501", f"=Referensi!$A${ranges['kategori'][0]}:$A${ranges['kategori'][1]}",
                "Kategori", "Opsional; pilih kategori label.")
    add_list_dv(ws_data, "E2:E501", f"=Referensi!$A${ranges['indikator'][0]}:$A${ranges['indikator'][1]}",
                "Kode Indikator", "Pilih indikator yang terhubung ke kode_iku.")
    add_list_dv(ws_data, "F2:F501", '"DESA,KECAMATAN,KABUPATEN"',
                "Scope Type", "Pilih level wilayah data; faskes bukan scope.")
    add_list_dv(ws_data, "G2:G501", f"=Referensi!$A${ranges['scope'][0]}:$A${ranges['scope'][1]}",
                "Scope Code", "Pilih kode wilayah sesuai scope_type.")
    add_decimal_dv(ws_data, "H2:I501", "Nilai", "Isi angka >= 0; denominator harus > 0 untuk rasio.")
    # Denominator specifically > 0.
    den_dv = DataValidation(type="decimal", operator="greaterThan", formula1="0", allow_blank=True)
    den_dv.errorTitle = "Denominator tidak valid"
    den_dv.error = "Denominator wajib > 0 untuk indikator rasio."
    den_dv.promptTitle = "Denominator"
    den_dv.prompt = "Isi angka > 0."
    ws_data.add_data_validation(den_dv)
    den_dv.add("J2:J501")

    red_fill = PatternFill(start_color=LIGHT_RED, end_color=LIGHT_RED, fill_type="solid")
    orange_fill = PatternFill(start_color=LIGHT_ORANGE, end_color=LIGHT_ORANGE, fill_type="solid")
    ws_data.conditional_formatting.add("J2:J501", CellIsRule(operator="lessThanOrEqual", formula=["0"], fill=red_fill))
    ws_data.conditional_formatting.add("I2:I501", FormulaRule(formula=['AND($E2<>"",I2="")'], fill=orange_fill))
    ws_data.conditional_formatting.add("J2:J501", FormulaRule(formula=['AND($E2<>"",J2="")'], fill=orange_fill))

    ws_data.freeze_panes = "A2"
    ws_data.auto_filter.ref = "A1:K501"
    ws_data.print_title_rows = "1:1"
    ws_data.page_setup.orientation = "landscape"
    ws_data.page_setup.fitToWidth = 1
    ws_data.sheet_properties.pageSetUpPr.fitToPage = True
    protect(ws_data)

    # ------------------------- Petunjuk -------------------------
    ws_petunjuk.sheet_view.showGridLines = False
    ws_petunjuk.sheet_properties.tabColor = BLUE
    ws_petunjuk.column_dimensions["A"].width = 3
    ws_petunjuk.column_dimensions["B"].width = 30
    ws_petunjuk.column_dimensions["C"].width = 105

    ws_petunjuk.merge_cells("B1:C1")
    t = ws_petunjuk["B1"]
    t.value = "TEMPLATE BNBA DINAMIS — MATA SIDIK"
    t.font = TITLE
    t.fill = fill(NAVY)
    t.alignment = Alignment(vertical="center")
    ws_petunjuk.row_dimensions[1].height = 36

    ws_petunjuk.merge_cells("B2:C2")
    c = ws_petunjuk["B2"]
    c.value = "Kontrak upload generik: indikator dan IKU dinamis dari master; scope wilayah eksplisit (desa/kecamatan/kabupaten); faskes adalah pengunggah."
    c.font = SUBTITLE
    c.fill = fill(LIGHT_BLUE)
    c.alignment = Alignment(vertical="center", wrap_text=True)
    ws_petunjuk.row_dimensions[2].height = 34

    ws_petunjuk.merge_cells("B3:C3")
    c = ws_petunjuk["B3"]
    c.value = "STATUS: KONTRAK UPLOAD PILOT — diterima parser/classifier untuk baris rasio. Filter/kalkulasi per desa/kecamatan masih menunggu kode."
    c.font = Font(name=FONT, size=10, bold=True, color=RED)
    c.fill = fill(YELLOW)
    c.alignment = Alignment(vertical="center", wrap_text=True)
    ws_petunjuk.row_dimensions[3].height = 30

    r = 5
    section(ws_petunjuk, r, "1. Prinsip dinamis")
    r += 1
    for label, text in [
        ("Indikator", "Tidak ada daftar indikator hardcode di sheet Data. Kode diambil dari master; tambah indikator lewat Master lalu regenerate template."),
        ("IKU", "kode_iku hanya leaf/direct-input. Boleh lebih dari satu, pisahkan dengan titik koma (;). Label pertama = utama/konteks."),
        ("Mapping", "kode_indikator harus terhubung ke kode_iku yang dipilih. Jika tidak, baris masuk karantina."),
        ("Many-to-many", "Satu indikator boleh punya banyak label. Isi beberapa kode di kode_iku dengan pemisah ; (contoh: IKU-A;IKK-B). Label pertama = utama/konteks; semua label tetap satu nilai. Jangan menduplikasi baris per label."),
        ("PJ", "Bidang/Timker adalah dimensi PJ, bukan label IKU. Jangan memasukkannya ke kode_iku."),
    ]:
        kv(ws_petunjuk, r, label, text, height=38); r += 1
    r += 1

    section(ws_petunjuk, r, "2. Level wilayah dan filter")
    r += 1
    for label, text in [
        ("DESA", "scope_type=DESA; scope_code=kode desa. Ini grain terendah untuk filter desa."),
        ("KECAMATAN", "scope_type=KECAMATAN; scope_code=kode kecamatan. Jika data desa sudah ada, kecamatan diturunkan; jangan input dua level sekaligus."),
        ("KABUPATEN", "scope_type=KABUPATEN; scope_code=kode kabupaten. Untuk data yang memang hanya level kabupaten."),
        ("Faskes", "Faskes bukan scope. Faskes adalah pengunggah/pemilik data (kode_faskes). Filter puskesmas dihitung dari kode_faskes, bukan scope_type."),
        ("Aturan", "Simpan satu grain terendah per indikator/periode; level atas diturunkan. Jangan menulis baris desa + kecamatan + kabupaten sekaligus."),
    ]:
        kv(ws_petunjuk, r, label, text, height=40); r += 1
    r += 1

    section(ws_petunjuk, r, "3. Cara isi")
    r += 1
    for label, text in [
        ("Langkah 1", "Buka sheet Data; jangan ubah header atau urutan kolom."),
        ("Langkah 2", "Isi satu baris = satu nilai indikator + satu scope + satu periode. Jangan mengulang baris untuk setiap label."),
        ("Langkah 3", "Pilih kode_indikator dari dropdown. Isi kode_iku satu atau beberapa, pisahkan dengan ; (contoh: IKU-A;IKK-B). Label pertama = utama/konteks."),
        ("Langkah 4", "Pilih scope_type dan scope_code sesuai level filter yang diinginkan."),
        ("Langkah 5", "Isi numerator dan denominator; biarkan nilai kosong. Denominator wajib > 0."),
        ("Langkah 6", "Jangan menambahkan total/subtotal, merged cell, atau formula."),
        ("Langkah 7", "Simpan .xlsx, unggah, lalu periksa preview dan baris bermasalah."),
    ]:
        kv(ws_petunjuk, r, label, text, height=34); r += 1
    r += 1

    section(ws_petunjuk, r, "4. Aturan wajib")
    r += 1
    for label, text in [
        ("Satu file", "Satu faskes + satu periode (bulan)."),
        ("Rasio", "Isi numerator/denominator; jangan isi persentase jadi."),
        ("Blank ≠ 0", "Blank = belum ada data; 0 = nol yang dilaporkan."),
        ("Tanpa PII", "Jangan masukkan nama, NIK, atau identitas individu."),
        ("Tanpa formula", "Jangan mengubah header atau menambah formula di Data."),
    ]:
        kv(ws_petunjuk, r, label, text, height=30); r += 1
    r += 1

    section(ws_petunjuk, r, "5. Kesiapan kode pilot (fakta)")
    r += 1
    for label, text in [
        ("Parser", "Menerima header ini; sheet Data harus berisi minimal satu baris."),
        ("Klasifikasi", "Setiap kode_iku dipisah dengan ; lalu dicocokkan ke master; semua harus leaf dan terhubung ke indikator; baris ambigu masuk karantina."),
        ("Scope upload", "Hanya DESA/KECAMATAN/KABUPATEN. Faskes bukan scope; faskes adalah pengunggah."),
        ("Kalkulasi", "Pilot hanya menghitung KABUPATEN/FASKES internal; nilai DESA/KECAMATAN tersimpan tetapi belum dihitung per level."),
        ("Filter dashboard", "Filter desa/kecamatan belum aktif; filter faskes memakai kode_faskes pengunggah."),
        ("Count/scalar", "Belum didukung; mesin pilot menghitung semua input langsung sebagai rasio."),
        ("Satuan", "Template ini untuk rasio; satuan hasil %. Count/tahun/poin perlu penyesuaian kode."),
    ]:
        kv(ws_petunjuk, r, label, text, height=38); r += 1
    r += 1

    section(ws_petunjuk, r, "6. Catatan")
    r += 1
    for label, text in [
        ("Dinamis", "Template tidak mengunci indikator. Setelah master berubah, regenerate template; jangan mengedit kode di file. Tabel mapping akan menampilkan semua label per indikator."),
        ("Label", "Satu indikator banyak label; nilai tetap satu. Beberapa kode di kode_iku dipisahkan ;. Label owner (Bidang/Timker) tidak dimasukkan ke kode_iku."),
        ("Scope", "Gunakan kode wilayah resmi. Faskes diisi di kode_faskes, bukan di scope_type."),
        ("Contoh", "Angka di sheet Contoh fiktif; jangan dianggap data produksi."),
        ("Tanpa kode", "Template ini tidak mengubah kode aplikasi. Penyesuaian kode ada di docs/05."),
    ]:
        kv(ws_petunjuk, r, label, text, height=34); r += 1

    # ------------------------- Contoh -------------------------
    ws_contoh.sheet_view.showGridLines = False
    ws_contoh.sheet_properties.tabColor = GREEN
    ws_contoh.merge_cells("A1:K1")
    t = ws_contoh["A1"]
    t.value = "CONTOH — TEMPLATE DINAMIS — ANGKA FIKTIF"
    t.font = TITLE
    t.fill = fill(GREEN)
    t.alignment = Alignment(vertical="center")
    ws_contoh.row_dimensions[1].height = 34

    ws_contoh.merge_cells("A2:K2")
    c = ws_contoh["A2"]
    c.value = "Contoh level scope: DESA, KECAMATAN, KABUPATEN. Faskes ada di kolom kode_faskes, bukan scope. Jangan mengunggah sheet Contoh; salin polanya ke Data."
    c.font = SUBTITLE
    c.fill = fill(LIGHT_GREEN)
    c.alignment = Alignment(vertical="center", wrap_text=True)
    ws_contoh.row_dimensions[2].height = 30

    header_row(ws_contoh, 4, headers, GREEN, 34)
    for h in headers:
        ws_contoh.column_dimensions[cl[h]].width = widths[h]
    examples = [
        ["2026-05", "PKM-TALIWANG", "IKU", "IKU-2026-001.01", "IND-001", "DESA", "5201.01.01", None, 92, 120, "ILUSTRASI — angka fiktif"],
        ["2026-05", "PKM-TALIWANG", "SPM", "SPM-2026-04", "IND-004", "KECAMATAN", "5201.01", None, 60, 80, "ILUSTRASI — angka fiktif"],
        ["2026-05", "PKM-TALIWANG", "IKU", "IKU-2026-001.02", "IND-002", "KABUPATEN", "5201", None, 64, 80, "ILUSTRASI — angka fiktif"],
    ]
    for rr, rowvals in enumerate(examples, 5):
        for cc, val in enumerate(rowvals, 1):
            cell = ws_contoh.cell(row=rr, column=cc, value=val)
            body_cell(cell, wrap=True, align="right" if cc in (8, 9, 10) else "left", color=LIGHT_GREEN)
            if cc in (8, 9, 10):
                cell.number_format = "0.###"
        ws_contoh.row_dimensions[rr].height = 24

    note(
        ws_contoh, 11,
        "Satu baris hanya untuk satu level scope. Jangan memasukkan desa + kecamatan + kabupaten sekaligus untuk indikator/periode yang sama. "
        "Faskes bukan scope; faskes adalah pengunggah (kode_faskes). Scope hanya wilayah. "
        "kode_iku boleh lebih dari satu, pisahkan dengan ; (contoh: IKU-A;IKK-B); label pertama = utama/konteks.",
        color=YELLOW, span=11, height=36,
    )

    wb.properties.title = "Template BNBA Dinamis — MATA SIDIK"
    wb.properties.subject = "Kontrak upload generik: IKU/indikator dinamis + scope wilayah"
    wb.properties.creator = "MATA SIDIK"
    wb.properties.description = (
        "Template generik untuk parser pilot. Indikator/IKU dinamis dari master; scope desa/kecamatan/faskes/kabupaten. "
        "Hanya baris rasio (numerator/denominator) yang didukung mesin pilot saat ini."
    )
    wb.properties.keywords = "BNBA, MATA SIDIK, dinamis, scope, IKU"
    wb.active = 0
    OUT.parent.mkdir(parents=True, exist_ok=True)
    wb.save(OUT)
    static = ROOT / "frontend/static/templates"
    static.mkdir(parents=True, exist_ok=True)
    shutil.copy2(OUT, static / OUT.name)
    print(f"Wrote {OUT} ({OUT.stat().st_size} bytes)")


if __name__ == "__main__":
    main()
