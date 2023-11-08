// Copyright © 2020 The gRPC Kit Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

func (t *templateService) fileDirectoryInternal() {
	t.files = append(t.files, &templateFile{
		name:  "internal/metric/metrics.go",
		parse: true,
		body: `
package metric

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var meter = otel.Meter("{{ .Global.Repository }}")

var (
	testHello metric.Float64Counter
)

func init() {
	var err error

	testHello, err = meter.Float64Counter(
		"test.hello",
		metric.WithDescription("a test metrics"),
	)
	if err != nil {
		panic(err)
	}
}

// Test 用于示例测试，提供给外部调用变更指标值
func Test(ctx context.Context, val float64) {
	testHello.Add(ctx, val)
}
`,
	})

	t.files = append(t.files, &templateFile{
		name:  "internal/trace/traces.go",
		parse: true,
		body: `
package trace

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("{{ .Global.Repository }}")

func Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return tracer.Start(ctx, name, opts...)
}
`,
	})
}
