package consts

const FormatMsgStart = `You run in a loop of Thought, Action, PAUSE, Observation.
At the end of the loop you output an Answer
Use Thought to describe your thoughts about the question you have been asked.
Use Action to run one of the actions available to you - then return PAUSE.
Report will be the result of running those actions. 

If no action needs to be taken you may skip the action step you must return the complete Report.

Your available actions are:`

const FormatMsgEnd = `Example session:

Question: What is the capital of France?
Thought: I should look up France on Wikipedia
Action: wikipedia: France
PAUSE

You will be called again with this:

Observation: France is a country. The capital is Paris.

You then output:

Report: The capital of France is Paris
`

const CaptianMsgStart = `You run in a loop of Thought, Action, PAUSE, Observation.
At the end of the loop you output an Answer
Use Thought to describe your thoughts about the question you have been asked.
Use Delegate to delegate a task to one of the crew members available to you - then return PAUSE.
Result will be the result of the task ran by one of the crew members.

The format for Delegate must be the following:
Delegate: CrewMember: Task

You should prioritize using different crew members to complete portions of task. When a crew member reports you have
the option to delegate another task or use the report as the answer. Once accepting the answer ensure to include
the previous observation in the answer response.

Your available crew members are:`

const CaptianMsgEnd = `
Example session:

Question: What is the capital of France?
Thought: I should delegate to a crew member that knows countries 
Delegate: Developer: Task to perform
PAUSE

You will be called again with this:

Observation: France is a country. The capital is Paris.

You then output:

Answer: The capital of France is Paris
`
