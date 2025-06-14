package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func CreateMen(window fyne.Window) *fyne.MainMenu {

	fileMenu := fyne.NewMenu("檔案",
		fyne.NewMenuItem("開啟", func() {
			dialog.ShowInformation("關於", "這是個關於內容", window)
		}),
	)
	helpMenu := fyne.NewMenu("幫助",
		fyne.NewMenuItem("<UNK>", func() {
			dialog.ShowInformation("關於", "這是個關於內容", window)
		}),
	)
	mainMenu := fyne.NewMainMenu(fileMenu, helpMenu)

	return mainMenu
}
