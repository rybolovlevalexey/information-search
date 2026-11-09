from __future__ import annotations

import re
from dataclasses import dataclass


@dataclass
class Heading:
    clause: str | None
    level: int | None
    title: str


@dataclass
class Paragraph:
    text: str


@dataclass
class SpecMeta:
    spec: str = ""
    version: str = ""
    spec_title: str = ""


# regex vars
# Заголовки Markdown в стиле ATX: ## Heading, ### Heading и т.д.
# Ищет символы #, за которыми следует текст (например, "## 5.3.5 Configuration")
ATX_RE = re.compile(r"^(#+)\s+(.*?)\s*#*\s*$")

# Подчёркивания заголовков в стиле Setext: === или --- (3+ символов)
# Используется для распознавания заголовков, оформленных подчёркиванием снизу
UNDERLINE_RE = re.compile(r"^(?:={3,}|-{3,})$")

# Разделители строк в сетчатых таблицах: +---+, +===+, +:::+, и т.д.
# Распознаёт горизонтальные линии-разделители в сетчатых Markdown-таблицах
GRID_SEP_RE = re.compile(r"^\+[-=:+]+\+\s*$")

# Границы простых таблиц: строки, состоящие только из тире, разделённых пробелами (например, "----  ----")
# Используется для распознавания и пропуска строк-границ простых таблиц
SIMPLE_BORDER_RE = re.compile(r"^\s*-{3,}(?:\s+-{3,})+\s*$")

# Маркер начала ASN.1-кода: -- ASN1START -- или \-- ASN1START --
ASN1_START_RE = re.compile(r"^\\?--\s*ASN1START\s*$")

# Маркер окончания ASN.1-кода: -- ASN1STOP -- или \-- ASN1STOP --
ASN1_STOP_RE = re.compile(r"^\\?--\s*ASN1STOP\s*$")

# Распознавание строк таблиц: непробельный символ, 2+ пробела, непробельный символ
# Идентифицирует строки таблиц по нескольким идущим подряд пробелам в пределах строки
TABLE_ROW_RE = re.compile(r"\S\s{2,}\S")  # несколько пробелов внутри строки = колонки

# Встроенные атрибуты в конце строки: {#id .Class}, {scope="foo"}
# Убирает атрибуты Markdown из заголовков (например, "{#sec-5.3.5 .Heading1}")
ATTR_RE = re.compile(r"\s*\{[^{}]*\}\s*$")  # {#scope .Heading1}

# Изображения Markdown: ![alt](url) и ![alt](url){attributes}
# Удаляет ссылки на изображения из текста, заменяя их на пустую строку
IMAGE_RE = re.compile(r"!\[[^\]]*\]\([^)]*\)(?:\{[^}]*\})?")

# Встроенные ссылки Markdown: [text](url)
# Извлекает текст ссылки, отбрасывая часть с URL
LINK_RE = re.compile(r"\[([^\]]*)\]\([^)]*\)")

# Встроенные span-элементы Markdown: [text]{.class}, [text]{#id}
# Извлекает текст из встроенных элементов уровня span с атрибутами
SPAN_RE = re.compile(r"\[([^\]]*)\]\{[^}]*\}")

# Примечания-ссылки: [11], [12, 13], [1, 2, 3]
# Удаляет числовые ссылки на литературы из текста (например, сноски/ссылки)
REF_RE = re.compile(r"\s*\[\d+(?:\s*,\s*\d+)*\]")  # [11]  [12, 13]

# Маркеры верхнего индекса: ^text^ -> text
# Извлекает содержимое верхнего индекса (например, ^TM^ -> TM)
SUP_RE = re.compile(r"\^([^^\s]+)\^")  # ^TM^

# Маркеры нижнего индекса: ~text~ -> text
# Извлекает содержимое нижнего индекса (например, ~1~ -> 1)
SUB_RE = re.compile(r"~([^~\s]+)~")

