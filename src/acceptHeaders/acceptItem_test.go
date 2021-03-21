package acceptHeaders

import "testing"

func TestParseAcceptItem(t *testing.T) {
	var input string
	var output acceptItem

	input = "en-US"
	output = parseAcceptItem(input)
	if output.value != input {
		t.Error(output.value)
	}
	if output.quality != 1000 {
		t.Error(output.quality)
	}

	input = "en-US;q=1"
	output = parseAcceptItem(input)
	if output.value != "en-US" {
		t.Error(output.value)
	}
	if output.quality != 1000 {
		t.Error(output.quality)
	}

	input = "en-US;q=1.0"
	output = parseAcceptItem(input)
	if output.value != "en-US" {
		t.Error(output.value)
	}
	if output.quality != 1000 {
		t.Error(output.quality)
	}

	input = "en-US;q=1."
	output = parseAcceptItem(input)
	if output.value != "en-US" {
		t.Error(output.value)
	}
	if output.quality != 1000 {
		t.Error(output.quality)
	}

	input = "en-US;q=0.97"
	output = parseAcceptItem(input)
	if output.value != "en-US" {
		t.Error(output.value)
	}
	if output.quality != 970 {
		t.Error(output.quality)
	}

	input = "en-US;q=1.1"
	output = parseAcceptItem(input)
	if output.value != "en-US" {
		t.Error(output.value)
	}
	if output.quality != 1000 {
		t.Error(output.quality)
	}

	input = "en-US;q=2.1"
	output = parseAcceptItem(input)
	if output.value != "en-US" {
		t.Error(output.value)
	}
	if output.quality != 1000 {
		t.Error(output.quality)
	}

	input = "en-US;q=2."
	output = parseAcceptItem(input)
	if output.value != "en-US" {
		t.Error(output.value)
	}
	if output.quality != 1000 {
		t.Error(output.quality)
	}

	input = "en-US;q=2.xyz"
	output = parseAcceptItem(input)
	if output.value != "en-US" {
		t.Error(output.value)
	}
	if output.quality != 1000 {
		t.Error(output.quality)
	}

	input = "en-US;q=2"
	output = parseAcceptItem(input)
	if output.value != "en-US" {
		t.Error(output.value)
	}
	if output.quality != 1000 {
		t.Error(output.quality)
	}

	input = "en-US;q=0.123456"
	output = parseAcceptItem(input)
	if output.value != "en-US" {
		t.Error(output.value)
	}
	if output.quality != 123 {
		t.Error(output.quality)
	}

	input = "en-US;q=0."
	output = parseAcceptItem(input)
	if output.value != "en-US" {
		t.Error(output.value)
	}
	if output.quality != 0 {
		t.Error(output.quality)
	}

	input = "en-US;q=0.xyz"
	output = parseAcceptItem(input)
	if output.value != "en-US" {
		t.Error(output.value)
	}
	if output.quality != 1000 {
		t.Error(output.quality)
	}

	input = "en-US;q=xyz"
	output = parseAcceptItem(input)
	if output.value != "en-US" {
		t.Error(output.value)
	}
	if output.quality != 1000 {
		t.Error(output.quality)
	}

	input = "en-US;q=0zyx"
	output = parseAcceptItem(input)
	if output.value != "en-US" {
		t.Error(output.value)
	}
	if output.quality != 1000 {
		t.Error(output.quality)
	}
}
