package task

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/Ladicle/tabwriter"
	"golang.org/x/sync/errgroup"

	"github.com/vikbert/taskr/v3/internal/editors"
	"github.com/vikbert/taskr/v3/internal/fingerprint"
	"github.com/vikbert/taskr/v3/internal/logger"
	tasksort "github.com/vikbert/taskr/v3/internal/sort"
	"github.com/vikbert/taskr/v3/taskfile/ast"
)

// Constants for task listing
const (
	DefaultTaskGroup = "General"
	MinColumnWidth   = 4
	TabWidth         = 8
	TabPadding       = 2
)

// ListOptions collects list-related options
type ListOptions struct {
	ListOnlyTasksWithDescriptions bool
	ListAllTasks                  bool
	FormatTaskListAsJSON          bool
	NoStatus                      bool
	Nested                        bool
}

// NewListOptions creates a new ListOptions instance
func NewListOptions(list, listAll, listAsJson, noStatus, nested bool) ListOptions {
	return ListOptions{
		ListOnlyTasksWithDescriptions: list,
		ListAllTasks:                  listAll,
		FormatTaskListAsJSON:          listAsJson,
		NoStatus:                      noStatus,
		Nested:                        nested,
	}
}

// ShouldListTasks returns true if one of the options to list tasks has been set to true
func (o ListOptions) ShouldListTasks() bool {
	return o.ListOnlyTasksWithDescriptions || o.ListAllTasks
}

// Filters returns the slice of FilterFunc which filters a list
// of ast.Task according to the given ListOptions
func (o ListOptions) Filters() []FilterFunc {
	filters := []FilterFunc{FilterOutInternal}

	if o.ListOnlyTasksWithDescriptions {
		filters = append(filters, FilterOutNoDesc)
	}

	return filters
}

// ListTasks prints a list of tasks.
// Tasks that match the given filters will be excluded from the list.
// The function returns a boolean indicating whether tasks were found
// and an error if one was encountered while preparing the output.
func (e *Executor) ListTasks(o ListOptions) (bool, error) {
	tasks, err := e.GetTaskList(o.Filters()...)
	if err != nil {
		return false, fmt.Errorf("failed to get task list: %w", err)
	}

	if o.FormatTaskListAsJSON {
		return e.listTasksAsJSON(tasks, o)
	}
	return e.listTasksAsTable(tasks, o)
}

