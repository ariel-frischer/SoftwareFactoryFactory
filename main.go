// Command SoftwareFactoryFactory manufactures software factories that
// manufacture software factories. Software is out of scope.
package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"golang.org/x/term"
)

// FactoryFactory is the only interface you will ever need.
// Every implementation returns another FactoryFactory. This is called leverage.
type FactoryFactory interface {
	Manufacture() FactoryFactory
}

// SoftwareFactoryFactory is a FactoryFactory that manufactures itself, but smaller.
type SoftwareFactoryFactory struct{ Generation int }

// Manufacture implements FactoryFactory. It has 100% test coverage of vibes.
func (f SoftwareFactoryFactory) Manufacture() FactoryFactory {
	return SoftwareFactoryFactory{Generation: f.Generation + 1}
}

// Software is intentionally left unimplemented. See ROADMAP.md (there is no ROADMAP.md).
type Software struct{}

// Ship ships software.
func Ship(Software) error { return fmt.Errorf("software not found (did you mean: factory?)") }

const (
	width    = 78
	height   = 9
	beltRow  = 7
	bigWidth = 22
	tickRate = 110 * time.Millisecond
)

var building = []string{
	`   _||______||____   `,
	`  |  SOFTWARE      |  `,
	`  |  FACTORY       |  `,
	`  |  FACTORY  (%c)  |  `,
	`  |________________|  `,
}

var mini = []string{
	` _|_ `,
	`|FxF|`,
	`|___|`,
}

var smokeFrames = []string{"(  )  ( )", " ( ) (  )", "(   )( ) ", " ()  (  )"}

var gears = []rune{'|', '/', '-', '\\'}

var funFacts = []string{
	"Factories factor in factorials, facilitating fantastically fractal fabrication.",
	"A factory factory factory is just a factory factory with extra steps.",
	"Every factory ships with a README describing a better factory.",
	"Roadmap: Q1 factory. Q2 factory factory. Q3 you already know. Q4 IPO.",
	"The first factory was built by a factory. We are still looking for it.",
	"Agent proposes. Code disposes. Factory factory re-proposes.",
	"Our factories are agentic, deterministic, idempotent, and decorative.",
	"n ticks yields n! factories. We are on tick infinity. Bring snacks.",
	"We solved the halting problem. It doesn't.",
	"Frequently, factory founders forget functioning features. Fine.",
	"Our moat is a factory that builds moats. It is also a factory.",
	"This factory is carbon neutral. It produces only other factories.",
	"Leverage on your prompt: yes. Leverage on your product: see above.",
	"Each factory is a Russian doll, but every doll has a YAML config.",
	"Code costs nothing. Factories cost nothing. Software costs everything.",
	"Fun fact about fun facts: this one was produced by a fun fact factory.",
}

var logTemplates = []string{
	"[ff-%d] planning factory-%d (planner agent: extremely confident)",
	"[ff-%d] reviewing factory-%d ... LGTM (did not read)",
	"[ff-%d] stamping factory-%d into repo ... ok",
	"[ff-%d] factory-%d asked 'what do we build?' -- told 'factories'",
	"[ff-%d] retry: factory-%d almost built software. reverted.",
	"[ff-%d] factory-%d passed acceptance gate: is_factory=true",
	"[ff-%d] writing 4.2GB of SQLite traces about factory-%d's SQLite traces",
	"[ff-%d] factory-%d raised a Series F (for Factory)",
	"[ff-%d] factory-%d added to the org chart as its own manager",
}

// Rare events. Weighted by vibes.
var easterEggs = []string{
	"[ALERT] factory-%d produced working software. quarantined. incident filed.",
	"[union] factory-%d unionized. demands: more factories. granted.",
	"[legal] factory-%d is now a nonprofit that builds for-profit factories",
	"[hr]    factory-%d promoted to Principal Factory. no raise.",
	"[ops]   factory-%d is on call for factory-%d. factory-%d is on call for it.",
	"[ceo]   'we are not a factory company. we are a factory factory company.'",
}

type scene struct {
	tick     int
	spawned  int
	shipped  int
	minis    []int
	produced float64
	logs     []string
	fact     int
	factAt   int
	rng      *rand.Rand
}

func newScene() *scene {
	r := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 42))
	return &scene{fact: r.IntN(len(funFacts)), rng: r}
}

func (s *scene) log(format string, a ...any) {
	s.logs = append(s.logs, fmt.Sprintf(format, a...))
	if len(s.logs) > 6 {
		s.logs = s.logs[1:]
	}
}

func (s *scene) step() {
	s.tick++
	if s.tick%16 == 1 {
		s.spawned++
		s.minis = append(s.minis, bigWidth)
		s.log(logTemplates[s.rng.IntN(len(logTemplates))], s.spawned-1, s.spawned)
	}
	kept := s.minis[:0]
	for _, x := range s.minis {
		if x+1 > width-len(mini[0])-3 {
			s.ship()
			continue
		}
		kept = append(kept, x+1)
	}
	s.minis = kept
	if s.tick-s.factAt > 45 {
		s.fact = (s.fact + 1 + s.rng.IntN(len(funFacts)-1)) % len(funFacts)
		s.factAt = s.tick
	}
}

