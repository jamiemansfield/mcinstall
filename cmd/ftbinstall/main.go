// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"errors"
	"fmt"
	"github.com/gosimple/slug"
	"github.com/jamiemansfield/ftbinstall/ftb"
	"github.com/jamiemansfield/ftbinstall/minecraft"
	"github.com/jamiemansfield/go-ftbmeta/ftbmeta"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name: "ftbinstall",
		Usage: "install packs from the modpacks.ch service",
		Version: "0.1.0-indev",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "target",
				Aliases:  []string{"t"},
				Usage:    "sets the install target",
				Value:    "client",
			},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() < 2 {
				return errors.New("usage: ftbinstall pack version")
			}
			packSlug := ctx.Args().Get(0)
			versionSlug := slug.MakeLang(ctx.Args().Get(1), "en")
			installTargetRaw := ctx.Value("target").(string)
			var installTarget minecraft.InstallTarget
			if installTargetRaw == "client" || installTargetRaw == "c" {
				installTarget = minecraft.Client
			} else
			if installTargetRaw == "server" || installTargetRaw == "s" {
				installTarget = minecraft.Server
			} else {
				return errors.New("unknown install target "+ installTargetRaw)
			}

			client := ftbmeta.NewClient(nil)

			pack, err := client.Packs.GetPack(packSlug)
			if err != nil {
				return err
			}

			version, err := client.Packs.GetVersion(packSlug, versionSlug)
			if err != nil {
				return err
			}

			fmt.Println("Installing " + pack.Name + " v" + version.Name + "...")
			return ftb.InstallPackVersion(installTarget, "", pack, version)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
