package main

import (
	"github.com/anilmisirlioglu/f1-telemetry-go/pkg/env/event"
	"github.com/anilmisirlioglu/f1-telemetry-go/pkg/packets"
	"github.com/anilmisirlioglu/f1-telemetry-go/pkg/telemetry"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func main() {
	// prometheus handler
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	client, err := telemetry.NewClientByCustomIpAddressAndPort("0.0.0.0", 20777)
	if err != nil {
		log.Fatal(err)
	}

	client.OnEventPacket(func(packet *packets.PacketEventData) {
		switch packet.EventCodeString() {
		case event.SpeedTrapTriggered:
			trap := packet.EventDetails.(*packets.SpeedTrap)
			if trap.VehicleIdx == packet.Header.PlayerCarIndex {
				log.Printf("Speed Trap: %f\n\n", trap.Speed)
				speedTrapMetric.Set(float64(trap.Speed))
			}
			break
		case event.FastestLap:
			fp := packet.EventDetails.(*packets.FastestLap)
			if fp.VehicleIdx == packet.Header.PlayerCarIndex {
				log.Printf("Fastest Lap: %f seconds\n", fp.LapTime)
				fastestLapMetric.Set(float64(fp.LapTime))
			}
			break
		}
	})

	client.OnCarTelemetryPacket(func(packet *packets.PacketCarTelemetryData) {
		car := packet.CarTelemetryData[packet.Header.PlayerCarIndex]
		speed := float64(car.Speed)
		engineRPM := float64(car.EngineRPM)
		log.Printf("Received Speed : %f \n", speed)
		speedMetric.Set(speed)
		engineRPMMetric.Set(engineRPM)

		for i, breaks := range []string{"rl", "rr", "fl", "fr"} {
			brakesTempMetric.WithLabelValues(breaks).Set(float64(car.BrakesTemperature[i]))
			log.Printf("Break temp : %f \n", float64(car.BrakesTemperature[i]))
		}
	})

	client.OnLapPacket(func(packet *packets.PacketLapData) {
		lapData := packet.LapData[packet.Header.PlayerCarIndex]
		lastLapTime := float64(lapData.LastLapTime)
		log.Printf("Last Lap : %f \n", lastLapTime)
		lastLapTimeMetric.Set(lastLapTime)
	})

	client.OnCarStatusPacket(func(packet *packets.PacketCarStatusData) {
		carStatus := packet.CarStatusData[packet.Header.PlayerCarIndex]
		tyresAgeLapsMetric.Set(float64(carStatus.TyresAgeLaps))

		for i, trye := range []string{"rl", "rr", "fl", "fr"} {
			tyreWearMetric.WithLabelValues(trye).Set(float64(carStatus.TyresWear[i]))
		}
	})

	client.Run()
}
