package agent

import (
	"testing"

	"bytemind/internal/llm"
	planpkg "bytemind/internal/plan"
)

func TestShouldRepairBuildHandoffTurnAllowsLegitimateExecutionBlockers(t *testing.T) {
	state := planpkg.State{
		Phase: planpkg.PhaseExecuting,
		Steps: []planpkg.Step{
			{Title: "Implement feature", Status: planpkg.StepInProgress},
		},
	}
	reply := llm.Message{
		Role:    llm.RoleAssistant,
		Content: "I need your API key and repository path before I can continue.",
	}
	messages := []llm.Message{
		llm.NewUserTextMessage("start execution"),
	}

	if shouldRepairBuildHandoffTurn(planpkg.ModeBuild, state, turnIntentAskUser, reply, messages) {
		t.Fatalf("expected legitimate dependency clarification not to trigger build-handoff repair")
	}
}

func TestShouldRepairBuildHandoffTurnRepairsModeConfusionReplies(t *testing.T) {
	state := planpkg.State{
		Phase: planpkg.PhaseExecuting,
		Steps: []planpkg.Step{
			{Title: "Implement feature", Status: planpkg.StepInProgress},
		},
	}
	reply := llm.Message{
		Role:    llm.RoleAssistant,
		Content: "We are still in plan confirmation. Please send continue execution again.",
	}
	messages := []llm.Message{
		llm.NewUserTextMessage("start execution"),
	}

	if !shouldRepairBuildHandoffTurn(planpkg.ModeBuild, state, turnIntentFinalize, reply, messages) {
		t.Fatalf("expected mode-confusion blocker to trigger build-handoff repair")
	}
}

func TestShouldRepairBuildHandoffTurnDoesNotRepairMixedMessageWithRealBlocker(t *testing.T) {
	state := planpkg.State{
		Phase: planpkg.PhaseExecuting,
		Steps: []planpkg.Step{
			{Title: "Implement feature", Status: planpkg.StepInProgress},
		},
	}
	reply := llm.Message{
		Role:    llm.RoleAssistant,
		Content: "Please switch to build mode, but first provide the missing token so I can access the repo.",
	}
	messages := []llm.Message{
		llm.NewUserTextMessage("继续执行"),
	}

	if shouldRepairBuildHandoffTurn(planpkg.ModeBuild, state, turnIntentAskUser, reply, messages) {
		t.Fatalf("expected real missing dependency signals to bypass build-handoff repair")
	}
}
