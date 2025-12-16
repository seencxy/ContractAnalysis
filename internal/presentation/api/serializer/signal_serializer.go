package serializer

import (
	"ContractAnalysis/internal/domain/entity"
	"ContractAnalysis/internal/presentation/api/dto"
)

// ToSignalResponse converts a Signal entity to SignalResponse DTO
func ToSignalResponse(signal *entity.Signal) *dto.SignalResponse {
	return ToSignalResponseWithOutcome(signal, nil)
}

// ToSignalResponseWithOutcome converts a Signal entity and optional SignalOutcome to SignalResponse DTO
func ToSignalResponseWithOutcome(signal *entity.Signal, outcome *entity.SignalOutcome) *dto.SignalResponse {
	resp := &dto.SignalResponse{
		SignalID:           signal.SignalID,
		Symbol:             signal.Symbol,
		Type:               string(signal.Type),
		StrategyName:       signal.StrategyName,
		GeneratedAt:        signal.GeneratedAt.Format("2006-01-02T15:04:05Z"),
		PriceAtSignal:      signal.PriceAtSignal.String(),
		LongAccountRatio:   signal.LongAccountRatio.String(),
		ShortAccountRatio:  signal.ShortAccountRatio.String(),
		LongPositionRatio:  signal.LongPositionRatio.String(),
		ShortPositionRatio: signal.ShortPositionRatio.String(),
		LongTraderCount:    signal.LongTraderCount,
		ShortTraderCount:   signal.ShortTraderCount,
		Status:             string(signal.Status),
		IsConfirmed:        signal.IsConfirmed,
		Reason:             signal.Reason,
		StrategyContext:    signal.ConfigSnapshot,
		CreatedAt:          signal.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:          signal.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if signal.ConfirmedAt != nil {
		confirmedAt := signal.ConfirmedAt.Format("2006-01-02T15:04:05Z")
		resp.ConfirmedAt = &confirmedAt
	}

	// Add outcome data if available (for CLOSED signals)
	if outcome != nil {
		finalPnl := outcome.FinalPriceChangePct.String()
		resp.FinalPnlPct = &finalPnl
		resp.Outcome = &outcome.Outcome
		resp.TotalTrackingHours = &outcome.TotalTrackingHours

		// Add closed time
		closedAt := outcome.ClosedAt.Format("2006-01-02T15:04:05Z")
		resp.ClosedAt = &closedAt
	}

	return resp
}

// ToSignalListResponse converts a slice of Signal entities to a slice of SignalResponse DTOs
func ToSignalListResponse(signals []*entity.Signal) []*dto.SignalResponse {
	responses := make([]*dto.SignalResponse, 0, len(signals))
	for _, signal := range signals {
		responses = append(responses, ToSignalResponse(signal))
	}
	return responses
}

// ToSignalTrackingResponse converts a SignalTracking entity to SignalTrackingResponse DTO
func ToSignalTrackingResponse(tracking *entity.SignalTracking) *dto.SignalTrackingResponse {
	highestPrice := tracking.HighestPrice.String()
	lowestPrice := tracking.LowestPrice.String()
	highestPricePct := tracking.HighestPricePct.String()
	lowestPricePct := tracking.LowestPricePct.String()

	resp := &dto.SignalTrackingResponse{
		ID:                tracking.ID,
		SignalID:          tracking.SignalID,
		TrackedAt:         tracking.TrackedAt.Format("2006-01-02T15:04:05Z"),
		CurrentPrice:      tracking.CurrentPrice.String(),
		PriceChangePct:    tracking.PriceChangePct.String(),
		HighestPrice:      &highestPrice,
		LowestPrice:       &lowestPrice,
		HighestChangePct:  &highestPricePct,
		LowestChangePct:   &lowestPricePct,
		HoursTracked:      int(tracking.HoursElapsed.IntPart()),
		IsProfitTargetHit: false, // This info comes from SignalOutcome
		IsStopLossHit:     false, // This info comes from SignalOutcome
	}

	return resp
}

// ToSignalTrackingListResponse converts a slice of SignalTracking entities
func ToSignalTrackingListResponse(trackings []*entity.SignalTracking) []*dto.SignalTrackingResponse {
	responses := make([]*dto.SignalTrackingResponse, 0, len(trackings))
	for _, tracking := range trackings {
		responses = append(responses, ToSignalTrackingResponse(tracking))
	}
	return responses
}

// ToSignalKlineTrackingResponse converts a SignalKlineTracking entity to DTO
func ToSignalKlineTrackingResponse(kline *entity.SignalKlineTracking) *dto.SignalKlineTrackingResponse {
	hourlyReturn := kline.HourlyReturnPct.String()

	resp := &dto.SignalKlineTrackingResponse{
		ID:                  kline.ID,
		SignalID:            kline.SignalID,
		KlineOpenTime:       kline.KlineOpenTime.Format("2006-01-02T15:04:05Z"),
		KlineCloseTime:      kline.KlineCloseTime.Format("2006-01-02T15:04:05Z"),
		OpenPrice:           kline.OpenPrice.String(),
		HighPrice:           kline.HighPrice.String(),
		LowPrice:            kline.LowPrice.String(),
		ClosePrice:          kline.ClosePrice.String(),
		Volume:              kline.Volume.String(),
		OpenChangePct:       kline.OpenChangePct.String(),
		HighChangePct:       kline.HighChangePct.String(),
		LowChangePct:        kline.LowChangePct.String(),
		CloseChangePct:      kline.CloseChangePct.String(),
		HourlyReturnPct:     &hourlyReturn,
		IsProfitableAtHigh:  kline.IsProfitableAtHigh,
		IsProfitableAtClose: kline.IsProfitableAtClose,
	}

	return resp
}

// ToSignalKlineTrackingListResponse converts a slice of SignalKlineTracking entities
func ToSignalKlineTrackingListResponse(klines []*entity.SignalKlineTracking) []*dto.SignalKlineTrackingResponse {
	responses := make([]*dto.SignalKlineTrackingResponse, 0, len(klines))
	for _, kline := range klines {
		responses = append(responses, ToSignalKlineTrackingResponse(kline))
	}
	return responses
}
