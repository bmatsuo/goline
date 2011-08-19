// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 *  Filename:    goline.go
 *  Package:     goline
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Sat Aug 13 02:28:54 PDT 2011
 *  Description: 
 */

// Package goline manages prompts in the spirit of highline.
package goline
import (
    "fmt"
    "os"
)

func Say(msg string) (int, os.Error) { return fmt.Println(msg) }
