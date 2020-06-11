package ui

import (
	"github.com/cinarmert/doclogs/cmd/docklogs/container"
	"github.com/gdamore/tcell"
	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
)

type LayoutManager struct {
	App      *tview.Application
	Grid     *tview.Grid
	Sessions []*container.Session
}

func NewManagerForSessions(sessions []*container.Session) (*LayoutManager, error) {
	if len(sessions) == 0 {
		return nil, errors.New("could not create layour manager: no sessions given")
	}

	lm := (&LayoutManager{}).
		SetApp().
		SetSessions(sessions).
		SetGrid()

	return lm, nil
}

func (lm *LayoutManager) SetApp() *LayoutManager {
	lm.App = tview.NewApplication()
	return lm
}

func (lm *LayoutManager) SetSessions(sessions []*container.Session) *LayoutManager {
	lm.Sessions = sessions
	return lm
}

func (lm *LayoutManager) SetGrid() *LayoutManager {
	n := len(lm.Sessions)
	grid := tview.NewGrid().
		SetRows(make([]int, n)...).
		SetBorders(true)

	lm.Grid = grid
	return lm
}

func (lm *LayoutManager) createTextView(title string) *tview.TextView {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		ScrollToEnd().
		SetChangedFunc(func() {
			lm.App.Draw()
		})

	textView.SetBorder(true)
	textView.SetTitle(" " + title + " ")
	textView.SetTitleColor(tcell.ColorRed)
	return textView
}

func (lm *LayoutManager) Run() {
	var wg sync.WaitGroup

	for i, session := range lm.Sessions {
		wg.Add(1)
		tv := lm.createTextView(session.Name)
		lm.Grid.AddItem(tv, i, 0, 1, 1, 0, 0, false)
		go session.ReadLogs(&wg, tv)
	}

	if !isatty.IsTerminal(os.Stdout.Fd()) {
		log.Warnf("doclogs is only available in a tty environment")
		return
	}

	if err := lm.App.SetRoot(lm.Grid, true).Run(); err != nil {
		log.Fatalf("could not init ui: %v", err)
		os.Exit(1)
	}
}
