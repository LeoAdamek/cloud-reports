package cmd

import (
	"github.com/urfave/cli"
	"log"
	"net"
	"os"
	"reports/printer"
	"time"
)

func App() *cli.App {
	app := cli.NewApp()

	app.Before = func(c *cli.Context) error {
		log.SetOutput(os.Stdout)
		log.SetFlags(log.LstdFlags)
		log.SetPrefix("reports ")

		return nil
	}

	app.Commands = []cli.Command{
		{
			Name: "cpu",
			Usage: "Print CPU Chart",
			Action: cpuChart,

			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name: "group",
				},


				cli.DurationFlag{
					Name: "d",
					Value: 24*time.Hour,
				},
			},
		},

		{
			Name: "memory",
			Usage: "Memory",
			Action: memoryChart,
		},

		{
			Name: "tests",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name :"d",
					Usage: "Set this to print data to STDOUT instead of to network",
				},
			},
			Subcommands: []cli.Command {
				{
					Name: "chart",
					Action: testChart,
				},
			},
		},
	}
	
	return app
}


func getPrinter(c *cli.Context) printer.Printer {
	
	if c.Bool("d") {
		return printer.New(os.Stdout)
	}
	
	conn, err := net.Dial("tcp", "10.1.0.4:30000")
	
	if err != nil {
		log.Fatal(err)
	}
	
	p := printer.New(conn)
	
	p.Debug = true
	
	return p
}