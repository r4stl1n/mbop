package managers

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/r4stl1n/mbop/pkg/api/llm"
	"github.com/r4stl1n/mbop/pkg/consts"
	"github.com/r4stl1n/mbop/pkg/structs"
	"github.com/r4stl1n/mbop/pkg/tools"
	"github.com/r4stl1n/mbop/pkg/tools/wiki"
	"github.com/r4stl1n/mbop/pkg/util"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type SailManager struct {
	model    string
	task     string
	agentDir string

	characterTrim int
	managerAgent  *structs.Agent

	agents []*structs.Agent
	tools  map[string]tools.Tool

	thoughtRegex *regexp.Regexp
	commandRegex *regexp.Regexp

	openaiClient *llm.OpenAIAPI
	utils        *util.Utils
}

func (s *SailManager) Init(model string, task string, agentDir string) (*SailManager, error) {

	openaiClient, openaiClientError := new(llm.OpenAIAPI).Init()

	thoughtReg, thoughtRegError := regexp.Compile(`^Thought: (.*)$`)

	if thoughtRegError != nil {
		return nil, thoughtRegError
	}

	commandReg, commandRegError := regexp.Compile(`^CrewResponse: (.*)$`)

	if commandRegError != nil {
		return nil, commandRegError
	}

	*s = SailManager{
		model:         model,
		task:          task,
		agentDir:      agentDir,
		characterTrim: 40,
		thoughtRegex:  thoughtReg,
		commandRegex:  commandReg,

		openaiClient: openaiClient,
		tools: map[string]tools.Tool{
			wiki.Wikipedia{}.Name(): wiki.Wikipedia{},
		},

		utils: new(util.Utils).Init(),
	}

	return s, openaiClientError
}

func (s *SailManager) loadAgents() error {

	folder, folderError := os.Open(s.agentDir)
	if folderError != nil {
		return folderError
	}

	files, filesError := folder.Readdir(0)
	if filesError != nil {
		return filesError
	}

	for _, v := range files {
		if v.IsDir() {
			continue
		}

		fileExt := filepath.Ext(v.Name())

		if fileExt != ".json" {
			continue
		}

		fileData, readFileError := os.ReadFile(s.agentDir + "/" + v.Name())
		if readFileError != nil {
			zap.L().Error("failed to read file", zap.String("file", s.agentDir+"/"+v.Name()), zap.Error(readFileError))
			continue
		}

		agent := structs.Agent{}
		unmarshallError := json.Unmarshal(fileData, &agent)

		if unmarshallError != nil {
			zap.L().Error("failed to unmarshall json", zap.String("file", s.agentDir+"/"+v.Name()), zap.Error(unmarshallError))
			continue
		}

		zap.L().Info("found agent", zap.String("file", s.agentDir+"/"+v.Name()), zap.String("name", agent.Role))

		if agent.IsCaptain {
			s.managerAgent = &agent
		} else {
			s.agents = append(s.agents, &agent)
		}
	}

	return nil

}

func (s *SailManager) findAgent(agentRole string) (*structs.Agent, error) {

	// extract tool and tool name

	for _, x := range s.agents {

		if strings.ToLower(x.Role) == strings.ToLower(agentRole) {
			return x, nil
		}
	}

	return nil, fmt.Errorf("no agent with role %s found", agentRole)
}

func (s *SailManager) runTool(toolName string, toolData string) (string, error) {

	// Run the specified tool
	zap.L().Info("attempting to use tool", zap.String("toolName", toolName))

	tool, ok := s.tools[toolName]

	if !ok {
		if !(strings.ToLower(toolName) == "none" || strings.ToLower(toolName) == "nil") {
			zap.L().Warn("attempt to use unknown tool", zap.String("toolName", toolName))
		}

		return "no tool output", nil
	}

	toolResponse, toolResponseError := tool.Run(toolData)

	if toolResponseError != nil {
		return "no tool output", toolResponseError
	}

	return toolResponse, nil
}

