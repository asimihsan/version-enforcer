/*
 * Copyright 2023 Asim Ihsan
 *
 * Licensed under the Apache License, SemverVersion 2.0 (the "License");
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

package main

import (
	"enforce-tool-versions/config"
	"enforce-tool-versions/identifier"
	"github.com/rs/zerolog"
	"os"
)

func main() {
	zlog := zerolog.New(os.Stdout).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	cfg, err := config.LoadConfig("tool-enforcer.hcl", &zlog)
	if err != nil {
		zlog.Error().Err(err).Msg("failed to load config")
		return
	}
	zlog.Debug().Interface("config", cfg).Msg("loaded config")

	for _, binary := range cfg.Binary {
		program, err := identifier.GetProgram(binary.Name)
		if err != nil {
			zlog.Error().Err(err).Interface("binary", binary).Msg("failed to get program")
			continue
		}

		version, err := identifier.Identify(*program, &zlog)
		if err != nil {
			zlog.Error().Err(err).Msg("failed to identify program")
			continue
		}

		if !identifier.Satisfies(string(version), binary.Version) {
			zlog.Debug().
				Interface("version", version).
				Interface("binary", binary).
				Msg("version does not satisfy requirement")
			continue
		} else {
			zlog.Debug().
				Interface("version", version).
				Interface("binary", binary).
				Msg("version satisfies requirement")
		}
	}
}
