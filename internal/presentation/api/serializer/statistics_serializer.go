package serializer

import (
	"ContractAnalysis/internal/domain/repository"
	"ContractAnalysis/internal/presentation/api/dto"
)

// ToStatisticsResponse converts StrategyStatistics entity to StatisticsResponse DTO
func ToStatisticsResponse(stats *repository.StrategyStatistics) *dto.StatisticsResponse {
	resp := &dto.StatisticsResponse{
		StrategyName:               stats.StrategyName,
		Symbol:                     stats.Symbol,
		PeriodLabel:                stats.PeriodLabel,
		PeriodStart:                stats.PeriodStart.Format("2006-01-02T15:04:05Z"),
		PeriodEnd:                  stats.PeriodEnd.Format("2006-01-02T15:04:05Z"),
		TotalSignals:               stats.TotalSignals,
		ConfirmedSignals:           stats.ConfirmedSignals,
		InvalidatedSignals:         stats.InvalidatedSignals,
		ProfitableSignals:          stats.ProfitableSignals,
		LosingSignals:              stats.LosingSignals,
		NeutralSignals:             stats.NeutralSignals,
		TotalKlineHours:            stats.TotalKlineHours,
		ProfitableKlineHoursHigh:   stats.ProfitableKlineHoursHigh,
		ProfitableKlineHoursClose:  stats.ProfitableKlineHoursClose,
		CalculatedAt:               stats.CalculatedAt.Format("2006-01-02T15:04:05Z"),
	}

	// Convert decimal pointers to string pointers
	if stats.WinRate != nil {
		winRate := stats.WinRate.String()
		resp.WinRate = &winRate
	}

	if stats.AvgProfitPct != nil {
		avgProfit := stats.AvgProfitPct.String()
		resp.AvgProfitPct = &avgProfit
	}

	if stats.AvgLossPct != nil {
		avgLoss := stats.AvgLossPct.String()
		resp.AvgLossPct = &avgLoss
	}

	if stats.AvgHoldingHours != nil {
		avgHolding := stats.AvgHoldingHours.String()
		resp.AvgHoldingHours = &avgHolding
	}

	if stats.BestSignalPct != nil {
		best := stats.BestSignalPct.String()
		resp.BestSignalPct = &best
	}

	if stats.WorstSignalPct != nil {
		worst := stats.WorstSignalPct.String()
		resp.WorstSignalPct = &worst
	}

	if stats.ProfitFactor != nil {
		profitFactor := stats.ProfitFactor.String()
		resp.ProfitFactor = &profitFactor
	}

	if stats.KlineTheoreticalWinRate != nil {
		theoreticalWinRate := stats.KlineTheoreticalWinRate.String()
		resp.KlineTheoreticalWinRate = &theoreticalWinRate
	}

	if stats.KlineCloseWinRate != nil {
		closeWinRate := stats.KlineCloseWinRate.String()
		resp.KlineCloseWinRate = &closeWinRate
	}

	if stats.AvgHourlyReturnPct != nil {
		avgHourlyReturn := stats.AvgHourlyReturnPct.String()
		resp.AvgHourlyReturnPct = &avgHourlyReturn
	}

	if stats.MaxHourlyReturnPct != nil {
		maxHourlyReturn := stats.MaxHourlyReturnPct.String()
		resp.MaxHourlyReturnPct = &maxHourlyReturn
	}

	if stats.MinHourlyReturnPct != nil {
		minHourlyReturn := stats.MinHourlyReturnPct.String()
		resp.MinHourlyReturnPct = &minHourlyReturn
	}

	if stats.AvgMaxPotentialProfitPct != nil {
		avgMaxProfit := stats.AvgMaxPotentialProfitPct.String()
		resp.AvgMaxPotentialProfitPct = &avgMaxProfit
	}

	if stats.AvgMaxPotentialLossPct != nil {
		avgMaxLoss := stats.AvgMaxPotentialLossPct.String()
		resp.AvgMaxPotentialLossPct = &avgMaxLoss
	}

	return resp
}

// ToStatisticsListResponse converts a slice of StrategyStatistics entities
func ToStatisticsListResponse(statsList []*repository.StrategyStatistics) []*dto.StatisticsResponse {
	responses := make([]*dto.StatisticsResponse, 0, len(statsList))
	for _, stats := range statsList {
		responses = append(responses, ToStatisticsResponse(stats))
	}
	return responses
}
