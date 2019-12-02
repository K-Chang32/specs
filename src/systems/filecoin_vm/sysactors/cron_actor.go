package sysactors

import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import exitcode "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/exitcode"
import util "github.com/filecoin-project/specs/util"
import vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"

func (a *CronActorCode_I) Constructor(rt vmr.Runtime) InvocOutput {
	// Nothing. intentionally left blank.
	return rt.SuccessReturn()
}

func (a *CronActorCode_I) EpochTick(rt vmr.Runtime) InvocOutput {
	// Hook period actions in here.

	// a.Entries is basically a static registry for now, loaded
	// in the interpreter static registry.
	for _, entry := range a.Entries() {
		rt.SendCatchingErrors(&vmr.InvocInput_I{
			To_:     entry.ToAddr(),
			Method_: entry.MethodNum(),
			Params_: []util.Serialization{},
			Value_:  actor.TokenAmount(0),
		})
	}

	return rt.SuccessReturn()
}

func (a *CronActorCode_I) InvokeMethod(rt Runtime, method actor.MethodNum, params actor.MethodParams) InvocOutput {
	switch method {
	case actor.MethodConstructor:
		rt.Assert(len(params) == 0)
		return a.Constructor(rt)

	case actor.MethodCron:
		rt.Assert(len(params) == 0)
		return a.EpochTick(rt)

	default:
		rt.Abort(exitcode.SystemError(exitcode.InvalidMethod), "Invalid method")
		panic("")
	}
}
