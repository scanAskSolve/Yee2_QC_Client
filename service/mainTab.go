package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"io/ioutil"
	"log"
	"net/http"
)

func CreateMainTab() *container.AppTabs {

	tabs := container.NewAppTabs(
		container.NewTabItem("首頁", createHomeTab()),
		//container.NewTabItem("表單", createFormTab()),
		//container.NewTabItem("列表", createListTab()),
		//container.NewTabItem("IP", createIPTab()),
	)

	return tabs
}

func createHomeTab() *fyne.Container {
	title := widget.NewCard("歡迎", "這是首頁",
		widget.NewLabel("歡迎使用 Fyne GUI 框架！"))

	var ips []string
	var selectedIP string

	ipSelect := widget.NewSelect(ips, func(selected string) {
		selectedIP = selected
		log.Println("選擇了 IP:", selected) // 當用戶選擇 IP 時，打印到控制台
	})

	ipSelect.PlaceHolder = "請選擇 IP 地址" // 設置預設提示文字

	spinner := widget.NewProgressBarInfinite()
	spinner.Hide()

	refreshButton := widget.NewButton("刷新IP", func() {
		selectedIP = ""
		spinner.Show()
		ipSelect.Options = []string{}

		go func() {
			ips = udpService()
			fmt.Println("find IP:", ips)
			fyne.Do(func() {

				fyne.CurrentApp().Driver().CanvasForObject(ipSelect).Refresh(ipSelect)
				ipSelect.Options = ips
				ipSelect.Refresh()
				spinner.Hide()

				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title:   "更新完成",
					Content: fmt.Sprintf("已搜尋到IP:%v", ips),
				})
			})
		}()
	})

	statusLabel := canvas.NewText("", nil)
	testButton := widget.NewButton("測試連線", func() {
		var url = "http://" + selectedIP + ":8080/health"

		go func() {
			resp, err := http.Get(url)

			if resp != nil {
				log.Printf("run url:%s resp:%+v err:%v\n", url, resp, err)
			} else {
				log.Printf("run url:%s resp:nil err:%v\n", url, err)
			}
			if err != nil {
				fyne.Do(func() {
					statusLabel.Text = "錯誤訊息:" + err.Error()
					statusLabel.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
				})
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fyne.Do(func() {
					statusLabel.Text = "讀取失敗"
					statusLabel.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
				})
				return
			}

			if string(body) != "data:1" {
				fyne.Do(func() {
					statusLabel.Text = "連線成功"
					statusLabel.Color = color.RGBA{R: 0, G: 255, B: 0, A: 255}
				})
			} else {
				fyne.Do(func() {
					statusLabel.Text = "連線回應錯誤:" + string(body)
					statusLabel.Color = color.RGBA{R: 0, G: 255, B: 0, A: 255}
				})
			}
		}()
	})

	row1 := container.NewHBox(refreshButton, spinner)
	row2 := container.NewHBox(ipSelect, testButton)
	row3 := container.NewHBox(statusLabel)
	return container.NewVBox(title, row1, row2, row3)
}

func createFormTab() *fyne.Container {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("請輸入姓名")

	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("請輸入電子郵件")

	ageSlider := widget.NewSlider(18, 100)
	ageSlider.Value = 25

	ageLabel := widget.NewLabel("年齡: 25")
	ageSlider.OnChanged = func(value float64) {
		ageLabel.SetText(fmt.Sprintf("年齡: %.0f", value))
	}

	submitButton := widget.NewButton("提交", func() {
		info := fmt.Sprintf("姓名: %s\n電子郵件: %s\n年齡: %.0f",
			nameEntry.Text, emailEntry.Text, ageSlider.Value)
		dialog.ShowInformation("表單資料", info, fyne.CurrentApp().Driver().AllWindows()[0])
	})

	form := container.NewVBox(
		widget.NewLabel("個人資訊表單"),
		nameEntry,
		emailEntry,
		ageLabel,
		ageSlider,
		submitButton,
	)

	return form
}

func createListTab() *fyne.Container {
	list := widget.NewList(
		func() int {
			return 10 // 列表項目數量
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("項目")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(fmt.Sprintf("列表項目 %d", id+1))
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		dialog.ShowInformation("選擇",
			fmt.Sprintf("你選擇了項目 %d", id+1),
			fyne.CurrentApp().Driver().AllWindows()[0])
	}

	return container.NewVBox(
		widget.NewLabel("項目列表"),
		list,
	)
}

// 建立 IP 標籤頁內容
func createIPTab() fyne.CanvasObject {
	// 生成下拉選單的選項（127.0.0.1 ~ 127.0.0.5）
	ipOptions := []string{}
	// 建立下拉選單
	ipSelect := widget.NewSelect(ipOptions, func(selected string) {
		fmt.Println("選擇了 IP:", selected) // 當用戶選擇 IP 時，打印到控制台
	})

	ipSelect.PlaceHolder = "請選擇 IP 地址" // 設置預設提示文字

	// 顯示SSH搜尋到的內容
	sshContent := widget.NewMultiLineEntry()
	sshContent.SetPlaceHolder("連線成功後的內容顯示在這裡....")

	// 輸入框 額外測試指令
	commandInput := widget.NewEntry()
	commandInput.SetPlaceHolder("輸入指令或是參數")

	refreshButton := widget.NewButton("刷新IP", func() {
		sshContent.SetText("正在掃描內網中")
		//go func() {
		//	subnet := "192.168.0"
		//	foundIPs := ScanNetwork(subnet)
		//	if len(foundIPs) == 0 {
		//		sshContent.SetText("未找到任何符合")
		//	} else {
		//		ipOptions = foundIPs
		//		ipSelect.Options = ipOptions
		//		ipSelect.Refresh()
		//		sshContent.SetText(fmt.Sprintf("找到 IP: \n%s", foundIPs))
		//	}
		//}()
	})

	connectButton := widget.NewButton("連線", func() {
		selectedIP := ipSelect.Selected
		if selectedIP == "" {
			sshContent.SetText("請先選擇一個IP地址")
			return
		}

		//go func() {
		//	sshContent.SetText(fmt.Sprintf("正在嘗試連線到 %s ...", selectedIP))
		//	if testSSHLogin(selectedIP) {
		//		sshContent.SetText(fmt.Sprintf("成功連線到 %s！", selectedIP))
		//	} else {
		//		sshContent.SetText(fmt.Sprintf("無法連線到 %s 遺失裝置請刷新IP...", selectedIP))
		//	}
		//}()

	})

	// 返回下拉選單作為標籤頁的內容
	row1 := container.NewHBox(
		widget.NewLabel("選擇一個 IP 地址："),
		ipSelect)

	row2 := container.NewHBox(
		refreshButton,
		connectButton)
	// 自動讀取ssh搜到內容
	row3 := container.NewHBox(
		widget.NewLabel("連線結果:"),
		container.NewMax(sshContent))
	// 輸入框
	row4 := container.NewHBox(
		widget.NewLabel("輸入指令:"),
		container.NewMax(commandInput))
	return container.NewVBox(row1, row2, row3, row4)
}
