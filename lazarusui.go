package main

import (
	// Do not reference this package if you use custom syso files

	"path/filepath"
	"strings"

	_ "github.com/ying32/govcl/pkgs/winappres"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

type TMainForm struct {
	*vcl.TForm
	Btn1         *vcl.TButton
	DlgOpen      *vcl.TOpenDialog
	SpinEdit     *vcl.TSpinEdit
	SpinEditSize *vcl.TSpinEdit
	Label        *vcl.TLabel
	LabelSize    *vcl.TLabel
	ProgressBar  *vcl.TProgressBar
	Memo         *vcl.TMemo
}

var (
	mainForm *TMainForm
)

func main() {
	vcl.RunApp(&mainForm)
}

// -- TMainForm

func (f *TMainForm) OnFormCreate(sender vcl.IObject) {
	f.SetCaption("Упаковщик нескольких страниц pdf в одну")
	f.SetWidth(500)
	f.SetHeight(265)
	f.SetPosition(types.PoScreenCenter)
	f.EnabledMaximize(false)
	f.SetBorderStyle(types.BsSingle)

	f.SetDoubleBuffered(true)

	f.Btn1 = vcl.NewButton(f)
	f.Btn1.SetParent(f)
	f.Btn1.SetBounds(10, 80, 480, 28)
	f.Btn1.SetCaption("Выбрать pdf для преобразования")
	f.Btn1.SetOnClick(f.OnBtn1Click)
	f.DlgOpen = vcl.NewOpenDialog(f)
	f.DlgOpen.SetFilter("PDF(*.pdf)|*.pdf")

	f.SpinEdit = vcl.NewSpinEdit(f)
	f.SpinEdit.SetParent(f)
	f.SpinEdit.SetIncrement(1)
	f.SpinEdit.SetValue(4)
	f.SpinEdit.SetMaxValue(300)
	f.SpinEdit.SetMinValue(2)
	f.SpinEdit.SetBounds(210, 47, 50, 20)

	f.SpinEditSize = vcl.NewSpinEdit(f)
	f.SpinEditSize.SetParent(f)
	f.SpinEditSize.SetIncrement(1)
	f.SpinEditSize.SetValue(594)
	f.SpinEditSize.SetMinValue(20)
	f.SpinEditSize.SetMaxValue(15000)
	f.SpinEditSize.SetBounds(210, 17, 50, 20)

	f.Label = vcl.NewLabel(f)
	f.Label.SetParent(f)
	f.Label.SetBounds(10, 50, 50, 5)
	f.Label.SetCaption("Количество листов на новый лист")

	f.LabelSize = vcl.NewLabel(f)
	f.LabelSize.SetParent(f)
	f.LabelSize.SetBounds(10, 20, 50, 5)
	f.LabelSize.SetCaption("Высота полосы в мм")

	f.ProgressBar = vcl.NewProgressBar(f)
	f.ProgressBar.SetParent(f)
	f.ProgressBar.SetBounds(10, 120, 480, 5)

	f.Memo = vcl.NewMemo(f)
	f.Memo.SetParent(f)
	f.Memo.SetBounds(10, 150, 480, 100)
	f.Memo.SetEnabled(false)
	f.Memo.SetScrollBars(2)

}

func (f *TMainForm) OnBtn1Click(sender vcl.IObject) {
	var filename, newfile string
	if f.DlgOpen.Execute() {
		filename = strings.TrimSuffix(filepath.Base(f.DlgOpen.FileName()), filepath.Ext(f.DlgOpen.FileName()))
	}

	rand, err := GenerateRandomString(6)
	if err != nil {
		return
	}

	if filename != "" {
		f.Memo.SetText("Преобразую " + filepath.Base(f.DlgOpen.FileName()) + "...")
		newfile = filepath.Join(filepath.Dir(f.DlgOpen.FileName()), filename+"-"+rand+filepath.Ext(f.DlgOpen.FileName()))

		f.Btn1.SetEnabled(false)
		f.SpinEdit.SetEnabled(false)
		f.SpinEditSize.SetEnabled(false)
		f.Refresh()

		err = repack(f, f.DlgOpen.FileName(), newfile, int(f.SpinEdit.Value()), 594)
		if err != nil {
			f.Memo.SetText(err.Error())
			f.Btn1.SetEnabled(true)
			return
		}
		f.ProgressBar.SetPosition(100)
		f.Memo.SetText(f.Memo.Text() + "\n" + "Готово! Там же файл " + filename + "-" + rand + filepath.Ext(f.DlgOpen.FileName()))
		MemoScrollDown(f.Memo)

		f.Btn1.SetEnabled(true)
		f.SpinEdit.SetEnabled(true)
		f.SpinEditSize.SetEnabled(true)
		f.Refresh()
	}
}
