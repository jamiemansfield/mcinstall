// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"errors"
	"log"
	"os"

	"github.com/jamiemansfield/go-technic/platform"
	"github.com/jamiemansfield/mcinstall/technic"
	"github.com/jamiemansfield/mcinstall/util"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "technicinstall",
		Usage:   "install packs from the Technic Pack",
		Version: "0.1.0-indev",
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() < 2 {
				return errors.New("usage: technicinstall pack version")
			}
			packSlug := ctx.Args().Get(0)
			version := ctx.Args().Get(1)

			client := platform.NewClient(nil)
			client.UserAgent = util.UserAgent
			client.Build = "mcinstall"

			pack, err := client.Modpack.GetModpack(packSlug)
			if err != nil {
				return err
			}

			return technic.InstallPackVersion("", pack, version)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
