package managers

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/r4stl1n/mbop/pkg/api/llm"
	"github.com/r4stl1n/mbop/pkg/consts"
	"github.com/r4stl1n/mbop/pkg/structs"
	"github.com/r4stl1n/mbop/pkg/tools"
	"github.com/r4stl1n/mbop/pkg/tools/directory"
	"github.com/r4stl1n/mbop/pkg/tools/duckduckgo"
	"github.com/r4stl1n/mbop/pkg/tools/file"
	mathtools "github.com/r4stl1n/mbop/pkg/tools/math"
	"github.com/r4stl1n/mbop/pkg/tools/terminal"
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
	debug    bool

	characterTrim int
	managerAgent  *structs.Agent

	agents []*structs.Agent
	tools  map[string]tools.Tool

	thoughtRegex *regexp.Regexp
	commandRegex *regexp.Regexp

	openaiClient *llm.OpenAIAPI
	utils        *util.Utils
}

func (s *SailManager) Init(debug bool, model string, task string, agentDir string) (*SailManager, error) {

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
		debug:         debug,
		model:         model,
		task:          task,
		agentDir:      agentDir,
		characterTrim: 40,
		thoughtRegex:  thoughtReg,
		commandRegex:  commandReg,

		openaiClient: openaiClient,
		tools: map[string]tools.Tool{
			wiki.Wikipedia{}.Name():      wiki.Wikipedia{},
			mathtools.Math{}.Name():      mathtools.Math{},
			directory.Directory{}.Name(): directory.Directory{},
			file.File{}.Name():           file.File{},
			terminal.Terminal{}.Name():   terminal.Terminal{},
			duckduckgo.DuckDuckGo{}.Name(): duckduckgo.DuckDuckGo{},
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
			zap.L().Error("agent file read failed", 
				zap.String("filename", v.Name()),
				zap.String("path", s.agentDir), 
				zap.Error(readFileError))
			continue
		}

		agent := structs.Agent{Context: llm.CompletionHistory{
			Model: s.model,
		}}

		unmarshallError := json.Unmarshal(fileData, &agent)

		if unmarshallError != nil {
			zap.L().Error("agent JSON parsing failed", 
				zap.String("filename", v.Name()),
				zap.String("path", s.agentDir), 
				zap.Error(unmarshallError))
			continue
		}

		zap.L().Info("agent loaded", 
			zap.String("name", agent.Role),
			zap.String("filename", v.Name()),
			zap.Bool("isCaptain", agent.IsCaptain))

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
	zap.L().Info("tool execution started", 
		zap.String("tool", toolName),
		zap.String("data", s.utils.EllipticalTruncate(toolData, s.characterTrim)))

	tool, ok := s.tools[toolName]

	if !ok {
		if !(strings.ToLower(toolName) == "none" || strings.ToLower(toolName) == "nil") {
			zap.L().Warn("unknown tool requested", 
				zap.String("tool", toolName))
		}

		return "no tool output", nil
	}

	toolResponse, toolResponseError := tool.Run(toolData)

	if toolResponseError != nil {
		zap.L().Error("tool execution failed",
			zap.String("tool", toolName),
			zap.String("data", s.utils.EllipticalTruncate(toolData, s.characterTrim)),
			zap.Error(toolResponseError))
		return "no tool output", toolResponseError
	}

	zap.L().Info("tool execution completed", 
		zap.String("tool", toolName),
		zap.Int("responseLength", len(toolResponse)))
	return toolResponse, nil
}

func (s *SailManager) extractJSON(input string) (string, error) {
	// Pattern to match any JSON object, regardless of surrounding content
	// Looks for content between { and } that could be valid JSON
	pattern := "{[\\s\\S]*?\"[\\s\\S]*?:[\\s\\S]*?\"[\\s\\S]*?}+"
	re := regexp.MustCompile(pattern)

	matches := re.FindString(input)
	if matches == "" {
		return "", fmt.Errorf("no JSON content found")
	}

	// Return the matched JSON
	return strings.TrimSpace(matches), nil
}

