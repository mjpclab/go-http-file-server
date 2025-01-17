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

	input = "en-US;level=3"
	output = parseAcceptItem(input)
	if output.value != "en-US" {
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

	input = "en-US; q=0.97"
	output = parseAcceptItem(input)
	if output.value != "en-US" {
		t.Error(output.value)
	}
	if output.quality != 970 {
		t.Error(output.quality)
	}

	input = "en-US;v=7;q=0.66"
	output = parseAcceptItem(input)
	if output.value != "en-US" {
		t.Error(output.value)
	}
	if output.quality != 660 {
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

func TestAcceptItemMatch(t *testing.T) {
	var item acceptItem

	item = acceptItem{"text/html", 1000}
	if !item.match("text/html") {
		t.Error()
	}
	if item.match("text/plain") {
		t.Error()
	}

	item = acceptItem{"text/*", 1000}
	if !item.match("text/*") {
		t.Error()
	}
	if !item.match("text/html") {
		t.Error()
	}
	if !item.match("text/plain") {
		t.Error()
	}
	if item.match("image/png") {
		t.Error()
	}

	item = acceptItem{"*/*", 1000}
	if !item.match("text/*") {
		t.Error()
	}
	if !item.match("text/html") {
		t.Error()
	}
	if !item.match("text/plain") {
		t.Error()
	}
	if !item.match("image/png") {
		t.Error()
	}

}
