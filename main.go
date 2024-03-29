package main

import (
	"embed"
	"fmt"
	"sliggooey/backend"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed frontend/dist
var assets embed.FS

func goPrint(a ...any) {
	// convenience method to distinguish backend logs from frontend logs
	s := make([]any, 1)
	s[0] = "[GO]"
	for _, e := range a {
		s = append(s, e)
	}
	fmt.Println(s...)
}

func main() {
	// Create an instance of the app structure
	app := backend.NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "gooey",
		MaxWidth:  2160,
		MaxHeight: 1440,
		Width:     1400,
		Height:    890,
		Assets:    assets,
		OnStartup: app.Startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		goPrint("Error:", err)
	}
}