# Снятие экранирования спецсимволов (непробельных, не алфавитно-цифровых)
# Убирает обратную косую черту из экранированных последовательностей вида \[ -> [, \* -> *
UNESCAPE_RE = re.compile(r"\\([^\sA-Za-z0-9])")

# Маркеры нумерованных списков: 1>, 2>, 10> и т.д.
# Убирает числовые префиксы нумерованных списков из текста
LIST_MARK_RE = re.compile(r"^\d+>\s*")  # "1> ", "2> "

# Маркеры маркированных списков: - или •
# Убирает маркеры маркированных списков из текста
BULLET_RE = re.compile(r"^[-•]\s+")

# Нормализация пробелов: один или более символов пробела
# Заменяет все последовательности пробелов на один пробел
WS_RE = re.compile(r"\s+")

# Нумерованные заголовки: "5.3.5 Some Title", "A.2 Annex Title"
# Ищет последовательности номеров через точку, за которыми следует текст (номера разделов вида 5.3.5.1)
NUM_HEADING_RE = re.compile(r"^(\d+(?:\.\d+)*)\s+(\S.*)$")

# Подзаголовки приложений: "A.1 Subtitle", "B.2.3 Deep Subtitle"
# Ищет подразделы приложений с префиксом-буквой (например, "A.1 Description here")
ANNEX_SUB_RE = re.compile(r"^([A-Z](?:\.\d+)+)\s+(\S.*)$")

# Заголовки приложений: "Annex A (normative): Title", "Annex B (informative): Title"
# Ищет заголовки приложений верхнего уровня с необязательной пометкой нормативный/информативный
ANNEX_RE = re.compile(
    r"^Annex\s+([A-Z])\b\s*(?:\((?:normative|informative)\))?\s*:?\s*(.*)$", re.I
)

# Сведения о релизе: "(Release 17)", "(Release 16)" и т.д.
# Распознаёт маркеры версий релизов 3GPP в тексте
RELEASE_RE = re.compile(r"^\W*\(Release\s+\d+\)")

# Идентификация спецификации 3GPP: "3GPP TS 38.331 V17.6.0"
# Извлекает номер спецификации (например, 38.331) и версию (например, 17.6.0) из заголовков спецификаций
SPEC_RE = re.compile(r"3GPP\s+(?:TS|TR)\s+(\d{2}\.\d{3})\s+V(\d+\.\d+\.\d+)")


def clean_inline(s: str) -> str:
    s = s.replace("\u00a0", " ")
    # Маркеры цитирования: "> text" -> убирает префиксы цитат Markdown
    s = re.sub(r"^(?:>\s?)+", "", s)
    # Обратная косая черта в конце строки = перенос строки, заменяем на пробел
    s = re.sub(r"\\$", " ", s)
    # Protect reference brackets like \[11\] by temporarily replacing with control chars
    s = s.replace("\\[", "\x01").replace("\\]", "\x02")
    s = IMAGE_RE.sub("", s)
    s = LINK_RE.sub(r"\1", s)
    s = SPAN_RE.sub(r"\1", s)
    s = s.replace("\x01", "[").replace("\x02", "]")
    s = REF_RE.sub("", s)
    s = SUP_RE.sub(r"\1", s)
    s = SUB_RE.sub(r"\1", s)
    s = s.replace("*", "")
    s = UNESCAPE_RE.sub(r"\1", s)
    s = WS_RE.sub(" ", s).strip()
    s = LIST_MARK_RE.sub("", s)
    s = BULLET_RE.sub("", s)
    return s.strip()


def make_heading(text: str) -> Heading:
    t = clean_inline(ATTR_RE.sub("", text))
    t = re.sub(r"^--\s*", "", t)
    m = ANNEX_RE.match(t)
    if m:
        return Heading(m.group(1), 1, m.group(2).strip() or t)
    m = ANNEX_SUB_RE.match(t)
    if m:
        return Heading(m.group(1), m.group(1).count(".") + 1, m.group(2).strip())
    m = NUM_HEADING_RE.match(t)
    if m:
        return Heading(m.group(1), m.group(1).count(".") + 1, m.group(2).strip())
    return Heading(None, None, t)


