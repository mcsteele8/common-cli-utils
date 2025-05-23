package spinner

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	cilUtilsColor "github.com/mcsteele8/common-cli-utils/color"

	"github.com/fatih/color"

	"golang.org/x/term"
)

// Spinner struct to hold the provided options.
type Spinner struct {
	mu            *sync.RWMutex                 //
	Delay         time.Duration                 // Delay is the speed of the indicator
	chars         []string                      // chars holds the chosen character set
	Prefix        string                        // Prefix is the text preppended to the indicator
	Suffix        string                        // Suffix is the text appended to the indicator
	FinalMSG      string                        // string displayed after Stop() is called
	lastOutput    string                        // last character(set) written
	color         func(a ...interface{}) string // default color is white
	Writer        io.Writer                     // to make testing better, exported so users have access. Use `WithWriter` to update after initialization.
	active        bool                          // active holds the state of the spinner
	stopChan      chan struct{}                 // stopChan is a channel used to stop the indicator
	HideCursor    bool                          // hideCursor determines if the cursor is visible
	PreUpdate     func(s *Spinner)              // will be triggered before every spinner update
	PostUpdate    func(s *Spinner)              // will be triggered after every spinner update
	DisableOutput bool                          //dynamically configure output. Sometime someone will want verbose output and this should disable spinner when the is true.
}

const (
	spinningDotsCharsetNum = 14
	spinInterval           = time.Millisecond * 100
)

var (
	state = New(CharSets[spinningDotsCharsetNum], spinInterval)

	// returns true if the OS is windows and the WT_SESSION env variable is set.
	isWindowsTerminalOnWindows = len(os.Getenv("WT_SESSION")) > 0 && runtime.GOOS == "windows"
)

// Start begins the spinner.
func Start() {
	state.HideCursor = true
	state.Color("yellow")

	state.Start()
}

// Stop ends the spinner.
func Stop() {
	state.Stop()
}

// IsActive will return if the global spinner is active
func IsActive() bool {
	return state.active
}

func Color(c string) {
	state.Color(c)
}

func Suffix(message string) {
	state.Suffix = message
}

func SetFinalMsg(message string) {
	state.FinalMSG = "\t" + cilUtilsColor.Yellow.Paint(message) + "\n"
}

// NewDefault provides a pointer to an instance of Spinner with our
// default charset and spin interval, along with any options provided
func NewDefault(options ...Option) *Spinner {
	s := New(CharSets[spinningDotsCharsetNum], spinInterval)
	for _, opt := range options {
		opt(s)
	}

	return s
}

// New provides a pointer to an instance of Spinner with the supplied options.
func New(cs []string, d time.Duration, options ...Option) *Spinner {
	s := &Spinner{
		Delay:    d,
		chars:    cs,
		color:    color.New(color.FgWhite).SprintFunc(),
		mu:       &sync.RWMutex{},
		Writer:   color.Output,
		active:   false,
		stopChan: make(chan struct{}, 1),
	}

	for _, option := range options {
		option(s)
	}
	return s
}

// Option is a function that takes a spinner and applies
// a given configuration.
type Option func(*Spinner)

// Options contains fields to configure the spinner.
type Options struct {
	Color      string
	Suffix     string
	FinalMSG   string
	HideCursor bool
}

// WithColor adds the given color to the spinner.
func WithColor(color string) Option {
	return func(s *Spinner) {
		err := s.Color(color)
		if err != nil {
			fmt.Printf("failed to set color to %s: %s", color, err)
		}
	}
}

// WithFinalMSG adds the given string ot the spinner
// as the final message to be written.
func WithFinalMSG(finalMsg string) Option {
	return func(s *Spinner) {
		s.FinalMSG = finalMsg
	}
}

// WithHiddenCursor hides the cursor
// if hideCursor = true given.
func WithHiddenCursor(hideCursor bool) Option {
	return func(s *Spinner) {
		s.HideCursor = hideCursor
	}
}

// WithSuffix adds the given string to the spinner
// as the suffix.
func WithSuffix(suffix string) Option {
	return func(s *Spinner) {
		s.Suffix = suffix
	}
}

// WithWriter adds the given writer to the spinner. This
// function should be favored over directly assigning to
// the struct value.
func WithWriter(w io.Writer) Option {
	return func(s *Spinner) {
		s.mu.Lock()
		s.Writer = w
		s.mu.Unlock()
	}
}

// Active will return whether or not the spinner is currently active.
func (s *Spinner) Active() bool {
	return s.active
}

