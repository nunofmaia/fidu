package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/nunofmaia/fidu/marker"
)

var (
	Version = "0.0.1"
	Build   = "gobuild"
)

func main() {
	app := cli.NewApp()
	app.Author = "Nuno Maia"
	app.Email = "nunofmaia@gmail.com"
	app.Name = "fidu"
	app.Usage = "2D fiducial markers generator for MultiTaction displays"
	app.Version = fmt.Sprintf("%s (%s)", Version, Build)
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "code, c",
			Value: 0,
			Usage: "code value to generate, a 32-bit integer value",
		},
		cli.IntFlag{
			Name:  "blocksize, b",
			Value: 32,
			Usage: "number of pixels per block",
		},
		cli.IntFlag{
			Name:  "division, d",
			Value: 5,
			Usage: "the number of blocks in X/Y directions, range from 3 to 8",
		},
		cli.StringFlag{
			Name:  "filename, o",
			Value: "",
			Usage: "the name of the output file, default 'code-[number].png'",
		},
		cli.BoolFlag{
			Name:  "no-border",
			Usage: "make the border of the marker transparent",
		},
	}
	app.Action = func(c *cli.Context) {
		code := c.Int("code")
		blocksize := c.Int("blocksize")
		division := c.Int("division")
		filename := c.String("filename")
		hasBorder := !c.Bool("no-border")

		mkr := marker.New(code, division, blocksize, filename, hasBorder)
		if err := mkr.Save(); err != nil {
			log.Fatal(err)
		}

		println(fmt.Sprintf("Saved code to %s", mkr.Name))
	}

	app.Run(os.Args)
}
