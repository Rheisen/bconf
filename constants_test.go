package bconf_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/rheisen/bconf"
)

func TestConstantsMatchReflectKinds(t *testing.T) {
	if bconf.Bool != reflect.Bool.String() {
		t.Errorf(
			"bconf '%s' does not match reflect kind '%s'",
			bconf.Bool,
			reflect.Bool.String(),
		)
	}

	if bconf.Bools != reflect.TypeOf([]bool{}).String() {
		t.Errorf(
			"bconf '%s' does not match reflect kind '%s'",
			bconf.Bools,
			reflect.TypeOf([]bool{}).String(),
		)
	}

	if bconf.String != reflect.String.String() {
		t.Errorf(
			"bconf '%s' does not match reflect kind '%s'",
			bconf.String,
			reflect.String.String(),
		)
	}

	if bconf.Strings != reflect.TypeOf([]string{}).String() {
		t.Errorf(
			"bconf '%s' does not match reflect kind '%s'",
			bconf.Strings,
			reflect.TypeOf([]string{}).String(),
		)
	}

	if bconf.Int != reflect.Int.String() {
		t.Errorf(
			"bconf '%s' does not match reflect kind '%s'",
			bconf.Int,
			reflect.Int.String(),
		)
	}

	if bconf.Ints != reflect.TypeOf([]int{}).String() {
		t.Errorf(
			"bconf '%s' does not match reflect kind '%s'",
			bconf.Ints,
			reflect.TypeOf([]int{}).String(),
		)
	}

	if bconf.Float != reflect.Float64.String() {
		t.Errorf(
			"bconf '%s' does not match reflect kind '%s'",
			bconf.Float,
			reflect.Float64.String(),
		)
	}

	if bconf.Floats != reflect.TypeOf([]float64{}).String() {
		t.Errorf(
			"bconf '%s' does not match reflect kind '%s'",
			bconf.Ints,
			reflect.TypeOf([]float64{}).String(),
		)
	}

	if bconf.Time != reflect.TypeOf(time.Time{}).String() {
		t.Errorf(
			"bconf '%s' does not match reflect type '%s'",
			bconf.Time,
			reflect.TypeOf(time.Time{}).String(),
		)
	}

	if bconf.Times != reflect.TypeOf([]time.Time{}).String() {
		t.Errorf(
			"bconf '%s' does not match reflect kind '%s'",
			bconf.Times,
			reflect.TypeOf([]time.Time{}).String(),
		)
	}

	if bconf.Duration != reflect.TypeOf(time.Nanosecond).String() {
		t.Errorf(
			"bconf '%s' does not match reflect type '%s'",
			bconf.Duration,
			reflect.TypeOf(time.Nanosecond).String(),
		)
	}

	if bconf.Durations != reflect.TypeOf([]time.Duration{}).String() {
		t.Errorf(
			"bconf '%s' does not match reflect kind '%s'",
			bconf.Durations,
			reflect.TypeOf([]time.Duration{}).String(),
		)
	}
}