func (s *SailManager) processAgents() error {

	activeAgent := s.managerAgent
	previousReport := ""

	activeAgent.Context.Add(llm.Message{
		Role:    "user",
		Content: fmt.Sprintf("%s\n%s", activeAgent.ConstructCaptainPrompt(s.agents), "Current Task: "+s.task),
	})

	if s.debug {
		color.Cyan(fmt.Sprintf("Role: %s\nContent: %s\n\n", activeAgent.Role, activeAgent.Context.Context[0].Content))
	}

	zap.L().Info("sail process started", 
		zap.String("captain", activeAgent.Role), 
		zap.String("task", s.task),
		zap.Int("agentCount", len(s.agents)))

	for i := 0; i < 100; i++ {
		failure := false

		rawCompletion, _, err := s.openaiClient.GetCompletion(activeAgent.Context)

		if err != nil {
			zap.L().Error("completion request failed", 
				zap.String("agent", activeAgent.Role),
				zap.Error(err))
			return err
		}

		completion, err := s.extractJSON(rawCompletion)

		if err != nil {
			zap.L().Debug("JSON extraction failed", 
				zap.String("agent", activeAgent.Role),
				zap.String("rawResponse", s.utils.EllipticalTruncate(rawCompletion, s.characterTrim)),
				zap.Error(err))
			failure = true
		}

		if !failure {
			activeAgent.Context.Add(llm.Message{
				Role:    "assistant",
				Content: completion,
			})

			if s.debug {
				color.Yellow(fmt.Sprintf("Response:\n%s\n", completion))
			}

			var command structs.CrewResponse
			unmarshallError := json.Unmarshal([]byte(completion), &command)

			if unmarshallError != nil {
				zap.L().Debug("command parsing failed", 
					zap.String("agent", activeAgent.Role),
					zap.String("completion", s.utils.EllipticalTruncate(completion, s.characterTrim)),
					zap.Error(unmarshallError))
				failure = true
			}

			if !failure {
				zap.L().Info("agent thinking", 
					zap.String("agent", activeAgent.Role), 
					zap.String("thought", command.Thought))

				switch strings.ToLower(command.Type) {
				case "action":
					zap.L().Info("agent action", 
						zap.String("agent", activeAgent.Role),
						zap.String("tool", command.Tool),
						zap.String("data", s.utils.EllipticalTruncate(command.Data, s.characterTrim)))

					toolResponse, toolError := s.runTool(command.Tool, command.Data)

					if toolError != nil {
						// Note: detailed error already logged in runTool
						failure = true
						break
					}

					activeAgent.Context.Add(llm.Message{
						Role:    "user",
						Content: fmt.Sprintf("Observation: %s", toolResponse),
					})

				case "delegate":
					zap.L().Info("agent delegation", 
						zap.String("from", activeAgent.Role),
						zap.String("to", command.Crew),
						zap.String("task", s.utils.EllipticalTruncate(command.Data, s.characterTrim)))

					// Find the agent
					foundAgent, foundAgentError := s.findAgent(command.Crew)

					if foundAgentError != nil {
						zap.L().Warn("agent delegation failed", 
							zap.String("from", activeAgent.Role),
							zap.String("to", command.Crew),
							zap.Error(foundAgentError))
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

				case "report":
					previousAgent := activeAgent

					// Swap over to the manager agent as active
					activeAgent = s.managerAgent

					previousReport = command.Response

					activeAgent.Context.Add(llm.Message{
						Role:    "user",
						Content: fmt.Sprintf("Result: %s", command.Response),
					})

					zap.L().Info("agent reporting", 
						zap.String("from", previousAgent.Role),
						zap.String("to", activeAgent.Role),
						zap.Int("responseLength", len(command.Response)))

				case "answer":
					zap.L().Info("task completed",
						zap.String("agent", activeAgent.Role),
						zap.Int("resultLength", len(command.Result)))
					color.Green(fmt.Sprintf("\nAnswer: \n%s", command.Result))
					color.Green(fmt.Sprintf("\nReport: \n%s", previousReport))
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
			zap.L().Debug("retry with format correction", 
				zap.String("agent", activeAgent.Role),
				zap.Int("attempt", i+1))
		}

		if s.debug {
			color.Cyan(activeAgent.Context.PrintLatestHistory())
		}
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
