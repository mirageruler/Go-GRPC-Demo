package main

import (
	"strings"
	"testing"

	pb "github.com/mirageruler/grpc-go-course/addressbook/addressbookpb"
)

func TestPromtForAddress_ReturnsAddress(t *testing.T) {
	in := `12345
Example Name
name@example.com
123-456-7890
home
222-222-2222
mobile
111-111-1111
work
777-777-7777
unknown

`

	got, err := promptForAddress(strings.NewReader(in))
	if err != nil {
		t.Fatalf("promptForAddress(%q) had unexpected error: %s", in, err.Error())
	}
	if got.Id != 12345 {
		t.Errorf("promptForAddress(%q) got %d, want ID %d", in, got.Id, 12345)
	}
	if got.Name != "Example Name" {
		t.Errorf("promptForAddress(%q) got %s, want ID %s", in, got.Name, "Example Name")
	}
	if got.Email != "name@example.com" {
		t.Errorf("promptForAddress(%q) got %s, want ID %s", in, got.Email, "name@example.com")
	}

	want := []*pb.Person_PhoneNumber{
		{Number: "123-456-7890", Type: pb.Person_HOME},
		{Number: "222-222-2222", Type: pb.Person_MOBILE},
		{Number: "111-111-1111", Type: pb.Person_WORK},
		{Number: "777-777-7777", Type: pb.Person_MOBILE},
	}

	if len(got.Phones) != len(want) {
		t.Errorf("want %d phone numbers, got %d", len(want), len(got.Phones))
	}
}
