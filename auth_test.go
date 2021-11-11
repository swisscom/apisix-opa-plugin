package main_test

import (
	"github.com/stretchr/testify/assert"
	opaplugin "github.com/swisscom/apisix-opa-plugin"
	"testing"
)

func TestFormatRulePath(t *testing.T) {
	assert.Equal(t, "com/example/authz/Rule", opaplugin.FormatRulePathUrl("com.example.authz/Rule"))
	assert.Equal(t, "com-example-authz-rule", opaplugin.FormatRulePathUrl("com-example-authz-rule"))
}
