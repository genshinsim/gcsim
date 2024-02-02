package backend

import (
	"regexp"
	"testing"
)

func TestValidateLink(t *testing.T) {
	s := &Store{
		Config: Config{
			LinkValidationRegex: regexp.MustCompile(`https://\S+.app/\S+/(\S+)$`),
		},
	}

	id, err := s.validateLink("https://gcsim.app/v3/viewer/share/23543b84-045a-47bb-be48-cdb242b413b6")
	if err != nil {
		t.Error(err)
	}
	if id != "23543b84-045a-47bb-be48-cdb242b413b6" {
		t.Errorf("expecting id 23543b84-045a-47bb-be48-cdb242b413b6, got %v", id)
	}

	_, err = s.validateLink("gcsim.app/v4/viewer/share/blah")
	if err == nil {
		t.Errorf("expecting invalid link error, got nil")
	}
}
