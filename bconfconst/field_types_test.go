package bconfconst_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/rheisen/bconf/bconfconst"
)

func TestConstantsMatchReflectKinds(t *testing.T) {
	if bconfconst.Bool != reflect.Bool.String() {
		t.Errorf(
			"bconfconst '%s' does not match reflect kind '%s'",
			bconfconst.Bool,
			reflect.Bool.String(),
		)
	}
	if bconfconst.Bools != reflect.TypeOf([]bool{}).String() {
		t.Errorf(
			"bconfconst '%s' does not match reflect kind '%s'",
			bconfconst.Bools,
			reflect.TypeOf([]bool{}).String(),
		)
	}
	if bconfconst.String != reflect.String.String() {
		t.Errorf(
			"bconfconst '%s' does not match reflect kind '%s'",
			bconfconst.String,
			reflect.String.String(),
		)
	}
	if bconfconst.Strings != reflect.TypeOf([]string{}).String() {
		t.Errorf(
			"bconfconst '%s' does not match reflect kind '%s'",
			bconfconst.Strings,
			reflect.TypeOf([]string{}).String(),
		)
	}
	if bconfconst.Int != reflect.Int.String() {
		t.Errorf(
			"bconfconst '%s' does not match reflect kind '%s'",
			bconfconst.Int,
			reflect.Int.String(),
		)
	}
	if bconfconst.Ints != reflect.TypeOf([]int{}).String() {
		t.Errorf(
			"bconfconst '%s' does not match reflect kind '%s'",
			bconfconst.Ints,
			reflect.TypeOf([]int{}).String(),
		)
	}
	if bconfconst.Float != reflect.Float64.String() {
		t.Errorf(
			"bconfconst '%s' does not match reflect kind '%s'",
			bconfconst.Float,
			reflect.Float64.String(),
		)
	}
	if bconfconst.Floats != reflect.TypeOf([]float64{}).String() {
		t.Errorf(
			"bconfconst '%s' does not match reflect kind '%s'",
			bconfconst.Ints,
			reflect.TypeOf([]float64{}).String(),
		)
	}
	if bconfconst.Time != reflect.TypeOf(time.Time{}).String() {
		t.Errorf(
			"bconfconst '%s' does not match reflect type '%s'",
			bconfconst.Time,
			reflect.TypeOf(time.Time{}).String(),
		)
	}
	if bconfconst.Times != reflect.TypeOf([]time.Time{}).String() {
		t.Errorf(
			"bconfconst '%s' does not match reflect kind '%s'",
			bconfconst.Times,
			reflect.TypeOf([]time.Time{}).String(),
		)
	}
	if bconfconst.Duration != reflect.TypeOf(time.Nanosecond).String() {
		t.Errorf(
			"bconfconst '%s' does not match reflect type '%s'",
			bconfconst.Duration,
			reflect.TypeOf(time.Nanosecond).String(),
		)
	}
	if bconfconst.Durations != reflect.TypeOf([]time.Duration{}).String() {
		t.Errorf(
			"bconfconst '%s' does not match reflect kind '%s'",
			bconfconst.Durations,
			reflect.TypeOf([]time.Duration{}).String(),
		)
	}
}
