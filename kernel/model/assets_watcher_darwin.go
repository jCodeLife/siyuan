// SiYuan - Build Your Eternal Digital Garden
// Copyright (c) 2020-present, b3log.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

//go:build darwin

package model

import (
	"path/filepath"
	"time"

	"github.com/radovskyb/watcher"
	"github.com/siyuan-note/siyuan/kernel/util"
)

var assetsWatcher *watcher.Watcher

func WatchAssets() {
	if "iOS" == util.Container {
		return
	}

	go func() {
		watchAssets()
	}()
}

func watchAssets() {
	if nil != assetsWatcher {
		assetsWatcher.Close()
	}
	assetsWatcher = watcher.New()

	assetsDir := filepath.Join(util.DataDir, "assets")

	go func() {
		for {
			select {
			case event, ok := <-assetsWatcher.Event:
				if !ok {
					return
				}

				//util.LogInfof("assets changed: %s", event)
				if watcher.Write == event.Op {
					IncWorkspaceDataVer()
				}
			case err, ok := <-assetsWatcher.Error:
				if !ok {
					return
				}
				util.LogErrorf("watch assets failed: %s", err)
			case <-assetsWatcher.Closed:
				return
			}
		}
	}()

	if err := assetsWatcher.Add(assetsDir); nil != err {
		util.LogErrorf("add assets watcher for folder [%s] failed: %s", assetsDir, err)
		return
	}

	//util.LogInfof("added file watcher [%s]", assetsDir)
	if err := assetsWatcher.Start(10 * time.Second); nil != err {
		util.LogErrorf("start assets watcher for folder [%s] failed: %s", assetsDir, err)
		return
	}
}

func CloseWatchAssets() {
	if nil != assetsWatcher {
		assetsWatcher.Close()
	}
}
