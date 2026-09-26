# /brag-slim plan — SoftwareFactoryFactory

- **What:** a software factory that builds software factories. Software shipped: 0.
- **For:** people drowning in "software factory" launches on X.
- **Hook:** the real factory already stamping factories, caption "Everyone built a software factory."
- **Angle / tone:** yc-parody played straight, with one chaotic beat (Konami alarm).
- **Identity:** the product's own page: #0d0f14 bg, #d2a8ff title, #79c0ff scene, #7ee787 stats, #e3b341 fun facts, #ff7b72 alarm. JetBrains Mono.
- **Real UI:** the exact scene/log/fun-fact renderer from `index.html`, driven deterministically by time.
- **Share caption:** "Everyone built a software factory. I built a software factory factory. Software shipped: 0."

## Storyboard (1920×1080, 30fps, 26s)

| Time | Scene | On screen | Sound |
|---|---|---|---|
| 0–3.2 | Hook | Factory running mid-shift; caption "Everyone built a software factory." | bed starts, soft stamps |
| 3.2–5.0 | Turn | Caption "We built the factory that builds them." | |
| 5.0–6.0 | Droste zoom | Camera dives into a mini `FxF` factory | riser whoosh |
| 6.0–8.2 | Title | "SoftwareFactoryFactory" / "we ship factories, not software" (hold) | boom |
| 8.2–12.5 | Factorial | Time accelerates, punch-in on counter exploding; caption "Growth: factorial." | stamps accelerate |
| 12.5–15.3 | Deadpan | Punch-in on "Software shipped: 0"; caption "Software shipped: 0. By design." | music drops to pad |
| 15.3–18.3 | Incident | Konami keycaps, red "Software shipped: 1", SOFTWARE DETECTED, shake → "rolled back" | keys, alarm, relief |
| 18.3–21.8 | Fun fact | Big typed fun fact (factorials line) | |
| 21.8–26 | Close | Name, tagline, `go install …`, hold | final chord |

## Rebuild

Renders (`brag.mp4`, `brag.jpg`, `work/`) are gitignored. Requires JetBrains Mono NL installed, headless Chromium, `uv`, `ffmpeg`.

1. Load `composition/index.html` at 1920×1080, await `window.ready`, save `cues()` to `work/cues.json`, and screenshot `renderAt(f/30)` for f = 0..779 into `work/frames/%05d.jpg`.
2. `cd composition && uv run mix.py ../work/cues.json ../work/audio.wav`
3. `ffmpeg -framerate 30 -i work/frames/%05d.jpg -i work/audio.wav -c:v libx264 -pix_fmt yuv420p -crf 18 -c:a aac -b:a 192k -shortest work/brag-raw.mp4`
4. Poster = frame 228 (title card) → `brag.jpg`; bake as frame 0 via `overlay=enable='eq(n\,0)'`.
