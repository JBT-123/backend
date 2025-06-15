package main

import (
	"database/sql"
	"fmt"
	"reflect"
)

type Seed struct {
	DB *sql.DB
}

func (s *Seed) Execute(methods ...string) error {
	seedType := reflect.TypeOf(s)
	seedValue := reflect.ValueOf(s)

	if len(methods) == 0 {
		// Run all methods
		for i := 0; i < seedType.NumMethod(); i++ {
			method := seedType.Method(i)
			if method.Name == "Execute" {
				continue
			}
			fmt.Printf("Running seeder: %s\n", method.Name)
			seedValue.MethodByName(method.Name).Call(nil)
		}
		return nil
	}

	// Run only specified methods
	for _, name := range methods {
		method := seedValue.MethodByName(name)
		if !method.IsValid() {
			return fmt.Errorf("seeder method %s not found", name)
		}
		fmt.Printf("Running seeder: %s\n", name)
		method.Call(nil)
	}
	return nil
}
