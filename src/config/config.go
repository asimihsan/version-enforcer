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

package config

import (
	"enforce-tool-versions/identifier"
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/rs/zerolog"
)

type Config struct {
	Binary []*Binary `hcl:"binary,block"`
}

type Binary struct {
	Name    string `hcl:"name,label"`
	Version string `hcl:"version"`
}

func LoadConfig(configPath string, zlog *zerolog.Logger) (*Config, error) {
	var cfg Config
	err := hclsimple.DecodeFile(configPath, nil, &cfg)
	if err != nil {
		if diagnostics, ok := err.(hcl.Diagnostics); ok {
			for _, diagnostic := range diagnostics {
				fmt.Println(diagnostic)
			}
		} else {
			zlog.Error().Stack().Err(err).Msg("Failed to decode config")
		}
		return nil, err
	}

	for _, binary := range cfg.Binary {
		_, err := identifier.GetProgram(binary.Name)
		if err != nil {
			zlog.Error().Err(err).Interface("binary", binary).Msg("failed to get program")
			return nil, err
		}

		_, err = identifier.NewRequirement(binary.Version)
		if err != nil {
			zlog.Error().Err(err).Interface("binary", binary).Msg("failed to parse requirement")
			return nil, err
		}
	}

	return &cfg, nil
}
