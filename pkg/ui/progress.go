package ui

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ProgressPhase represents different phases of git clone
type ProgressPhase int

const (
	PhaseStarting ProgressPhase = iota
	PhaseEnumerating
	PhaseCounting
	PhaseCompressing
	PhaseReceiving
	PhaseResolving
	PhaseDone
)

// ProgressInfo holds information about clone progress
type ProgressInfo struct {
	Phase      ProgressPhase
	Percentage int
	Current    int
	Total      int
	Speed      string
	Message    string
	mu         sync.RWMutex
}

// NewProgressInfo creates a new progress tracker
func NewProgressInfo() *ProgressInfo {
	return &ProgressInfo{
		Phase: PhaseStarting,
	}
}

// Update updates progress information
func (p *ProgressInfo) Update(phase ProgressPhase, percentage, current, total int, speed, message string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Phase = phase
	p.Percentage = percentage
	p.Current = current
	p.Total = total
	p.Speed = speed
	p.Message = message
}

// GetSnapshot returns a copy of current progress
func (p *ProgressInfo) GetSnapshot() (ProgressPhase, int, int, int, string, string) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Phase, p.Percentage, p.Current, p.Total, p.Speed, p.Message
}

// ProgressModel is the bubble tea model for displaying progress
type ProgressModel struct {
	spinner      spinner.Model
	progress     *ProgressInfo
	done         bool
	phaseHistory map[ProgressPhase]bool
}

// NewProgressModel creates a new progress display model
func NewProgressModel(progress *ProgressInfo) ProgressModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(colorBlue)

	return ProgressModel{
		spinner:      s,
		progress:     progress,
		phaseHistory: make(map[ProgressPhase]bool),
	}
}

func (m ProgressModel) Init() tea.Cmd {
	return m.spinner.Tick
}

type doneMsg struct{}

func (m ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.done = true
			return m, tea.Quit
		}
	case doneMsg:
		m.done = true
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m ProgressModel) View() string {
	if m.done {
		return ""
	}

	phase, percentage, current, total, speed, _ := m.progress.GetSnapshot()

	var output strings.Builder

	// Mark phase as seen
	m.phaseHistory[phase] = true

	// Render completed phases with checkmarks
	phases := []struct {
		phase ProgressPhase
		name  string
		emoji string
	}{
		{PhaseEnumerating, "Enumerating objects", "🔍"},
		{PhaseCounting, "Counting objects", "📊"},
		{PhaseCompressing, "Compressing objects", "⚡"},
		{PhaseReceiving, "Receiving objects", "↓"},
		{PhaseResolving, "Resolving deltas", "🔗"},
	}

	for _, p := range phases {
		if phase == p.phase {
			// Current phase - show spinner and progress
			if total > 0 && percentage > 0 {
				bar := renderProgressBar(percentage, 25)
				stats := fmt.Sprintf("%d%% (%s/%s)", percentage, formatNumber(current), formatNumber(total))
				if speed != "" {
					stats += fmt.Sprintf(", %s", speed)
				}
				output.WriteString(fmt.Sprintf("   %s %s %s %s\n",
					m.spinner.View(),
					dimStyle.Render(p.name+":"),
					bar,
					progressText.Render(stats)))
			} else {
				output.WriteString(fmt.Sprintf("   %s %s...\n",
					m.spinner.View(),
					dimStyle.Render(p.name)))
			}
		} else if m.phaseHistory[p.phase] && phase > p.phase {
			// Completed phase - show checkmark
			output.WriteString(fmt.Sprintf("   %s %s\n",
				successIcon.Render("✓"),
				dimStyle.Render(p.name+": done")))
		}
	}

	return output.String()
}

// renderProgressBar creates a styled progress bar
func renderProgressBar(percentage, width int) string {
	if percentage > 100 {
		percentage = 100
	}
	filled := (percentage * width) / 100
	if filled > width {
		filled = width
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return progressBar.Render(bar)
}

// formatNumber formats numbers with commas for readability
func formatNumber(n int) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}

	var result strings.Builder
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(c)
	}
	return result.String()
}

// GitProgressParser parses git output and updates progress
type GitProgressParser struct {
	progress *ProgressInfo
	program  *tea.Program
	lastLine string
}

// NewGitProgressParser creates a new parser
func NewGitProgressParser(progress *ProgressInfo, program *tea.Program) *GitProgressParser {
	return &GitProgressParser{
		progress: progress,
		program:  program,
	}
}

var (
	enumeratingRe = regexp.MustCompile(`Enumerating objects:\s*(\d+)`)
	countingRe    = regexp.MustCompile(`Counting objects:\s*(\d+)%\s*\((\d+)/(\d+)\)`)
	compressingRe = regexp.MustCompile(`Compressing objects:\s*(\d+)%\s*\((\d+)/(\d+)\)`)
	receivingRe   = regexp.MustCompile(`Receiving objects:\s*(\d+)%\s*\((\d+)/(\d+)\)(?:,\s*([0-9.]+\s*[KMG]iB))?(?:\s*\|\s*([0-9.]+\s*[KMG]iB/s))?`)
	resolvingRe   = regexp.MustCompile(`Resolving deltas:\s*(\d+)%\s*\((\d+)/(\d+)\)`)
	doneRe        = regexp.MustCompile(`done\.?\s*$`)
)

