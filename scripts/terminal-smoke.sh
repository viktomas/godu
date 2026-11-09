#!/bin/sh
# End-to-end terminal smoke test for godu.
#
# Runs the built binary in a real terminal (a tmux pane) against a fixture
# tree and asserts on what actually lands on the screen. This is the only
# layer that can catch regressions the headless render tests cannot: the
# blank-screen hang (draw loop spinning at 100% CPU) and terminal-emission
# bugs in tcell itself, which look fine through tcell's SimulationScreen —
# e.g. tcell v2.8/v2.9 dropped wide-rune names entirely on real terminals
# and v2.10+ clobbered every second wide rune ("日本語フォルダ" rendering as
# "日語ォダ").
#
# Usage: ./scripts/terminal-smoke.sh /path/to/godu-binary
# Requires: tmux, coreutils. Run via `make smoke` or directly.

set -eu

if [ "$#" -ne 1 ]; then
    echo "usage: $0 /path/to/godu-binary" >&2
    exit 2
fi

# tcell picks its charset from the locale; force a UTF-8 one so CJK names
# are actually drawn (with LC_ALL=C they become '?').
LC_ALL=C.UTF-8
export LC_ALL

BIN=$(cd "$(dirname "$1")" && pwd)/$(basename "$1")
FIXTURE=$(mktemp -d)
CAPTURE=$(mktemp)
OUT=$(mktemp)
ERR=$(mktemp)
SERVER="godu-smoke-$$"

# Private tmux server with no user config, so the test is isolated from
# ~/.tmux.conf, the default shell and any other session.
cleanup() {
    tmux -L "$SERVER" kill-server 2>/dev/null || true
    rm -rf "$FIXTURE" "$CAPTURE" "$OUT" "$ERR"
}
trap cleanup EXIT INT TERM

# Fixture: sizes, nested CJK names, a filename containing a newline.
mkdir -p "$FIXTURE/日本語フォルダ"
head -c 5242880 /dev/urandom > "$FIXTURE/big.bin"
head -c 2048 /dev/urandom > "$FIXTURE/日本語フォルダ/中文.txt"
printf 'x' > "$FIXTURE/$(printf 'new\nline.txt')"

# Wait until the screen contains a pattern (polling beats fixed sleeps on
# slow CI runners), up to a limit. Fails with the captured screen and godu's
# stderr if the session died or the pattern never appears.
wait_for() {
    pattern=$1
    description=$2
    command=$3 # key to send before waiting, or '' for none
    i=0
    while [ $i -lt 40 ]; do
        if [ -n "$command" ]; then
            tmux -L "$SERVER" send-keys -t godu-smoke "$command"
            command=''
        fi
        if ! tmux -L "$SERVER" has-session -t godu-smoke 2>/dev/null; then
            echo "FAIL: $description — godu exited early; stderr:"
            cat "$ERR"
            return 1
        fi
        tmux -L "$SERVER" capture-pane -t godu-smoke -p > "$CAPTURE" 2>/dev/null || true
        if grep -qF "$pattern" "$CAPTURE"; then
            echo "ok: $description"
            return 0
        fi
        i=$((i + 1))
        sleep 0.5
    done
    echo "FAIL: $description; captured screen:"
    cat "$CAPTURE"
    echo "--- godu stderr ---"
    cat "$ERR"
    return 1
}

start_godu() {
    # Pass paths via the environment (tmux -e) to survive quoting; godu's
    # stderr goes to a file so a crash is visible in failures.
    tmux -L "$SERVER" new-session -d -s godu-smoke -x 120 -y 40 \
        -e "GODU_BIN=$BIN" -e "GODU_FIXTURE=$FIXTURE" -e "GODU_OUT=$OUT" -e "GODU_ERR=$ERR" \
        'sh -c '\''"$GODU_BIN" -l 0 "$GODU_FIXTURE" > "$GODU_OUT" 2> "$GODU_ERR"'\'''
}

fail=0

# 1. Rendering: the tree must appear (not a blank screen), wide runes must
#    survive in full (not dropped or every-second-rune clobbered), runes
#    after the newline must survive, and the preview pane must render.
start_godu
wait_for 'big.bin' 'folder tree rendered (not a blank screen)' '' || fail=1
wait_for '日本語フォルダ/' 'wide-rune name rendered in full' '' || fail=1
wait_for 'line.txt' 'runes after a newline in a filename rendered' '' || fail=1
# Select the CJK folder so its contents show in the preview pane.
wait_for '中文.txt' 'preview pane rendered' 'j' || fail=1
tmux -L "$SERVER" kill-session -t godu-smoke 2>/dev/null || true

# 2. Interaction: mark a file, quit, and check it is printed on stdout.
start_godu
wait_for 'big.bin' 'folder tree rendered' '' || fail=1
tmux -L "$SERVER" send-keys -t godu-smoke Space # mark big.bin
sleep 1
tmux -L "$SERVER" send-keys -t godu-smoke q

# godu must quit on its own (a hang here means the render loop is stuck).
i=0
while [ $i -lt 20 ] && tmux -L "$SERVER" has-session -t godu-smoke 2>/dev/null; do
    i=$((i + 1))
    sleep 0.5
done
if tmux -L "$SERVER" has-session -t godu-smoke 2>/dev/null; then
    echo "FAIL: godu did not quit after 'q' (hung render loop?); captured screen:"
    tmux -L "$SERVER" capture-pane -t godu-smoke -p
    echo "--- godu stderr ---"
    cat "$ERR"
    fail=1
else
    echo "ok: godu quit on 'q'"
fi

if grep -qF 'big.bin' "$OUT"; then
    echo "ok: marked file printed on quit"
else
    echo "FAIL: marked file not printed on quit; stdout:"
    cat "$OUT"
    fail=1
fi

if [ "$fail" -ne 0 ]; then
    exit 1
fi
echo "terminal smoke test passed"
