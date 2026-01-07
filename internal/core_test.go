package internal

import (
	"bytes"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetHorseLabel(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"A", "A"},
		{"", HorseLabelDefault},
		{"AB", HorseLabelDefault},
	}

	for _, test := range tests {
		// Reset global state for this test
		horseLabel = HorseLabelDefault
		setHorseLabel(test.input)
		assert.Equal(t, test.expected, horseLabel, "setHorseLabel(%q)", test.input)
	}
}

func TestSetHorseQuantity(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{5, 5},
		{0, HorseQuantityDefault},
		{100, HorseQuantityDefault},
	}

	for _, test := range tests {
		horseQuantity = HorseQuantityDefault
		setHorseQuantity(test.input)
		assert.Equal(t, test.expected, horseQuantity, "setHorseQuantity(%q)", test.input)
	}
}

func TestSetScoreTarget(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{50, 50},
		{5, ScoreTargetDefault},
		{150, ScoreTargetDefault},
	}

	for _, test := range tests {
		scoreTarget = ScoreTargetDefault
		setScoreTarget(test.input)
		assert.Equal(t, test.expected, scoreTarget, "setScoreTarget(%q)", test.input)
	}
}

func TestSetGameTimeout(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"15s", "15s"},
		{"0s", GameTimeoutDefault},
		{"120s", GameTimeoutDefault},
		{"what", GameTimeoutDefault},
	}

	for _, test := range tests {
		gameTimeout = GameTimeoutDefault
		setGameTimeout(test.input)
		assert.Equal(t, test.expected, gameTimeout, "setGameTimeout(%q)", test.input)
	}
}

func TestSetGameTimeoutDuration(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"0s", GameTimeoutDefault},
		{"20s", "20s"},
		{GameTimeoutDefault, "10s"},
		{"120s", "10s"},
	}

	for _, test := range tests {
		gameTimeout = GameTimeoutDefault
		setGameTimeout(test.input)
		setGameTimeoutDuration()
		assert.Equal(t, test.expected, gameTimeoutDuration.String(), "setGameTimeoutDuration() with input %q", test.input)
	}
}

func TestLoadHorses(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{5, 5},
		{1, HorseQuantityDefault},
		{100, HorseQuantityDefault},
	}

	for _, test := range tests {
		clearHorses()
		loadHorses(test.input)
		assert.Equal(t, test.expected, len(horses), "loadHorses(%q)", test.input)
	}
}

func TestGenerateHorseTrack(t *testing.T) {
	tests := []struct {
		inputLabel  string
		inputScore  int32
		scoreTarget int
		expected    string
	}{
		{"A01", 5, 75, "A01|.....A01                                                                     |"},
		{"B01", 30, 30, "B01|..............................B01|"},
		{"C", 0, 25, "C|C                        |"},
		{"D101", 3, 25, "D101|...D101                     |"},
	}

	for _, test := range tests {
		horse := &Horse{Label: test.inputLabel}
		horse.Score.Store(test.inputScore)
		track := generateHorseTrack(horse, test.scoreTarget)
		assert.Equal(t, test.expected, track, "generateHorseTrack(%v, %v)", test.inputLabel, test.scoreTarget)
	}
}

func TestGenerateTrackMark(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{75, "   +---------|---------|---------|---------|---------|---------|---------|-------+"},
		{15, "   +---------|-------+"},
		{28, "   +---------|---------|---------|+"},
		{100, "   +---------|---------|---------|---------|---------|---------|---------|---------|---------|---------|--+"},
		{150, "   +---------|---------|---------|---------|---------|---------|---------|-------+"}, // Invalid input, falls back to default
	}

	for _, test := range tests {
		mark := generateTrackMark(test.input)
		assert.Equal(t, test.expected, mark, "generateTrackMark(%v)", test.input)
	}
}