func (s *SailManager) processAgents() error {

	activeAgent := s.managerAgent
	previousReport := ""

	activeAgent.Context.Add(llm.Message{
		Role:    "user",
		Content: fmt.Sprintf("%s\n%s", activeAgent.ConstructCaptainPrompt(s.agents), "Current Task: "+s.task),
	})

	color.Cyan(fmt.Sprintf("Role: %s\nContent: %s\n\n", activeAgent.Role, activeAgent.Context.Context[0].Content))
	zap.L().Info("sail process started", zap.String("agent", activeAgent.Role), zap.String("task", s.task))

	for i := 0; i < 20; i++ {

		completion, _, err := s.openaiClient.GetCompletion(activeAgent.Context)

		if err != nil {
			return err
		}

		activeAgent.Context.Add(llm.Message{
			Role:    "assistant",
			Content: completion,
		})

		// Quick housekeeping clean up
		completion = strings.Split(completion, "PAUSE")[0]
		completion = strings.TrimSpace(completion)
		after, _ := strings.CutPrefix(completion, "{")
		before, _ := strings.CutSuffix(after, "}")
		completion = "{" + before + "}"
		completion = strings.ReplaceAll(completion, "\n", "\\n")

		color.Yellow(fmt.Sprintf("Response:\n%s\n", completion))
		failure := false

		var command structs.CrewResponse
		unmarshallError := json.Unmarshal([]byte(completion), &command)

		if unmarshallError != nil {
			zap.L().Debug("invalid command format received", zap.String("command",
				s.utils.EllipticalTruncate(completion, s.characterTrim)))
			failure = true
		}

		switch strings.ToLower(command.Type) {

		case "action":

			zap.L().Info("action", zap.String("agent", activeAgent.Role),
				zap.String("data", s.utils.EllipticalTruncate(completion, s.characterTrim)))

			toolResponse, toolError := s.runTool(command.Tool, command.Data)

			if toolError != nil {
				zap.L().Warn("invalid tool requested", zap.String("tool", command.Tool))
				failure = true
				break
			}

			activeAgent.Context.Add(llm.Message{
				Role:    "user",
				Content: fmt.Sprintf("Observation: %s", toolResponse),
			})

		case "delegate":
			zap.L().Info("delegate", zap.String("agent", activeAgent.Role),
				zap.String("data", s.utils.EllipticalTruncate(completion, s.characterTrim)))

			// Find the agent
			foundAgent, foundAgentError := s.findAgent(command.Crew)

			if foundAgentError != nil {
				zap.L().Debug("invalid agent requested", zap.String("agent", command.Crew))
				failure = true
				break
			}

			// Make it the active agent
			activeAgent = foundAgent

			// If no context we want to set the initial
			activeAgent.Context.Add(llm.Message{
				Role: "user",
				Content: fmt.Sprintf("%s\n%s\n\n%s", activeAgent.ConstructAgentPrompt(s.tools),
					"Relevant Information: "+previousReport,
					"Current Task: "+command.Data),
			})

			zap.L().Info("changing", zap.String("agent", activeAgent.Role))

		case "report":

			previousAgent := activeAgent

			// Swap over to the manager agent as active
			activeAgent = s.managerAgent

			previousReport = command.Response

			activeAgent.Context.Add(llm.Message{
				Role:    "user",
				Content: fmt.Sprintf("Result: %s", command.Response),
			})

			zap.L().Info("reporting", zap.String("previousAgent", previousAgent.Role),
				zap.String("newAgent", activeAgent.Role))

		case "answer":
			color.Green(fmt.Sprintf("\nAnswer: \n%s", command.Result))
			color.Green(fmt.Sprintf("\nReport: \n%s", previousReport))
			return nil

		}

		// Failed to get valid response
		if failure == true {
			activeAgent.Context.Add(llm.Message{
				Role:    "user",
				Content: consts.IncorrectFormatMsg,
			})
			zap.L().Debug("attempting to query again, desired response format invalid")
		}

		color.Cyan(activeAgent.Context.PrintLatestHistory())

	}

	return nil
}

func (s *SailManager) Run() error {

	loadAgentsError := s.loadAgents()

	if loadAgentsError != nil {
		return loadAgentsError
	}

	return s.processAgents()
}
