package plugins

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	frameworkruntime "k8s.io/kubernetes/pkg/scheduler/framework/runtime"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
)

// 插件名称定义
const (
	Name              = "sample-plugin"
	preFilterStateKey = "PreFilter" + Name
)

// sample实现PreFilterPlugin插件,实现了PreFilter和Filter扩展点
var _ framework.PreFilterPlugin = &Sample{}
var _ framework.FilterPlugin = &Sample{}

// 定义传递的参数
type SampleArgs struct {
	FavoriteColor  string `json:"favoriteColor,omitempty"`
	FavoriteNumber int    `json:"FavoriteNumber,omitempty"`
	ThanksTo       string `json:"ThanksTo,omitempty"`
}

type PreFilterState struct {
	framework.Resource // 获取request,limits
}

// 定义pluginFactory
type Sample struct {
	args    *SampleArgs
	handler framework.FrameworkHandle
}

// sample实现prefilter插件函数定义
func (s *Sample) Name() string {
	return Name
}

func (s *Sample) PreFilter(ctx context.Context, state *framework.CycleState, pod *v1.Pod) *framework.Status {
	klog.V(3).Infof("prefilter pod: %v", pod.Name)
	state.Write(preFilterStateKey, computePodResourceLimit(pod))
	return nil
}

func (s *Sample) PreFilterExtensions() framework.PreFilterExtensions {
	return nil
}

func (s *Sample) Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	preState, err := getPreFilterState(state)
	if err != nil {
		return framework.NewStatus(framework.Error, err.Error())
	}

	// logic 真正处理
	klog.V(3).Infof("Filter Pod: %v, node: %v, pre state: %v", pod.Name, nodeInfo.Node().Name, preState)
	return framework.NewStatus(framework.Success, "")
}

// new一个plugin,plugin返回Name方法，由sample实现过了
func New(object runtime.Object, f framework.FrameworkHandle) (framework.Plugin, error) {
	// 获取插件配置的参数
	args, err := getSampleArgs(object)
	if err != nil {
		return nil, err
	}
	// 校验参数是否正确
	klog.V(3).Infof("get plugin config args: %+v", args)
	return &Sample{
		args:    args,
		handler: f,
	}, nil
}

func getSampleArgs(object runtime.Object) (*SampleArgs, error) {
	sa := &SampleArgs{}
	if err := frameworkruntime.DecodeInto(object, sa); err != nil {
		return nil, err
	}
	return sa, nil
}

// 实现clone函数
//type StateData interface {
//	Clone() StateData
//}
func (p *PreFilterState) Clone() framework.StateData {
	return p
}

// 获取pod的resource值
func computePodResourceLimit(pod *v1.Pod) *PreFilterState {
	result := &PreFilterState{}
	for _, container := range pod.Spec.Containers {
		// 添加
		result.Add(container.Resources.Limits)
	}

	return result
}

//
func getPreFilterState(state *framework.CycleState) (*PreFilterState, error) {
	data, err := state.Read(preFilterStateKey)
	if err != nil {
		return nil, err
	}
	s, ok := data.(*PreFilterState)
	if !ok {
		return nil, fmt.Errorf("%+v convent to samplaePlugin prefilterstate error", data)
	}
	return s, nil
}
