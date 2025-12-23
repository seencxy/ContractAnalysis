package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"ContractAnalysis/config"
	"ContractAnalysis/internal/domain/entity"
	"ContractAnalysis/internal/infrastructure/logger"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// Client wraps the Binance Futures API client
type Client struct {
	client     *futures.Client
	httpClient *http.Client
	baseURL    string
	apiKey     string
	apiSecret  string
	timeout    time.Duration
	logger     *logger.Logger
}

// NewClient creates a new Binance API client
func NewClient(cfg config.BinanceConfig) (*Client, error) {
	// Create Binance futures client
	futuresClient := futures.NewClient(cfg.APIKey, cfg.APISecret)

	// Set base URL if custom
	if cfg.APIURL != "" && cfg.APIURL != "https://fapi.binance.com" {
		futures.UseTestnet = false
		futuresClient.BaseURL = cfg.APIURL
	}

	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: cfg.Timeout,
	}

	client := &Client{
		client:     futuresClient,
		httpClient: httpClient,
		baseURL:    cfg.APIURL,
		apiKey:     cfg.APIKey,
		apiSecret:  cfg.APISecret,
		timeout:    cfg.Timeout,
		logger:     logger.WithComponent("binance-client"),
	}

	return client, nil
}

// GetAllUSDTFuturesPairs retrieves all USDT-margined futures trading pairs
func (c *Client) GetAllUSDTFuturesPairs(ctx context.Context) ([]string, error) {
	c.logger.Info("Fetching all USDT futures pairs")

	exchangeInfo, err := c.client.NewExchangeInfoService().Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchange info: %w", err)
	}

	var usdtPairs []string
	for _, symbol := range exchangeInfo.Symbols {
		if symbol.QuoteAsset == "USDT" && symbol.Status == "TRADING" {
			usdtPairs = append(usdtPairs, symbol.Symbol)
		}
	}

	c.logger.Info("Fetched USDT futures pairs",
		zap.Int("count", len(usdtPairs)),
	)

	return usdtPairs, nil
}