// Start will start the indicator.
func (s *Spinner) Start() {
	if s.DisableOutput {
		return
	}

	s.mu.Lock()
	if s.active {
		s.mu.Unlock()
		return
	}
	if s.HideCursor && !isWindowsTerminalOnWindows {
		// hides the cursor
		fmt.Fprint(s.Writer, "\033[?25l")
	}
	s.active = true
	s.mu.Unlock()

	go func() {
		for {
			for i := 0; i < len(s.chars); i++ {
				select {
				case <-s.stopChan:
					return
				default:
					s.mu.Lock()
					if !s.active {
						s.mu.Unlock()
						return
					}
					if !isWindowsTerminalOnWindows {
						s.erase()
					}

					if s.PreUpdate != nil {
						s.PreUpdate(s)
					}

					var outColor string
					if runtime.GOOS == "windows" {
						if s.Writer == os.Stderr {
							outColor = fmt.Sprintf("\r%s%s %s ", s.Prefix, s.chars[i], s.Suffix)
						} else {
							outColor = fmt.Sprintf("\r%s%s %s ", s.Prefix, s.color(s.chars[i]), s.Suffix)
						}
					} else {
						outColor = fmt.Sprintf("\r%s%s %s ", s.Prefix, s.color(s.chars[i]), s.Suffix)
					}
					maxWidth, _, _ := term.GetSize(0)
					if maxWidth != 0 && len(outColor) > maxWidth {
						outColor = outColor[:maxWidth]
					}
					outPlain := fmt.Sprintf("\r%s%s %s ", s.Prefix, s.chars[i], s.Suffix)
					fmt.Fprint(s.Writer, outColor)
					s.lastOutput = outPlain
					delay := s.Delay

					if s.PostUpdate != nil {
						s.PostUpdate(s)
					}

					s.mu.Unlock()
					time.Sleep(delay)
				}
			}
		}
	}()
}

// Stop stops the indicator.
func (s *Spinner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.active {
		s.active = false
		if s.HideCursor && !isWindowsTerminalOnWindows {
			// makes the cursor visible
			fmt.Fprint(s.Writer, "\033[?25h")
		}
		s.erase()
		if s.FinalMSG != "" {
			if isWindowsTerminalOnWindows {
				fmt.Fprint(s.Writer, "\r", s.FinalMSG)
			} else {
				fmt.Fprint(s.Writer, s.FinalMSG)
			}
		}
		s.stopChan <- struct{}{}
	}
}

// Restart will stop and start the indicator.
func (s *Spinner) Restart() {
	s.Stop()
	s.Start()
}

// Reverse will reverse the order of the slice assigned to the indicator.
func (s *Spinner) Reverse() {
	s.mu.Lock()
	for i, j := 0, len(s.chars)-1; i < j; i, j = i+1, j-1 {
		s.chars[i], s.chars[j] = s.chars[j], s.chars[i]
	}
	s.mu.Unlock()
}

// Color will set the struct field for the given color to be used. The spinner
// will need to be explicitly restarted.
func (s *Spinner) Color(colors ...string) error {
	colorAttributes := make([]color.Attribute, len(colors))

	// Verify colours are valid and place the appropriate attribute in the array
	for index, c := range colors {
		if !validColor(c) {
			return errInvalidColor
		}
		colorAttributes[index] = colorAttributeMap[c]
	}

	s.mu.Lock()
	s.color = color.New(colorAttributes...).SprintFunc()
	s.mu.Unlock()
	return nil
}

// UpdateSpeed will set the indicator delay to the given value.
func (s *Spinner) UpdateSpeed(d time.Duration) {
	s.mu.Lock()
	s.Delay = d
	s.mu.Unlock()
}

// UpdateCharSet will change the current character set to the given one.
func (s *Spinner) UpdateCharSet(cs []string) {
	s.mu.Lock()
	s.chars = cs
	s.mu.Unlock()
}

// erase deletes written characters.
// Caller must already hold s.lock.
func (s *Spinner) erase() {
	n := utf8.RuneCountInString(s.lastOutput)
	if runtime.GOOS == "windows" && !isWindowsTerminalOnWindows {
		clearString := "\r" + strings.Repeat(" ", n) + "\r"
		fmt.Fprint(s.Writer, clearString)
		s.lastOutput = ""
		return
	}

	fmt.Fprintf(s.Writer, "\r\033[K")
	s.lastOutput = ""
}

// Lock allows for manual control to lock the spinner.
func (s *Spinner) Lock() {
	s.mu.Lock()
}

// Unlock allows for manual control to unlock the spinner.
func (s *Spinner) Unlock() {
	s.mu.Unlock()
}

// GenerateNumberSequence will generate a slice of integers at the
// provided length and convert them each to a string.
func GenerateNumberSequence(length int) []string {
	numSeq := make([]string, length)
	for i := 0; i < length; i++ {
		numSeq[i] = strconv.Itoa(i)
	}
	return numSeq
}