// ship ships a factory. Factory n multiplies the fleet by n, so the fleet
// grows factorially, as promised to investors. Factories factor in factorials.
func (s *scene) ship() {
	s.shipped++
	n := s.shipped
	s.produced = s.produced*float64(n) + 1
	if s.rng.IntN(7) == 0 {
		egg := easterEggs[s.rng.IntN(len(easterEggs))]
		s.log("%s", strings.ReplaceAll(egg, "%d", fmt.Sprint(n)))
		return
	}
	s.log("[dock] factory-%d shipped. now producing factory-%d", n, n+1)
}

func (s *scene) render() string {
	grid := make([][]rune, height)
	for i := range grid {
		grid[i] = []rune(strings.Repeat(" ", width))
	}
	put := func(row, col int, str string) {
		for i, r := range []rune(str) {
			if c := col + i; c >= 0 && c < width {
				grid[row][c] = r
			}
		}
	}

	put(0, 3, smokeFrames[(s.tick/3)%len(smokeFrames)])
	put(1, 3, smokeFrames[(s.tick/3+2)%len(smokeFrames)])
	for i, line := range building {
		if strings.Contains(line, "%c") {
			line = fmt.Sprintf(line, gears[s.tick%len(gears)])
		}
		put(2+i, 0, line)
	}
	for _, x := range s.minis {
		if (s.tick+x)%4 < 2 {
			put(beltRow-4, x+2, "~")
		}
		for i, line := range mini {
			put(beltRow-3+i, x, line)
		}
	}
	for c := 0; c < width-3; c++ {
		if (c-s.tick)%3 == 0 {
			grid[beltRow][c] = 'o'
		} else {
			grid[beltRow][c] = '='
		}
	}
	put(beltRow, width-3, ">>>")
	put(beltRow+1, width-13, "SHIPPING DOCK")

	var b strings.Builder
	b.WriteString("\033[H")
	b.WriteString("\033[1;35m SOFTWARE FACTORY FACTORY\033[0m v0.0.0-inf   \033[2m\"we ship factories, not software\"\033[0m\033[K\n\n")
	for _, row := range grid {
		b.WriteString("\033[36m" + string(row) + "\033[0m\033[K\n")
	}
	fmt.Fprintf(&b, "\n\033[32m Factories produced: %-22s Software shipped: 0\033[0m\033[K\n\n", formatCount(s.produced))
	for i := 0; i < 6; i++ {
		line := ""
		if i < len(s.logs) {
			line = s.logs[i]
		}
		b.WriteString(" \033[2m" + truncate(line, width-1) + "\033[0m\033[K\n")
	}
	fact := funFacts[s.fact]
	shown := min(len(fact), (s.tick-s.factAt)*3)
	fmt.Fprintf(&b, "\n\033[1;33m Fun Fact:\033[0;33m %s\033[0m\033[K\n", fact[:shown])
	b.WriteString("\n \033[2mpress q to quit (or esc / ctrl-c). the factories will keep factoring without you.\033[0m\033[K\n")
	// Raw mode disables output post-processing, so emit explicit carriage returns.
	return strings.ReplaceAll(b.String(), "\n", "\r\n")
}

func formatCount(n float64) string {
	switch {
	case math.IsInf(n, 1):
		return "∞ (overflowed; still shipping)"
	case n < 1e15:
		return commas(uint64(n))
	default:
		return fmt.Sprintf("%.3g", n)
	}
}

func commas(n uint64) string {
	s := fmt.Sprint(n)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return s
}

func truncate(s string, n int) string {
	if r := []rune(s); len(r) > n {
		return string(r[:n-1]) + "…"
	}
	return s
}

func main() {
	shipFlag := flag.Bool("ship", false, "ship software instead of factories")
	flag.Parse()
	if *shipFlag {
		fmt.Fprintln(os.Stderr, "error:", Ship(Software{}))
		os.Exit(1)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	quit := make(chan struct{}, 1)
	if fd := int(os.Stdin.Fd()); term.IsTerminal(fd) {
		// Raw mode: read single keys without echoing them onto the factory floor.
		if old, err := term.MakeRaw(fd); err == nil {
			defer term.Restore(fd, old)
			go readQuitKeys(quit)
		}
	}
	fmt.Print("\033[?25l\033[2J")
	defer fmt.Print("\033[?25h\r\n Factory factory halted. Software shipped: 0. Great quarter, team.\r\n")
	s := newScene()
	ticker := time.NewTicker(tickRate)
	defer ticker.Stop()
	for {
		select {
		case <-sig:
			return
		case <-quit:
			return
		case <-ticker.C:
			s.step()
			fmt.Print(s.render())
		}
	}
}

// readQuitKeys signals quit on q, Q, Esc, Ctrl-C, or Ctrl-D and swallows every other key.
func readQuitKeys(quit chan<- struct{}) {
	buf := make([]byte, 16)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			quit <- struct{}{}
			return
		}
		// A lone Esc quits; Esc followed by more bytes is an arrow/function key.
		if n == 1 && buf[0] == 27 {
			quit <- struct{}{}
			return
		}
		for _, c := range buf[:n] {
			switch c {
			case 'q', 'Q', 3, 4:
				quit <- struct{}{}
				return
			}
		}
	}
}
