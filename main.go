package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

// const (
// 	programExe = `M:\ZaprUI\ZaprUiApp\zaprUI\build\bin\OUTPUT.exe`
// )

func main() {

	// wg := &sync.WaitGroup{}

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	if !isAdmin() {
	// 		rerunAsAdmin()
	// 	}

	// 	runProgram()
	// }()

	// wg.Wait()

	// запускается main ->

	// ===========

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "zaprLauncher",
		Width:  200,
		Height: 250,

		MinWidth:  200,
		MaxWidth:  200,
		MinHeight: 250,
		MaxHeight: 250,
		Frameless: true,
		CSSDragProperty: "widows",
		CSSDragValue:    "1",

		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}

}

// ======================================================

// func isAdmin() bool {
// 	h, err := syscall.GetCurrentProcess()
// 	if err != nil {
// 		return false
// 	}

// 	var token syscall.Token
// 	err = syscall.OpenProcessToken(h, syscall.TOKEN_QUERY, &token)
// 	if err != nil {
// 		return false
// 	}
// 	defer token.Close()

// 	var elevation uint32
// 	var outLen uint32

// 	err = syscall.GetTokenInformation(
// 		token,
// 		syscall.TokenElevation,
// 		(*byte)(unsafe.Pointer(&elevation)),
// 		uint32(unsafe.Sizeof(elevation)),
// 		&outLen,
// 	)
// 	return err == nil && elevation != 0
// }

// func rerunAsAdmin() {

// 	verb, _ := syscall.UTF16PtrFromString("runas")
// 	file, _ := syscall.UTF16PtrFromString(programExe)

// 	shell32 := syscall.NewLazyDLL("shell32.dll")
// 	shellExecute := shell32.NewProc("ShellExecuteW")

// 	shellExecute.Call(
// 		0,
// 		uintptr(unsafe.Pointer(verb)),
// 		uintptr(unsafe.Pointer(file)),
// 		0,
// 		0,
// 		1,
// 	)
// }

// func runProgram() {
// 	cmd := exec.Command(programExe)
// 	cmd.Start()
// }
