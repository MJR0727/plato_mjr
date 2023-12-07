package client

import (
	"log"

	"github.com/jroimartin/gocui"
)

func main() {
	// 1、创建一个新的Gui对象
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	// 2、重置视图
	g.SetManagerFunc(layout)

	// 3、设置键位绑定
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	// 4、GUI界面循环
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	err := ViewHead(g, maxX/2-5, maxY/2, maxX/2+5, maxY/2+2)
	if err != nil {
		return err
	}
	return nil
}

func ViewHead(g *gocui.Gui, x0, y0, x1, y1 int) error {

	if v, err := g.SetView("hello", x0, y0, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Plato"
		v.Write([]byte("This is Plato.A IM system."))
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
