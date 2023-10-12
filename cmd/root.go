package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/mitchellh/go-homedir"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
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

var clipboardHistoryFile = "~/.clipboard/history"
var clipboardHistory []string

func init() {
	clipboardHistory = []string{}
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
		clipboardList.SetCellSimple(i+1, 0, item)
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
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if line != "" {
			clipboardHistory = append(clipboardHistory, line)
		}
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
