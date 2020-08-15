package main

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"io"
	"syscall/js"
)

func save(s *Scene) (string, error) {
	//deselect all before saving
	curSelection.deselectAll(s)
	var data bytes.Buffer
	for _, vox := range s.voxels {
		enc := gob.NewEncoder(&data)
		err := enc.Encode(vox.public())
		if err != nil {
			return "", err
		}
	}

	return base64.StdEncoding.EncodeToString(data.Bytes()), nil
}

func load(s *Scene, dataStr string) error {
	s.voxels = []*Voxel{}
	data64, err := base64.StdEncoding.DecodeString(dataStr)
	if err != nil {
		return err
	}
	data := bytes.NewBuffer(data64)
	if err != nil {
		return err
	}

	dec := gob.NewDecoder(data)
	ok := true
	for ok {
		var v publicVoxel
		err := dec.Decode(&v)
		if err == io.EOF {
			fmt.Printf("Reached the end of the data")
			ok = false
		}
		if err != nil {
			continue
		}
		s.addVoxel(v.private())
	}

	s.update = true
	return nil
}

func initSaveBtn(s *Scene) {
	js.Global().Set("saveCurrentScene", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		data, err := save(s)
		if err != nil {
			return fmt.Sprintf("ERROR: %v\n", err.Error())
		}
		return data
	}))
}

func initLoadBtn(s *Scene) {
	fmt.Printf("Loading Data\n")
	js.Global().Set("loadNewScene", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		data := args[0].String()
		err := load(s, data)
		if err != nil {
			return fmt.Sprintf("ERROR: %v\n", err.Error())
		}
		return ""
	}))
}
