# Golang	Interview	Task

# Task	Description
You	are	tasked	with	creating	a	job	processing	server that	accepts	and	processes	jobs	
asynchronously.	Each	job	has	a	random	processing	time	between	5	and	30	seconds.	The	server	
should	handle	multiple	concurrent	clients	gracefully,	process	jobs	using	a	worker	pool,	and	provide	
status	tracking	for	each	job.

# Requirements

## Server	Features

1. HTTP	Endpoints:
• POST	/job:
• Accepts	a	job	payload	(string)	in	JSON	format	and	adds	it	to	a	queue.
• Example	request:	{	"payload":	"Process	this	job!"	}
• Respond	with	202	Accepted and	the	job	ID:	{	"job_id":	1	}.
• GET	/status/{job_id}:
• Returns	the	status	of	a	job	(pending,	processing,	completed).
• Example	response:	{	"job_id":	1,	"status":	"completed"	}.

2. Concurrency	with	Channels:
• Use	a	channel to	manage	the	job	queue.
• Create	a	worker	pool to	process	jobs	concurrently.	The	number	of	workers	should	be	
configurable.

3. Interfaces:
• Define	an	interface	JobProcessor with	a	method	Process(job	Job)	error.
• Implement	a	concrete	StringJobProcessor that	processes	string-based	jobs.

4. Random	Job	Processing	Time:
• Each	job	should	have	a	random	processing	time	between	5	and	30	seconds.

5. Status	Tracking:
• Maintain	job	statuses	(pending,	processing,	completed)	in	a	thread-safe manner	(e.g.,	using	
sync.Map).

6. Graceful	Shutdown:
• Ensure	the	server	can	shut	down	cleanly,	completing	all	in-progress	jobs.

## Client	Simulation
1. Simulate	Clients:
• Create	a	function	to	simulate	multiple	clients	(e.g.,	10	clients)	sending	requests	to	the	
server.
• Each	client	should:
• Send	multiple	POST	/job requests.
• Periodically	query	job	statuses	using	GET	/status/{job_id}.

2. Load	Testing:
• Test	the	server	with	a	high	volume	of	concurrent	requests	to	evaluate	its	robustness.

# Deliverables
1. Server	Code:
• Complete	implementation	of	the	HTTP	server	with	the	specified	features as	an	executable

2. Client	Simulation:
• An	executable to	simulate	multiple	concurrent	clients	interacting	with	the	server and	a	way	
to	query	job	statuses

3. Documentation:
• Provide	instructions	for	running	the	server	and	simulation.
• Explain	assumptions	and	trade-offs	made	during	implementation.

4. Error	Handling:
• Ensure	the	server	handles	invalid	inputs	and	errors	gracefully