package pubnub

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"unicode"
)

// operationStringUnset is the OperationType.String default. Validation errors
// embed it ("pubnub/validation: pubnub: No Category Matched: …"), so a new
// constant that falls through changes no error text until a case is added,
// and the missing case is easy to miss.
const operationStringUnset = "No Category Matched"

// operationStringGaps are constants that already fall through to
// operationStringUnset. Their validation tests assert that text. Do not add
// a new constant here; add a String() case instead. Removing a name from
// this map is required once that constant grows a case.
var operationStringGaps = map[string]struct{}{
	"PNMessageCountsOperation":   {},
	"PNAccessManagerRevokeToken": {},
}

// objectsCreateUserValue is the iota value of PNCreateUserOperation.
// Inserting or deleting a constant earlier in the OperationType block
// renumbers every Objects constant. enums.go says that order is part of
// the Objects API.
const objectsCreateUserValue = 26

// objectsOperationOrder is the Objects API section of the OperationType
// const block, from PNCreateUserOperation through PNGetChannelMetadataOperation.
// New Objects constants are appended after this list. They must not be
// inserted into it or placed before it.
var objectsOperationOrder = []string{
	"PNCreateUserOperation",
	"PNGetUsersOperation",
	"PNGetUserOperation",
	"PNUpdateUserOperation",
	"PNDeleteUserOperation",
	"PNGetSpaceOperation",
	"PNGetSpacesOperation",
	"PNCreateSpaceOperation",
	"PNDeleteSpaceOperation",
	"PNUpdateSpaceOperation",
	"PNGetMembershipsOperation",
	"PNGetChannelMembersOperation",
	"PNManageMembershipsOperation",
	"PNManageMembersOperation",
	"PNSetChannelMembersOperation",
	"PNSetMembershipsOperation",
	"PNRemoveChannelMetadataOperation",
	"PNRemoveUUIDMetadataOperation",
	"PNGetAllChannelMetadataOperation",
	"PNGetAllUUIDMetadataOperation",
	"PNGetUUIDMetadataOperation",
	"PNRemoveMembershipsOperation",
	"PNRemoveChannelMembersOperation",
	"PNSetUUIDMetadataOperation",
	"PNSetChannelMetadataOperation",
	"PNGetChannelMetadataOperation",
}

type operationConst struct {
	name  string
	value int
}

// TestOperationTypeInvariants checks two mistakes that otherwise show up
// only in review: a new OperationType constant with no String() case, and
// a reordered or renumbered Objects constant.
func TestOperationTypeInvariants(t *testing.T) {
	ops, err := parseOperationTypeConstants(enumsSource(t))
	if err != nil {
		t.Fatal(err)
	}

	seen := make(map[string]struct{}, len(ops))
	for _, op := range ops {
		if _, dup := seen[op.name]; dup {
			t.Errorf("duplicate OperationType constant %s", op.name)
		}
		seen[op.name] = struct{}{}

		got := OperationType(op.value).String()
		if _, gap := operationStringGaps[op.name]; gap {
			if got != operationStringUnset {
				t.Errorf("%s is listed in operationStringGaps but String() returns %q; remove it from operationStringGaps", op.name, got)
			}
			continue
		}
		if got == operationStringUnset {
			t.Errorf("OperationType %s = %d has no String() case (returns %q). Add a case in OperationType.String. Validation errors embed this text.", op.name, op.value, got)
		}
	}
	for name := range operationStringGaps {
		if _, ok := seen[name]; !ok {
			t.Errorf("operationStringGaps lists %s, which is not an OperationType constant", name)
		}
	}

	assertObjectsOperationOrder(t, ops)
}

func enumsSource(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	path := filepath.Join(filepath.Dir(file), "enums.go")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// parseOperationTypeConstants reads the iota block that starts at
// PNSubscribeOperation. Each following name repeats the expression
// "1 + iota", so values are 1, 2, 3, ... in source order.
func parseOperationTypeConstants(src string) ([]operationConst, error) {
	const marker = "PNSubscribeOperation OperationType = 1 + iota"
	idx := strings.Index(src, marker)
	if idx < 0 {
		return nil, fmt.Errorf("enums.go: %q not found", marker)
	}
	open := strings.LastIndex(src[:idx], "const (")
	if open < 0 {
		return nil, fmt.Errorf("enums.go: OperationType const block not found")
	}

	var ops []operationConst
	value := 1
	for _, line := range strings.Split(src[open+len("const ("):], "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == ")" {
			break
		}
		name, ok := operationConstName(trimmed)
		if !ok {
			continue
		}
		ops = append(ops, operationConst{name: name, value: value})
		value++
	}
	if len(ops) == 0 || ops[0].name != "PNSubscribeOperation" || ops[0].value != 1 {
		return nil, fmt.Errorf("enums.go: OperationType block did not start at PNSubscribeOperation = 1")
	}
	return ops, nil
}

func operationConstName(trimmed string) (string, bool) {
	if trimmed == "" || strings.HasPrefix(trimmed, "//") {
		return "", false
	}
	field := trimmed
	if i := strings.IndexAny(field, " \t"); i >= 0 {
		field = field[:i]
	}
	if !strings.HasPrefix(field, "PN") {
		return "", false
	}
	for _, r := range field {
		if r != '_' && !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return "", false
		}
	}
	return field, true
}

func assertObjectsOperationOrder(t *testing.T, ops []operationConst) {
	t.Helper()
	start := -1
	for i, op := range ops {
		if op.name == objectsOperationOrder[0] {
			start = i
			break
		}
	}
	if start < 0 {
		t.Errorf("Objects constant %s is missing from the OperationType block", objectsOperationOrder[0])
		return
	}
	if ops[start].value != objectsCreateUserValue {
		t.Errorf("PNCreateUserOperation value is %d, want %d. A constant was inserted or removed before the Objects block, which renumbers the Objects API. Append new constants after %s.",
			ops[start].value, objectsCreateUserValue, objectsOperationOrder[len(objectsOperationOrder)-1])
	}
	if start+len(objectsOperationOrder) > len(ops) {
		t.Fatalf("Objects constant list is longer than the parsed OperationType block")
	}
	for i, want := range objectsOperationOrder {
		got := ops[start+i]
		if got.name != want {
			t.Errorf("Objects constant %d is %s, want %s. Do not reorder or insert inside the Objects constants. Append new Objects constants after %s.",
				i, got.name, want, objectsOperationOrder[len(objectsOperationOrder)-1])
			return
		}
		if got.value != objectsCreateUserValue+i {
			t.Errorf("%s value is %d, want %d", got.name, got.value, objectsCreateUserValue+i)
		}
	}
}
