/*
 * Copyright 2023 Asim Ihsan
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package cmd

import (
	"fmt"
	"github.com/asimihsan/version-enforcer/src/config"
	"github.com/asimihsan/version-enforcer/src/identifier"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:  "enforce --config <config file>",
	Long: "Enforce tool versions",
	Run: func(cmd *cobra.Command, args []string) {
		zlog := zerolog.New(os.Stdout).With().Timestamp().Logger()

		if verbose {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}

		cfg, err := config.LoadConfig(cfgFile, &zlog)
		if err != nil {
			zlog.Error().Err(err).Msg("failed to load config")
			os.Exit(1)
		}
		zlog.Debug().Interface("config", cfg).Msg("loaded config")

		anyFailures := false
		for _, binary := range cfg.Binary {
			program, err := identifier.GetProgram(binary.Name)
			if err != nil {
				zlog.Error().Err(err).Interface("binary", binary).Msg("failed to get program")
				os.Exit(1)
			}

			version, err := identifier.Identify(*program, &zlog)
			if err != nil {
				zlog.Error().Err(err).Msg("failed to identify program")
				os.Exit(1)
			}

			if !identifier.Satisfies(string(version), binary.Version) {
				zlog.Debug().
					Interface("version", version).
					Interface("binary", binary).
					Msg("version does not satisfy requirement")
				msg := fmt.Sprintf("%s version %s does not satisfy requirement %s", binary.Name, version, binary.Version)
				PrintErrorLine(msg)
				anyFailures = true

				continue
			} else {
				zlog.Debug().
					Interface("version", version).
					Interface("binary", binary).
					Msg("version satisfies requirement")
				if verbose {
					msg := fmt.Sprintf("%s version %s satisfies requirement %s", binary.Name, version, binary.Version)
					PrintSuccessLine(msg)
				}
			}
		}

		if anyFailures {
			os.Exit(1)
		}
	},
}

// PrintErrorLine prints an error message in bright red.
func PrintErrorLine(message string) {
	fmt.Printf("\033[31;1m%s\033[0m %s\n", "Error:", message)
}

// PrintSuccessLine prints a success message in bright green.
func PrintSuccessLine(message string) {
	fmt.Printf("\033[32;1m%s\033[0m %s\n", "Success:", message)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
