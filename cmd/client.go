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
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"rp-bot-client/src/bot"
	"rp-bot-client/src/captcha"
	"rp-bot-client/src/captcha/recognizer"
	"rp-bot-client/src/cropper"
	"rp-bot-client/src/event"
	"rp-bot-client/src/repainter"
	"rp-bot-client/src/saver"
	"rp-bot-client/src/storage"
	"rp-bot-client/src/window"
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
	withStorage, err := cmd.Flags().GetBool("storage")
	if err != nil {
		return err
	}

	botType, err := getBotType(args)
	if err != nil {
		return err
	}

	pid, err := window.FindGtaPid("GTA5.exe")
	if err != nil {
		return err
	}

	r, err := recognizer.Get(recognizer.GOCR)
	if err != nil {
		return err
	}

	captchaChecker := captcha.NewChecker()
	captchaSolver := captcha.NewSolver(pid, r,
		captcha.NewScreenshotProcessor(cropper.NewImage(repainter.NewImage()), saver.NewImg()))
	storageManipulator := storage.NewManipulator(pid)

	w, err := worker.GetWorker(pid, botType, "e", captchaChecker, captchaSolver, storageManipulator, withStorage)
	if err != nil {
		return err
	}

	return bot.NewBot(pid, w, event.NewEventListener()).Start()
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
