package runtime

import actor "github.com/filecoin-project/specs/actors"
import abi "github.com/filecoin-project/specs/actors/abi"
import crypto "github.com/filecoin-project/specs/actors/crypto"
import exitcode "github.com/filecoin-project/specs/actors/runtime/exitcode"
import addr "github.com/filecoin-project/go-address"
import indices "github.com/filecoin-project/specs/actors/runtime/indices"
import cid "github.com/ipfs/go-cid"

// Runtime is the VM's internal runtime object.
// this is everything that is accessible to actors, beyond parameters.
type Runtime interface {
	CurrEpoch() abi.ChainEpoch

	// Randomness returns a (pseudo)random string for the given epoch and tag.
	GetRandomness(epoch abi.ChainEpoch) abi.RandomnessSeed

	// The address of the immediate calling actor.
	// Not necessarily the actor in the From field of the initial on-chain Message.
	// Always an ID-address.
	ImmediateCaller() addr.Address
	ValidateImmediateCallerIs(caller addr.Address)
	ValidateImmediateCallerInSet(callers []addr.Address)
	ValidateImmediateCallerAcceptAnyOfType(type_ abi.ActorCodeID)
	ValidateImmediateCallerAcceptAnyOfTypes(types []abi.ActorCodeID)
	ValidateImmediateCallerAcceptAny()
	ValidateImmediateCallerMatches(CallerPattern)

	// The address of the actor receiving the message. Always an ID-address.
	CurrReceiver() addr.Address

	// The actor who mined the block in which the initial on-chain message appears.
	// Always an ID-address.
	ToplevelBlockWinner() addr.Address

	AcquireState() ActorStateHandle

	SuccessReturn() InvocOutput
	ValueReturn([]byte) InvocOutput

	// Throw an error indicating a failure condition has occurred, from which the given actor
	// code is unable to recover.
	Abort(errExitCode exitcode.ExitCode, msg string)

	// Calls Abort with InvalidArguments_User.
	AbortArgMsg(msg string)
	AbortArg()

	// Calls Abort with InconsistentState_User.
	AbortStateMsg(msg string)
	AbortState()

	// Calls Abort with InsufficientFunds_User.
	AbortFundsMsg(msg string)
	AbortFunds()

	// Calls Abort with RuntimeAPIError.
	// For internal use only (not in actor code).
	AbortAPI(msg string)

	// Check that the given condition is true (and call Abort if not).
	Assert(bool)

	CurrentBalance() abi.TokenAmount
	ValueReceived() abi.TokenAmount

	// Look up the current values of several system-wide economic indices.
	CurrIndices() indices.Indices

	// Look up the code ID of a given actor address.
	GetActorCodeID(addr addr.Address) (ret abi.ActorCodeID, ok bool)

	// Run a (pure function) computation, consuming the gas cost associated with that function.
	// This mechanism is intended to capture the notion of an ABI between the VM and native
	// functions, and should be used for any function whose computation is expensive.
	Compute(ComputeFunctionID, args []interface{}) interface{}

	// Sends a message to another actor.
	// If the invoked method does not return successfully, this caller will be aborted too.
	SendPropagatingErrors(input InvocInput) InvocOutput
	Send(
		toAddr addr.Address,
		methodNum abi.MethodNum,
		params abi.MethodParams,
		value abi.TokenAmount,
	) InvocOutput
	SendQuery(
		toAddr addr.Address,
		methodNum abi.MethodNum,
		params abi.MethodParams,
	) []byte
	SendFunds(toAddr addr.Address, value abi.TokenAmount)

	// Sends a message to another actor, trapping an unsuccessful execution.
	// This may only be invoked by the singleton Cron actor.
	SendCatchingErrors(input InvocInput) (output InvocOutput, exitCode exitcode.ExitCode)

	// Computes an address for a new actor. The returned address is intended to uniquely refer to
	// the actor even in the event of a chain re-org (whereas an ID-address might refer to a
	// different actor after messages are re-ordered).
	// Always an ActorExec address.
	NewActorAddress() addr.Address

	// Creates an actor in the state tree, with empty state. May only be called by InitActor.
	CreateActor(
		// The new actor's code identifier.
		codeId abi.ActorCodeID,
		// Address under which the new actor's state will be stored. Must be an ID-address.
		address addr.Address,
	)

	// Deletes an actor in the state tree. May only be called by the actor itself,
	// or by StoragePowerActor in the case of StorageMinerActors.
	DeleteActor(address addr.Address)

	// Retrieves and deserializes an object from the store into o. Returns whether successful.
	IpldGet(c cid.Cid, o interface{}) bool
	// Serializes and stores an object, returning its CID.
	IpldPut(x interface{}) cid.Cid

	// Provides the system call interface.
	Syscalls() Syscalls
}

type Syscalls interface {
	// Verifies that a signature is valid for an address and plaintext.
	VerifySignature(
		signature crypto.Signature,
		signer addr.Address,
		plaintext []byte,
	) bool
	// Computes an unsealed sector CID (CommD) from its constituent piece CIDs (CommPs) and sizes.
	ComputeUnsealedSectorCID(sectorSize abi.SectorSize, pieces []abi.PieceInfo) (abi.UnsealedSectorCID, error)
	// Verifies a sector seal proof.
	VerifySeal(sectorSize abi.SectorSize, vi abi.SealVerifyInfo) bool
	// Verifies a proof of spacetime.
	VerifyPoSt(sectorSize abi.SectorSize, vi abi.PoStVerifyInfo) bool
}

type InvocInput struct {
	To     addr.Address
	Method abi.MethodNum
	Params abi.MethodParams
	Value  abi.TokenAmount
}

type InvocOutput struct {
	ReturnValue []byte
}

type ActorStateHandle interface {
	UpdateRelease(newStateCID actor.ActorSubstateCID)
	Release(checkStateCID actor.ActorSubstateCID)
	Take() actor.ActorSubstateCID
}

type ComputeFunctionID int64

const (
	Compute_VerifySignature = ComputeFunctionID(1)
)
