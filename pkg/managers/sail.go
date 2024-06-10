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

	commandReg, commandRegError := regexp.Compile(`^Command: (.*)$`)

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

		if x.Role == agentRole {
			return x, nil
		}
	}

	return nil, fmt.Errorf("no agent with role %s found", agentRole)
}

func (s *SailManager) runTool(toolName string, toolData string) (string, error) {

	// Run the specified tool
	zap.L().Info("attempting to use tool", zap.String("toolName", toolName))

	toolResponse, toolResponseError := s.tools[toolName].Run(toolData)

	if toolResponseError != nil {
		return "", toolResponseError
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

		// Remove the second half
		completion = strings.Split(completion, "PAUSE")[0]

		activeAgent.Context.Add(llm.Message{
			Role:    "assistant",
			Content: completion,
		})

		completionStrings := strings.Split(completion, "\n")

		color.Yellow(fmt.Sprintf("Response:\n%s\n", completion))

		failure := false
		done := false
		// Make sure we have a PAUSE or Answer in our response
		for _, x := range completionStrings {

			if failure || done {
				break
			}

			if !strings.Contains(completion, "Command:") {
				failure = true
				break
			}

			if s.thoughtRegex.Match([]byte(x)) {
				thoughtSplit := strings.Split(x, "Thought: ")
				zap.L().Info("process", zap.String("agent", activeAgent.Role), zap.String("thought", thoughtSplit[1]))
			}

			if s.commandRegex.Match([]byte(x)) {

				commandSplit := strings.Split(x, "Command:")

				var command structs.Command
				unmarshallError := json.Unmarshal([]byte(commandSplit[1]), &command)

				if unmarshallError != nil {
					zap.L().Debug("invalid command format received", zap.String("command",
						s.utils.EllipticalTruncate(x, s.characterTrim)))
					failure = true
					continue
				}

				switch strings.ToLower(command.Type) {

				case "action":

					zap.L().Info("action", zap.String("agent", activeAgent.Role),
						zap.String("data", s.utils.EllipticalTruncate(x, s.characterTrim)))

					toolResponse, toolError := s.runTool(command.Tool, command.Data)

					if toolError != nil {
						zap.L().Warn("invalid tool requested", zap.String("tool", command.Tool))
						failure = true
						continue
					}

					activeAgent.Context.Add(llm.Message{
						Role:    "user",
						Content: fmt.Sprintf("Observation: %s", toolResponse),
					})

					done = true

				case "delegate":
					zap.L().Info("delegate", zap.String("agent", activeAgent.Role),
						zap.String("data", s.utils.EllipticalTruncate(x, s.characterTrim)))

					// Find the agent
					foundAgent, foundAgentError := s.findAgent(command.Crew)

					if foundAgentError != nil {
						zap.L().Debug("invalid agent requested", zap.String("agent", command.Crew))
						failure = true
						continue
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

					done = true

				case "report":

					previousAgent := activeAgent

					// Swap over to the manager agent as active
					activeAgent = s.managerAgent

					previousReport = command.Data

					activeAgent.Context.Add(llm.Message{
						Role:    "user",
						Content: fmt.Sprintf("Observation: %s", command.Data),
					})

					zap.L().Info("reporting", zap.String("previousAgent", previousAgent.Role),
						zap.String("newAgent", activeAgent.Role))

					done = true

				case "answer":
					color.Green(fmt.Sprintf("\n\nAnswer: %s", command.Data))
					return nil

				}

			}
		}

		// Failed to get valid response
		if failure {
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