// ParseLine parses a single line of git output
func (p *GitProgressParser) ParseLine(line string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}

	// Check for "done" completion
	if doneRe.MatchString(line) {
		return
	}

	// Parse different progress patterns
	if matches := receivingRe.FindStringSubmatch(line); matches != nil {
		percentage := parseInt(matches[1])
		current := parseInt(matches[2])
		total := parseInt(matches[3])
		speed := ""
		if len(matches) > 4 && matches[4] != "" {
			speed = matches[4]
		}
		p.progress.Update(PhaseReceiving, percentage, current, total, speed, "")
	} else if matches := resolvingRe.FindStringSubmatch(line); matches != nil {
		percentage := parseInt(matches[1])
		current := parseInt(matches[2])
		total := parseInt(matches[3])
		p.progress.Update(PhaseResolving, percentage, current, total, "", "")
	} else if matches := compressingRe.FindStringSubmatch(line); matches != nil {
		percentage := parseInt(matches[1])
		current := parseInt(matches[2])
		total := parseInt(matches[3])
		p.progress.Update(PhaseCompressing, percentage, current, total, "", "")
	} else if matches := countingRe.FindStringSubmatch(line); matches != nil {
		percentage := parseInt(matches[1])
		current := parseInt(matches[2])
		total := parseInt(matches[3])
		p.progress.Update(PhaseCounting, percentage, current, total, "", "")
	} else if matches := enumeratingRe.FindStringSubmatch(line); matches != nil {
		total := parseInt(matches[1])
		p.progress.Update(PhaseEnumerating, 0, 0, total, "", "")
	} else if strings.Contains(line, "Enumerating") {
		p.progress.Update(PhaseEnumerating, 0, 0, 0, "", line)
	} else if strings.Contains(line, "Counting") {
		p.progress.Update(PhaseCounting, 0, 0, 0, "", line)
	} else if strings.Contains(line, "Compressing") {
		p.progress.Update(PhaseCompressing, 0, 0, 0, "", line)
	}

	p.lastLine = line
}

// parseInt safely converts string to int
func parseInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}

// ProgressWriter wraps an io.Writer and parses git progress
type ProgressWriter struct {
	parser *GitProgressParser
	buf    []byte
}

// NewProgressWriter creates a new progress writer
func NewProgressWriter(parser *GitProgressParser) *ProgressWriter {
	return &ProgressWriter{
		parser: parser,
		buf:    make([]byte, 0, 4096),
	}
}

// Write implements io.Writer
func (w *ProgressWriter) Write(p []byte) (n int, err error) {
	// Accumulate data in buffer
	w.buf = append(w.buf, p...)

	// Process complete lines
	for {
		// Look for line endings (\n or \r)
		idx := -1
		for i, b := range w.buf {
			if b == '\n' || b == '\r' {
				idx = i
				break
			}
		}

		if idx == -1 {
			// No complete line yet
			break
		}

		// Extract line
		line := string(w.buf[:idx])

		// Skip past the line ending
		skip := idx + 1
		if skip < len(w.buf) && w.buf[idx] == '\r' && w.buf[skip] == '\n' {
			skip++ // Skip \r\n
		}

		// Remove processed data from buffer
		w.buf = w.buf[skip:]

		// Parse the line
		w.parser.ParseLine(line)
	}

	return len(p), nil
}

// CloneProgressDisplay manages the progress display for git clone
type CloneProgressDisplay struct {
	progress *ProgressInfo
	program  *tea.Program
	parser   *GitProgressParser
	writer   *ProgressWriter
}

// NewCloneProgressDisplay creates a new progress display
func NewCloneProgressDisplay() *CloneProgressDisplay {
	progress := NewProgressInfo()
	model := NewProgressModel(progress)
	program := tea.NewProgram(model)
	parser := NewGitProgressParser(progress, program)
	writer := NewProgressWriter(parser)

	return &CloneProgressDisplay{
		progress: progress,
		program:  program,
		parser:   parser,
		writer:   writer,
	}
}

// Start starts the progress display
func (d *CloneProgressDisplay) Start() {
	go func() {
		d.program.Run()
	}()
	time.Sleep(50 * time.Millisecond) // Give UI time to initialize
}

// Writer returns the io.Writer for git output
func (d *CloneProgressDisplay) Writer() io.Writer {
	return d.writer
}

// Stop stops the progress display
func (d *CloneProgressDisplay) Stop() {
	d.progress.Update(PhaseDone, 100, 0, 0, "", "")
	time.Sleep(100 * time.Millisecond) // Let UI update
	d.program.Send(doneMsg{})
	time.Sleep(50 * time.Millisecond) // Give time to clean up
}

// StreamGitOutput streams git output through the progress parser
func StreamGitOutput(reader io.Reader, display *CloneProgressDisplay) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		display.Writer().Write([]byte(line + "\n"))
	}
	return scanner.Err()
}

// ShowCloneProgress displays the initial clone message
func ShowCloneProgress(repoUrl string) {
	fmt.Printf("\n%s  %s %s\n\n",
		infoIcon.Render("↓"),
		Dim("Cloning"),
		Path(repoUrl))
}
