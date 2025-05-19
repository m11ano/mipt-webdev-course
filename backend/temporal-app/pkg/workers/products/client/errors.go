package productstc

import "github.com/m11ano/e"

var ErrWorkflowCantStart = e.NewErrorFrom(e.ErrInternal).SetMessage("workflow can't start")
var ErrWorkflowResutError = e.NewErrorFrom(e.ErrInternal).SetMessage("workflow result error")
