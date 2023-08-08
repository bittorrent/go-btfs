package s3action

import "testing"

func TestAction_IsValid(t *testing.T) {
	testCases := []struct {
		action         Action
		expectedResult bool
	}{
		{Action("*"), true},
		{Action(PutObjectAction), true},
		{Action("abcd"), false},
		{Action(PutObjectAction + "*"), true},
	}
	for _, testCase := range testCases {
		if testCase.action.IsValid() != testCase.expectedResult {
			t.Errorf("Test case failed: %s", testCase.action)
		}
	}
}
func TestAction_Match(t *testing.T) {
	testCases := []struct {
		name           string
		action         Action
		resource       Action
		expectedResult bool
	}{
		{"test1", Action("*"), Action(""), true},
		{"test1", Action("*"), Action(PutObjectAction), true},
		{"test1", Action("*"), Action("abcd"), true},
		{"test2", Action(PutObjectAction), Action(""), false},
		{"test2", Action(PutObjectAction), Action(PutObjectAction), true},
		{"test2", Action(PutObjectAction), Action("abcd"), false},
		{"test3", Action(""), Action("*"), false},
		{"test3", Action(""), Action(PutObjectAction), false},
		{"test3", Action(""), Action("abcd"), false},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.action.Match(testCase.resource) != testCase.expectedResult {
				t.Errorf("Test case failed: %s", testCase.action)
			}
		})
	}
}
