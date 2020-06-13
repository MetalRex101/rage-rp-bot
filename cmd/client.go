/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"rp-bot-client/src/bot"
	"rp-bot-client/src/captcha"
	"rp-bot-client/src/cropper"
	"rp-bot-client/src/event"
	"rp-bot-client/src/miner"
	"rp-bot-client/src/saver"
)

// #cgo windows LDFLAGS: -lgdi32 -luser32
// #cgo windows,amd64 LDFLAGS: -L${SRCDIR}/cdeps/win64 -lpng -lz
// #cgo windows,386 LDFLAGS: -L${SRCDIR}/cdeps/win32 -lpng -lz

// //#include "screen/goScreen.h"
// //#include "mouse/goMouse.h"
// //#include "key/goKey.h"
// //#include "bitmap/goBitmap.h"
// //#include "event/goEvent.h"
// // #include "github.com/go-vgo/robotgo/window/goWindow.h"
import "C"

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: runClient,
}

func runClient(cmd *cobra.Command, args []string) error {
	pid, err := findGtaPid("GTA5.exe")
	if err != nil {
		return err
	}

	minerWorker := miner.NewWorker("e")
	captchaSolver := captcha.NewSolver(
		pid,
		captcha.NewRecognizerClient(),
		captcha.NewScreenshotProcessor(cropper.NewImage(), saver.NewImg()),
	)

	return bot.NewBot(pid, minerWorker, captchaSolver, captcha.NewMouseManipulator(pid), event.NewEventListener()).Start()
}

func findGtaPid(name string) (int32, error) {
	var pid int32 = -1

	pc, err := robotgo.Process()
	if err != nil {
		return pid, err
	}

	for _, proc := range pc {
		if proc.Name == name {
			pid = proc.Pid
			fmt.Println(fmt.Sprintf("%+v", proc))
		}
	}

	if pid == -1 {
		return pid, errors.New(fmt.Sprintf("process not found, %s", name))
	}

	return pid, nil
}

func init() {
	rootCmd.AddCommand(clientCmd)
}
