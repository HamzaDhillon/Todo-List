package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"github.com/mergestat/timediff"
	"todo/internal/todo"
)
func printTasks(tasks []todo.Task, showAll bool) {

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	
	if showAll {
		fmt.Fprintln(w, "ID\tTask\tCreated\tDone")
	} else {
		fmt.Fprintln(w, "ID\tTask\tCreated")
	}
	for _, t := range tasks {
		created := timediff.TimeDiff(t.CreatedAt)
		if showAll {
			fmt.Fprintf(w, "%d\t%s\t%s\t%v\n",
				t.ID, t.Description, created, t.IsComplete)
		} else {
			fmt.Fprintf(w, "%d\t%s\t%s\n",
				t.ID, t.Description, created)
		}
	}
	w.Flush() // required — columns won't align without this
}