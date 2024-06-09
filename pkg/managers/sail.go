package managers

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/r4stl1n/mbop/pkg/api/llm"
	"github.com/r4stl1n/mbop/pkg/structs"
	"github.com/r4stl1n/mbop/pkg/tools"
	"github.com/r4stl1n/mbop/pkg/tools/wiki"
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

	managerAgent *structs.Agent

	agents []*structs.Agent
	tools  map[string]tools.Tool

	actionRegex   *regexp.Regexp
	delegateRegex *regexp.Regexp
	reportRegex   *regexp.Regexp

	openaiClient *llm.OpenAIAPI
}

func (s *SailManager) Init(model string, task string, agentDir string) (*SailManager, error) {

	openaiClient, openaiClientError := new(llm.OpenAIAPI).Init()

	actionReg, actionRegError := regexp.Compile(`^Action: (\w+): (.*)$`)

	if actionRegError != nil {
		return nil, actionRegError
	}

	delegateReg, delegateRegError := regexp.Compile(`^Delegate: (.*): (.*)$`)

	if delegateRegError != nil {
		return nil, delegateRegError
	}

	reportReg, reportRegError := regexp.Compile(`^Report: (.*)$`)

	if reportRegError != nil {
		return nil, reportRegError
	}

	*s = SailManager{
		model:    model,
		task:     task,
		agentDir: agentDir,

		actionRegex:   actionReg,
		delegateRegex: delegateReg,
		reportRegex:   reportReg,

		openaiClient: openaiClient,
		tools: map[string]tools.Tool{
			wiki.Wikipedia{}.Name(): wiki.Wikipedia{},
		},
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

func (s *SailManager) findAgent(input string) (*structs.Agent, string, error) {

	// extract tool and tool name
	agentSplit := strings.Split(input, ": ")

	for _, x := range s.agents {

		if x.Role == agentSplit[1] {
			return x, agentSplit[2], nil
		}
	}

	return nil, "", fmt.Errorf("no agent with role %s found", agentSplit[1])
}

func (s *SailManager) runTool(input string) (string, error) {
	// extract tool and tool name
	toolSplit := strings.Split(input, ": ")

	// Run the specified tool
	toolResponse, toolResponseError := s.tools[toolSplit[1]].Run(toolSplit[2])
	if toolResponseError != nil {
		return "", toolResponseError
	}

	return toolResponse, nil
}

func (s *SailManager) processAgents() error {

	activeAgent := s.managerAgent

	activeAgent.Context.Add(llm.Message{
		Role:    "user",
		Content: fmt.Sprintf("%s\n%s", activeAgent.ConstructCaptainPrompt(s.agents), "Current Task: "+s.task),
	})

	//color.Red(fmt.Sprintf("Role: %s\nContent: %s\n\n", activeAgent.Role, activeAgent.Context.Context[0].Content))
	zap.L().Info("sail process started", zap.String("agent", activeAgent.Role), zap.String("task", s.task))

	for i := 0; i < 5; i++ {

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

		//color.Yellow(fmt.Sprintf("Response:\n%s\n", completion))

		// Make sure we have a PAUSE or Answer in our response
		for _, x := range completionStrings {

			actionMatch := s.actionRegex.Match([]byte(x))
			delegateMatch := s.delegateRegex.Match([]byte(x))
			reportMatch := s.reportRegex.Match([]byte(x))

			// Check if there is a answer in the response
			if strings.Contains(completion, "Answer:") {
				answerSplit := strings.Split(completion, "Answer:")
				color.Green(fmt.Sprintf("\n\nAnswer: %s", answerSplit[1]))
				return nil
			}

			// Found action to perform
			if actionMatch {

				zap.L().Info("update", zap.String("agent", activeAgent.Role), zap.String("data", x))

				toolResponse, toolError := s.runTool(x)

				if toolError != nil {
					return toolError
				}

				activeAgent.Context.Add(llm.Message{
					Role:    "user",
					Content: fmt.Sprintf("Observation: %s", toolResponse),
				})

			}

			if delegateMatch {

				zap.L().Info("update", zap.String("agent", activeAgent.Role), zap.String("data", x))

				// Find the agent
				foundAgent, foundAgentTask, foundAgentError := s.findAgent(x)

				if foundAgentError != nil {
					return foundAgentError
				}

				// Make it the active agent
				activeAgent = foundAgent

				// If no context we want to set the initial
				activeAgent.Context.Add(llm.Message{
					Role: "user",
					Content: fmt.Sprintf("%s\n%s", activeAgent.ConstructAgentPrompt(s.tools),
						"Current Task: "+foundAgentTask),
				})

			}

			if reportMatch {

				zap.L().Info("update", zap.String("agent", activeAgent.Role), zap.String("data", x))

				// Swap over to the manager agent as active
				activeAgent = s.managerAgent

				// extract tool and tool name
				toolSplit := strings.Split(x, ": ")

				activeAgent.Context.Add(llm.Message{
					Role:    "user",
					Content: fmt.Sprintf("Observation: %s", toolSplit[1]),
				})
			}

		}

		//color.Cyan(activeAgent.Context.PrintHistory())

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
