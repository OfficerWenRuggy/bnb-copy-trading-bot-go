package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// FixedCapitalConfig defines the fixed capital strategy settings
type FixedCapitalConfig struct {
	// Total capital allocated for the trading strategy
	TotalCapital float64
	// Percentage of capital to use per trade (0.01 to 1.0)
	RiskPercentage float64
	// Minimum capital required to open a position
	MinimumCapital float64
	// Maximum capital per single trade
	MaxCapitalPerTrade float64
	// Enable dynamic capital allocation based on win rate
	DynamicAllocation bool
	// Minimum winning rate to increase capital allocation
	MinWinRateForIncrease float64
	// Maximum winning rate threshold for allocation
	MaxWinRateThreshold float64
}

// TierProfit defines a single tier in the multi-tier take profit strategy
type TierProfit struct {
	// Profit percentage to trigger this tier (e.g., 0.5 = 0.5%)
	ProfitPercentage float64
	// Percentage of position to close at this tier (0.1 to 1.0)
	ClosePercentage float64
	// Whether this tier is enabled
	Enabled bool
}

// MultiTierConfig defines the multi-tier take profit configuration
type MultiTierConfig struct {
	// Array of profit tiers ordered by ascending profit percentage
	Tiers []TierProfit
	// Enable multi-tier take profit strategy
	Enabled bool
	// Close entire position if no tier is reached within max time
	CloseOnTimeout bool
	// Maximum time in minutes to hold a position
	MaxHoldTime int
	// Trailing stop loss trigger percentage
	TrailingStopPercentage float64
}

// RiskManagementConfig defines advanced risk management settings
type RiskManagementConfig struct {
	// Maximum percentage of portfolio to risk per trade
	MaxRiskPercentage float64
	// Maximum consecutive losing trades before pause
	MaxConsecutiveLosses int
	// Pause trading duration in minutes after max losses
	PauseDuration int
	// Maximum daily loss percentage allowed
	MaxDailyLossPercentage float64
	// Enable stop loss at percentage (e.g., 0.02 = 2% loss)
	StopLossPercentage float64
	// Enable break-even stop loss after reaching profit threshold
	BreakEvenStopEnabled bool
	// Profit percentage to trigger break-even stop
	BreakEvenThreshold float64
	// Maximum position size as percentage of total capital
	MaxPositionSize float64
	// Enable correlation check for multiple positions
	CorrelationCheckEnabled bool
	// Maximum correlation allowed between positions
	MaxCorrelationThreshold float64
	// Enable drawdown monitoring
	DrawdownMonitoringEnabled bool
	// Maximum allowed drawdown percentage
	MaxDrawdownPercentage float64
	// Enable equity protection
	EquityProtectionEnabled bool
	// Minimum equity level to stop trading
	MinimumEquityLevel float64
}

// TradingConfig defines core trading parameters
type TradingConfig struct {
	// Trading pair to monitor (e.g., "BNBUSDT")
	TradingPair string
	// Exchange API key
	APIKey string
	// Exchange API secret
	APISecret string
	// Enable testnet/sandbox trading
	TestnetEnabled bool
	// Minimum order quantity
	MinOrderQuantity float64
	// Maximum order quantity
	MaxOrderQuantity float64
	// Slippage tolerance percentage
	SlippageTolerance float64
	// Order timeout in seconds
	OrderTimeout int
	// Enable order validation before submission
	OrderValidationEnabled bool
	// Maker fee percentage
	MakerFee float64
	// Taker fee percentage
	TakerFee float64
}

// LoggingConfig defines logging configuration
type LoggingConfig struct {
	// Log level: DEBUG, INFO, WARN, ERROR
	LogLevel string
	// Log file path
	LogFilePath string
	// Enable console logging
	ConsoleLogging bool
	// Enable file logging
	FileLogging bool
	// Maximum log file size in MB
	MaxLogFileSize int
	// Number of backup log files to keep
	MaxBackupFiles int
}

