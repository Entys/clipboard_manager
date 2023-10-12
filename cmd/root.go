package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mitchellh/go-homedir"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	Version    string
	daemonFlag bool
)

var rootCmd = &cobra.Command{
	Use:     "cpmanage",
	Short:   "Clipboard manager",
	Version: Version,
	Run:     run,
}

type Copy struct {
	Date string `json:"date"`
	Copy string `json:"copy"`
}

var clipboardHistoryFile = "/tmp/clipboard/history"
var clipboardHistory []Copy

func init() {
	clipboardHistory = []Copy{}
	rootCmd.Flags().BoolVar(&daemonFlag, "daemon", false, "Run as a daemon")
}

func saveClipboardHistoryToFile(filename string, history []Copy) {
	data, err := json.Marshal(history)
	if err != nil {
		log.Printf("Error marshaling clipboard history: %v\n", err)
		return
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		log.Printf("Error writing clipboard history to file: %v\n", err)
	}
}

func checkClipboardChanges() {
	cmd := exec.Command("wl-paste")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error reading clipboard with wl-paste: %v\n", err)
		return
	}

	currentClipboardContent := strings.TrimSpace(string(output))

	if len(clipboardHistory) > 0 && currentClipboardContent == clipboardHistory[0].Copy {
		return
	}

	newEntry := Copy{
		Date: time.Now().Format("2006-01-02 15:04:05"),
		Copy: currentClipboardContent,
	}

	clipboardHistory = append([]Copy{newEntry}, clipboardHistory...)

	saveClipboardHistoryToFile(clipboardHistoryFile, clipboardHistory)
}

func run(cmd *cobra.Command, args []string) {

	if daemonFlag {
		for {
			checkClipboardChanges()
			time.Sleep(2 * time.Second)
		}
		return
	}

	clipboardHistoryFile, e := homedir.Expand(clipboardHistoryFile)
	if e != nil {
		log.Fatal(e)
	}

	if _, e := os.Stat(clipboardHistoryFile); os.IsNotExist(e) {
		if e := createFileRecursive(clipboardHistoryFile); e != nil {
			log.Fatal(e)
		}
	}

	app := tview.NewApplication()

	loadClipboardHistoryFromFile(clipboardHistoryFile)

	clipboardList := tview.NewTable().SetSelectable(true, false)
	clipboardList.SetBackgroundColor(tcell.ColorReset)
	clipboardList.SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorBlue))

	timeOfCopyStyle := tcell.StyleDefault.Foreground(tcell.ColorYellow)
	messageCopiedStyle := tcell.StyleDefault.Foreground(tcell.ColorBlue)

	clipboardList.SetCell(0, 0, tview.NewTableCell("Time of Copy").SetAlign(tview.AlignCenter).SetStyle(timeOfCopyStyle))
	clipboardList.SetCell(0, 1, tview.NewTableCell("|").SetAlign(tview.AlignCenter).SetStyle(tcell.StyleDefault))
	clipboardList.SetCell(0, 2, tview.NewTableCell("Message Copied").SetAlign(tview.AlignCenter).SetStyle(messageCopiedStyle))

	clipboardList.SetTitle("Clipboard History").SetTitleAlign(tview.AlignCenter)
	clipboardList.SetBorder(true)

	for i, item := range clipboardHistory {
		clipboardList.SetCellSimple(i+1, 0, item.Date)
		clipboardList.SetCellSimple(i+1, 1, "|")
		clipboardList.SetCellSimple(i+1, 2, item.Copy)
	}
	clipboardList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			row, _ := clipboardList.GetSelection()
			if row == 0 {
				return nil
			}
			text := clipboardList.GetCell(row, 2).Text
			err := exec.Command("wl-copy", text).Run()
			if err != nil {
				return nil
			}
		}

		return event
	})
	if e := app.SetRoot(clipboardList, true).SetFocus(clipboardList).Run(); e != nil {
		panic(e)
	}
}

func createFileRecursive(filename string) error {
	if _, e := os.Stat(filename); os.IsNotExist(e) {
		if e := os.MkdirAll(filepath.Dir(filename), os.ModePerm); e != nil {
			return e
		}

		f, e := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0644)
		if e != nil {
			return e
		}
		f.Close()
	}

	return nil
}

func loadClipboardHistoryFromFile(clipboardHistoryFile string) {
	data, err := os.ReadFile(clipboardHistoryFile)
	if err != nil {
		log.Printf("Error reading clipboard history file: %v\n", err)
		return
	}
	err = json.Unmarshal(data, &clipboardHistory)
	if err != nil {
		log.Printf("Error unmarshaling clipboard history file: %v\n", err)
		return
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
