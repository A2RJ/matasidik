#!/usr/bin/env python3
"""Generate the single BNBA template for all units and period granularities.

The same file supports:
- Rasio (numerator/denominator), period YYYY-MM-DD/YYYY-MM/YYYY
- Nilai for count (Kasus/Dokumen/Laporan/Berita Acara/orang) and scalar
  (UHH/Poin/AKIP), same period formats

Output: assets/templates/bnba/final/BNBA-Semua-Satuan.xlsx
and a copy in frontend/static/templates/.

Requires: python3 -m pip install openpyxl
Run from any directory: python3 assets/templates/bnba/final/gen_template.py
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
WHITE = "FFFFFF"
RED = "C00000"
LIGHT_RED = "F8CBAD"
MID_GRAY = "BFBFBF"

THIN = Side(style="thin", color=MID_GRAY)
BORDER = Border(left=THIN, right=THIN, top=THIN, bottom=THIN)
FONT = "Calibri"
BODY = Font(name=FONT, size=10)
BOLD = Font(name=FONT, size=10, bold=True)
HEADER = Font(name=FONT, size=10, bold=True, color=WHITE)
TITLE = Font(name=FONT, size=18, bold=True, color=WHITE)
SUB = Font(name=FONT, size=10, color="404040")


def fill(c):
    return PatternFill("solid", fgColor=c)


def hdr(ws, row, headers, color=NAVY, height=34):
    for i, h in enumerate(headers, 1):
        c = ws.cell(row=row, column=i, value=h)
        c.font = HEADER
        c.fill = fill(color)
        c.alignment = Alignment(horizontal="center", vertical="center", wrap_text=True)
        c.border = BORDER
    ws.row_dimensions[row].height = height


def cell_style(c, align="left", color=None, wrap=False):
    c.font = BODY
    c.border = BORDER
    c.alignment = Alignment(horizontal=align, vertical="top", wrap_text=wrap)
    if color:
        c.fill = fill(color)


def section(ws, row, title, color=BLUE):
    ws.merge_cells(start_row=row, start_column=2, end_row=row, end_column=3)
    c = ws.cell(row=row, column=2, value=title)
    c.font = Font(name=FONT, size=11, bold=True, color=WHITE)
    c.fill = fill(color)
    c.alignment = Alignment(vertical="center")
    ws.row_dimensions[row].height = 22


def kv(ws, row, label, value, height=24):
    a = ws.cell(row=row, column=2, value=label)
    b = ws.cell(row=row, column=3, value=value)
    a.font = BOLD
    a.fill = fill(PALE_BLUE)
    a.border = BORDER
    a.alignment = Alignment(vertical="top", wrap_text=True)
    b.font = BODY
    b.border = BORDER
    b.alignment = Alignment(vertical="top", wrap_text=True)
    ws.row_dimensions[row].height = height


def list_dv(ws, rng, formula, title, prompt):
    dv = DataValidation(type="list", formula1=formula, allow_blank=True, showDropDown=False)
    dv.errorTitle = "Nilai tidak valid"
    dv.error = "Pilih nilai dari sheet Referensi atau ikuti format."
    dv.promptTitle = title
    dv.prompt = prompt
    ws.add_data_validation(dv)
    dv.add(rng)


def num_dv(ws, rng, title, prompt, positive=False):
    dv = DataValidation(type="decimal", operator="greaterThan" if positive else "greaterThanOrEqual", formula1="0", allow_blank=True)
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
    categories = [dict(r) for r in q("SELECT kode,nama,jenis FROM dashboard_category WHERE aktif=1 ORDER BY urutan,kode")]
    indicators = [dict(r) for r in q(
        """SELECT i.kode,i.nama,i.satuan,i.tipe_nilai,l.kode AS iku_kode,l.nama AS iku_nama,
                  c.kode AS kategori_kode,l.aggregation_method,
                  pb.nama AS pj_bidang_nama,pt.nama AS pj_timker_nama
           FROM indikator i
           JOIN indikator_label il ON il.indikator_id=i.id AND il.is_primary=1
           JOIN label l ON l.id=il.label_id
           LEFT JOIN dashboard_category c ON c.id=l.category_id
           LEFT JOIN bidang pb ON pb.id=i.pj_bidang_id
           LEFT JOIN bidang pt ON pt.id=i.pj_timker_id
           WHERE i.aktif=1 AND l.aktif=1 ORDER BY l.kode"""
    )]
    links = [dict(r) for r in q(
        """SELECT i.kode AS indikator_kode,i.nama AS indikator_nama,i.satuan,
                  l.kode AS iku_kode,l.nama AS iku_nama,c.kode AS kategori_kode,
                  il.is_primary,l.allow_direct_input,l.is_rollup_only
           FROM indikator_label il
           JOIN indikator i ON i.id=il.indikator_id
           JOIN label l ON l.id=il.label_id
           LEFT JOIN dashboard_category c ON c.id=l.category_id
           WHERE i.aktif=1 AND l.aktif=1
           ORDER BY i.kode,il.is_primary DESC,l.kode"""
    )]
    wilayah = [dict(r) for r in q("SELECT kode,nama,level,parent_id FROM wilayah ORDER BY kode")]
    faskes = [dict(r) for r in q("SELECT id,kode,nama FROM faskes WHERE aktif=1 ORDER BY kode")]
    con.close()
    return categories, indicators, links, wilayah, faskes


def ref_table(ws, start, title, headers, rows, color=BLUE):
    ws.merge_cells(start_row=start, start_column=1, end_row=start, end_column=len(headers))
    c = ws.cell(row=start, column=1, value=title)
    c.font = Font(name=FONT, size=11, bold=True, color=WHITE)
    c.fill = fill(color)
    c.alignment = Alignment(vertical="center")
    ws.row_dimensions[start].height = 22
    hdr(ws, start + 1, headers, color, 28)
    for rr, vals in enumerate(rows, start + 2):
        for cc, val in enumerate(vals, 1):
            cell_style(ws.cell(row=rr, column=cc, value=val), wrap=True)
        ws.row_dimensions[rr].height = 22
    return start + 2, start + 1 + len(rows)


UNIT_ROWS = [
    ["(tahun)", "SCALAR", "nilai", "LAST_VALUE", "UHH; data 1x/tahun, kabupaten"],
    ["Kasus", "COUNT", "nilai", "SUM", "Kematian ibu/balita; kejadian"],
    ["%", "RATIO", "numerator + denominator", "RATIO_OF_SUMS", "Prevalensi/cakupan"],
    ["Persentase", "RATIO", "numerator + denominator", "RATIO_OF_SUMS", "Cakupan"],
    ["(%)", "RATIO", "numerator + denominator", "RATIO_OF_SUMS", "Cakupan"],
    ["Poin", "SCALAR", "nilai", "LAST_VALUE", "AKIP/indeks"],
    ["Rasio", "RATIO", "numerator + denominator", "RATIO_OF_SUMS", "Rasio non-persen; hasil num/den"],
    ["Dokumen", "COUNT", "nilai", "SUM", "Jumlah dokumen"],
    ["Laporan", "COUNT", "nilai", "SUM", "Jumlah laporan"],
    ["Berita Acara", "COUNT", "nilai", "SUM", "Jumlah berita acara"],
    ["orang", "COUNT/RATIO", "nilai atau numerator/denominator", "SUM/RATIO_OF_SUMS", "Tergantung indikator"],
]

PERIODE_ROWS = [
    ["2026-05-15", "Harian: YYYY-MM-DD", "Data harian"],
    ["2026-05", "Bulanan: YYYY-MM", "Data bulanan"],
    ["2026", "Tahunan: YYYY", "Data 1x/tahun"],
]

POLA_ROWS = [
    ["RASIO", "Rasio/persentase", "Sheet Rasio", "numerator + denominator"],
    ["NILAI", "Jumlah/kejadian/scalar", "Sheet Nilai", "nilai"],
]


def write_ref(ws, master):
    categories, indicators, links, wilayah, faskes = master
    ws.sheet_view.showGridLines = False
    ws.sheet_properties.tabColor = GREEN
    for col, w in {"A":22,"B":36,"C":22,"D":18,"E":24,"F":18,"G":18,"H":24,"I":16,"J":16}.items():
        ws.column_dimensions[col].width = w

    ws.merge_cells("A1:J1")
    c = ws["A1"]; c.value = "REFERENSI MASTER, SATUAN & PERIODE"; c.font = TITLE; c.fill = fill(GREEN); c.alignment = Alignment(vertical="center")
    ws.row_dimensions[1].height = 32
    ws.merge_cells("A2:J2")
    c = ws["A2"]; c.value = "Sumber master pilot `backend/data/matasidik.db` per 2026-11-19. Satuan mengikuti sheet 2 `Proses Bisnis Mata Sidik.xlsx`. Regenerate setelah master berubah."
    c.font = SUB; c.fill = fill(LIGHT_GREEN); c.alignment = Alignment(vertical="center", wrap_text=True); ws.row_dimensions[2].height = 34

    ranges = {}
    row = 4
    start, end = ref_table(ws, row, "1. Satuan dan tipe nilai", ["satuan","tipe_nilai","field_wajib","agregasi","contoh"], UNIT_ROWS)
    ranges["satuan"] = (start, end)

    row = end + 3
    start, end = ref_table(ws, row, "2. Pola input", ["kode","nama","sheet","field"], POLA_ROWS)
    ranges["pola"] = (start, end)

    row = end + 3
    start, end = ref_table(ws, row, "3. Format periode (satu kolom, tiga granularitas)", ["contoh","format","penggunaan"], PERIODE_ROWS)
    ranges["periode_format"] = (start, end)

    row = end + 3
    start, end = ref_table(ws, row, "4. Kategori dashboard", ["kode","nama","jenis"], [[c["kode"],c["nama"],c["jenis"]] for c in categories])
    ranges["kategori"] = (start, end)

    row = end + 3
    start, end = ref_table(ws, row, "5. Faskes pengunggah", ["kode","nama"], [[f["kode"],f["nama"]] for f in faskes])
    ranges["faskes"] = (start, end)

    row = end + 3
    scope_rows = []
    for w in wilayah:
        scope_type = "KABUPATEN" if w["level"] == "kabupaten" else w["level"].upper()
        scope_rows.append([w["kode"], w["nama"], scope_type, w["level"]])
    start, end = ref_table(ws, row, "6. Scope wilayah (DESA/KECAMATAN/KABUPATEN)", ["kode","nama","scope_type","level"], scope_rows)
    ranges["scope"] = (start, end)

    row = end + 3
    ind_rows = []
    for i in indicators:
        satuan = i["satuan"]
        tipe = (i.get("tipe_nilai") or "RATIO").upper()
        ind_rows.append([i["kode"],i["nama"],i["iku_kode"],i["iku_nama"],i["kategori_kode"] or "",satuan,"%" if tipe=="RATIO" else satuan,tipe,"RASIO" if tipe=="RATIO" else "NILAI",i["aggregation_method"]])
    start, end = ref_table(ws, row, "7. Indikator master dan tipe nilai", ["kode","nama","kode_iku","nama_iku","kategori","satuan_master","satuan_hasil","tipe_nilai","pola","agregasi"], ind_rows)
    ranges["indikator"] = (start, end)

    row = end + 3
    mapping = []
    for x in links:
        mapping.append([x["indikator_kode"],x["indikator_nama"],x["iku_kode"],x["iku_nama"],x["kategori_kode"] or "","Ya" if x["is_primary"] else "Tidak","Ya" if (x["allow_direct_input"] and not x["is_rollup_only"]) else "Tidak",x["satuan"]])
    start, end = ref_table(ws, row, "8. Mapping semua label IKU ↔ indikator (many-to-many)", ["kode_indikator","nama_indikator","kode_iku","nama_iku","kategori","is_primary","is_input","satuan"], mapping)
    ranges["mapping"] = (start, end)

    row = end + 3
    start, end = ref_table(ws, row, "9. PJ Bidang/Timker (bukan label)", ["kode_indikator","nama_indikator","pj_bidang","pj_timker"], [[i["kode"],i["nama"],i["pj_bidang_nama"] or "—",i["pj_timker_nama"] or "—"] for i in indicators])
    ranges["pj"] = (start, end)

    row = end + 3
    ws.merge_cells(start_row=row,start_column=1,end_row=row,end_column=10)
    c = ws.cell(row=row,column=1,value=(
        "Catatan: satu template ini menangani periode harian, bulanan, dan tahunan melalui kolom `periode`. "
        "Jangan membuat file tahunan terpisah. Data tahunan memakai `YYYY`; bulanan `YYYY-MM`; harian `YYYY-MM-DD`. "
        "Rasio memakai numerator/denominator; count/scalar memakai nilai. Tambahkan indikator dan tipe nilai lewat Master."
    ))
    c.font = BODY; c.fill = fill(YELLOW); c.border = BORDER; c.alignment = Alignment(vertical="top", wrap_text=True); ws.row_dimensions[row].height = 58
    return ranges


def write_data(ws, mode, ranges, has_indicator, headers):
    hdr(ws, 1, headers, NAVY, 36)
    ws.sheet_view.showGridLines = False
    ws.sheet_properties.tabColor = NAVY
    widths = {"periode":14,"kode_faskes":18,"kode_kategori":14,"kode_iku":22,"kode_indikator":22,"scope_type":14,"scope_code":16,"numerator":12,"denominator":14,"nilai":12,"catatan":34}
    cl = {h:get_column_letter(i) for i,h in enumerate(headers,1)}
    for h in headers:
        ws.column_dimensions[cl[h]].width = widths.get(h,14)

    comments = {
        "periode":"Wajib. Harian YYYY-MM-DD; bulanan YYYY-MM; tahunan YYYY. Satu kolom untuk semua granularitas.",
        "kode_faskes":"Kode faskes pengunggah/pemilik data.",
        "kode_kategori":"Opsional; harus cocok dengan kategori label utama.",
        "kode_iku":"Satu/beberapa kode IKU leaf, pisahkan ;. Label pertama = utama/konteks.",
        "kode_indikator":"Kode indikator dari master; menentukan satuan/pola.",
        "scope_type":"DESA/KECAMATAN/KABUPATEN. Faskes bukan scope.",
        "scope_code":"Kode wilayah resmi sesuai scope_type.",
        "numerator":"Untuk rasio: pembilang.",
        "denominator":"Untuk rasio: wajib > 0.",
        "nilai":"Untuk count/scalar: angka jumlah/nilai sesuai satuan indikator.",
        "catatan":"Catatan baris, boleh kosong.",
    }
    for i,h in enumerate(headers,1):
        if h in comments:
            ws.cell(row=1,column=i).comment = Comment(comments[h],"MATA SIDIK")

    numeric = {"numerator","denominator","nilai"}
    for r in range(2,502):
        for i,h in enumerate(headers,1):
            c = ws.cell(row=r,column=i)
            cell_style(c, align="right" if h in numeric else "left", color=WHITE if h=="catatan" else YELLOW)
            c.number_format = "@" if h in ("periode","kode_faskes","kode_kategori","kode_iku","kode_indikator","scope_type","scope_code") else ("0.###" if h in numeric else "General")
            c.protection = Protection(locked=False)
        ws.row_dimensions[r].height = 18

    # No dropdown for periode: one column supports YYYY-MM-DD, YYYY-MM, YYYY.
    list_dv(ws,"B2:B501",f"=Referensi!$A${ranges['faskes'][0]}:$A${ranges['faskes'][1]}","Faskes","Pilih kode faskes.")
    list_dv(ws,"C2:C501",f"=Referensi!$A${ranges['kategori'][0]}:$A${ranges['kategori'][1]}","Kategori","Opsional.")
    list_dv(ws,"E2:E501",f"=Referensi!$A${ranges['indikator'][0]}:$A${ranges['indikator'][1]}","Indikator","Pilih indikator dari master.")
    list_dv(ws,"F2:F501",'"DESA,KECAMATAN,KABUPATEN"',"Scope Type","Pilih level wilayah.")
    list_dv(ws,"G2:G501",f"=Referensi!$A${ranges['scope'][0]}:$A${ranges['scope'][1]}","Scope","Pilih kode wilayah.")
    if mode=="rasio":
        num_dv(ws,"H2:H501","Numerator","Isi angka >= 0.")
        num_dv(ws,"I2:I501","Denominator","Wajib > 0 untuk rasio.",positive=True)
    else:
        num_dv(ws,"H2:H501","Nilai","Isi angka >= 0; satuan mengikuti indikator.")

    red = PatternFill(start_color=LIGHT_RED,end_color=LIGHT_RED,fill_type="solid")
    orange = PatternFill(start_color=LIGHT_ORANGE,end_color=LIGHT_ORANGE,fill_type="solid")
    if mode=="rasio":
        ws.conditional_formatting.add("I2:I501",CellIsRule(operator="lessThanOrEqual",formula=["0"],fill=red))
        ws.conditional_formatting.add("H2:H501",FormulaRule(formula=['AND($E2<>"",H2="")'],fill=orange))
        ws.conditional_formatting.add("I2:I501",FormulaRule(formula=['AND($E2<>"",I2="")'],fill=orange))
    else:
        ws.conditional_formatting.add("H2:H501",CellIsRule(operator="lessThan",formula=["0"],fill=red))
        ws.conditional_formatting.add("H2:H501",FormulaRule(formula=['AND($E2<>"",H2="")'],fill=orange))

    ws.freeze_panes="A2"; ws.auto_filter.ref=f"A1:{get_column_letter(len(headers))}501"; ws.print_title_rows="1:1"
    ws.page_setup.orientation="landscape"; ws.page_setup.fitToWidth=1; ws.sheet_properties.pageSetUpPr.fitToPage=True
    protect(ws)


def write_petunjuk(ws):
    ws.sheet_view.showGridLines=False; ws.sheet_properties.tabColor=BLUE
    ws.column_dimensions["A"].width=3; ws.column_dimensions["B"].width=30; ws.column_dimensions["C"].width=112
    ws.merge_cells("B1:C1"); c=ws["B1"]; c.value="BNBA — SEMUA SATUAN & PERIODE"; c.font=TITLE; c.fill=fill(NAVY); c.alignment=Alignment(vertical="center"); ws.row_dimensions[1].height=36
    ws.merge_cells("B2:C2"); c=ws["B2"]; c.value="Satu template untuk data harian/bulanan/tahunan. Pilih sheet sesuai pola nilai: Rasio atau Nilai."; c.font=SUB; c.fill=fill(LIGHT_BLUE); c.alignment=Alignment(vertical="center",wrap_text=True); ws.row_dimensions[2].height=36
    ws.merge_cells("B3:C3"); c=ws["B3"]; c.value="STATUS: FINAL — rasio, count/scalar, dan periode bulanan/tahunan didukung. Format harian tersedia; perhitungan harian belum aktif."; c.font=Font(name=FONT,size=10,bold=True,color=RED); c.fill=fill(YELLOW); c.alignment=Alignment(vertical="center",wrap_text=True); ws.row_dimensions[3].height=36

    r=5
    section(ws,r,"1. Dua pola nilai"); r+=1
    for label,text in [
        ("Sheet Rasio","Untuk %, Persentase, (%), dan Rasio non-persen: numerator + denominator."),
        ("Sheet Nilai","Untuk Kasus, Dokumen, Laporan, Berita Acara, orang (count), UHH (tahun), AKIP (Poin), dan indeks/scalar: nilai."),
        ("Pilih satu","Satu baris hanya memakai satu sheet/pola. Jangan mengisi pola berbeda untuk baris yang sama."),
    ]:
        kv(ws,r,label,text,34); r+=1
    r+=1
    section(ws,r,"2. Satu kolom periode untuk tanggal/bulan/tahun"); r+=1
    for label,text in [
        ("Harian","periode = YYYY-MM-DD, contoh 2026-05-15."),
        ("Bulanan","periode = YYYY-MM, contoh 2026-05."),
        ("Tahunan","periode = YYYY, contoh 2026. Tidak perlu file template tahunan terpisah."),
        ("Aturan","Jangan menambahkan bulan pada data tahunan; satu baris = satu periode tahunan."),
    ]:
        kv(ws,r,label,text,32); r+=1
    r+=1
    section(ws,r,"3. Satuan"); r+=1
    for label,text in [
        ("Rasio","%, Persentase, (%) → persen; satuan Rasio → num/den tanpa ×100."),
        ("Count","Kasus, Dokumen, Laporan, Berita Acara → nilai, agregasi SUM."),
        ("Scalar","(tahun), Poin, indeks → nilai, LAST_VALUE."),
        ("orang","Bisa count orang atau numerator cakupan; master indikator menentukan."),
    ]:
        kv(ws,r,label,text,34); r+=1
    r+=1
    section(ws,r,"4. Cara isi"); r+=1
    for label,text in [
        ("Langkah 1","Pilih sheet Rasio atau Nilai sesuai indikator."),
        ("Langkah 2","Isi periode (tanggal/bulan/tahun), faskes, kode_iku, kode_indikator, scope_type, scope_code."),
        ("Langkah 3","Isi field nilai sesuai pola: numerator/denominator atau nilai."),
        ("Langkah 4","kode_iku boleh beberapa kode dengan ;. Label pertama = utama/konteks."),
        ("Langkah 5","Jangan mengisi total/subtotal, formula, atau merged cell."),
    ]:
        kv(ws,r,label,text,32); r+=1
    r+=1
    section(ws,r,"5. Aturan wajib"); r+=1
    for label,text in [
        ("Scope","Hanya DESA/KECAMATAN/KABUPATEN. Faskes adalah pengunggah, bukan scope."),
        ("Satu grain","Simpan level terendah; level atas diturunkan. Jangan tulis desa + kecamatan + kabupaten sekaligus."),
        ("Blank ≠ 0","Blank = belum ada data; 0 = nol yang dilaporkan."),
        ("Tanpa PII","Jangan masukkan nama, NIK, atau identitas individu."),
        ("Tanpa formula","Jangan mengubah header atau menambah formula."),
    ]:
        kv(ws,r,label,text,32); r+=1
    r+=1
    section(ws,r,"6. Kesiapan kode"); r+=1
    for label,text in [
        ("Rasio/count/scalar","Sudah didukung upload → kalkulasi → dashboard."),
        ("Periode bulanan/tahunan","Didukung; pilih periode yang sama saat upload."),
        ("Master","Tambahkan indikator dan tipe nilai lewat Master sebelum memakai kode baru."),
    ]:
        kv(ws,r,label,text,34); r+=1
    r+=1
    section(ws,r,"7. Catatan"); r+=1
    for label,text in [
        ("Contoh","Angka/kode CONTOH-* hanya ilustrasi; bukan data produksi."),
        ("Regenerate","Setelah master berubah, regenerate template agar Referensi ikut berubah."),
    ]:
        kv(ws,r,label,text,30); r+=1


def write_contoh(ws):
    ws.sheet_view.showGridLines=False; ws.sheet_properties.tabColor=GREEN
    ws.merge_cells("A1:J1"); c=ws["A1"]; c.value="CONTOH — ANGKA FIKTIF"; c.font=TITLE; c.fill=fill(GREEN); c.alignment=Alignment(vertical="center"); ws.row_dimensions[1].height=34
    ws.merge_cells("A2:J2"); c=ws["A2"]; c.value="Satu template untuk periode harian/bulanan/tahunan. Kode CONTOH-* bukan master."; c.font=SUB; c.fill=fill(LIGHT_GREEN); c.alignment=Alignment(vertical="center",wrap_text=True); ws.row_dimensions[2].height=30

    rasio_h=["periode","kode_faskes","kode_kategori","kode_iku","kode_indikator","scope_type","scope_code","numerator","denominator","catatan"]
    nilai_h=["periode","kode_faskes","kode_kategori","kode_iku","kode_indikator","scope_type","scope_code","nilai","catatan"]

    ws.merge_cells("A4:J4"); c=ws["A4"]; c.value="Contoh Rasio (bulanan)"; c.font=Font(name=FONT,size=11,bold=True,color=WHITE); c.fill=fill(BLUE); c.alignment=Alignment(vertical="center")
    hdr(ws,5,rasio_h,BLUE,30)
    rasio=[["2026-05","PKM-TALIWANG","IKU","IKU-2026-001.01","IND-001","DESA","5201.01.01",92,120,"ILUSTRASI — angka fiktif"]]
    for rr,vals in enumerate(rasio,6):
        for cc,v in enumerate(vals,1):
            cell_style(ws.cell(row=rr,column=cc,value=v),align="right" if cc in (8,9) else "left",color=LIGHT_GREEN)
    ws.merge_cells("A9:J9"); c=ws["A9"]; c.value="Contoh Nilai — count (bulanan)"; c.font=Font(name=FONT,size=11,bold=True,color=WHITE); c.fill=fill("7F7F7F"); c.alignment=Alignment(vertical="center")
    hdr(ws,10,nilai_h,"7F7F7F",30)
    nilai_count=[["2026-05","PKM-TALIWANG","IKU","CONTOH-IKU-KEMATIAN","CONTOH-KASUS-IBU","DESA","5201.01.01",2,"ILUSTRASI — kode contoh"]]
    for rr,vals in enumerate(nilai_count,11):
        for cc,v in enumerate(vals,1):
            cell_style(ws.cell(row=rr,column=cc,value=v),align="right" if cc==8 else "left",color=LIGHT_GREEN)
    ws.merge_cells("A14:J14"); c=ws["A14"]; c.value="Contoh Nilai — scalar (tahunan)"; c.font=Font(name=FONT,size=11,bold=True,color=WHITE); c.fill=fill(GREEN); c.alignment=Alignment(vertical="center")
    hdr(ws,15,nilai_h,GREEN,30)
    nilai=[["2026","PKM-TALIWANG","IKU","CONTOH-IKU-UHH","CONTOH-UHH","KABUPATEN","5201",73.5,"ILUSTRASI — kode contoh"]]
    for rr,vals in enumerate(nilai,16):
        for cc,v in enumerate(vals,1):
            cell_style(ws.cell(row=rr,column=cc,value=v),align="right" if cc==8 else "left",color=LIGHT_GREEN)
    ws.merge_cells("A19:J19"); c=ws["A19"]; c.value="Jangan mengunggah sheet Contoh. Gunakan template ini untuk upload; pastikan periode file sama dengan periode submission."; c.font=BODY; c.fill=fill(YELLOW); c.border=BORDER; c.alignment=Alignment(vertical="top",wrap_text=True); ws.row_dimensions[19].height=36


def main():
    master = load_master()
    categories, indicators, links, wilayah, faskes = master
    wb=Workbook()
    ws_pet=wb.active; ws_pet.title="Petunjuk"
    ws_rasio=wb.create_sheet("Rasio")
    ws_nilai=wb.create_sheet("Nilai")
    ws_ref=wb.create_sheet("Referensi")
    ws_contoh=wb.create_sheet("Contoh")
    ranges=write_ref(ws_ref,master)
    write_data(ws_rasio,"rasio",ranges,True,["periode","kode_faskes","kode_kategori","kode_iku","kode_indikator","scope_type","scope_code","numerator","denominator","catatan"])
    write_data(ws_nilai,"nilai",ranges,True,["periode","kode_faskes","kode_kategori","kode_iku","kode_indikator","scope_type","scope_code","nilai","catatan"])
    write_petunjuk(ws_pet)
    write_contoh(ws_contoh)
    wb.active=0
    OUT.mkdir(parents=True,exist_ok=True); STATIC.mkdir(parents=True,exist_ok=True)
    out=OUT/"BNBA-Semua-Satuan.xlsx"; wb.save(out); shutil.copy2(out,STATIC/out.name)
    print("wrote",out)


if __name__=="__main__":
    main()
