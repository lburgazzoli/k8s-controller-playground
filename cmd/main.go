/*
Copyright 2023.

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

package main

import (
	"flag"
	"os"

	"github.com/lburgazzoli/k8s-controller-playground/cmd/modelschema"
	"github.com/lburgazzoli/k8s-controller-playground/cmd/run"
	"github.com/lburgazzoli/k8s-controller-playground/pkg/logger"

	"github.com/spf13/cobra"

	"k8s.io/klog/v2"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "playground",
		Short: "playground",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	rootCmd.AddCommand(modelschema.NewCmd())
	rootCmd.AddCommand(run.NewCmd())

	fs := flag.NewFlagSet("", flag.PanicOnError)

	klog.InitFlags(fs)
	logger.Options.BindFlags(fs)

	rootCmd.PersistentFlags().AddGoFlagSet(fs)

	if err := rootCmd.Execute(); err != nil {
		klog.ErrorS(err, "problem running command")
		os.Exit(1)
	}
}