// listTasksAsJSON formats and outputs tasks as JSON
func (e *Executor) listTasksAsJSON(tasks []*ast.Task, o ListOptions) (bool, error) {
	output, err := e.ToEditorOutput(tasks, o.NoStatus, o.Nested)
	if err != nil {
		return false, fmt.Errorf("failed to generate editor output: %w", err)
	}

	encoder := json.NewEncoder(e.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(output); err != nil {
		return false, fmt.Errorf("failed to encode JSON output: %w", err)
	}

	return len(tasks) > 0, nil
}

// listTasksAsTable formats and outputs tasks as a formatted table
func (e *Executor) listTasksAsTable(tasks []*ast.Task, o ListOptions) (bool, error) {
	if len(tasks) == 0 {
		return e.printEmptyTaskListMessage(o), nil
	}

	// Display banner if enabled
	if e.Taskfile.Banner {
		e.Logger.PrintBannerWithProject(e.Taskfile.Project)
	}

	// Group and sort tasks
	grouper := newTaskGrouper(DefaultTaskGroup, e.Taskfile.Categories)
	groupedTasks := grouper.group(tasks)
	sortedGroups := grouper.sortedGroups(groupedTasks)

	// Calculate maximum task name length for alignment
	maxTaskNameLen := calculateMaxTaskNameLength(tasks)

	// Build and print the table
	builder := newTaskTableBuilder(e.Stdout, e.Logger, maxTaskNameLen)
	if err := builder.build(groupedTasks, sortedGroups); err != nil {
		return false, fmt.Errorf("failed to build task table: %w", err)
	}

	return true, nil
}

// printEmptyTaskListMessage prints appropriate message when no tasks are found
func (e *Executor) printEmptyTaskListMessage(o ListOptions) bool {
	if o.ListAllTasks {
		e.Logger.Outf(logger.Yellow, "task: No tasks available\n")
	} else {
		e.Logger.Outf(logger.Yellow, "task: No tasks with description available. Try --list-all to list all tasks\n")
	}
	return false
}

// calculateMaxTaskNameLength finds the longest task name for column alignment
func calculateMaxTaskNameLength(tasks []*ast.Task) int {
	maxLen := 0
	for _, task := range tasks {
		if len(task.Task) > maxLen {
			maxLen = len(task.Task)
		}
	}
	return maxLen
}

// taskGrouper handles task grouping logic
type taskGrouper struct {
	defaultGroup string
	categories   []string
}

// newTaskGrouper creates a new task grouper
func newTaskGrouper(defaultGroup string, categories []string) *taskGrouper {
	return &taskGrouper{defaultGroup: defaultGroup, categories: categories}
}

// group organizes tasks by their group field and sorts tasks within each group by line number
func (g *taskGrouper) group(tasks []*ast.Task) map[string][]*ast.Task {
	grouped := make(map[string][]*ast.Task, len(tasks)/2) // Estimate initial capacity

	for _, task := range tasks {
		group := task.Category
		if group == "" {
			group = g.defaultGroup
		}
		grouped[group] = append(grouped[group], task)
	}

	// Sort tasks within each group by index, then by line number
	for _, groupTasks := range grouped {
		sort.Slice(groupTasks, func(i, j int) bool {
			// Get index values, defaulting to 999 for tasks without index
			indexI := 999
			if groupTasks[i].Index.IsSet() {
				indexI = groupTasks[i].Index.Get()
			}
			indexJ := 999
			if groupTasks[j].Index.IsSet() {
				indexJ = groupTasks[j].Index.Get()
			}

			// First compare by index
			if indexI != indexJ {
				return indexI < indexJ
			}

			// If indices are equal, compare by line number
			lineI := 0
			if groupTasks[i].Location != nil {
				lineI = groupTasks[i].Location.Line
			}
			lineJ := 0
			if groupTasks[j].Location != nil {
				lineJ = groupTasks[j].Location.Line
			}
			return lineI < lineJ
		})
	}

	return grouped
}

// sortedGroups returns group names in the order specified by categories, with default group prioritized
func (g *taskGrouper) sortedGroups(grouped map[string][]*ast.Task) []string {
	var result []string

	// If categories are specified, use their order
	if len(g.categories) > 0 {
		// Add categories that exist in the grouped tasks
		for _, category := range g.categories {
			if _, exists := grouped[category]; exists {
				result = append(result, category)
			}
		}
		// Add any remaining groups not in categories, sorted alphabetically
		remainingGroups := make([]string, 0)
		for group := range grouped {
			found := false
			for _, category := range g.categories {
				if group == category {
					found = true
					break
				}
			}
			if !found {
				remainingGroups = append(remainingGroups, group)
			}
		}
		sort.Strings(remainingGroups)
		result = append(result, remainingGroups...)
	} else {
		// Default behavior: default group first, then alphabetical
		groups := make([]string, 0, len(grouped))
		hasDefault := false

		for group := range grouped {
			if group == g.defaultGroup {
				hasDefault = true
			} else {
				groups = append(groups, group)
			}
		}

		sort.Strings(groups)

		if hasDefault {
			result = append([]string{g.defaultGroup}, groups...)
		} else {
			result = groups
		}
	}

	return result
}

// taskTableBuilder builds formatted task tables
type taskTableBuilder struct {
	writer      io.Writer
	logger      *logger.Logger
	minWidth    int
	tabWidth    int
	padding     int
	showAliases bool
}

// newTaskTableBuilder creates a new table builder
func newTaskTableBuilder(w io.Writer, log *logger.Logger, maxTaskNameLen int) *taskTableBuilder {
	return &taskTableBuilder{
		writer:      w,
		logger:      log,
		minWidth:    maxTaskNameLen + MinColumnWidth,
		tabWidth:    TabWidth,
		padding:     TabPadding,
		showAliases: true,
	}
}

// build constructs and prints the complete task table
func (b *taskTableBuilder) build(groupedTasks map[string][]*ast.Task, groups []string) error {
	for _, group := range groups {
		if err := b.printGroup(group, groupedTasks[group]); err != nil {
			return err
		}
	}
	return nil
}

// printGroup prints a single group of tasks
func (b *taskTableBuilder) printGroup(groupName string, tasks []*ast.Task) error {
	// Print group header
	b.logger.Outf(logger.BoldYellow, "\n%s\n", strings.ToUpper(groupName))

	// Create tabwriter for aligned columns
	w := tabwriter.NewWriter(b.writer, b.minWidth, b.tabWidth, b.padding, ' ', 0)

	// Print each task in the group
	for _, task := range tasks {
		b.printTask(w, task)
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("failed to flush output for group %s: %w", groupName, err)
	}

	return nil
}

// printTask prints a single task row
func (b *taskTableBuilder) printTask(w io.Writer, task *ast.Task) {
	// Task name
	b.logger.FOutf(w, logger.Green, task.Task)

	// Task description (normalize newlines)
	desc := strings.ReplaceAll(task.Desc, "\n", " ")
	b.logger.FOutf(w, logger.Default, "\t%s", desc)

	// Task aliases (if any)
	if b.showAliases && len(task.Aliases) > 0 {
		aliasStr := strings.Join(task.Aliases, ", ")
		b.logger.FOutf(w, logger.Cyan, " (aliases: %s)", aliasStr)
	}

	fmt.Fprint(w, "\n")
}

// ListTaskNames prints only the task names in a Taskfile.
// Only tasks with a non-empty description are printed if allTasks is false.
// Otherwise, all task names are printed.
func (e *Executor) ListTaskNames(allTasks bool) error {
	w := e.Stdout
	if w == nil {
		w = os.Stdout
	}

	// Ensure task sorter is initialized
	if e.TaskSorter == nil {
		e.TaskSorter = tasksort.AlphaNumericWithRootTasksFirst
	}

	// Collect task names
	taskNames := e.collectTaskNames(allTasks)

	// Print each task name
	for _, name := range taskNames {
		if _, err := fmt.Fprintln(w, name); err != nil {
			return fmt.Errorf("failed to write task name: %w", err)
		}
	}

	return nil
}

// collectTaskNames extracts task names based on filter criteria
func (e *Executor) collectTaskNames(allTasks bool) []string {
	taskNames := make([]string, 0, e.Taskfile.Tasks.Len())

	for task := range e.Taskfile.Tasks.Values(e.TaskSorter) {
		// Filter internal tasks and tasks without description (if needed)
		if task.Internal || (!allTasks && task.Desc == "") {
			continue
		}

		// Add main task name
		taskNames = append(taskNames, strings.TrimRight(task.Task, ":"))

		// Add aliases
		for _, alias := range task.Aliases {
			taskNames = append(taskNames, strings.TrimRight(alias, ":"))
		}
	}

	return taskNames
}

// ToEditorOutput converts tasks to editor-friendly format with optional status checking
func (e *Executor) ToEditorOutput(tasks []*ast.Task, noStatus bool, nested bool) (*editors.Namespace, error) {
	editorTasks := make([]editors.Task, len(tasks))

	// Fast path: no status checking needed
	if noStatus {
		for i, task := range tasks {
			editorTasks[i] = editors.NewTask(task)
		}
		return e.buildNamespace(editorTasks, nested), nil
	}

	// Concurrent status checking with controlled parallelism
	if err := e.checkTasksStatusConcurrently(tasks, editorTasks); err != nil {
		return nil, err
	}

	return e.buildNamespace(editorTasks, nested), nil
}

// checkTasksStatusConcurrently checks task status in parallel with controlled concurrency
func (e *Executor) checkTasksStatusConcurrently(tasks []*ast.Task, editorTasks []editors.Task) error {
	var g errgroup.Group

	// Limit concurrent goroutines to avoid resource exhaustion
	g.SetLimit(runtime.NumCPU())

	for i := range tasks {
		i := i // Capture loop variable (though Go 1.22+ doesn't require this)
		task := tasks[i]

		g.Go(func() error {
			editorTask := editors.NewTask(task)

			upToDate, err := e.checkTaskStatus(task)
			if err != nil {
				return fmt.Errorf("failed to check status for task %s: %w", task.Task, err)
			}

			editorTask.UpToDate = &upToDate
			editorTasks[i] = editorTask
			return nil
		})
	}

	return g.Wait()
}

// checkTaskStatus determines if a single task is up-to-date
func (e *Executor) checkTaskStatus(task *ast.Task) (bool, error) {
	method := e.Taskfile.Method
	if task.Method != "" {
		method = task.Method
	}

	return fingerprint.IsTaskUpToDate(
		context.Background(),
		task,
		fingerprint.WithMethod(method),
		fingerprint.WithTempDir(e.TempDir.Fingerprint),
		fingerprint.WithDry(e.Dry),
		fingerprint.WithLogger(e.Logger),
	)
}

// buildNamespace constructs the namespace hierarchy for editor integration
func (e *Executor) buildNamespace(editorTasks []editors.Task, nested bool) *editors.Namespace {
	// Determine initial capacity
	var tasksLen int
	if !nested {
		tasksLen = len(editorTasks)
	}

	rootNamespace := &editors.Namespace{
		Tasks:    make([]editors.Task, tasksLen),
		Location: e.Taskfile.Location,
	}

	// Build namespace structure
	for i, task := range editorTasks {
		taskNamespacePath := strings.Split(task.Task, ast.NamespaceSeparator)

		if nested {
			rootNamespace.AddNamespace(taskNamespacePath, task)
		} else {
			rootNamespace.Tasks[i] = task
		}
	}

	return rootNamespace
}
