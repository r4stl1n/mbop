package structs

import (
	"fmt"
	"github.com/r4stl1n/mbop/pkg/api/llm"
	"github.com/r4stl1n/mbop/pkg/consts"
	"github.com/r4stl1n/mbop/pkg/tools"
)

type Agent struct {
	Role      string
	Goal      string
	Persona   string
	IsCaptain bool
	Context   llm.CompletionHistory
}

func ConvertAgentArrayToPrompt(agents []*Agent) string {

	prompt := ""

	for _, agent := range agents {
		prompt += "\n"
		prompt += agent.Role
	}

	return prompt
}

func (a *Agent) ConstructCaptainPrompt(agentList []*Agent) string {

	agentsConv := ConvertAgentArrayToPrompt(agentList)

	systemPrompt := fmt.Sprintf("%s\n%s\nYour personal goal is: %s\n\n"+
		"%s\n%s\n%s", a.Role, a.Persona, a.Goal, consts.CaptianMsgStart, agentsConv, consts.CaptianMsgEnd)

	return systemPrompt
}

func (a *Agent) ConstructAgentPrompt(toolList map[string]tools.Tool) string {

	toolsConv := tools.ConvertToolArrayToPrompt(toolList)

	systemPrompt := fmt.Sprintf("%s\n%s\nYour personal goal is: %s\n\n"+
		"%s\n%s\n%s", a.Role, a.Persona, a.Goal, consts.FormatMsgStart, toolsConv, consts.FormatMsgEnd)

	return systemPrompt
}
