package main

import (
	"log"
	"os"
	"pbed/bed"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/flopp/go-findfont"
)

var (
	a            fyne.App
	w            fyne.Window
	listTree     *widget.List
	infoTextBind binding.String
	strListBind  binding.StringList
)

var (
	xb         bed.Bed
	uploadPath string
)

func init() {
	initFont()
}

func startGUI(bed bed.Bed) {
	xb = bed

	a = app.New()
	w = a.NewWindow("kazma图床")
	w.CenterOnScreen()

	w.Resize(fyne.NewSize(621, 552))

	infoTextBind = binding.NewString()
	_ = infoTextBind.Set("信息面板")

	strListBind = binding.NewStringList()
	listTree = widget.NewListWithData(strListBind, func() fyne.CanvasObject {
		return widget.NewEntry()
	}, func(item binding.DataItem, object fyne.CanvasObject) {
		val, err := item.(binding.String).Get()
		if err != nil {
			log.Printf("refresh list error: %v", err)
			return
		}

		object.(*widget.Entry).SetText(val)
	})

	w.SetContent(container.NewVBox(
		// 选择文件
		container.NewGridWithColumns(2,
			widget.NewButton("选择文件", func() {
				fileDialog := dialog.NewFileOpen(func(u fyne.URIReadCloser, e error) {
					if e != nil {
						log.Printf("选取文件[夹]异常: %v", e.Error())
					} else if u != nil {
						fileDialogCallback(u.URI().Path())
					}
				}, w)

				fileDialog.Resize(w.Canvas().Size())
				fileDialog.Show()
			}),
			widget.NewButton("选择目录", func() {
				fileDialog := dialog.NewFolderOpen(func(lu fyne.ListableURI, e error) {
					if e != nil {
						log.Printf("选取目录异常: %v", e.Error())
					} else if lu != nil {
						dirDialogCallback(lu.Path())
					}
				}, w)

				fileDialog.Resize(w.Canvas().Size())
				fileDialog.Show()
			}),
		),
		widget.NewLabelWithData(infoTextBind),
		container.NewGridWithColumns(2,
			widget.NewSelect(BedType, func(s string) {

			}),
			widget.NewButton("上传", uploadAction),
		),
		container.NewGridWithRows(1, listTree),
	))

	w.SetCloseIntercept(func() {
		w.Close()
		os.Exit(0)
	})

	w.ShowAndRun()

	os.Exit(0)
}

func uploadAction() {
	if uploadPath == "" {
		return
	}

	up := uploadPath

	info, err := os.Stat(up)
	if err != nil {
		log.Printf("打开文件/目录失败: %v", err)
		return
	}

	if info.IsDir() {
		for _, p := range uploadDir(xb, up) {
			_ = strListBind.Append(p)
		}
	} else {
		p := upload(xb, up)
		if p != "" {
			_ = strListBind.Append(p)
		}
	}

	listTree.Resize(fyne.NewSize(listTree.Size().Width, 360))
}

func fileDialogCallback(up string) {
	log.Printf("upload file path: %s", up)

	if up == "" {
		return
	}
	uploadPath = up
	_ = infoTextBind.Set("选中文件: " + up)
}

func dirDialogCallback(up string) {
	log.Printf("upload file path: %s", up)

	if up == "" {
		return
	}

	uploadPath = up
	_ = infoTextBind.Set("选中目录: " + up)
}

func initFont() {

	for _, fn := range findfont.List() {
		lowerFn := strings.ToLower(fn)
		if strings.Contains(lowerFn, "sarasa") {
			log.Printf("find font: %s", fn)
			err := os.Setenv("FYNE_FONT", fn)
			if err != nil {
				log.Printf("set env FYNE_FONT error: %v", err)
			}

			break
		}
	}
}