// GetGlobalLongShortRatio retrieves global long/short account ratio
func (c *Client) GetGlobalLongShortRatio(ctx context.Context, symbol string, period string) (*GlobalLongShortAccountRatio, error) {
	endpoint := fmt.Sprintf("%s/futures/data/globalLongShortAccountRatio", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("period", period) // 5m, 15m, 30m, 1h, 2h, 4h, 6h, 12h, 1d
	q.Add("limit", "1")     // Get only the latest
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var ratios []GlobalLongShortAccountRatio
	if err := json.NewDecoder(resp.Body).Decode(&ratios); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(ratios) == 0 {
		return nil, fmt.Errorf("no data returned for symbol %s", symbol)
	}

	return &ratios[0], nil
}

// GetTopLongShortPositionRatio retrieves top trader long/short position ratio
func (c *Client) GetTopLongShortPositionRatio(ctx context.Context, symbol string, period string) (*TopLongShortPositionRatio, error) {
	endpoint := fmt.Sprintf("%s/futures/data/topLongShortPositionRatio", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("period", period)
	q.Add("limit", "1")
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var ratios []TopLongShortPositionRatio
	if err := json.NewDecoder(resp.Body).Decode(&ratios); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(ratios) == 0 {
		return nil, fmt.Errorf("no data returned for symbol %s", symbol)
	}

	return &ratios[0], nil
}

// GetTopLongShortAccountRatio retrieves top trader long/short account ratio
func (c *Client) GetTopLongShortAccountRatio(ctx context.Context, symbol string, period string) (*TopLongShortAccountRatio, error) {
	endpoint := fmt.Sprintf("%s/futures/data/topLongShortAccountRatio", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("period", period)
	q.Add("limit", "1")
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var ratios []TopLongShortAccountRatio
	if err := json.NewDecoder(resp.Body).Decode(&ratios); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(ratios) == 0 {
		return nil, fmt.Errorf("no data returned for symbol %s", symbol)
	}

	return &ratios[0], nil
}

// GetOpenInterest retrieves open interest for a symbol
func (c *Client) GetOpenInterest(ctx context.Context, symbol string) (*OpenInterest, error) {
	endpoint := fmt.Sprintf("%s/futures/data/openInterestHist", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("period", "5m") // Get latest 5m period
	q.Add("limit", "1")
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var interests []OpenInterest
	if err := json.NewDecoder(resp.Body).Decode(&interests); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(interests) == 0 {
		return nil, fmt.Errorf("no open interest data for symbol %s", symbol)
	}

	return &interests[0], nil
}

// GetFundingRate retrieves the current funding rate for a symbol
func (c *Client) GetFundingRate(ctx context.Context, symbol string) (*FundingRate, error) {
	endpoint := fmt.Sprintf("%s/fapi/v1/premiumIndex", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	q.Add("symbol", symbol)
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var fundingRate FundingRate
	if err := json.NewDecoder(resp.Body).Decode(&fundingRate); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &fundingRate, nil
}

// GetPrice retrieves the current price for a symbol
func (c *Client) GetPrice(ctx context.Context, symbol string) (float64, error) {
	prices, err := c.client.NewListPricesService().Symbol(symbol).Do(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get price: %w", err)
	}

	if len(prices) == 0 {
		return 0, fmt.Errorf("no price data for symbol %s", symbol)
	}

	price := 0.0
	fmt.Sscanf(prices[0].Price, "%f", &price)

	return price, nil
}

// Get24hrTicker retrieves 24-hour ticker statistics
func (c *Client) Get24hrTicker(ctx context.Context, symbol string) (*Ticker24hr, error) {
	tickers, err := c.client.NewListPriceChangeStatsService().Symbol(symbol).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get 24hr ticker: %w", err)
	}

	if len(tickers) == 0 {
		return nil, fmt.Errorf("no ticker data for symbol %s", symbol)
	}

	ticker := tickers[0]

	// Parse values
	var lastPrice, volume, quoteVolume float64
	fmt.Sscanf(ticker.LastPrice, "%f", &lastPrice)
	fmt.Sscanf(ticker.Volume, "%f", &volume)
	fmt.Sscanf(ticker.QuoteVolume, "%f", &quoteVolume)

	return &Ticker24hr{
		Symbol:      ticker.Symbol,
		LastPrice:   lastPrice,
		Volume:      volume,
		QuoteVolume: quoteVolume,
		OpenTime:    ticker.OpenTime,
		CloseTime:   ticker.CloseTime,
		Count:       ticker.Count,
	}, nil
}

// GetMarketData retrieves comprehensive market data for a symbol
func (c *Client) GetMarketData(ctx context.Context, symbol string) (*MarketData, error) {
	c.logger.Debug("Fetching market data", zap.String("symbol", symbol))

	now := time.Now()

	// Fetch long/short account ratio (global)
	accountRatio, err := c.GetGlobalLongShortRatio(ctx, symbol, "5m")
	if err != nil {
		return nil, fmt.Errorf("failed to get account ratio: %w", err)
	}

	// Fetch top trader position ratio (optional - some pairs may not have this data)
	var longPositionPct, shortPositionPct float64
	var positionRatioAvailable bool = true

	positionRatio, err := c.GetTopLongShortPositionRatio(ctx, symbol, "5m")
	if err != nil {
		// Log warning but continue - position ratio is optional
		c.logger.Warn("Position ratio not available for symbol",
			zap.String("symbol", symbol),
			zap.Error(err),
		)
		longPositionPct = 0
		shortPositionPct = 0
		positionRatioAvailable = false
	} else {
		// Convert position ratios from 0-1 to percentages 0-100
		longPositionPct = positionRatio.LongAccount * 100
		shortPositionPct = positionRatio.ShortAccount * 100
	}

	// Fetch current price and volume
	ticker, err := c.Get24hrTicker(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticker: %w", err)
	}

	// Fetch open interest (optional)
	var openInterest float64
	oi, err := c.GetOpenInterest(ctx, symbol)
	if err != nil {
		c.logger.Debug("Open interest not available", zap.String("symbol", symbol), zap.Error(err))
		openInterest = 0
	} else {
		openInterest = oi.Value // Use Value (in USDT) or OpenInterest (in coins)? Usually Value is more comparable.
		// Wait, user requirement says "Open Interest, OI".
		// binance API: sumOpenInterest (coins), sumOpenInterestValue (USDT).
		// Let's use USDT value as it's more standard across pairs.
	}

	// Fetch funding rate
	var fundingRate float64
	fr, err := c.GetFundingRate(ctx, symbol)
	if err != nil {
		c.logger.Debug("Funding rate not available", zap.String("symbol", symbol), zap.Error(err))
		fundingRate = 0
	} else {
		fundingRate = fr.FundingRate
	}

	// Convert account ratios from 0-1 to percentages 0-100
	longAccountPct := accountRatio.LongAccount * 100
	shortAccountPct := accountRatio.ShortAccount * 100

	// Calculate data quality score
	dataQualityScore := 100
	if !positionRatioAvailable {
		dataQualityScore = 80 // Deduct 20 points for missing position data
	}

	marketData := &MarketData{
		Symbol:                 symbol,
		Timestamp:              now,
		LongAccountRatio:       longAccountPct,
		ShortAccountRatio:      shortAccountPct,
		LongPositionRatio:      longPositionPct,
		ShortPositionRatio:     shortPositionPct,
		PositionRatioAvailable: positionRatioAvailable,
		DataQualityScore:       dataQualityScore,
		Price:                  ticker.LastPrice,
		Volume24h:              ticker.QuoteVolume,
		OpenInterest:           openInterest,
		FundingRate:            fundingRate,
	}

	c.logger.Debug("Fetched market data successfully",
		zap.String("symbol", symbol),
	)

	return marketData, nil
}

// GetMarketDataBatch retrieves market data for multiple symbols
func (c *Client) GetMarketDataBatch(ctx context.Context, symbols []string) ([]*MarketData, error) {
	c.logger.Info("Fetching market data batch",
		zap.Int("count", len(symbols)),
	)

	var results []*MarketData
	var errors []error

	for _, symbol := range symbols {
		data, err := c.GetMarketData(ctx, symbol)
		if err != nil {
			c.logger.WithError(err).WithSymbol(symbol).Warn("Failed to fetch market data for symbol")
			errors = append(errors, fmt.Errorf("%s: %w", symbol, err))
			continue
		}

		results = append(results, data)

		// Small delay to avoid rate limiting
		time.Sleep(100 * time.Millisecond)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("failed to fetch data for all symbols: %v", errors)
	}

	c.logger.Info("Fetched market data batch",
		zap.Int("success", len(results)),
		zap.Int("failed", len(errors)),
	)

	return results, nil
}

// GetKlines retrieves kline/candlestick data for a symbol
// interval: "1m", "5m", "15m", "1h", "4h", "1d", etc.
// limit: maximum 1000 (default 500)
func (c *Client) GetKlines(ctx context.Context, symbol string, interval string, limit int) ([]*entity.Kline, error) {
	c.logger.Debug("Fetching klines",
		zap.String("symbol", symbol),
		zap.String("interval", interval),
		zap.Int("limit", limit),
	)

	// Validate limit
	if limit <= 0 || limit > 1000 {
		limit = 500
	}

	klines, err := c.client.NewKlinesService().
		Symbol(symbol).
		Interval(interval).
		Limit(limit).
		Do(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get klines: %w", err)
	}

	// Convert to entity.Kline
	result := make([]*entity.Kline, len(klines))
	for i, k := range klines {
		kline := convertToKline(k)
		result[i] = kline
	}

	c.logger.Debug("Fetched klines successfully",
		zap.String("symbol", symbol),
		zap.Int("count", len(result)),
	)

	return result, nil
}

// GetKlinesSince retrieves kline data since a specific time
func (c *Client) GetKlinesSince(ctx context.Context, symbol string, interval string, startTime time.Time) ([]*entity.Kline, error) {
	c.logger.Debug("Fetching klines since",
		zap.String("symbol", symbol),
		zap.String("interval", interval),
		zap.Time("start_time", startTime),
	)

	klines, err := c.client.NewKlinesService().
		Symbol(symbol).
		Interval(interval).
		StartTime(startTime.UnixMilli()).
		Limit(1000). // Max limit
		Do(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get klines since %v: %w", startTime, err)
	}

	// Convert to entity.Kline
	result := make([]*entity.Kline, len(klines))
	for i, k := range klines {
		kline := convertToKline(k)
		result[i] = kline
	}

	c.logger.Debug("Fetched klines since successfully",
		zap.String("symbol", symbol),
		zap.Int("count", len(result)),
	)

	return result, nil
}

// convertToKline converts a Binance API kline to entity.Kline
func convertToKline(k *futures.Kline) *entity.Kline {
	// Parse decimal values
	open, _ := decimal.NewFromString(k.Open)
	high, _ := decimal.NewFromString(k.High)
	low, _ := decimal.NewFromString(k.Low)
	close, _ := decimal.NewFromString(k.Close)
	volume, _ := decimal.NewFromString(k.Volume)
	quoteVolume, _ := decimal.NewFromString(k.QuoteAssetVolume)

	return &entity.Kline{
		OpenTime:    time.Unix(0, k.OpenTime*int64(time.Millisecond)),
		CloseTime:   time.Unix(0, k.CloseTime*int64(time.Millisecond)),
		Open:        open,
		High:        high,
		Low:         low,
		Close:       close,
		Volume:      volume,
		QuoteVolume: quoteVolume,
	}
}
