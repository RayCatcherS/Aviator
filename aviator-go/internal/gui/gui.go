package gui

import (
	"aviator/internal/config"
	"aviator/internal/server"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	qrcode "github.com/skip2/go-qrcode"
)

type GUI struct {
	mw          *walk.MainWindow
	ni          *walk.NotifyIcon
	cfg         *config.ConfigManager
	srv         *server.Server
	appURL      string
	
	// Models
	appModel    *AppModel
	
	// Widgets
	table       *walk.TableView
	qrImage     *walk.ImageView
	statusLabel *walk.Label
	startBtn    *walk.PushButton
	stopBtn     *walk.PushButton
	
	// State
	serverRunning bool
}

type AppModel struct {
	walk.TableModelBase
	items []config.App
}

func (m *AppModel) RowCount() int {
	return len(m.items)
}

func (m *AppModel) Value(row, col int) interface{} {
	item := m.items[row]
	switch col {
	case 0:
		return item.Name
	case 1:
		return item.Path
	case 2:
		return item.Args
	}
	return ""
}

func NewGUI(defaultUrl string, cfg *config.ConfigManager, srv *server.Server) (*GUI, error) {
	localIP := GetOutboundIP()
	finalURL := fmt.Sprintf("http://%s:8000", localIP)

	g := &GUI{
		cfg:           cfg,
		srv:           srv,
		appURL:        finalURL,
		appModel:      &AppModel{items: cfg.GetApps()},
		serverRunning: true, // Assuming it starts running
	}

	// Generate QR
	var qrBitmap *walk.Bitmap
	// Create a temp file for QR to ensure Walk can load it reliably
	tmpQR := filepath.Join(os.TempDir(), "aviator_qr.png")
	if err := qrcode.WriteFile(finalURL, qrcode.Medium, 256, tmpQR); err == nil {
		if img, err := walk.NewBitmapFromFile(tmpQR); err == nil {
			qrBitmap = img
		}
	}

	if err := (MainWindow{
		AssignTo: &g.mw,
		Title:    "Aviator Server",
		Size:     Size{Width: 650, Height: 700},
		Layout:   VBox{},
		Children: []Widget{
			// Top Section: 3 Columns
			Composite{
				Layout: HBox{Alignment: AlignHNearVCenter},
				MaxSize: Size{Height: 180}, // Fixed height
				Children: []Widget{
					// Column 1: Status (Left)
					Composite{
						Layout: VBox{Alignment: AlignHNearVNear},
						Children: []Widget{
							Label{
								Text: "Server Status:",
								Font: Font{PointSize: 10, Bold: true},
							},
							Label{
								AssignTo: &g.statusLabel,
								Text:     "RUNNING",
								Font:     Font{PointSize: 12, Bold: true},
								TextColor: walk.RGB(0, 128, 0),
							},
							VSpacer{Size: 10},
							PushButton{
								AssignTo: &g.startBtn,
								Text:     "Start Server",
								Enabled:  false,
								OnClicked: func() {
									g.startServer()
								},
							},
							PushButton{
								AssignTo: &g.stopBtn,
								Text:     "Stop Server",
								Enabled:  true,
								OnClicked: func() {
									g.stopServer()
								},
							},
						},
					},
					
					HSpacer{Size: 20}, // Spacer

					// Column 2: URLs (Center)
					Composite{
						Layout: VBox{Alignment: AlignHNearVNear},
						Children: []Widget{
							Label{Text: "Local URL:"},
							LinkLabel{
								Text: "<a href=\"http://localhost:8000\">http://localhost:8000</a>",
								OnLinkActivated: func(link *walk.LinkLabelLink) {
									openBrowser("http://localhost:8000")
								},
							},
							VSpacer{Size: 5},
							Label{Text: "Network URL:"},
							LinkLabel{
								Text: fmt.Sprintf("<a href=\"%s\">%s</a>", finalURL, finalURL),
								OnLinkActivated: func(link *walk.LinkLabelLink) {
									openBrowser(finalURL)
								},
							},
						},
					},

					HSpacer{Size: 20}, // Spacer

					// Column 3: QR Code (Right)
					ImageView{
						AssignTo: &g.qrImage,
						Image:    qrBitmap,
						MinSize:  Size{Width: 160, Height: 160},
						MaxSize:  Size{Width: 160, Height: 160},
						Mode:     ImageViewModeZoom,
					},
				},
			},

			VSpacer{Size: 10},
			Label{Text: "Applications:", Font: Font{Bold: true}},

			// App List
			TableView{
				AssignTo:         &g.table,
				AlternatingRowBG: true,
				// This allows the table to take all available space
				VerticalStretch: 1,
				Columns: []TableViewColumn{
					{Title: "Name", Width: 150},
					{Title: "Path", Width: 300},
					{Title: "Args", Width: 100},
				},
				Model: g.appModel,
			},

			// App Controls
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						Text: "Add App...",
						OnClicked: func() {
							g.appDialog(nil)
						},
					},
					PushButton{
						Text: "Edit App...",
						OnClicked: func() {
							idx := g.table.CurrentIndex()
							if idx >= 0 {
								app := g.appModel.items[idx]
								g.appDialog(&app)
							}
						},
					},
					PushButton{
						Text: "Remove",
						OnClicked: func() {
							idx := g.table.CurrentIndex()
							if idx >= 0 {
								app := g.appModel.items[idx]
								cfg.RemoveApp(app.ID)
								g.refreshList()
							}
						},
					},
				},
			},
		},
	}.Create()); err != nil {
		return nil, err
	}

	// Minimizing instead of closing on X
	g.mw.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
		*canceled = true
		g.mw.SetVisible(false)
		g.showTrayMessage("Aviator Minimized", "The server is still running in the background.")
	})

	// Tray Icon
	ni, err := walk.NewNotifyIcon(g.mw)
	if err == nil {
		g.ni = ni
		ni.SetVisible(true)
		if icon, err := walk.Resources.Icon("O"); err == nil {
			ni.SetIcon(icon)
		}

		ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
			if button == walk.LeftButton {
				g.mw.SetVisible(true)
				// Activate/BringToFront
				if win := g.mw; win != nil {
					// win.BringToFront() // If checking
				}
			}
		})

		// Context Menu
		showAction := walk.NewAction()
		showAction.SetText("Show Interface")
		showAction.Triggered().Attach(func() { 
			g.mw.SetVisible(true) 
		})
		ni.ContextMenu().Actions().Add(showAction)

		stopAction := walk.NewAction()
		stopAction.SetText("Stop Server")
		stopAction.Triggered().Attach(func() { g.stopServer() })
		ni.ContextMenu().Actions().Add(stopAction)

		exitAction := walk.NewAction()
		exitAction.SetText("Exit Aviator")
		exitAction.Triggered().Attach(func() { g.exitApp() })
		ni.ContextMenu().Actions().Add(exitAction)
	}

	g.mw.SetVisible(true)
	return g, nil
}

