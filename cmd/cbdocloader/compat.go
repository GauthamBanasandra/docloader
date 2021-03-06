/*   Copyright 2016 Couchbase, Inc.
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package main

/***
 * The compatibility functions below are for supporting the old command line format
 * for the cbdocloader data path specification. In the old version we allowed the
 * user to specify the data path without any preceeding flags. This data path could
 * be specified anywhere in the command as long as it wasn't between an option flag
 * and parameter. These functions modify the command line arguments for conform to
 * the new format. This feature should be deprecated as soon as possible to reduce
 * complexity in cbdocloader.
 ***/

import (
	"fmt"
)

func compatMode(args []string) []string {
	hasDatasetOption := false
	compatDataset := ""
	outArgs := make([]string, 0)

	outArgs = append(outArgs, args[0])
	for i := 1; i < len(args); i++ {
		if isOpt, isDataset := isOption(args[i]); isOpt {
			if isDataset {
				hasDatasetOption = true
			}
			outArgs = append(outArgs, args[i])
			i++
			outArgs = append(outArgs, args[i])
		} else if isFlag(args[i]) {
			outArgs = append(outArgs, args[i])
		} else {
			if compatDataset == "" {
				compatDataset = args[i]
			}
		}
	}

	if !hasDatasetOption && compatDataset != "" {
		fmt.Printf("Warning: Specifying the dataset without the -d/--dataset option is deprecated\n")
		outArgs = append(outArgs, "-d")
		outArgs = append(outArgs, compatDataset)
	}

	return outArgs
}

func isOption(arg string) (bool, bool) {
	options := []string{"-c", "--cluster", "-n", "-u", "--username", "-p", "--password",
		"-b", "--bucket", "-m", "--bucket-quota", "-s", "-d", "--dataset", "-t", "--threads"}
	for _, option := range options {
		if option == arg {
			if arg == "-d" || arg == "--dataset" {
				return true, true
			} else {
				return true, false
			}
		}
	}

	return false, false
}

func isFlag(arg string) bool {
	flags := []string{"-v", "--verbose", "-h", "--help"}
	for _, flag := range flags {
		if flag == arg {
			return true
		}
	}
	return false
}