def parse_meta(lines: list[str], filename: str = "") -> SpecMeta:
    head = " ".join(l.strip(" |") for l in lines[:80])
    meta = SpecMeta()
    m = SPEC_RE.search(head)
    if m:
        meta.spec, meta.version = m.group(1), m.group(2)
    else:
        fm = re.match(r"(\d{2})(\d{3})-([a-z0-9]{3})", filename)
        if fm:
            meta.spec, meta.version = f"{fm.group(1)}.{fm.group(2)}", fm.group(3)

    for j, line in enumerate(lines[:80]):
        if RELEASE_RE.match(line.strip(" |")):
            prev = [clean_inline(l.strip(" |")) for l in lines[:j]]
            prev = [p for p in prev if p and not p.startswith(("+-", "+="))]
            if prev:
                title = prev[-1]
                if (
                        len(prev) > 1
                        and prev[-2].endswith(";")
                        and not prev[-2].startswith(("Technical Specification", "3rd Generation"))
                ):
                    title = prev[-2] + " " + title
                meta.spec_title = title
            break
    return meta


def parse(text: str, filename: str = "") -> tuple[SpecMeta, list[Heading | Paragraph]]:
    lines = text.replace("\r\n", "\n").replace("\u00a0", " ").split("\n")
    meta = parse_meta(lines, filename)

    blocks: list[Heading | Paragraph] = []
    para: list[str] = []

    def flush() -> None:
        if para:
            t = WS_RE.sub(" ", " ".join(para)).strip()
            if re.search(r"\w", t):
                blocks.append(Paragraph(t))
            para.clear()

    def emit(t: str) -> None:
        if t and re.search(r"\w", t):
            blocks.append(Paragraph(t))

    def grid_row(rows: list[str]) -> None:
        cells = []
        for r in rows:
            r = r.strip()
            r = r[1:] if r.startswith("|") else r
            r = r[:-1] if r.endswith("|") else r
            cells.append(clean_inline(r.replace("|", " ")))
        emit(WS_RE.sub(" ", " ".join(c for c in cells if c)).strip())

    n = len(lines)
    i = 0
    in_asn1 = False
    while i < n:
        s = lines[i].rstrip()
        st = s.strip()

        if in_asn1:
            if ASN1_STOP_RE.match(st) or ATX_RE.match(s):
                in_asn1 = False
                if ATX_RE.match(s):
                    continue
            i += 1
            continue
        if ASN1_START_RE.match(st):
            flush()
            in_asn1 = True
            i += 1
            continue

        if GRID_SEP_RE.match(s):
            flush()
            rows: list[str] = []
            while i < n:
                t = lines[i].rstrip()
                if GRID_SEP_RE.match(t):
                    grid_row(rows)
                    rows = []
                elif t.startswith("|"):
                    rows.append(t)
                else:
                    break
                i += 1
            grid_row(rows)
            continue

        if not st:
            flush()
            i += 1
            continue

        m = ATX_RE.match(s)
        if m:
            flush()
            blocks.append(make_heading(m.group(2)))
            i += 1
            continue

        if (
                not para
                and i + 1 < n
                and UNDERLINE_RE.match(lines[i + 1].strip())
                and "|" not in s
                and not SIMPLE_BORDER_RE.match(s)
        ):
            blocks.append(make_heading(s))
            i += 2
            continue

        if SIMPLE_BORDER_RE.match(s) or UNDERLINE_RE.match(st):
            flush()
            i += 1
            continue

        cleaned = clean_inline(s)
        if TABLE_ROW_RE.search(st):
            flush()
            emit(cleaned)
        elif cleaned:
            para.append(cleaned)
        i += 1

    flush()
    return meta, blocks