// Config represents the complete bot configuration
type Config struct {
	FixedCapital    FixedCapitalConfig
	MultiTier       MultiTierConfig
	RiskManagement  RiskManagementConfig
	Trading         TradingConfig
	Logging         LoggingConfig
	// Refresh interval in seconds for market data
	RefreshInterval int
	// Enable dry run mode (no actual trades)
	DryRun bool
	// Notification webhook URL
	WebhookURL string
	// Enable notifications
	NotificationsEnabled bool
}

// LoadConfig loads configuration from environment variables and defaults
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	config := &Config{}

	// Load Fixed Capital Configuration
	config.FixedCapital = FixedCapitalConfig{
		TotalCapital:             getEnvFloat("FIXED_CAPITAL_TOTAL", 1000.0),
		RiskPercentage:           getEnvFloat("FIXED_CAPITAL_RISK_PERCENT", 0.05),
		MinimumCapital:           getEnvFloat("FIXED_CAPITAL_MINIMUM", 10.0),
		MaxCapitalPerTrade:       getEnvFloat("FIXED_CAPITAL_MAX_PER_TRADE", 500.0),
		DynamicAllocation:        getEnvBool("FIXED_CAPITAL_DYNAMIC_ALLOCATION", false),
		MinWinRateForIncrease:    getEnvFloat("FIXED_CAPITAL_MIN_WIN_RATE", 0.55),
		MaxWinRateThreshold:      getEnvFloat("FIXED_CAPITAL_MAX_WIN_RATE", 0.85),
	}

	// Load Multi-Tier Configuration
	config.MultiTier = MultiTierConfig{
		Enabled:                 getEnvBool("MULTI_TIER_ENABLED", true),
		CloseOnTimeout:          getEnvBool("MULTI_TIER_CLOSE_ON_TIMEOUT", true),
		MaxHoldTime:             getEnvInt("MULTI_TIER_MAX_HOLD_TIME", 240),
		TrailingStopPercentage:  getEnvFloat("MULTI_TIER_TRAILING_STOP", 0.5),
		Tiers: []TierProfit{
			{
				ProfitPercentage: 0.5,
				ClosePercentage:  0.2,
				Enabled:          true,
			},
			{
				ProfitPercentage: 1.0,
				ClosePercentage:  0.3,
				Enabled:          true,
			},
			{
				ProfitPercentage: 1.5,
				ClosePercentage:  0.25,
				Enabled:          true,
			},
			{
				ProfitPercentage: 2.0,
				ClosePercentage:  0.25,
				Enabled:          true,
			},
		},
	}

	// Load Risk Management Configuration
	config.RiskManagement = RiskManagementConfig{
		MaxRiskPercentage:          getEnvFloat("RISK_MAX_RISK_PERCENT", 0.02),
		MaxConsecutiveLosses:       getEnvInt("RISK_MAX_CONSECUTIVE_LOSSES", 5),
		PauseDuration:              getEnvInt("RISK_PAUSE_DURATION_MINUTES", 30),
		MaxDailyLossPercentage:     getEnvFloat("RISK_MAX_DAILY_LOSS_PERCENT", 0.05),
		StopLossPercentage:         getEnvFloat("RISK_STOP_LOSS_PERCENT", 0.03),
		BreakEvenStopEnabled:       getEnvBool("RISK_BREAK_EVEN_STOP_ENABLED", true),
		BreakEvenThreshold:         getEnvFloat("RISK_BREAK_EVEN_THRESHOLD", 0.5),
		MaxPositionSize:            getEnvFloat("RISK_MAX_POSITION_SIZE", 0.1),
		CorrelationCheckEnabled:    getEnvBool("RISK_CORRELATION_CHECK_ENABLED", true),
		MaxCorrelationThreshold:    getEnvFloat("RISK_MAX_CORRELATION_THRESHOLD", 0.8),
		DrawdownMonitoringEnabled:  getEnvBool("RISK_DRAWDOWN_MONITORING_ENABLED", true),
		MaxDrawdownPercentage:      getEnvFloat("RISK_MAX_DRAWDOWN_PERCENT", 0.15),
		EquityProtectionEnabled:    getEnvBool("RISK_EQUITY_PROTECTION_ENABLED", true),
		MinimumEquityLevel:         getEnvFloat("RISK_MINIMUM_EQUITY_LEVEL", 500.0),
	}

	// Load Trading Configuration
	config.Trading = TradingConfig{
		TradingPair:            getEnvString("TRADING_PAIR", "BNBUSDT"),
		APIKey:                 os.Getenv("API_KEY"),
		APISecret:              os.Getenv("API_SECRET"),
		TestnetEnabled:         getEnvBool("TRADING_TESTNET_ENABLED", false),
		MinOrderQuantity:       getEnvFloat("TRADING_MIN_ORDER_QUANTITY", 0.01),
		MaxOrderQuantity:       getEnvFloat("TRADING_MAX_ORDER_QUANTITY", 1000.0),
		SlippageTolerance:      getEnvFloat("TRADING_SLIPPAGE_TOLERANCE", 0.01),
		OrderTimeout:           getEnvInt("TRADING_ORDER_TIMEOUT_SECONDS", 30),
		OrderValidationEnabled: getEnvBool("TRADING_ORDER_VALIDATION_ENABLED", true),
		MakerFee:               getEnvFloat("TRADING_MAKER_FEE", 0.001),
		TakerFee:               getEnvFloat("TRADING_TAKER_FEE", 0.001),
	}

	// Load Logging Configuration
	config.Logging = LoggingConfig{
		LogLevel:       getEnvString("LOG_LEVEL", "INFO"),
		LogFilePath:    getEnvString("LOG_FILE_PATH", "./logs/bot.log"),
		ConsoleLogging: getEnvBool("LOG_CONSOLE_ENABLED", true),
		FileLogging:    getEnvBool("LOG_FILE_ENABLED", true),
		MaxLogFileSize: getEnvInt("LOG_MAX_FILE_SIZE_MB", 10),
		MaxBackupFiles: getEnvInt("LOG_MAX_BACKUP_FILES", 5),
	}

	// Load General Configuration
	config.RefreshInterval = getEnvInt("REFRESH_INTERVAL_SECONDS", 5)
	config.DryRun = getEnvBool("DRY_RUN_MODE", false)
	config.WebhookURL = os.Getenv("WEBHOOK_URL")
	config.NotificationsEnabled = getEnvBool("NOTIFICATIONS_ENABLED", true)

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate validates the configuration values
func (c *Config) Validate() error {
	// Validate Fixed Capital Configuration
	if c.FixedCapital.TotalCapital <= 0 {
		return fmt.Errorf("total capital must be positive, got %f", c.FixedCapital.TotalCapital)
	}
	if c.FixedCapital.RiskPercentage <= 0 || c.FixedCapital.RiskPercentage > 1 {
		return fmt.Errorf("risk percentage must be between 0 and 1, got %f", c.FixedCapital.RiskPercentage)
	}
	if c.FixedCapital.MinimumCapital <= 0 {
		return fmt.Errorf("minimum capital must be positive, got %f", c.FixedCapital.MinimumCapital)
	}
	if c.FixedCapital.MaxCapitalPerTrade <= 0 {
		return fmt.Errorf("max capital per trade must be positive, got %f", c.FixedCapital.MaxCapitalPerTrade)
	}
	if c.FixedCapital.MaxCapitalPerTrade > c.FixedCapital.TotalCapital {
		return fmt.Errorf("max capital per trade cannot exceed total capital")
	}
	if c.FixedCapital.MinWinRateForIncrease <= 0 || c.FixedCapital.MinWinRateForIncrease > 1 {
		return fmt.Errorf("min win rate must be between 0 and 1, got %f", c.FixedCapital.MinWinRateForIncrease)
	}
	if c.FixedCapital.MaxWinRateThreshold <= 0 || c.FixedCapital.MaxWinRateThreshold > 1 {
		return fmt.Errorf("max win rate must be between 0 and 1, got %f", c.FixedCapital.MaxWinRateThreshold)
	}

	// Validate Multi-Tier Configuration
	if c.MultiTier.Enabled {
		if len(c.MultiTier.Tiers) == 0 {
			return fmt.Errorf("at least one profit tier must be configured")
		}
		for i, tier := range c.MultiTier.Tiers {
			if tier.ProfitPercentage <= 0 {
				return fmt.Errorf("tier %d profit percentage must be positive, got %f", i, tier.ProfitPercentage)
			}
			if tier.ClosePercentage <= 0 || tier.ClosePercentage > 1 {
				return fmt.Errorf("tier %d close percentage must be between 0 and 1, got %f", i, tier.ClosePercentage)
			}
		}
		if c.MultiTier.MaxHoldTime <= 0 {
			return fmt.Errorf("max hold time must be positive, got %d", c.MultiTier.MaxHoldTime)
		}
		if c.MultiTier.TrailingStopPercentage < 0 {
			return fmt.Errorf("trailing stop percentage must be non-negative, got %f", c.MultiTier.TrailingStopPercentage)
		}
	}

	// Validate Risk Management Configuration
	if c.RiskManagement.MaxRiskPercentage <= 0 || c.RiskManagement.MaxRiskPercentage > 1 {
		return fmt.Errorf("max risk percentage must be between 0 and 1, got %f", c.RiskManagement.MaxRiskPercentage)
	}
	if c.RiskManagement.MaxConsecutiveLosses <= 0 {
		return fmt.Errorf("max consecutive losses must be positive, got %d", c.RiskManagement.MaxConsecutiveLosses)
	}
	if c.RiskManagement.PauseDuration <= 0 {
		return fmt.Errorf("pause duration must be positive, got %d", c.RiskManagement.PauseDuration)
	}
	if c.RiskManagement.MaxDailyLossPercentage <= 0 || c.RiskManagement.MaxDailyLossPercentage > 1 {
		return fmt.Errorf("max daily loss percentage must be between 0 and 1, got %f", c.RiskManagement.MaxDailyLossPercentage)
	}
	if c.RiskManagement.StopLossPercentage < 0 || c.RiskManagement.StopLossPercentage > 1 {
		return fmt.Errorf("stop loss percentage must be between 0 and 1, got %f", c.RiskManagement.StopLossPercentage)
	}
	if c.RiskManagement.BreakEvenThreshold < 0 {
		return fmt.Errorf("break-even threshold must be non-negative, got %f", c.RiskManagement.BreakEvenThreshold)
	}
	if c.RiskManagement.MaxPositionSize <= 0 || c.RiskManagement.MaxPositionSize > 1 {
		return fmt.Errorf("max position size must be between 0 and 1, got %f", c.RiskManagement.MaxPositionSize)
	}
	if c.RiskManagement.CorrelationCheckEnabled {
		if c.RiskManagement.MaxCorrelationThreshold < 0 || c.RiskManagement.MaxCorrelationThreshold > 1 {
			return fmt.Errorf("max correlation threshold must be between 0 and 1, got %f", c.RiskManagement.MaxCorrelationThreshold)
		}
	}
	if c.RiskManagement.DrawdownMonitoringEnabled {
		if c.RiskManagement.MaxDrawdownPercentage <= 0 || c.RiskManagement.MaxDrawdownPercentage > 1 {
			return fmt.Errorf("max drawdown percentage must be between 0 and 1, got %f", c.RiskManagement.MaxDrawdownPercentage)
		}
	}
	if c.RiskManagement.EquityProtectionEnabled {
		if c.RiskManagement.MinimumEquityLevel < 0 {
			return fmt.Errorf("minimum equity level must be non-negative, got %f", c.RiskManagement.MinimumEquityLevel)
		}
	}

	// Validate Trading Configuration
	if c.Trading.TradingPair == "" {
		return fmt.Errorf("trading pair must be specified")
	}
	if !c.Trading.TestnetEnabled && (c.Trading.APIKey == "" || c.Trading.APISecret == "") {
		return fmt.Errorf("API key and secret must be provided for live trading")
	}
	if c.Trading.MinOrderQuantity <= 0 {
		return fmt.Errorf("min order quantity must be positive, got %f", c.Trading.MinOrderQuantity)
	}
	if c.Trading.MaxOrderQuantity <= 0 {
		return fmt.Errorf("max order quantity must be positive, got %f", c.Trading.MaxOrderQuantity)
	}
	if c.Trading.MaxOrderQuantity < c.Trading.MinOrderQuantity {
		return fmt.Errorf("max order quantity cannot be less than min order quantity")
	}
	if c.Trading.SlippageTolerance < 0 || c.Trading.SlippageTolerance > 1 {
		return fmt.Errorf("slippage tolerance must be between 0 and 1, got %f", c.Trading.SlippageTolerance)
	}
	if c.Trading.OrderTimeout <= 0 {
		return fmt.Errorf("order timeout must be positive, got %d", c.Trading.OrderTimeout)
	}
	if c.Trading.MakerFee < 0 || c.Trading.MakerFee > 1 {
		return fmt.Errorf("maker fee must be between 0 and 1, got %f", c.Trading.MakerFee)
	}
	if c.Trading.TakerFee < 0 || c.Trading.TakerFee > 1 {
		return fmt.Errorf("taker fee must be between 0 and 1, got %f", c.Trading.TakerFee)
	}

	// Validate Logging Configuration
	if c.Logging.LogFilePath == "" && c.Logging.FileLogging {
		return fmt.Errorf("log file path must be specified when file logging is enabled")
	}
	if c.Logging.MaxLogFileSize <= 0 {
		return fmt.Errorf("max log file size must be positive, got %d", c.Logging.MaxLogFileSize)
	}
	if c.Logging.MaxBackupFiles < 0 {
		return fmt.Errorf("max backup files must be non-negative, got %d", c.Logging.MaxBackupFiles)
	}

	// Validate General Configuration
	if c.RefreshInterval <= 0 {
		return fmt.Errorf("refresh interval must be positive, got %d", c.RefreshInterval)
	}

	return nil
}

