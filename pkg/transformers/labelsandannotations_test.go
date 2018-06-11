/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package transformers

import (
	"reflect"
	"testing"

	"github.com/kubernetes-sigs/kustomize/pkg/resmap"
	"github.com/kubernetes-sigs/kustomize/pkg/resource"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var service = schema.GroupVersionKind{Version: "v1", Kind: "Service"}
var secret = schema.GroupVersionKind{Version: "v1", Kind: "Secret"}
var cmap = schema.GroupVersionKind{Version: "v1", Kind: "ConfigMap"}
var ns = schema.GroupVersionKind{Version: "v1", Kind: "Namespace"}
var deploy = schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}
var statefulset = schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "StatefulSet"}
var foo = schema.GroupVersionKind{Group: "example.com", Version: "v1", Kind: "Foo"}
var crd = schema.GroupVersionKind{Group: "apiwctensions.k8s.io", Version: "v1beta1", Kind: "CustomResourceDefinition"}

func TestLabelsRun(t *testing.T) {
	m := resmap.ResMap{
		resource.NewResId(cmap, "cm1"): resource.NewResourceFromMap(
			map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "ConfigMap",
				"metadata": map[string]interface{}{
					"name": "cm1",
				},
			}),
		resource.NewResId(deploy, "deploy1"): resource.NewResourceFromMap(
			map[string]interface{}{
				"group":      "apps",
				"apiVersion": "v1",
				"kind":       "Deployment",
				"metadata": map[string]interface{}{
					"name": "deploy1",
				},
				"spec": map[string]interface{}{
					"template": map[string]interface{}{
						"metadata": map[string]interface{}{
							"labels": map[string]interface{}{
								"old-label": "old-value",
							},
						},
						"spec": map[string]interface{}{
							"containers": []interface{}{
								map[string]interface{}{
									"name":  "nginx",
									"image": "nginx:1.7.9",
								},
							},
						},
					},
				},
			}),
		resource.NewResId(service, "svc1"): resource.NewResourceFromMap(
			map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Service",
				"metadata": map[string]interface{}{
					"name": "svc1",
				},
				"spec": map[string]interface{}{
					"ports": []interface{}{
						map[string]interface{}{
							"name": "port1",
							"port": "12345",
						},
					},
				},
			}),
	}
	expected := resmap.ResMap{
		resource.NewResId(cmap, "cm1"): resource.NewResourceFromMap(
			map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "ConfigMap",
				"metadata": map[string]interface{}{
					"name": "cm1",
					"labels": map[string]interface{}{
						"label-key1": "label-value1",
						"label-key2": "label-value2",
					},
				},
			}),
		resource.NewResId(deploy, "deploy1"): resource.NewResourceFromMap(
			map[string]interface{}{
				"group":      "apps",
				"apiVersion": "v1",
				"kind":       "Deployment",
				"metadata": map[string]interface{}{
					"name": "deploy1",
					"labels": map[string]interface{}{
						"label-key1": "label-value1",
						"label-key2": "label-value2",
					},
				},
				"spec": map[string]interface{}{
					"selector": map[string]interface{}{
						"matchLabels": map[string]interface{}{
							"label-key1": "label-value1",
							"label-key2": "label-value2",
						},
					},
					"template": map[string]interface{}{
						"metadata": map[string]interface{}{
							"labels": map[string]interface{}{
								"old-label":  "old-value",
								"label-key1": "label-value1",
								"label-key2": "label-value2",
							},
						},
						"spec": map[string]interface{}{
							"containers": []interface{}{
								map[string]interface{}{
									"name":  "nginx",
									"image": "nginx:1.7.9",
								},
							},
						},
					},
				},
			}),
		resource.NewResId(service, "svc1"): resource.NewResourceFromMap(
			map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Service",
				"metadata": map[string]interface{}{
					"name": "svc1",
					"labels": map[string]interface{}{
						"label-key1": "label-value1",
						"label-key2": "label-value2",
					},
				},
				"spec": map[string]interface{}{
					"ports": []interface{}{
						map[string]interface{}{
							"name": "port1",
							"port": "12345",
						},
					},
					"selector": map[string]interface{}{
						"label-key1": "label-value1",
						"label-key2": "label-value2",
					},
				},
			}),
	}

	lt, err := NewDefaultingLabelsMapTransformer(map[string]string{"label-key1": "label-value1", "label-key2": "label-value2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = lt.Transform(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(m, expected) {
		err = expected.ErrorIfNotEqual(m)
		t.Fatalf("actual doesn't match expected: %v", err)
	}
}

