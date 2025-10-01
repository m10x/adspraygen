#!/usr/bin/env python3
"""
Simple unique placeholders generator.

- [WORD], [NUMBER], [SPECIAL] placeholders
- Values of the same type in the same pattern cannot repeat
- Multiple patterns supported
"""

import itertools
import argparse
import re
import sys
from typing import List

NUMBER_TOKENS = ["{YY}", "{YYYY}", "1", "2", "3", "12", "123"]
SPECIAL_TOKENS = ["!", ".", "#", "-", "_"]
NOUN_TOKENS = ["{monthsGerman}", "{seasonsGerman}", "{sn}", "{givenName}"]

PLACEHOLDER_REGEX = re.compile(r"(\[WORD\]|\[NUMBER\]|\[SPECIAL\])")

def read_lines_file(path: str) -> List[str]:
    try:
        with open(path, encoding="utf-8") as f:
            return [line.strip() for line in f if line.strip()]
    except FileNotFoundError:
        print(f"File not found: {path}", file=sys.stderr)
        return []

def build_word_list(nouns: List[str]) -> List[str]:
    combined = list(nouns) + NOUN_TOKENS
    seen = set()
    out = []
    for x in combined:
        if x not in seen:
            seen.add(x)
            out.append(x)
    return out

def expand_pattern(pattern: str):
    tokens = []
    last = 0
    for m in PLACEHOLDER_REGEX.finditer(pattern):
        start, end = m.span()
        if start > last:
            tokens.append(("TEXT", pattern[last:start]))
        tokens.append(("PH", m.group(1)))
        last = end
    if last < len(pattern):
        tokens.append(("TEXT", pattern[last:]))
    return tokens

def generate_combinations(tokens, word_list):
    """
    Generate all combinations.
    Skip any combination where the same placeholder type repeats a value.
    """
    # Build value lists
    value_lists = []
    placeholder_types = []
    for kind, val in tokens:
        if kind == "TEXT":
            value_lists.append([val])
            placeholder_types.append(None)
        elif val == "[WORD]":
            value_lists.append(word_list)
            placeholder_types.append("WORD")
        elif val == "[NUMBER]":
            value_lists.append(NUMBER_TOKENS)
            placeholder_types.append("NUMBER")
        elif val == "[SPECIAL]":
            value_lists.append(SPECIAL_TOKENS)
            placeholder_types.append("SPECIAL")
        else:
            value_lists.append([val])
            placeholder_types.append(None)

    for combo in itertools.product(*value_lists):
        valid = True
        seen = {"WORD": set(), "NUMBER": set(), "SPECIAL": set()}
        for typ, val in zip(placeholder_types, combo):
            if typ:
                if val in seen[typ]:
                    valid = False
                    break
                seen[typ].add(val)
        if valid:
            yield "".join(combo)

def main():
    parser = argparse.ArgumentParser(description="Generate strings with unique placeholders in pattern.")
    group = parser.add_mutually_exclusive_group(required=True)
    group.add_argument("--pattern", type=str)
    group.add_argument("--patterns-file", type=str)
    parser.add_argument("--nouns", type=str, required=True)
    parser.add_argument("--out", type=str, default=None)
    parser.add_argument("--limit", type=int, default=None)
    args = parser.parse_args()

    nouns = read_lines_file(args.nouns)
    word_list = build_word_list(nouns)

    patterns = [args.pattern] if args.pattern else read_lines_file(args.patterns_file)

    out_f = open(args.out, "w", encoding="utf-8") if args.out else None
    produced = 0

    for patt in patterns:
        tokens = expand_pattern(patt)
        for combo in generate_combinations(tokens, word_list):
            if args.limit and produced >= args.limit:
                if out_f:
                    out_f.close()
                return
            if out_f:
                out_f.write(combo + "\n")
            else:
                print(combo)
            produced += 1

    if out_f:
        out_f.close()

if __name__ == "__main__":
    main()
