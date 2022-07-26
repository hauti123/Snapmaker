package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/hauti123/Snapmaker/snapmaker"
)

type UploadCommand struct {
	flags    *flag.FlagSet
	filePath string
	printer  snapmaker.Snapmaker
}

func NewUploadCommand() *UploadCommand {

	uploadCommand := &UploadCommand{
		flags: flag.NewFlagSet("upload", flag.ContinueOnError),
	}

	return uploadCommand
}

func (uc *UploadCommand) Flags() *flag.FlagSet {
	return uc.flags
}

func (uc *UploadCommand) Name() string {
	return uc.flags.Name()
}

func (uc *UploadCommand) Init(printer snapmaker.Snapmaker) error {
	if !uc.flags.Parsed() {
		return fmt.Errorf("Command line arguments not yet parsed.")
	}

	if uc.flags.NArg() < 1 {
		return fmt.Errorf("Missing file path to .gcode file.")
	}

	uc.filePath = uc.flags.Args()[0]
	if len(uc.filePath) == 0 {
		return errors.New("Missing file path to .gcode file.")
	}

	uc.printer = printer

	fmt.Printf("%v", uc)
	return nil
}

func (uc *UploadCommand) Run() error {
	err := (&uc.printer).Connect()
	if err != nil {
		return err
	}
	err = uc.printer.SendGcodeFile(uc.filePath)
	if err != nil {
		return err
	}
	return nil
}