func TestAnnotationsRun(t *testing.T) {
	m := resmap.ResMap{
		resource.NewResId(cmap, "cm1"): resource.NewResourceFromMap(
			map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "ConfigMap",
				"metadata": map[string]interface{}{
					"name": "cm1",
				},
			}),
		resource.NewResId(deploy, "deploy1"): resource.NewResourceFromMap(
			map[string]interface{}{
				"group":      "apps",
				"apiVersion": "v1",
				"kind":       "Deployment",
				"metadata": map[string]interface{}{
					"name": "deploy1",
				},
				"spec": map[string]interface{}{
					"template": map[string]interface{}{
						"metadata": map[string]interface{}{
							"labels": map[string]interface{}{
								"old-label": "old-value",
							},
						},
						"spec": map[string]interface{}{
							"containers": []interface{}{
								map[string]interface{}{
									"name":  "nginx",
									"image": "nginx:1.7.9",
								},
							},
						},
					},
				},
			}),
		resource.NewResId(service, "svc1"): resource.NewResourceFromMap(
			map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Service",
				"metadata": map[string]interface{}{
					"name": "svc1",
				},
				"spec": map[string]interface{}{
					"ports": []interface{}{
						map[string]interface{}{
							"name": "port1",
							"port": "12345",
						},
					},
				},
			}),
	}
	expected := resmap.ResMap{
		resource.NewResId(cmap, "cm1"): resource.NewResourceFromMap(
			map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "ConfigMap",
				"metadata": map[string]interface{}{
					"name": "cm1",
					"annotations": map[string]interface{}{
						"anno-key1": "anno-value1",
						"anno-key2": "anno-value2",
					},
				},
			}),
		resource.NewResId(deploy, "deploy1"): resource.NewResourceFromMap(
			map[string]interface{}{
				"group":      "apps",
				"apiVersion": "v1",
				"kind":       "Deployment",
				"metadata": map[string]interface{}{
					"name": "deploy1",
					"annotations": map[string]interface{}{
						"anno-key1": "anno-value1",
						"anno-key2": "anno-value2",
					},
				},
				"spec": map[string]interface{}{
					"template": map[string]interface{}{
						"metadata": map[string]interface{}{
							"annotations": map[string]interface{}{
								"anno-key1": "anno-value1",
								"anno-key2": "anno-value2",
							},
							"labels": map[string]interface{}{
								"old-label": "old-value",
							},
						},
						"spec": map[string]interface{}{
							"containers": []interface{}{
								map[string]interface{}{
									"name":  "nginx",
									"image": "nginx:1.7.9",
								},
							},
						},
					},
				},
			}),
		resource.NewResId(service, "svc1"): resource.NewResourceFromMap(
			map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Service",
				"metadata": map[string]interface{}{
					"name": "svc1",
					"annotations": map[string]interface{}{
						"anno-key1": "anno-value1",
						"anno-key2": "anno-value2",
					},
				},
				"spec": map[string]interface{}{
					"ports": []interface{}{
						map[string]interface{}{
							"name": "port1",
							"port": "12345",
						},
					},
				},
			}),
	}
	at, err := NewDefaultingAnnotationsMapTransformer(map[string]string{"anno-key1": "anno-value1", "anno-key2": "anno-value2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = at.Transform(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(m, expected) {
		err = expected.ErrorIfNotEqual(m)
		t.Fatalf("actual doesn't match expected: %v", err)
	}
}