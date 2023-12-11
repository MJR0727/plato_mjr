package client

import (
	"fmt"
	"log"
	"strings"

	"github.com/CaoYnag/gocui"
	"github.com/MJR0727/plato_mjr/client/sdk"
	"github.com/gookit/color"
)

var (
	buf     string
	chat    *sdk.Chat
	step    int
	verbose bool
	pos     int //指向当前output识图的行
)

// 消息显示结构
type VOT struct {
	name, msg, sep string
}

func RunMain() {
	chat = sdk.NewChat("127.0.0.1:5666", "jack", "22249", "123321")

	// 1、创建一个新的Gui对象
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	// 2、重置视图
	g.SetManagerFunc(layout)

	// 3、设置键位绑定
	if err := g.SetKeybinding("input", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("input", gocui.KeyPgup, gocui.ModNone, viewUpScroll); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("input", gocui.KeyPgdn, gocui.ModNone, viewDownScroll); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("input", gocui.KeyArrowUp, gocui.ModNone, pasteUP); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("input", gocui.KeyArrowDown, gocui.ModNone, pasteDown); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, ViewUpdate); err != nil {
		log.Panicln(err)
	}

	go doRecv(g)

	// 4、GUI界面循环
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if err := ViewHead(g, 1, 1, maxX-1, 4); err != nil {
		return err
	}
	if err := ViewOutput(g, 1, 5, maxX-1, maxY-4); err != nil {
		return err
	}
	if err := ViewInput(g, 1, maxY-3, maxX-1, maxY-1); err != nil {
		return err
	}
	return nil
}

func ViewHead(g *gocui.Gui, x0, y0, x1, y1 int) error {

	if v, err := g.SetView("head", x0, y0, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = false
		v.Overwrite = true
		v.Clear()
		v.Write([]byte("This is Plato.A IM system.\n可以开始聊天了!"))
	}
	return nil
}

func ViewOutput(g *gocui.Gui, x0, y0, x1, y1 int) error {

	if v, err := g.SetView("output", x0, y0, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Overwrite = true
		v.Autoscroll = true
		v.SelBgColor = gocui.ColorGreen
		v.Title = "Message"
	}
	return nil
}

func ViewInput(g *gocui.Gui, x0, y0, x1, y1 int) error {

	if v, err := g.SetView("input", x0, y0, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Overwrite = false
		v.Editable = true
		v.SelBgColor = gocui.ColorBlue
		// 焦点转移
		if _, err := g.SetCurrentView("input"); err != nil {
			return err
		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	chat.Close()
	ov, _ := g.View("output")
	buf = ov.Buffer()
	g.Close()
	return gocui.ErrQuit
}

func viewUpScroll(g *gocui.Gui, cv *gocui.View) error {
	v, err := g.View("output")
	v.Autoscroll = false
	ox, oy := v.Origin()
	if err == nil {
		v.SetOrigin(ox, oy-1)
	}
	return err
}

func viewDownScroll(g *gocui.Gui, cv *gocui.View) error {
	v, err := g.View("output")
	// 视图缓冲行数，可视范围内的行数
	_, y := v.Size()
	lnum := len(v.BufferLines())
	// 原点坐标
	ox, oy := v.Origin()
	if err == nil {
		if oy < lnum-y-1 {
			v.Autoscroll = true
		} else {
			v.SetOrigin(ox, oy+1)
		}

	}
	return err
}

func pasteUP(g *gocui.Gui, cv *gocui.View) error {
	v, err := g.View("output")
	if err != nil {
		fmt.Fprintf(cv, "error:%s", err)
		return nil
	}
	bls := v.BufferLines()
	lnum := len(bls)
	if pos < lnum-1 {
		pos++
	}
	cv.Clear()
	strArr := strings.Split(bls[lnum-pos-1], ":")
	fmt.Fprintf(cv, "%s", strArr[len(strArr)-1])
	// fmt.Fprintf(cv, "%s", bls[lnum-pos-1])
	return nil
}

func pasteDown(g *gocui.Gui, cv *gocui.View) error {
	v, err := g.View("output")
	if err != nil {
		fmt.Fprintf(cv, "error:%s", err)
		return nil
	}
	if pos > 0 {
		pos--
	}
	bls := v.BufferLines()
	lnum := len(bls)
	cv.Clear()
	strArr := strings.Split(bls[lnum-pos-1], ":")
	fmt.Fprintf(cv, "%s", strArr[len(strArr)-1])
	// fmt.Fprintf(cv, "%s", bls[lnum-pos-1])
	return nil
}

// cv：input视图
func ViewUpdate(g *gocui.Gui, cv *gocui.View) error {
	// 更新视图并发送消息
	doSay(g, cv)
	l := len(cv.Buffer())
	cv.MoveCursor(0-l, 0, true)
	cv.Clear()
	return nil
}

// cv:input视图对象指针
func doSay(g *gocui.Gui, cv *gocui.View) error {
	// 1、获取input Editor对象的字符串
	b := make([]byte, 300)
	var content string
	if cv != nil {
		n, err := cv.Read(b)
		if n > 0 {
			content = string(b[:n])
		} else {
			return err
		}
	}

	// 2、更新output视图
	outputV, err := g.View("output")
	if err == nil {
		viewPrint(g, "me", content, false)
		outputV.Autoscroll = true

		// 3、将消息进行发送
		message := &sdk.Message{
			Type:       sdk.MsgType_Text,
			Name:       "mjr",
			ToUserId:   "22222",
			FromUserID: "22249",
			Content:    content,
			Session:    "123321",
		}
		chat.Send(message)
	} else {
		return err
	}
	return nil
}

func viewPrint(g *gocui.Gui, name, content string, newLine bool) {
	out := VOT{
		name: name,
		msg:  content,
	}
	if newLine {
		out.sep = "\n"
	} else {
		out.sep = " "
	}
	g.Update(out.Show)
}

func (out VOT) Show(g *gocui.Gui) error {
	v, err := g.View("output")
	if err == nil {
		// 将消息打印到视图中
		fmt.Fprintf(v, "%v:%v%v", color.FgGreen.Text(out.name), out.sep, color.FgYellow.Text(out.msg))
	} else {
		return nil
	}

	return nil
}

// 轮询接受消息
func doRecv(g *gocui.Gui) {
	recvChannel := chat.Recv()
	for msg := range recvChannel {
		// 将msg更新到output视图
		switch msg.Type {
		case sdk.MsgType_Text:
			viewPrint(g, msg.Name, msg.Content, false)
		}
	}
	g.Close()
}
