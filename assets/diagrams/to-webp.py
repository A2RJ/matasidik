#!/usr/bin/env python3
"""Convert a rendered PNG to lossless WebP for the diagram/docs build.

Usage:
  python3 to-webp.py INPUT.png OUTPUT.webp

Lossless WebP is used because diagrams and UI screenshots contain small text.
The caller is responsible for deleting the temporary PNG.
"""
import sys
from pathlib import Path

from PIL import Image


def main() -> None:
    if len(sys.argv) != 3:
        raise SystemExit(__doc__)
    src = Path(sys.argv[1])
    dst = Path(sys.argv[2])
    if not src.exists():
        raise SystemExit(f"input tidak ditemukan: {src}")
    with Image.open(src) as image:
        image.convert("RGBA").save(dst, "WEBP", lossless=True, method=6)
    before = src.stat().st_size
    after = dst.stat().st_size
    print(f"webp {src.name} -> {dst.name} ({before // 1024} KB -> {after // 1024} KB)")


if __name__ == "__main__":
    main()
