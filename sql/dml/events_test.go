// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dml

import (
	"context"
	"testing"

	"github.com/corestoreio/pkg/util/assert"
)

func TestQueryOptions(t *testing.T) {
	ctx := WithContextQueryOptions(context.Background(), QueryOptions{
		SkipEvents:     true,
		SkipTimestamps: true,
		SkipRelations:  true,
	})
	assert.True(t, FromContextQueryOptions(ctx).SkipEvents)
}
