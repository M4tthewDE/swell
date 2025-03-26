package jvm

import (
	"testing"

	"github.com/m4tthewde/swell/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestRunnerMain(t *testing.T) {
	log, err := logger.NewLogger()
	assert.Nil(t, err)

	ctx := logger.OnContext(t.Context(), log)

	runner := NewRunner([]string{"../../classes"})

	err = runner.RunMain(ctx, "Main")

	assert.Equal(t, "unknown instruction 1a\n\tjava.lang.StringUTF16.checkIndex()\n\tjava.lang.StringUTF16.charAt()\n\tjava.lang.String.charAt()\n\tjava.security.BasicPermission.init()\n\tjava.security.BasicPermission.<init>()\n\tjava.lang.reflect.ReflectPermission.<init>()\n\tjava.lang.invoke.MethodHandles.<clinit>()\n\tjava.lang.invoke.MethodHandles.<clinit>()\n\tjava.lang.invoke.MethodHandles.<clinit>()\n\tjdk.internal.access.SharedSecrets.<clinit>()\n\tjava.lang.System.setJavaLangAccess()\n\tjava.lang.System.initPhase1()\n\tjava.lang.System.<clinit>()\n\tMain.main()", err.Error())
	assert.Equal(t, "java/lang/StringUTF16", runner.classBeingInitialized)
	assert.Equal(t, map[string]struct{}{
		"Main":             {},
		"java/lang/Class":  {},
		"java/lang/Object": {},
		"java/lang/String": {},
		"java/lang/String$CaseInsensitiveComparator": {},
		"java/lang/StringUTF16":                      {},
		"java/lang/System$2":                         {},
		"java/lang/invoke/MemberName":                {},
		"java/lang/invoke/MemberName$Factory":        {},
		"java/lang/reflect/ReflectPermission":        {},
		"java/security/BasicPermission":              {},
		"java/security/Permission":                   {},
	}, runner.initializedClasses)
	assert.Equal(t, 0, runner.pc)
}