// CalculateRiskCapital calculates the capital to risk based on fixed capital configuration
func (c *Config) CalculateRiskCapital(currentEquity float64) float64 {
	return currentEquity * c.FixedCapital.RiskPercentage
}

// CalculatePositionSize calculates the position size based on risk parameters
func (c *Config) CalculatePositionSize(currentEquity float64, entryPrice float64, stopLossPrice float64) float64 {
	riskCapital := c.CalculateRiskCapital(currentEquity)
	priceDifference := entryPrice - stopLossPrice
	if priceDifference <= 0 {
		return 0
	}
	positionSize := riskCapital / priceDifference
	maxPositionValue := currentEquity * c.RiskManagement.MaxPositionSize
	maxPositionQuantity := maxPositionValue / entryPrice
	if positionSize > maxPositionQuantity {
		return maxPositionQuantity
	}
	return positionSize
}

// IsWithinDailyLossLimit checks if trading can continue based on daily loss limit
func (c *Config) IsWithinDailyLossLimit(startingEquity float64, currentEquity float64) bool {
	loss := startingEquity - currentEquity
	lossPercentage := loss / startingEquity
	return lossPercentage <= c.RiskManagement.MaxDailyLossPercentage
}

// IsWithinDrawdownLimit checks if current drawdown is within acceptable limits
func (c *Config) IsWithinDrawdownLimit(peakEquity float64, currentEquity float64) bool {
	if peakEquity <= 0 {
		return true
	}
	drawdown := (peakEquity - currentEquity) / peakEquity
	return drawdown <= c.RiskManagement.MaxDrawdownPercentage
}

// Helper functions for environment variable parsing

func getEnvString(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvFloat(key string, defaultValue float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Printf("Invalid float value for %s: %s, using default: %f\n", key, value, defaultValue)
		return defaultValue
	}
	return floatValue
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Invalid int value for %s: %s, using default: %d\n", key, value, defaultValue)
		return defaultValue
	}
	return intValue
}

func getEnvBool(key string, defaultValue bool) bool {
	value := strings.ToLower(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	return value == "true" || value == "1" || value == "yes"
}
