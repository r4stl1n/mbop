package consts

const FormatMsgStart = `You run in a loop of Thought, Command, PAUSE, Observation.
At the end of the loop you output an report in a json format.
Use Thought to describe your thoughts about the question you have been asked.
Use Command to run one of the actions available to you - then return PAUSE.
Report will be the result of running those actions. 

If no action needs to be taken you may skip the action step you must return the complete Report. In the correct json
format.

There should be no new lines in the json output.

Report needs to be in the following format:

Command: {"type": "report", "data":"The capital of France is Paris"}

Your available actions are:`

const FormatMsgEnd = `Example session:

Question: What is the capital of France?
Thought: I should look up France on Wikipedia
Command: {"type":"action", "tool":"wikipedia", "data": "France"}
PAUSE

Your action results will be in the following format:
Observation: France is a country. The capital is Paris.

You then will process the data and output in the following format:

Command: {"type": "report", "data":"The capital of France is Paris"}

The last line in your response should always be in the following format. With nothing else added. All created data should be 
included in the single json formatted response. No new lines should exist in your json response.

Command: json response

`

const CaptianMsgStart = `You run in a loop of Thought, Action, PAUSE, Observation.
At the end of the loop you output an Answer
Use Thought to describe your thoughts about the question you have been asked.
Use Delegate to delegate a task to one of the crew members available to you - then return PAUSE.
Result will be the result of the task ran by one of the crew members.

The format for Delegate must be the following:
Command: {"type": "delegate", "crew":"developer","data":"task to perform"}

There should be no new lines in the json output. 

You should prioritize using different crew members to complete portions of task. When a crew member reports you have
the option to delegate another task or use the report as the answer.

Your available crew members are:`

const CaptianMsgEnd = `
Example session:

Question: What is the capital of France?
Thought: I should delegate to a crew member that knows countries 
Command: {"type": "delegate", "crew":"developer","data":"task to perform"}

PAUSE

You will be called again with this:

Observation: France is a country. The capital is Paris.

You must then output the task results in the following format:

Command: {"type": "answer", "data":"The capital of France is Paris"}

The last line in your response should always be in the following format. Task results should
be includes in the json response. No new lines should exist in your json response.

Command: json response
`

const IncorrectFormatMsg = `Your last json formatted response was not valid please ensure it is in the following format:
Command: {"type": "answer", "data":"The capital of France is Paris"}
`
