package main

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"os"

	pc "github.com/InfinityTools/go-binpack2d"
	"github.com/oliverpool/unipdf/v3/creator"
	pdf "github.com/oliverpool/unipdf/v3/model"
	"github.com/oliverpool/unipdf/v3/model/optimize"
	"github.com/ying32/govcl/vcl"
)

// ------------------------------------------------------------------------------------
//  основано на https://github.com/unidoc/unipdf-examples/blob/v3/pages/pdf_4up.go
//  основано на https://github.com/unidoc/unipdf-examples/blob/master/pages/pdf_rotate_flatten.go
// ------------------------------------------------------------------------------------

func repack(f *TMainForm, inpath string, outpath string, divto int, polHeight int) error {
	if _, err := os.Stat(inpath); err != nil {
		return err
	}

	data, err := os.Open(inpath)
	if err != nil {
		return err
	}
	defer data.Close()

	//читаем
	pdfReader, err := pdf.NewPdfReader(data)
	if err != nil {
		return err
	}

	//Получаем количество страниц в файле
	colPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	divpages := divto
	rds := float64(colPages) / float64(divpages)
	rounds := int(math.Ceil(rds))

	currPage := 1

	//Создаем компановщик pdf
	c := creator.New()
	var currentPage *pdf.PdfPage

	for i := 0; i < rounds; i++ {
		packer := pc.Create(80000, polHeight)
		var pagesAlign []pc.Rectangle
		var blocks []*creator.Block

		for j := 0; j < divpages; j++ {
			if currPage+j <= colPages {

				currentPage, err = pdfReader.GetPage(currPage + j)

				var rotateDeg int64
				if currentPage.Rotate != nil && *currentPage.Rotate != 0 {
					rotateDeg = -*currentPage.Rotate
				}

				block, err := creator.NewBlockFromPage(currentPage)
				if err != nil {
					return err
				}
				if rotateDeg != 0 {
					block.SetAngle(float64(rotateDeg))
				}
				w, h := block.RotatedSize()

				blocks = append(blocks, block)

				t := fmt.Sprint("Новая страница ", i+1, " <- старая страница ", currPage+j)
				fmt.Println(t)
				f.Memo.SetText(f.Memo.Text() + "\n" + t)
				MemoScrollDown(f.Memo)
				f.Refresh()

				f.ProgressBar.SetPosition(int32(math.Ceil(((float64(currPage + j)) * 100 / float64(colPages)))))

				rect, ok := packer.Insert(int(math.Round(Px2mm(w))), int(math.Round(Px2mm(h))), 2)
				if ok {
					//fmt.Printf("Accepted: %v\n", rect)
					pagesAlign = append(pagesAlign, rect)
				} else {
					fmt.Printf("Rejected: %v\n", rect)
					f.Memo.SetText(f.Memo.Text() + "\n" + "Ошибка, страница не упаковывается: " + string(currPage))
					MemoScrollDown(f.Memo)
					f.Refresh()
				}
			}
		}
		packer.ShrinkBin(false)
		//fmt.Println("packer", packer.GetHeight(), packer.GetWidth(), packer.GetOccupancy())

		c.SetPageSize(creator.PageSize{Mm2px(float64(packer.GetWidth())), Mm2px(float64(packer.GetHeight()))})
		c.NewPage()

		for j := 0; j < len(pagesAlign); j++ {
			if currPage+j <= colPages {
				if err != nil {
					return err
				}
				rect := pagesAlign[j]
				grect := c.NewRectangle(float64(Mm2px(float64(rect.X))), float64(Mm2px(float64(rect.Y))), float64(Mm2px(float64(rect.W))), float64(Mm2px(float64(rect.H))))
				grect.SetBorderWidth(Mm2px(0.1))
				grect.SetBorderColor(creator.ColorBlack)

				block := blocks[j]
				w, h := block.RotatedSize()
				block.SetPos((w-block.Width())/2+float64(Mm2px(float64(rect.X))), (h-block.Height())/2+float64(Mm2px(float64(rect.Y))))

				err = c.Draw(block)
				if err != nil {
					return err
				}

				err = c.Draw(grect)
				if err != nil {
					return err
				}
			}
		}
		currPage += divpages
	}
	f.Memo.SetText(f.Memo.Text() + "\n" + "Сохраняю результат в " + outpath + ", ждите..")
	MemoScrollDown(f.Memo)
	f.Refresh()
	fmt.Println("Сохраняю результат в " + outpath)

	c.SetOptimizer(optimize.New(optimize.Options{
		CombineDuplicateDirectObjects:   true,
		CombineIdenticalIndirectObjects: true,
		CombineDuplicateStreams:         true,
		CompressStreams:                 true,
		CleanFonts:                      true,
		SubsetFonts:                     true,
	}))
	c.Finalize()
	err = c.WriteToFile(outpath)

	if err != nil {
		fmt.Println("Ошибка " + err.Error())
		return err
	}
	fmt.Print("Успешно!")
	return nil
}

func Px2mm(px float64) float64 {
	return px * 0.35277777777
}

func Mm2px(mm float64) float64 {
	return mm / 0.35277777777
}

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		// rand.Int generates a cryptographically secure random integer
		// in the range [0, max).
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()] // Use Int66 as big.NewInt(int64) returns Int66
	}
	return string(ret), nil
}

func MemoScrollDown(m *vcl.TMemo) {
	m.SetSelStart(m.GetTextLen() - 1)
	m.SetSelLength(1)
	m.ClearSelection()
}
