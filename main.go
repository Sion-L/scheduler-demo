package main

import (
	"github.com/Sion-L/scheculer-demo/pkg/plugins"
	"github.com/spf13/pflag"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"
	_ "k8s.io/component-base/metrics/prometheus/clientgo"
	_ "k8s.io/component-base/metrics/prometheus/version"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// 把自定义的插件注册到整体的调度器中去
	command := app.NewSchedulerCommand(app.WithPlugin(plugins.Name, plugins.New))

	pflag.CommandLine.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)

	logs.InitLogs()
	defer logs.FlushLogs()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
