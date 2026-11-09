from __future__ import annotations

import re
from dataclasses import dataclass, field

from extract import Heading, SpecMeta

EXCLUDED_TITLE_RE = re.compile(r"change history|^foreword$", re.I)


@dataclass
class Section:
    clause: str
    level: int
    title: str
    path: list[str]
    numbered: bool
    paras: list[str] = field(default_factory=list)


def _words(paras: list[str]) -> int:
    return sum(len(p.split()) for p in paras)


def _slug(s: str) -> str:
    return re.sub(r"[^A-Za-z0-9_.-]+", "-", s).strip("-")[:60]


def cut_preamble(blocks: list) -> list:
    for idx, b in enumerate(blocks):
        if isinstance(b, Heading) and (
                b.clause == "1" or (b.clause is None and b.title.lower() == "scope")
        ):
            return blocks[idx:]
    return blocks


def build_sections(blocks: list) -> list[Section]:
    sections: list[Section] = []
    stack: list[tuple[int, str]] = []
    cur: Section | None = None
    last_clause, last_level = "", 0
    skip_level: int | None = None

    for b in blocks:
        if isinstance(b, Heading):
            if b.clause:
                level, clause = b.level, b.clause
                last_clause, last_level = clause, level
                label = f"{clause} {b.title}"
            else:
                level, clause = last_level + 1, last_clause
                label = b.title
            while stack and stack[-1][0] >= level:
                stack.pop()
            stack.append((level, label))

            if skip_level is not None and level <= skip_level:
                skip_level = None
            if skip_level is None and EXCLUDED_TITLE_RE.search(b.title):
                skip_level = level

            cur = None
            if skip_level is None:
                cur = Section(clause, level, b.title, [l for _, l in stack], bool(b.clause))
                sections.append(cur)
        elif cur is not None:
            cur.paras.append(b.text)
    return sections


def split_paras(paras: list[str], max_words: int, min_words: int) -> list[str]:
    pieces: list[list[str]] = []
    cur: list[str] = []
    cw = 0
    for p in paras:
        w = len(p.split())
        if w > max_words:
            if cur:
                pieces.append(cur)
                cur, cw = [], 0
            ws = p.split()
            for k in range(0, len(ws), max_words):
                pieces.append([" ".join(ws[k: k + max_words])])
            continue
        if cur and cw + w > max_words:
            pieces.append(cur)
            cur, cw = [], 0
        cur.append(p)
        cw += w
    if cur:
        pieces.append(cur)
    if len(pieces) > 2 and _words(pieces[-1]) < min_words:
        pieces[-2] += pieces.pop()
    return ["\n".join(p) for p in pieces]


def build_chunks(
        meta: SpecMeta, blocks: list, max_words: int = 350, min_words: int = 20
) -> list[dict]:
    sections = build_sections(cut_preamble(blocks))

    ready: list[Section] = []
    pending: Section | None = None
    for sec in sections:
        if not sec.paras:
            continue
        if pending is not None:
            sec.paras = [pending.title] + pending.paras + sec.paras
            pending = None
        if _words(sec.paras) < min_words:
            pending = sec
        else:
            ready.append(sec)
    if pending is not None:
        ready.append(pending)

    chunks: list[dict] = []
    seen: dict[str, int] = {}
    for sec in ready:
        base = f"{meta.spec}@{meta.version}#{sec.clause}"
        if not sec.numbered:
            base += f":{_slug(sec.title)}"
        for part, text in enumerate(split_paras(sec.paras, max_words, min_words)):
            cid = f"{base}/{part}"
            if cid in seen:
                seen[cid] += 1
                cid = f"{cid}~{seen[cid]}"
            else:
                seen[cid] = 0
            chunks.append(
                {
                    "id": cid,
                    "spec": meta.spec,
                    "spec_title": meta.spec_title,
                    "version": meta.version,
                    "clause": sec.clause,
                    "path": sec.path,
                    "title": sec.title,
                    "text": text,
                    "part": part,
                }
            )
    return chunks
