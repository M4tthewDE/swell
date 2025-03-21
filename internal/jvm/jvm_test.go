package jvm

import (
	"context"
	"testing"

	"github.com/m4tthewde/swell/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestRunnerMain(t *testing.T) {
	log, err := logger.NewLogger()
	assert.Nil(t, err)

	ctx := logger.OnContext(context.Background(), log)

	runner := NewRunner([]string{"../../classes"})

	err = runner.RunMain(ctx, "Main")

	assert.Equal(t, "native method desiredAssertionStatus0 in java/lang/Class not implemented\n\tjava.lang.Class.desiredAssertionStatus()\n\tjava.lang.invoke.MethodHandles.<clinit>()\n\tjdk.internal.access.SharedSecrets.<clinit>()\n\tjava.lang.System.setJavaLangAccess()\n\tjava.lang.System.initPhase1()\n\tjava.lang.System.<clinit>()\n\tMain.main()", err.Error())
	assert.Equal(t, "java/lang/Class", runner.classBeingInitialized)
	assert.Equal(t, map[string]struct{}{
		"Main":               {},
		"java/lang/Class":    {},
		"java/lang/Object":   {},
		"java/lang/System$2": {},
	}, runner.initializedClasses)
	assert.Equal(t, 13, runner.pc)
}
