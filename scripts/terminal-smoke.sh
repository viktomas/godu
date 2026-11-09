#!/bin/sh
# End-to-end terminal smoke test for godu.
#
# Runs the built binary in a real terminal (a tmux pane) against a fixture
# tree and asserts on what actually lands on the screen. This is the only
# layer that can catch regressions the headless render tests cannot: the
# blank-screen hang (draw loop spinning at 100% CPU) and terminal-emission
# bugs in tcell itself, e.g. the v2.8+ wide-rune corruption that made
# "日本語フォルダ" render as "日語ォダ" or disappear entirely.
#
# Usage: ./scripts/terminal-smoke.sh /path/to/godu-binary
# Requires: tmux, coreutils.

set -eu

BIN=$(cd "$(dirname "$1")" && pwd)/$(basename "$1")
FIXTURE=$(mktemp -d)
CAPTURE=$(mktemp)
OUT=$(mktemp)
trap 'tmux kill-session -t godu-smoke 2>/dev/null || true; rm -rf "$FIXTURE" "$CAPTURE" "$OUT"' EXIT

# Fixture: sizes, nested CJK names, a filename containing a newline.
mkdir -p "$FIXTURE/日本語フォルダ"
head -c 5242880 /dev/urandom > "$FIXTURE/big.bin"
head -c 2048 /dev/urandom > "$FIXTURE/日本語フォルダ/中文.txt"
printf 'x' > "$FIXTURE/$(printf 'new\nline.txt')"

# 1. Rendering: tree, wide runes intact (not dropped, not every-other-rune),
#    runes after the newline survive, preview pane shows the selected folder.
tmux new-session -d -s godu-smoke -x 120 -y 40 "sh -c '$BIN -l 0 \"$FIXTURE\"'"
sleep 3
tmux send-keys -t godu-smoke j # select the CJK folder -> preview pane shows 中文.txt
sleep 1
tmux capture-pane -t godu-smoke -p > "$CAPTURE"
tmux kill-session -t godu-smoke

fail=0
check() {
    if grep -q "$1" "$CAPTURE"; then
        echo "ok: $2"
    else
        echo "FAIL: $2"
        fail=1
    fi
}
check '5M big.bin' 'folder tree with sizes rendered (not a blank screen)'
check '日本語フォルダ/' 'wide-rune name rendered in full (no dropped or overwritten glyphs)'
check '中文.txt' 'preview pane rendered'
check 'line.txt' 'runes after a newline in a filename rendered'

# 2. Interaction: mark a file, quit, and check it is printed on stdout.
tmux new-session -d -s godu-smoke -x 120 -y 40 "sh -c '$BIN -l 0 \"$FIXTURE\" > \"$OUT\"'"
sleep 3
tmux send-keys -t godu-smoke Space # mark big.bin
sleep 1
tmux send-keys -t godu-smoke q
sleep 1
tmux kill-session -t godu-smoke 2>/dev/null || true

if grep -q "big.bin" "$OUT"; then
    echo "ok: marked file printed on quit"
else
    echo "FAIL: marked file not printed on quit"
    fail=1
fi

if [ "$fail" -ne 0 ]; then
    echo "--- captured screen ---"
    cat "$CAPTURE"
    exit 1
fi
echo "terminal smoke test passed"