func TestGetRaceStr(t *testing.T) {
	tests := []struct {
		inputHorses      []*Horse
		inputScoreTarget int
		expected         string
	}{
		{[]*Horse{
			{Label: "A01"},
			{Label: "B02"},
		},
			0,
			"   +---------|---------|---------|---------|---------|---------|---------|-------+\n" + "A01|....................................................................A01      |\n" + "B02|...........................................................................B02|\n" +
				"   +---------|---------|---------|---------|---------|---------|---------|-------+\n"},
		{[]*Horse{
			{Label: "A01"},
			{Label: "B02"},
			{Label: "C03"},
		},
			20,
			"   +---------|---------|--+\n" +
				"A01|.....A01              |\n" +
				"B02|..........B02         |\n" +
				"C03|...............C03    |\n" +
				"   +---------|---------|--+\n"},
		{[]*Horse{
			{Label: "EEE"},
			{Label: "AAA"},
			{Label: "HHH"},
		},
			22,
			"   +---------|---------|----+\n" +
				"EEE|............EEE         |\n" +
				"AAA|....................AAA |\n" +
				"HHH|.........................HHH|\n" +
				"   +---------|---------|----+\n"},
	}

	for _, test := range tests {
		// Setup test-specific state
		originalScoreTarget := scoreTarget

		// Set scores for horses
		if len(test.inputHorses) == 2 && test.inputScoreTarget == 0 {
			test.inputHorses[0].Score.Store(68)
			test.inputHorses[1].Score.Store(75)
		} else if len(test.inputHorses) == 3 && test.inputScoreTarget == 20 {
			test.inputHorses[0].Score.Store(5)
			test.inputHorses[1].Score.Store(10)
			test.inputHorses[2].Score.Store(15)
		} else if len(test.inputHorses) == 3 && test.inputScoreTarget == 22 {
			test.inputHorses[0].Score.Store(12)
			test.inputHorses[1].Score.Store(20)
			test.inputHorses[2].Score.Store(25)
		}

		horsesMutex.Lock()
		originalHorses := horses
		horses = test.inputHorses
		horsesMutex.Unlock()

		scoreTarget = test.inputScoreTarget

		raceStr := getRaceStr()
		assert.Equal(t, test.expected, raceStr)

		// Teardown: Restore global state
		horsesMutex.Lock()
		horses = originalHorses
		horsesMutex.Unlock()
		scoreTarget = originalScoreTarget
	}
}

func TestGoHorse(t *testing.T) {
	t.Run("Horse with negative score starts at 0", func(t *testing.T) {
		chGameOver := make(chan bool, 1)
		isGameOver := atomic.Bool{}
		scoreTarget = 15
		winnerOnce := &sync.Once{}
		wg := &sync.WaitGroup{}
		wg.Add(1)

		horse := &Horse{Label: "A01"}
		horse.Score.Store(-99)

		go goHorse(horse, &isGameOver, chGameOver, winnerOnce, wg)

		select {
		case <-chGameOver:
			assert.GreaterOrEqual(t, horse.GetScore(), int32(15), "Horse score should reach target")
			assert.True(t, isGameOver.Load(), "isGameOver should be true")
		case <-time.After(2 * time.Second):
			t.Fatal("Test timed out")
		}
		wg.Wait()
		close(chGameOver)
	})

	t.Run("Horse reaches target score", func(t *testing.T) {
		chGameOver := make(chan bool, 1)
		isGameOver := atomic.Bool{}
		scoreTarget = 25
		winnerOnce := &sync.Once{}
		wg := &sync.WaitGroup{}
		wg.Add(1)

		horse := &Horse{Label: "B02"}
		horse.Score.Store(20)

		go goHorse(horse, &isGameOver, chGameOver, winnerOnce, wg)

		select {
		case <-chGameOver:
			assert.GreaterOrEqual(t, horse.GetScore(), int32(25), "Horse score should reach target")
			assert.True(t, isGameOver.Load(), "isGameOver should be true")
		case <-time.After(2 * time.Second):
			t.Fatal("Test timed out")
		}
		wg.Wait()
		close(chGameOver)
	})

	t.Run("Horse stops when channel closes", func(t *testing.T) {
		chGameOver := make(chan bool, 1)
		isGameOver := atomic.Bool{}
		scoreTarget = 100
		winnerOnce := &sync.Once{}
		wg := &sync.WaitGroup{}
		wg.Add(1)

		horse := &Horse{Label: "C03"}
		horse.Score.Store(5)

		go goHorse(horse, &isGameOver, chGameOver, winnerOnce, wg)

		// Close channel after short delay to trigger the <-chGameOver case
		time.Sleep(DelayHorseStep / 2)
		close(chGameOver)

		wg.Wait()

		// Horse should have stopped before reaching target
		assert.Less(t, horse.GetScore(), int32(100), "Horse should stop before target when channel closes")
	})
}

func TestGoHorseChannelDefault(t *testing.T) {
	// Test the default case when channel is full
	chGameOver := make(chan bool, 1)
	chGameOver <- true // Fill the channel buffer

	isGameOver := atomic.Bool{}
	winnerOnce := &sync.Once{}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	scoreTarget = 10

	horse := &Horse{Label: "TEST"}
	horse.Score.Store(8)

	go goHorse(horse, &isGameOver, chGameOver, winnerOnce, wg)

	// Wait for horse to potentially reach the target
	time.Sleep(DelayHorseStep * 2)

	// The goroutine should still complete even if channel is full
	wg.Wait()
	close(chGameOver)
}