func (g *GUI) Run() {
	g.mw.Run()
}

func (g *GUI) showTrayMessage(title, msg string) {
	if g.ni != nil {
		g.ni.ShowInfo(title, msg)
	}
}

func (g *GUI) startServer() {
	if g.serverRunning {
		return
	}
	
	go func() {
		if err := g.srv.Start(8000); err != nil {
			log.Printf("Server start error: %v", err)
		}
	}()
	
	g.serverRunning = true
	g.updateStatusUI()
}

func (g *GUI) stopServer() {
	// Confirmation? Using MsgBox
	if walk.MsgBox(g.mw, "Stop Server", "Are you sure you want to stop the server?", walk.MsgBoxYesNo|walk.MsgBoxIconQuestion) == walk.DlgCmdNo {
		return
	}

	if g.srv != nil {
		g.srv.Stop()
	}
	g.serverRunning = false
	g.updateStatusUI()
}

func (g *GUI) exitApp() {
	if walk.MsgBox(g.mw, "Exit Aviator", "Are you sure you want to quit Aviator?\nThis will stop the server.", walk.MsgBoxYesNo|walk.MsgBoxIconQuestion) == walk.DlgCmdYes {
		g.srv.Stop()
		g.ni.Dispose()
		walk.App().Exit(0)
	}
}

