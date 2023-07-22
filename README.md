# Todoist ager

Very basic utility to apply age-based labels to all tasks in Todoist.

The following labels are applied:

* Recurring tasks: _none_
* Completed tasks: _none_
* Tasks created within two weeks: _none_
* Tasks created between two and six weeks ago: `age-weeks`
* Tasks created between 6 weeks and 1 year ago: `age-months`
* Tasks older than 1 year: `age-years`

This gives a quick-and-dirty way to look through things you've been putting
off or have forgotten about, and delete/do them as desired.

## Usage

Obtain an API token from `Settings` > `Integrations` > `Developer`, then:

```bash
go run main.go -api-token <token>
```

You can also specify the token in an `API_TOKEN` env var if you prefer.
