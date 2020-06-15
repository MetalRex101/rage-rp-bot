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
	"rp-bot-client/src/repainter"
	"rp-bot-client/src/saver"
	"rp-bot-client/src/worker"
)

type botType string

func (b botType) isValid() bool {
	switch b {
	case oilMan, miner:
		return true
	default:
		return false
	}
}

const (
	oilMan botType = "oil"
	miner  botType = "miner"
)

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
	botType, err := getBotType(args)
	if err != nil {
		return err
	}

	pid, err := findGtaPid("GTA5.exe")
	if err != nil {
		return err
	}

	w, err := worker.GetWorker(pid, botType, "e")
	if err != nil {
		return err
	}

	captchaSolver := captcha.NewSolver(
		pid,
		captcha.NewRecognizer(),
		captcha.NewScreenshotProcessor(cropper.NewImage(repainter.NewImage()), saver.NewImg()),
	)

	return bot.NewBot(pid, w, captchaSolver, captcha.NewMouseManipulator(pid), event.NewEventListener()).Start()
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

func getBotType(args []string) (string, error) {
	b := "oil"

	if len(args) > 0 {
		b = args[0]
	}

	if !botType(b).isValid() {
		return "", errors.Errorf("invalid bot type. Waiting for 'oil' or 'mine', '%s' given", b)
	}

	return b, nil
}

func init() {
	rootCmd.AddCommand(clientCmd)
}
