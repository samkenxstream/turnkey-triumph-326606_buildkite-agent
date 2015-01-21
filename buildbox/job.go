package buildkite

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// The Job struct uses strings for StartedAt and FinishedAt because
// if they were actual date objects, then when this struct is
// initialized they would have a default value of: 00:00:00.000000000.
// This causes problems for the Buildkite Agent API because it looks for
// the presence of values in these properties to determine if the build
// has finished.
type Job struct {
	ID string

	State string

	Env map[string]string

	Output string `json:"output,omitempty"`

	ExitStatus string `json:"exit_status,omitempty"`

	StartedAt string `json:"started_at,omitempty"`

	FinishedAt string `json:"finished_at,omitempty"`

	// If the job is currently being cancelled
	cancelled bool

	// The currently running process of the job
	process *Process
}

func (b Job) String() string {
	return fmt.Sprintf("Job{ID: %s, State: %s, StartedAt: %s, FinishedAt: %s, Process: %s}", b.ID, b.State, b.StartedAt, b.FinishedAt, b.process)
}

func (c *Client) JobNext() (*Job, error) {
	// Create a new instance of a job that will be populated
	// by the client.
	var job Job

	// Return the job.
	return &job, c.Get(&job, "jobs/next")
}

func (c *Client) JobFind(id string) (*Job, error) {
	// Create a new instance of a job that will be populated
	// by the client.
	var job Job

	// Find the job
	return &job, c.Get(&job, "jobs/"+id)
}

func (c *Client) JobAccept(job *Job) (*Job, error) {
	// Accept the job
	return job, c.Put(job, "jobs/"+job.ID+"/accept", job)
}

func (c *Client) JobUpdate(job *Job) (*Job, error) {
	// Create a new instance of a job that will be populated
	// with the updated data by the client
	var updatedJob Job

	// Return the job.
	return &updatedJob, c.Put(&updatedJob, "jobs/"+job.ID, job)
}

func (j *Job) Kill() error {
	if j.cancelled {
		// Already cancelled
	} else {
		Logger.Infof("Cancelling job %s", j.ID)
		j.cancelled = true

		if j.process != nil {
			j.process.Kill()
		} else {
			Logger.Errorf("No process to kill")
		}
	}

	return nil
}

func (j *Job) Run(agent *Agent) error {
	Logger.Infof("Starting job %s", j.ID)

	// Create the environment that the script will use
	env := []string{}

	// These are client specific
	env = append(env, fmt.Sprintf("BUILDKITE_AGENT_ENDPOINT=%s", agent.Client.URL))
	env = append(env, fmt.Sprintf("BUILDKITE_AGENT_DEBUG=%t", InDebugMode()))
	env = append(env, fmt.Sprintf("BUILDKITE_AGENT_ACCESS_TOKEN=%s", agent.Client.AuthorizationToken))
	env = append(env, fmt.Sprintf("BUILDKITE_AGENT_VERSION=%s", Version()))

	// We know the BUILDKITE_BIN_DIR dir, because it's the path to the
	// currently running file (there is only 1 binary)
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	env = append(env, fmt.Sprintf("BUILDKITE_BIN_DIR=%s", dir))

	// Add the rest environment variables from the API to the process
	for key, value := range j.Env {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}

	// Mark the build as started
	j.StartedAt = time.Now().Format(time.RFC3339)
	_, err := agent.Client.JobUpdate(j)
	if err != nil {
		// We don't care if the HTTP request failed here. We hope that it
		// starts working during the actual build.
	}

	// This callback is called every second the build is running. This lets
	// us do a lazy-person's method of streaming data to Buildkite.
	callback := func(process *Process) {
		j.Output = process.Output

		// Post the update to the API
		updatedJob, err := agent.Client.JobUpdate(j)
		if err != nil {
			// We don't really care if the job couldn't update at this point.
			// This is just a partial update. We'll just let the job run
			// and hopefully the host will fix itself before we finish.
			Logger.Warnf("Problem with updating job %s (%s)", j.ID, err)
		} else if updatedJob.State == "canceled" {
			j.Kill()
		}
	}

	// Initialze our process to run
	process := InitProcess(agent.BootstrapScript, env, agent.RunInPty, callback)

	// Store the process so we can cancel it later.
	j.process = process

	// Start the process. This will block until it finishes.
	err = process.Start()
	if err == nil {
		// Store the final output
		j.Output = j.process.Output
	} else {
		j.Output = fmt.Sprintf("%s", err)
	}

	// Mark the build as finished
	j.FinishedAt = time.Now().Format(time.RFC3339)
	j.ExitStatus = j.process.ExitStatus

	// Keep trying this call until it works. This is the most important one.
	for {
		_, err = agent.Client.JobUpdate(j)
		if err != nil {
			Logger.Errorf("Problem with updating final job information %s (%s)", j.ID, err)

			// How long should we wait until we try again?
			idleSeconds := 5

			// Sleep for a while
			sleepTime := time.Duration(idleSeconds*1000) * time.Millisecond
			time.Sleep(sleepTime)
		} else {
			break
		}
	}

	Logger.Infof("Finished job %s", j.ID)

	return nil
}
