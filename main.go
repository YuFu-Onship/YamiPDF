package main

import (
	"app/src/trunk"
	"log"
	"os"
	"path"
	"runtime/debug"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	debug.SetGCPercent(100)

	exe_path, _ := os.Executable()
	project_path := path.Dir(exe_path)
	trunk := trunk.NewTrunk(project_path)
	trunk.Run()
}
