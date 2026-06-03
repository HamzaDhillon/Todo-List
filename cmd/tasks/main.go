package main

import (
	"fmt"
	"os"
	"strings"
	"todo/internal/todo"
	"strconv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Manage tasks from the terminal",
}

const dataFile = "data/tasks.csv"

func loadService() (*todo.Service, error) {
	svc := todo.NewService()
	if err := svc.LoadFromCSV(dataFile); err != nil {
		return nil, err
	}
	return svc, nil
}

func exitErr(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

var showAll bool
var showCompleted bool

var listCmd = &cobra.Command{
	Use:   "list [status]",
	Short: "List tasks",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mode := "pending"
		if showCompleted {
			mode = "completed"
		} else if len(args) == 1 {
			mode = strings.ToLower(args[0])
		}
		if showAll && mode == "pending" {
			mode = "all"
		}

		svc, err := loadService()
		if err != nil {
			exitErr(err)
		}

		var tasks []todo.Task
		switch mode {
		case "pending":
			tasks, err = svc.List(false)
		case "all":
			tasks, err = svc.List(true)
		case "completed":
			tasks, err = svc.List(true)
			tasks = onlyCompleted(tasks)
		default:
			exitErr(fmt.Errorf("invalid status %q: use pending, completed, or all", mode))
		}
		if err != nil {
			exitErr(err)
		}

		printTasks(tasks, mode == "all" || mode == "completed")
	},
}
var force bool
var completeCmd = &cobra.Command{
	Use:   "complete [id]",
	Short: "Complete task",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			exitErr(fmt.Errorf("please provide a task ID"))
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			exitErr(fmt.Errorf("invalid task ID: %v", err))
		}

		svc, err := loadService()
		_ = svc.Complete(id)

		if err != nil {
			exitErr(err)
		}
		_ = svc.Complete(1)
		println("Task completed successfully")
		_ = svc.SaveToCSV(dataFile)
	},
}

var del bool
var deleteCmd = &cobra.Command{
	Use: "delete [id]",
	Short: "Delete ID",
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string){
		if len(args) != 1 {
			exitErr(fmt.Errorf("Please provide an id"))
		} 
		id,err := strconv.Atoi(args[0])
		if err != nil{
			exitErr(fmt.Errorf("invalid task id %w",err))
		}
		svc,err := loadService()
		if err != nil {
			exitErr(err)
		}
	    err = svc.Delete(id)
		if err != nil{
			println("Task with such an id doesn't exist",)
		}else{
			println("Task deleted successfully",)
		}
		_ = svc.SaveToCSV(dataFile)
	},
}

var add bool
var addCmd = &cobra.Command{
	Use: "add[description]",
	Short: "Add a task to the list",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command,args []string)  {
		svc,err := loadService()
		if err != nil {
			exitErr(err)
		}
		_ = svc.Add(args[0])
		
	},
}

func onlyCompleted(tasks []todo.Task) []todo.Task {
	filtered := make([]todo.Task, 0, len(tasks))
	for _, t := range tasks {
		if t.IsComplete {
			filtered = append(filtered, t)
		}
	}
	return filtered
}


func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&showAll, "all", "a", false, "show completed tasks too")
	listCmd.Flags().BoolVarP(&showCompleted, "completed", "c", false, "show only completed tasks")

	rootCmd.AddCommand(completeCmd)
	completeCmd.Flags().BoolVarP(&force, "force", "f", false, "force completion even if task is already completed")

	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVarP(&del, "delete","d", false, "deleting a task given an id" )

}

func main() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
