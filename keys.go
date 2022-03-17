package vslparser

const (
	// KindRequest is a Kind string identifying Varnish Request object.
	KindRequest = "Request"
	// KindBeReq is a Kind string identifying Varnish BeReq object.
	KindBeReq = "BeReq"
	// KindSession is a Kind string identifying Varnish Session objects.
	KindSession = "Session"
)

const (
	// TagBegin is tag key informing about begin of tags section (start of
	// VXID).
	TagBegin = "Begin"
	// TagEnd is tag key informing about end of tags section (end of VXID).
	TagEnd = "End"
	// TagVSL is a tag key identifying any VSL API warning or error.
	TagVSL = "VSL"

	// TagReqURL is a tag key identifying request URL.
	TagReqURL = "ReqURL"
	// TagReqProtocol is a tag key identifying HTTP protocol version.
	TagReqProtocol = "ReqProtocol"
	// TagReqMethod is a tag key identifying HTTP request method verb.
	TagReqMethod = "ReqMethod"
	// TagRespStatus is a tag key identifying HTTP response status code.
	TagRespStatus = "RespStatus"

	// TagReqHeader is a tag informing that Request header was set.
	TagReqHeader = "ReqHeader"
	// TagReqUnset is a tag informing that Request header was unset.
	TagReqUnset = "ReqUnset"
	// TagRespHeader is a tag informing that Response header was set.
	TagRespHeader = "RespHeader"
	// TagRespUnset is a tag informing that Response header was unset.
	TagRespUnset = "RespUnset"
	// TagBeReqHeader is a tag informing that BeReq header was set.
	TagBeReqHeader = "BereqHeader"
	// TagBeReqUnset is a tag informing that BeReq header was unset.
	TagBeReqUnset = "BereqUnset"

	// TagBackendOpen is a tag informing that backend connection has been
	// open.
	TagBackendOpen = "BackendOpen"
	// TagFetchError is a tag informing about reason of backend fetch
	// operation failure.
	TagFetchError = "FetchError"
	// TagBerespStatus is a tag informing about backend (BeResp) response
	// status code.
	TagBerespStatus = "BerespStatus"

	// TagTimestamp is a tag containing timing information for the Varnish
	// worker thread.
	TagTimestamp = "Timestamp"
)

const (
	// TimestampReqEventStart is a Request-level timestamp which identifies
	// timestamp of request processing start.
	TimestampReqEventStart = "Start"
	// TimestampReqEventReq is a Request-level timestamp which identifies
	// timestamp of receiving a complete client request.
	TimestampReqEventReq = "Req"
	// TimestampReqEventReqBody  is a Request-level timestamp which
	// identifies timestamp when the client request body has been processed
	// (discarded, cached or passed to the backend).
	TimestampReqEventReqBody = "ReqBody"
	// TimestampReqEventWaitinglist  is a Request-level timestamp which
	// identifies timestamp when request came off waitinglist.
	TimestampReqEventWaitinglist = "Waitinglist"
	// TimestampReqEventFetch  is a Request-level timestamp which identifies
	// timestamp of fetch completion.
	TimestampReqEventFetch = "Fetch"
	// TimestampReqEventProcess  is a Request-level timestamp which
	// identifies timestamp when processing was finished. After this event,
	// the response is ready to be delivered.
	TimestampReqEventProcess = "Process"
	// TimestampReqEventResp  is a Request-level timestamp which identifies
	// timestamp when delivery of a response to a client finished.
	TimestampReqEventResp = "Resp"
	// TimestampReqEventRestart  is a Request-level timestamp which
	// identifies timestamp of request processing restart.
	TimestampReqEventRestart = "Restart"
)

// https://book.varnish-software.com/4.0/chapters/Examining_Varnish_Server_s_Output.html#transactions
const (
	// ReasonESI is a reason of VSL transaction (Request, BeReq, etc.) begin
	// which states that the transaction evaluates ESI (Edge Side Includes -
	// basically a scripting language).
	ReasonESI = "esi"
	// ReasonFetch is a reason of VSL transaction (Request, BeReq, etc.)
	// begin which states that transaction started to fetch data from
	// backend.
	ReasonFetch = "fetch"
	// ReasonRestart is a reason of VSL transaction (Request, BeReq, etc.)
	// begin which states that transaction started because of request
	// processing restart.
	ReasonRestart = "restart"
	// ReasonRxreq is a reason of VSL transaction (Request, BeReq, etc.)
	// begin which states that a new client request is the cause of
	// transaction start.
	ReasonRxreq = "rxreq"
)

const (
	// VSLStoreOverflow is a value of VSL tag informing that varnishlog is
	// not consuming VSL logs fast enough from the shared memory circular
	// buffer and that old (not yet consumed logs) are being overwritten by
	// new logs by Varnish.
	VSLStoreOverflow = "store overflow"
	// VSLFlush is a value of VSL tag informing that varnishlog has been
	// forced to immediately terminate log output. In contrast to
	// VSLStoreOverflow, this error happens if for example Varnish instance
	// dies or if the varnishlog is forced to immediately exit.
	VSLFlush = "flush"
)

// EndNoteSynth is a tag value of End tag in case something went wrong and the
// varnishlog is incomplete. In such a case, there is typically a VCL tag logged
// which provides further details.
const EndNoteSynth = "synth"
