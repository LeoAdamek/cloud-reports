package cmd

import (
	"github.com/urfave/cli"
	"reports/printer"
)

func testTable(c *cli.Context) {

	p := getPrinter(c)

	table := printer.Table{
		Boarders: false,
		Columns: []printer.TableColumn{
			{
				Title: "ID",
				Width: 40,
			},

			{
				Title: "Name",
				Width: 220,
			},

			{
				Title: "Value",
				Width: 220,
			},
		},

		Cells: [][]string{
			{"01","SomeBODY", "once told me"},
			{"02","That the","world is gonna roll me"},
			{"04","I ain't the", "sharpest tool in the shed"},
			{"05","She was looking", "kinda dumb"},
			{"06","With her finger and", "her thumb in"},
			{"07","The shape of an 'L'", "on her forehead."},
			{"08","Well...",""},
			{"09","The years start coming", "and they don't stop coming"},
			{"10","Fed to the wolves","and I hit the ground running."},
		},
	}

	p.DrawTable(table)
}
