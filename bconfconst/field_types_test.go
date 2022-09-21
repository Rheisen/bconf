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
	if bconfconst.String != reflect.String.String() {
		t.Errorf(
			"bconfconst '%s' does not match reflect kind '%s'",
			bconfconst.String,
			reflect.String.String(),
		)
	}
	if bconfconst.Int != reflect.Int.String() {
		t.Errorf(
			"bconfconst '%s' does not match reflect kind '%s'",
			bconfconst.Int,
			reflect.Int.String(),
		)
	}
	if bconfconst.Int16 != reflect.Int16.String() {
		t.Errorf(
			"bconfconst '%s' does not match reflect kind '%s'",
			bconfconst.Int16,
			reflect.Int16.String(),
		)
	}
	if bconfconst.Int32 != reflect.Int32.String() {
		t.Errorf(
			"bconfconst '%s' does not match reflect kind '%s'",
			bconfconst.Int32,
			reflect.Int32.String(),
		)
	}
	if bconfconst.Int64 != reflect.Int64.String() {
		t.Errorf(
			"bconfconst '%s' does not match reflect kind '%s'",
			bconfconst.Int64,
			reflect.Int64.String(),
		)
	}
	if bconfconst.Float32 != reflect.Float32.String() {
		t.Errorf(
			"bconfconst '%s' does not match reflect kind '%s'",
			bconfconst.Float32,
			reflect.Float32.String(),
		)
	}
	if bconfconst.Float64 != reflect.Float64.String() {
		t.Errorf(
			"bconfconst '%s' does not match reflect kind '%s'",
			bconfconst.Float64,
			reflect.Float64.String(),
		)
	}
	if bconfconst.Time != reflect.TypeOf(time.Time{}).String() {
		t.Errorf(
			"bconfconst '%s' does not match reflect type '%s'",
			bconfconst.Time,
			reflect.TypeOf(time.Time{}).String(),
		)
	}
	if bconfconst.Duration != reflect.TypeOf(time.Nanosecond).String() {
		t.Errorf(
			"bconfconst '%s' does not match reflect type '%s'",
			bconfconst.Duration,
			reflect.TypeOf(time.Nanosecond).String(),
		)
	}
}
