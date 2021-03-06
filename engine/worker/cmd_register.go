package main

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/ovh/cds/sdk"
	"github.com/ovh/cds/sdk/log"
)

func cmdRegister(w *currentWorker) *cobra.Command {
	var cmdRegister = &cobra.Command{
		Use:    "register",
		Long:   "worker register is a subcommand used by hatchery. This is not directly useful for end user",
		Hidden: true, // user should not use this command directly
		Run:    cmdRegisterRun(w),
	}
	initFlagsRun(cmdRegister)
	return cmdRegister
}

func cmdRegisterRun(w *currentWorker) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		initFlags(cmd, w)
		form := sdk.WorkerRegistrationForm{
			Name:         w.status.Name,
			Token:        w.token,
			HatcheryName: w.hatchery.name,
			ModelID:      w.model.ID,
		}

		if err := w.register(form); err != nil {
			log.Error("Unable to register worker %s: %v", w.status.Name, err)
		}
		if err := w.unregister(); err != nil {
			log.Error("Unable to unregister worker %s: %v", w.status.Name, err)
		}

		if FlagBool(cmd, flagForceExit) {
			log.Info("Exiting worker with force_exit true")
			return
		}

		if w.hatchery.name != "" {
			log.Info("Waiting 30min to be killed by hatchery, if not killed, worker will exit")
			time.Sleep(30 * time.Minute)
		}
	}
}
