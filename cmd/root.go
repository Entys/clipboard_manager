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
	"path/filepath"
)

var (
	Version string
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

var clipboardHistoryFile = "~/.clipboard/history"
var clipboardHistory []Copy

func init() {
	clipboardHistory = []Copy{}
}

func run(cmd *cobra.Command, args []string) {
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
	clipboardList.SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlue))

	clipboardList.SetCell(0, 0, tview.NewTableCell("Time of Copy").SetAlign(tview.AlignCenter))
	clipboardList.SetCell(0, 1, tview.NewTableCell("Message Copied").SetAlign(tview.AlignCenter))

	clipboardList.SetTitle("Clipboard History").SetTitleAlign(tview.AlignCenter)
	clipboardList.SetBorder(true)

	for i, item := range clipboardHistory {
		clipboardList.SetCellSimple(i+1, 0, item.Date)
		clipboardList.SetCellSimple(i+1, 1, item.Copy)
	}

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