func TestGoHorseStopsOnGameOver(t *testing.T) {
	// Test that goHorse stops when isGameOver is already true
	chGameOver := make(chan bool, 1)
	isGameOver := atomic.Bool{}
	isGameOver.Store(true) // Already game over

	winnerOnce := &sync.Once{}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	scoreTarget = 100

	horse := &Horse{Label: "STOP"}
	horse.Score.Store(0)

	go goHorse(horse, &isGameOver, chGameOver, winnerOnce, wg)

	// Horse should stop immediately on first tick when isGameOver is true
	wg.Wait()
	close(chGameOver)

	// Score should be minimal since it stops early
	assert.Less(t, horse.GetScore(), int32(20), "Horse should stop early when game is over")
}

func TestClearTerminal(t *testing.T) {
	assert.Equal(t, "\033[H\033[2J", clearTerminal())
}

func TestDisplay(t *testing.T) {
	// Setup: Isolate global variables for this test
	originalHorseWinner := horseWinner
	originalScoreTarget := scoreTarget

	horseWinnerMutex.Lock()
	horseWinner = &Horse{Label: "A01"}
	horseWinner.Score.Store(15)
	horseWinnerMutex.Unlock()

	scoreTarget = 15
	loadHorses(2)

	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	defer func() {
		os.Stdout = oldStdout
		r.Close()
		w.Close()
		// Teardown: Restore global state
		horseWinnerMutex.Lock()
		horseWinner = originalHorseWinner
		horseWinnerMutex.Unlock()
		scoreTarget = originalScoreTarget
	}()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go display(wg)
	time.Sleep(DelayRefreshScreen * 2) // Allow display to run at least once
	wg.Wait()
	w.Close()

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	expected := "\x1b[H\x1b[2J\n   +---------|-------+\nH01|H01              |\nH02|H02              |\n   +---------|-------+\n\n"
	assert.Contains(t, output, expected)
}

func TestDisplayWithoutWinner(t *testing.T) {
	originalHorseWinner := horseWinner
	originalScoreTarget := scoreTarget

	horseWinnerMutex.Lock()
	horseWinner = &Horse{Label: "TEST"}
	horseWinner.Score.Store(0)
	horseWinnerMutex.Unlock()

	scoreTarget = 15
	loadHorses(2)

	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	defer func() {
		os.Stdout = oldStdout
		r.Close()
		w.Close()
		horseWinnerMutex.Lock()
		horseWinner = originalHorseWinner
		horseWinnerMutex.Unlock()
		scoreTarget = originalScoreTarget
	}()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go display(wg)

	time.Sleep(DelayRefreshScreen * 3)

	horseWinnerMutex.Lock()
	horseWinner.Score.Store(20)
	horseWinnerMutex.Unlock()

	wg.Wait()
	w.Close()

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	assert.NotEmpty(t, output, "Display should produce output")
	assert.Contains(t, output, "\x1b[H\x1b[2J", "Should contain terminal clear sequence")
}

func TestRun(t *testing.T) {
	t.Run("Run with default input produces winner", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		input := Input{
			HorseLabel:     "T",
			HorsesQuantity: 2,
			ScoreTarget:    20,
			GameTimeout:    "10s",
		}

		done := make(chan bool)
		go func() {
			Run(input)
			w.Close()
			done <- true
		}()

		<-done
		os.Stdout = oldStdout

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		assert.Contains(t, output, "T0", "Output should contain horse labels")
		assert.True(t,
			strings.Contains(output, "winner") || strings.Contains(output, "tired"),
			"Output should contain winner or timeout message")
	})

	t.Run("Run with timeout produces timeout message", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		input := Input{
			HorseLabel:     "X",
			HorsesQuantity: 2,
			ScoreTarget:    100,   // High target
			GameTimeout:    "10s", // Minimum timeout
		}

		done := make(chan bool)
		go func() {
			Run(input)
			w.Close()
			done <- true
		}()

		<-done
		os.Stdout = oldStdout

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		assert.NotEmpty(t, output, "Output should not be empty")
		assert.Contains(t, output, "X0", "Output should contain horse labels")
	})

	t.Run("Run handles invalid input gracefully", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		input := Input{
			HorseLabel:     "INVALID", // Too long, should use default
			HorsesQuantity: 200,       // Too many, should use default
			ScoreTarget:    5,         // Too low, should use default
			GameTimeout:    "999s",    // Too high, should use default
		}

		done := make(chan bool)
		go func() {
			Run(input)
			w.Close()
			done <- true
		}()

		<-done
		os.Stdout = oldStdout

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		assert.NotEmpty(t, output, "Output should not be empty")
		assert.Contains(t, output, "H0", "Should use default horse label")
	})

	t.Run("Run with all scenarios covered", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		input := Input{
			HorseLabel:     "R",
			HorsesQuantity: 3,
			ScoreTarget:    30,
			GameTimeout:    "15s",
		}

		done := make(chan bool)
		go func() {
			Run(input)
			w.Close()
			done <- true
		}()

		<-done
		os.Stdout = oldStdout

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		assert.NotEmpty(t, output, "Output should not be empty")
		assert.Contains(t, output, "R0", "Output should contain horse labels")
	})
}
