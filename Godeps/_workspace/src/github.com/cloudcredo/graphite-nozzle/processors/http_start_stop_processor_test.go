package processors_test

import (
	. "github.com/cloudcredo/graphite-nozzle/processors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/noaa/events"
)

var _ = Describe("HttpStartStopProcessor", func() {
	var (
		processor          *HttpStartStopProcessor
		event              *events.Envelope
		httpStartStopEvent *events.HttpStartStop
	)

	BeforeEach(func() {
		processor = NewHttpStartStopProcessor()

		startTimestamp := int64(1425881484152112140)
		stopTimestamp := int64(1425881484161498528)
		method := events.Method_GET
		uri := "api.10.244.0.34.xip.io/v2/info"
		statusCode := int32(200)
		peerType := events.PeerType_Client

		httpStartStopEvent = &events.HttpStartStop{
			StartTimestamp: &startTimestamp,
			StopTimestamp:  &stopTimestamp,
			Method:         &method,
			Uri:            &uri,
			StatusCode:     &statusCode,
			PeerType:       &peerType,
		}

		event = &events.Envelope{
			HttpStartStop: httpStartStopEvent,
		}
	})

	Describe("#Process", func() {
		It("returns a Metric for each of the ProcessHttpStartStop* methods", func() {
			processedMetrics := processor.Process(event)

			Expect(processedMetrics).To(HaveLen(4))
		})
	})

	Describe("#ProcessHttpStartStopResponseTime", func() {
		It("formats the Stat string to include the hostname", func() {
			metric := processor.ProcessHttpStartStopResponseTime(httpStartStopEvent)

			Expect(metric.Stat).To(Equal("http.responsetimes.api_10_244_0_34_xip_io"))
		})

		It("calculates the HTTP response time in milliseconds", func() {
			metric := processor.ProcessHttpStartStopResponseTime(httpStartStopEvent)

			Expect(metric.Value).To(Equal(int64(9)))
		})
	})

	Describe("#ProcessHttpStartStopStatusCodeCount", func() {
		Context("with a HTTP 200 status code", func() {
			It("formats the Stat string to include the hostname and the status code", func() {
				metric := processor.ProcessHttpStartStopStatusCodeCount(httpStartStopEvent)

				Expect(metric.Stat).To(Equal("http.statuscodes.api_10_244_0_34_xip_io.200"))
			})

			It("it does not increment the error counter by one", func() {
				metric := processor.ProcessHttpStartStopHttpErrorCount(httpStartStopEvent)
				Expect(metric.Stat).To(Equal("http.errors.api_10_244_0_34_xip_io"))
				Expect(metric.Value).To(Equal(int64(0)))
			})

			It("increments the requests counter by one", func() {
				metric := processor.ProcessHttpStartStopHttpRequestCount(httpStartStopEvent)
				Expect(metric.Stat).To(Equal("http.requests.api_10_244_0_34_xip_io"))
				Expect(metric.Value).To(Equal(int64(1)))
			})
		})

		Context("with a HTTP 404 status code", func() {
			It("formats the Stat string to include the hostname and the status code", func() {
				statusCode := int32(404)
				httpStartStopEvent.StatusCode = &statusCode
				metric := processor.ProcessHttpStartStopStatusCodeCount(httpStartStopEvent)

				Expect(metric.Stat).To(Equal("http.statuscodes.api_10_244_0_34_xip_io.404"))
			})

			It("increments the error counter by one", func() {
				statusCode := int32(404)
				httpStartStopEvent.StatusCode = &statusCode
				metric := processor.ProcessHttpStartStopHttpErrorCount(httpStartStopEvent)
				Expect(metric.Stat).To(Equal("http.errors.api_10_244_0_34_xip_io"))
				Expect(metric.Value).To(Equal(int64(1)))
			})

			It("increments the requests counter by one", func() {
				statusCode := int32(404)
				httpStartStopEvent.StatusCode = &statusCode
				metric := processor.ProcessHttpStartStopHttpRequestCount(httpStartStopEvent)
				Expect(metric.Stat).To(Equal("http.requests.api_10_244_0_34_xip_io"))
				Expect(metric.Value).To(Equal(int64(1)))
			})
		})

		Context("when PeerType == PeerType_Client", func() {
			It("sets the increment value for the CounterMetric to 1", func() {
				metric := processor.ProcessHttpStartStopStatusCodeCount(httpStartStopEvent)

				Expect(metric.Value).To(Equal(int64(1)))
			})
		})

		Context("when PeerType == PeerType_Server", func() {
			It("sets the increment value for the CounterMetric to 0", func() {
				peerType := events.PeerType_Server
				httpStartStopEvent.PeerType = &peerType
				metric := processor.ProcessHttpStartStopStatusCodeCount(httpStartStopEvent)

				Expect(metric.Value).To(Equal(int64(0)))
			})
		})

	})
})
