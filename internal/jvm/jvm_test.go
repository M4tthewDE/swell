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

	assert.Equal(t, "invalid variable type: Int=16\n\tjava.lang.AbstractStringBuilder.<init>()\n\tjava.lang.StringBuilder.<init>()\n\tjava.lang.String.checkIndex()\n\tjava.lang.StringUTF16.checkIndex()\n\tjava.lang.StringUTF16.charAt()\n\tjava.lang.String.charAt()\n\tjava.security.BasicPermission.init()\n\tjava.security.BasicPermission.<init>()\n\tjava.lang.reflect.ReflectPermission.<init>()\n\tjava.lang.invoke.MethodHandles.<clinit>()\n\tjava.lang.invoke.MethodHandles.<clinit>()\n\tjava.lang.invoke.MethodHandles.<clinit>()\n\tjdk.internal.access.SharedSecrets.<clinit>()\n\tjava.lang.System.setJavaLangAccess()\n\tjava.lang.System.initPhase1()\n\tjava.lang.System.<clinit>()\n\tMain.main()", err.Error())
	assert.Equal(t, "java/lang/AbstractStringBuilder", runner.classBeingInitialized)
	assert.Equal(t, 1, runner.pc)
}
