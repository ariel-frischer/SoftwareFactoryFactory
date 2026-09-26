# /// script
# dependencies = ["numpy"]
# ///
"""Synthesize the soundtrack: A-minor bed at 120 bpm plus cue-synced SFX.

Usage: uv run mix.py ../work/cues.json ../work/audio.wav
"""
import json
import sys
import wave

import numpy as np

SR, DUR, BEAT = 48000, 26.0, 0.5
N = int(SR * DUR)
t_all = np.arange(N) / SR
rng = np.random.default_rng(7)


def hz(midi):
    return 440.0 * 2 ** ((midi - 69) / 12)


def env(n, a=0.01, r=0.2):
    e = np.ones(n)
    na, nr = int(a * SR), int(r * SR)
    e[:na] = np.linspace(0, 1, na)
    if nr:
        e[-nr:] *= np.linspace(1, 0, nr)
    return e


def lowpass(x, cutoff):
    a = np.broadcast_to(np.exp(-2 * np.pi * np.asarray(cutoff, float) / SR), x.shape)
    y = np.empty_like(x)
    acc = 0.0
    for i, v in enumerate(x):
        acc = (1 - a[i]) * v + a[i] * acc
        y[i] = acc
    return y


def add(buf, start, sig, gain=1.0):
    i = int(start * SR)
    if i >= N:
        return
    sig = sig[: N - i]
    buf[i : i + len(sig)] += sig * gain


music = np.zeros(N)
sfx = np.zeros(N)

# Chords: Am, F, C, G (one bar = 2s)
chords = [[57, 60, 64], [53, 57, 60], [48, 52, 55, 60], [55, 59, 62]]
roots = [45, 41, 48, 43]


def pad_gain(t):
    return 0.6 if 12.5 <= t < 15.3 else 1.0


for bar in range(13):
    start = bar * 2.0
    ch = chords[bar % 4]
    n = int(2.0 * SR)
    tt = np.arange(n) / SR
    sig = np.zeros(n)
    for m in ch:
        f = hz(m)
        for det in (-0.12, 0.12):
            sig += np.sign(np.sin(2 * np.pi * f * (1 + det / 100) * tt)) * 0.5 + np.sin(2 * np.pi * f * tt)
    sig = lowpass(sig, 1100) * env(n, 0.08, 0.25) * 0.035
    add(music, start, sig, pad_gain(start))

# Bass on beats 1 and 3, drums, hats — shaped by section.
def section_drums(t):
    if 5.0 <= t < 6.0 or 12.5 <= t < 15.3 or 16.0 <= t < 17.3 or t >= 21.8:
        return 0.0
    if 18.3 <= t < 21.8:
        return 0.5
    return 1.0


for b in range(int(DUR / BEAT)):
    t = b * BEAT
    g = section_drums(t)
    if g == 0:
        continue
    # kick
    n = int(0.25 * SR)
    tt = np.arange(n) / SR
    freq = 50 + 90 * np.exp(-tt * 30)
    kick = np.sin(2 * np.pi * np.cumsum(freq) / SR) * np.exp(-tt * 14)
    add(music, t, kick, 0.32 * g)
    # bass
    if b % 2 == 0:
        root = hz(roots[int(t // 2.0) % 4] - 12)
        n = int(0.45 * SR)
        tt = np.arange(n) / SR
        bass = (np.sin(2 * np.pi * root * tt) + 0.3 * np.sin(4 * np.pi * root * tt)) * env(n, 0.005, 0.2)
        add(music, t, bass, 0.2 * g)
    # hats: 8ths, 16ths during factorial ramp
    sub = 4 if 9.0 <= t < 12.3 else 2
    for k in range(sub):
        n = int(0.04 * SR)
        hat = rng.standard_normal(n)
        hat = hat - lowpass(hat, 6000)
        hat *= np.exp(-np.arange(n) / SR * 90)
        add(music, t + k * BEAT / sub, hat, 0.05 * g * (1.0 if k % 2 == 0 else 0.6))

# Riser into the Droste dive (5.0–6.0)
n = int(1.0 * SR)
tt = np.arange(n) / SR
noise = rng.standard_normal(n)
riser = (noise - lowpass(noise, 800 + 5000 * tt)) * (tt ** 2) * 0.12
riser += np.sin(2 * np.pi * np.cumsum(220 * 2 ** (tt * 2)) / SR) * (tt ** 2) * 0.05
add(sfx, 5.0, riser)

# Title boom at 6.0 + Am hit
n = int(2.0 * SR)
tt = np.arange(n) / SR
boom = np.sin(2 * np.pi * np.cumsum(40 + 60 * np.exp(-tt * 8)) / SR) * np.exp(-tt * 2.5)
hit = sum(np.sin(2 * np.pi * hz(m) * tt) for m in (57, 60, 64, 69)) * np.exp(-tt * 1.6) * 0.25
add(sfx, 6.0, boom * 0.45 + hit * 0.25)

cues = json.load(open(sys.argv[1]))

# Factory stamps: short tuned thunk (A2) with a tick.
for i, t in enumerate(cues["stamps"]):
    n = int(0.12 * SR)
    tt = np.arange(n) / SR
    thunk = np.sin(2 * np.pi * hz(45 + (7 if i % 2 else 0)) * tt) * np.exp(-tt * 40)
    tick = rng.standard_normal(n) * np.exp(-tt * 300) * 0.3
    add(sfx, t, (thunk + tick) * 0.22)

# Konami keys
for t in cues["keys"]:
    n = int(0.03 * SR)
    tt = np.arange(n) / SR
    click = rng.standard_normal(n) * np.exp(-tt * 250)
    add(sfx, t, lowpass(click, 3500) * 0.35)

# Alarm 16.0–17.3: two-tone A5/E5, softened square
for k in range(5):
    t = 16.0 + k * 0.26
    n = int(0.24 * SR)
    tt = np.arange(n) / SR
    f = hz(81 if k % 2 == 0 else 76)
    tone = lowpass(np.sign(np.sin(2 * np.pi * f * tt)), 2500) * env(n, 0.005, 0.05)
    add(sfx, t, tone, 0.06)

# Relief chord at 17.3 and final chord at 21.8
for start, length, gain in ((17.3, 1.5, 0.12), (21.8, 4.2, 0.2)):
    n = int(length * SR)
    tt = np.arange(n) / SR
    chord = sum(np.sin(2 * np.pi * hz(m) * tt) + 0.2 * np.sin(4 * np.pi * hz(m) * tt) for m in (45, 57, 60, 64, 69))
    add(sfx, start, chord * np.exp(-tt * (1.2 if length < 2 else 0.7)) * env(n, 0.01, 0.8), gain)

# Duck music under the alarm and fade out at the end.
duck = np.ones(N)
duck[(t_all >= 15.9) & (t_all < 17.3)] = 0.45
duck *= np.clip((DUR - t_all) / 1.5, 0, 1)
mix = music * duck + sfx
mix /= np.max(np.abs(mix)) / 0.89
stereo = np.stack([mix, mix], axis=1)
with wave.open(sys.argv[2], "wb") as w:
    w.setnchannels(2)
    w.setsampwidth(2)
    w.setframerate(SR)
    w.writeframes((stereo * 32767).astype("<i2").tobytes())
print("wrote", sys.argv[2], f"{DUR}s")
