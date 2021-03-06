package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v1"

	"github.com/Southclaws/sampctl/download"
	"github.com/Southclaws/sampctl/print"
	"github.com/Southclaws/sampctl/rook"
	"github.com/Southclaws/sampctl/util"
	"github.com/Southclaws/sampctl/versioning"
)

var packageUninstallFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "dir",
		Value: ".",
		Usage: "working directory for the project - by default, uses the current directory",
	},
	cli.BoolFlag{
		Name:  "dev",
		Usage: "for specifying dependencies only necessary for development or testing of the package",
	},
}

func packageUninstall(c *cli.Context) error {
	if c.Bool("verbose") {
		print.SetVerbose()
	}

	dir := util.FullPath(c.String("dir"))
	development := c.Bool("dev")

	if len(c.Args()) == 0 {
		cli.ShowCommandHelpAndExit(c, "uninstall", 0)
		return nil
	}

	cacheDir, err := download.GetCacheDir()
	if err != nil {
		print.Erro("Failed to retrieve cache directory path (attempted <user folder>/.samp) ")
		return err
	}

	deps := []versioning.DependencyString{}
	for _, dep := range c.Args() {
		deps = append(deps, versioning.DependencyString(dep))
	}

	pkg, err := rook.PackageFromDir(true, dir, runtime.GOOS, "")
	if err != nil {
		return errors.Wrap(err, "failed to interpret directory as Pawn package")
	}

	err = rook.Uninstall(context.Background(), gh, pkg, deps, development, gitAuth, runtime.GOOS, cacheDir)
	if err != nil {
		return err
	}

	print.Info("successfully removed dependency")

	return nil
}

func packageUninstallBash(c *cli.Context) {
	cacheDir, err := download.GetCacheDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to retrieve cache directory path (attempted <user folder>/.samp) ", err)
		return
	}

	packages, err := download.GetPackageList(cacheDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get package list:", err)
		return
	}

	query := c.Args().First()
	for _, pkg := range packages {
		if strings.HasPrefix(pkg.String(), query) {
			fmt.Println(pkg)
		}
	}
}
