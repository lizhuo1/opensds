// Copyright (c) 2016 Huawei Technologies Co., Ltd. All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

/*
This module implements a standard SouthBound interface of share resource to
storage plugins.

*/

package share

import (
	"bufio"
	"errors"
	"log"
	"os"
	"os/exec"
	// "time"
)

func isMounted(mountDir string) bool {
	findmntCmd := exec.Command("findmnt", "-n", mountDir)
	findmntStdout, err := findmntCmd.StdoutPipe()
	if err != nil {
		log.Println("Could not get findmount stdout pipe:", err.Error())
	}

	if err = findmntCmd.Start(); err != nil {
		log.Println("findmnt failed to start:", err.Error())
	}

	findmntScanner := bufio.NewScanner(findmntStdout)
	findmntScanner.Split(bufio.ScanWords)
	findmntScanner.Scan()
	if findmntScanner.Err() != nil {
		log.Println("Couldn't read findnmnt output:", findmntScanner.Err().Error())
	}

	findmntText := findmntScanner.Text()
	if err = findmntCmd.Wait(); err != nil {
		_, isExitError := err.(*exec.ExitError)
		if !isExitError {
			log.Println("findmnt failed:", err.Error())
		}
	}

	return findmntText == mountDir
}

func MountShare(mountDir, device, fsType string) (string, error) {
	if isMounted(mountDir) {
		err := errors.New("This path has been mounted!")
		return "", err
	}

	if err := os.MkdirAll(mountDir, 0777); err != nil {
		log.Println("Could not create directory:", err.Error())
		return "", err
	}

	mountCmd := exec.Command("mount.nfs", device, mountDir)
	if mountOut, err := mountCmd.CombinedOutput(); err != nil {
		log.Println("Could not mount:", err.Error(), "Output:", string(mountOut))
		return "", err
	}

	return "Mount share success!", nil
}

func UnmountShare(mountDir string) (string, error) {
	if !isMounted(mountDir) {
		err := errors.New("This path is not mounted!")
		return "", err
	}

	umountCmd := exec.Command("umount", "-l", mountDir)
	if umountOut, err := umountCmd.CombinedOutput(); err != nil {
		log.Println("Could not unmount:", err.Error(), "Output:", string(umountOut))
		return "", err
	}

	return "Unmount share success!", nil
}