# python main.py --input ../data/3gpp-specs --output ../data/corpus.jsonl
from __future__ import annotations

import argparse
import json
import re
import statistics
from pathlib import Path

from chunk import build_chunks
from extract import parse

NAME_RE = re.compile(r"^(\d{5})-([a-z0-9]{3})\.md$")


def pick_latest(files: list[Path]) -> list[Path]:
    best: dict[str, Path] = {}
    for f in files:
        m = NAME_RE.match(f.name)
        key = m.group(1) if m else f.name
        if key in best:
            print(f"[warn] несколько версий {key}: {best[key].name}, {f.name}")
            if f.name > best[key].name:
                best[key] = f
        else:
            best[key] = f
    return sorted(best.values())


def main() -> None:
    ap = argparse.ArgumentParser()
    ap.add_argument("--input", required=True, type=Path)
    ap.add_argument("--output", required=True, type=Path)
    ap.add_argument("--max-words", type=int, default=350)
    ap.add_argument("--min-words", type=int, default=20)
    ap.add_argument("--only", default="", help="подстрока в имени файла, например 38331")
    args = ap.parse_args()

    files = sorted(args.input.rglob("*.md"))
    if args.only:
        files = [f for f in files if args.only in f.name]
    files = pick_latest(files)

    args.output.parent.mkdir(parents=True, exist_ok=True)
    sizes: list[int] = []
    empty: list[str] = []
    with args.output.open("w", encoding="utf-8") as out:
        for f in files:
            meta, blocks = parse(f.read_text(encoding="utf-8", errors="replace"), f.name)
            chunks = build_chunks(meta, blocks, args.max_words, args.min_words)
            if not chunks:
                empty.append(f.name)
            for c in chunks:
                out.write(json.dumps(c, ensure_ascii=False) + "\n")
                sizes.append(len(c["text"].split()))

    print(f"спецификаций: {len(files)}, кусков: {len(sizes)}, слов всего: {sum(sizes)}")
    if sizes:
        q = statistics.quantiles(sizes, n=20)
        print(f"слов в куске: median={statistics.median(sizes)}, p95={q[18]:.0f}, max={max(sizes)}")
    if empty:
        print(f"[warn] без кусков ({len(empty)}): {', '.join(empty[:10])}")


if __name__ == "__main__":
    main()
