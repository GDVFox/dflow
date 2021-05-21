package about

import "github.com/pterm/pterm"

// HandleAbout выводит сообщение с помощью.
func HandleAbout() {
	title, _ := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("D", pterm.NewStyle(pterm.FgLightMagenta)),
		pterm.NewLettersFromStringWithStyle("Flow", pterm.NewStyle(pterm.FgCyan))).
		Srender()

	pterm.DefaultCenter.Println(title)
	pterm.DefaultCenter.WithCenterEachLineSeparately().Println(
		"Distributed dataflow managing tool\n" +
			"by Daniil Gavrilovsky.\n" +
			"GitHub repo: 'https://github.com/GDVFox/dflow'\n" +
			"DFlow is licensed under MIT License.\n" +
			"2021")
}