func (g *GUI) updateStatusUI() {
	// Must run on main thread
	g.mw.Synchronize(func() {
		if g.serverRunning {
			g.statusLabel.SetText("RUNNING")
			g.statusLabel.SetTextColor(walk.RGB(0, 128, 0))
			g.startBtn.SetEnabled(false)
			g.stopBtn.SetEnabled(true)
		} else {
			g.statusLabel.SetText("STOPPED")
			g.statusLabel.SetTextColor(walk.RGB(255, 0, 0))
			g.startBtn.SetEnabled(true)
			g.stopBtn.SetEnabled(false)
		}
	})
}

func (g *GUI) refreshList() {
	g.appModel.items = g.cfg.GetApps()
	g.table.SetModel(g.appModel)
}

// Unified Dialog for Adding/Editing App
func (g *GUI) appDialog(existingApp *config.App) {
	var name, path, args string
	var dlg *walk.Dialog
	var nameEdit, pathEdit, argsEdit *walk.LineEdit
	title := "Add Application"
	
	if existingApp != nil {
		title = "Edit Application"
		name = existingApp.Name
		path = existingApp.Path
		args = existingApp.Args
	}

	// Pre-fill path logic
	selectPath := func() {
		dlgFile := new(walk.FileDialog)
		dlgFile.Title = "Select Executable"
		dlgFile.Filter = "Executables (*.exe)|*.exe|All Files (*.*)|*.*"
		if ok, _ := dlgFile.ShowOpen(g.mw); ok {
			pathEdit.SetText(dlgFile.FilePath)
			if nameEdit.Text() == "" {
				base := filepath.Base(dlgFile.FilePath)
				ext := filepath.Ext(base)
				nameEdit.SetText(strings.TrimSuffix(base, ext))
			}
		}
	}

	Dialog{
		AssignTo: &dlg,
		Title:    title,
		MinSize:  Size{Width: 400, Height: 200},
		Layout:   VBox{},
		Children: []Widget{
			Label{Text: "Executable Path:"},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					LineEdit{AssignTo: &pathEdit, Text: path},
					PushButton{
						Text:      "...",
						OnClicked: selectPath,
						MaxSize:   Size{Width: 30},
					},
				},
			},
			Label{Text: "Application Name (Display):"},
			LineEdit{AssignTo: &nameEdit, Text: name},
			Label{Text: "Arguments (Optional):"},
			LineEdit{AssignTo: &argsEdit, Text: args},
			VSpacer{},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						Text: "Cancel",
						OnClicked: func() { dlg.Cancel() },
					},
					PushButton{
						Text: "Save",
						OnClicked: func() {
							path = pathEdit.Text()
							name = nameEdit.Text()
							args = argsEdit.Text()
							if path != "" && name != "" {
								dlg.Accept()
							}
						},
					},
				},
			},
		},
	}.Run(g.mw)

	if dlg.Result() == 1 { // Accepted
		if existingApp != nil {
			// Update logic: remove old, add new (or just update fields if we had ID)
			// Since ConfigManager uses ID, we can update.
			g.cfg.UpdateApp(existingApp.ID, name, path, args)
		} else {
			g.cfg.AddApp(name, path, args)
		}
		g.refreshList()
	}
}

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = exec.Command("cmd", "/c", "start", url).Start()
	}
	if err != nil {
		log.Printf("Error opening browser: %v", err)
	}
}
